package admincourse

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
)

type Service struct {
	courses adminrepo.AdminCourseRepository
	db      *gorm.DB
	rdb     *redis.Client
}

func New(courses adminrepo.AdminCourseRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{courses: courses, db: db, rdb: rdb}
}

type courseRow struct {
	ID            int            `gorm:"column:id;primaryKey"`
	Title         string         `gorm:"column:title"`
	Thumb         string         `gorm:"column:thumb"`
	Desc          string         `gorm:"column:desc"`
	Type          int            `gorm:"column:type"`
	Content       *string        `gorm:"column:content"`
	IsShow        int            `gorm:"column:is_show"`
	CategoryID    int            `gorm:"column:category_id"`
	IsRec         int            `gorm:"column:is_rec"`
	MinNeedCredit int            `gorm:"column:min_need_credit"`
	Credit        int            `gorm:"column:credit"`
	PaperID       *int           `gorm:"column:paper_id"`
	CertID        *int           `gorm:"column:cert_id"`
	CreatedID     int            `gorm:"column:created_id"`
	CreatedAt     *time.Time     `gorm:"column:created_at"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (courseRow) TableName() string { return "courses" }

type attachRow struct {
	ID            int64      `gorm:"column:id;primaryKey"`
	CourseID      int        `gorm:"column:course_id"`
	Name          string     `gorm:"column:name"`
	Path          string     `gorm:"column:path"`
	Size          int        `gorm:"column:size"`
	DownloadTimes int        `gorm:"column:download_times"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

func (attachRow) TableName() string { return "course_attach" }

type CertSaveMeta struct {
	ID               int  `json:"id"`
	CertChanged      bool `json:"cert_changed"`
	CertSwapped      bool `json:"cert_swapped"`
	PendingCertUsers int  `json:"pending_cert_users"`
	IssuedCertUsers  int  `json:"issued_cert_users"`
}

func (s *Service) GetCourseList(ctx context.Context, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	sort := asString(raw["sort"])
	if sort == "" || !allowedSort[sort] {
		sort = "created_at"
	}
	order := strings.ToLower(asString(raw["order"]))
	if order != "asc" {
		order = "desc"
	}

	query := s.db.WithContext(ctx).Model(&courseRow{})
	if id := asInt(raw["id"]); id > 0 {
		query = query.Where("id = ?", id)
	}
	if title := strings.TrimSpace(asString(raw["title"])); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if cid := asInt(raw["category_id"]); cid > 0 {
		query = query.Where("category_id = ?", cid)
	}
	if typ := asInt(raw["type"]); typ > 0 {
		query = query.Where("type = ?", typ)
	}
	if _, ok := raw["is_show"]; ok {
		if asString(raw["is_show"]) != "" {
			query = query.Where("is_show = ?", asInt(raw["is_show"]))
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []courseRow
	if err := query.Order(sort + " " + order).Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}

	ids := make([]int, 0, len(rows))
	cateIDs := map[int]struct{}{}
	certIDs := map[int]struct{}{}
	for _, row := range rows {
		ids = append(ids, row.ID)
		if row.CategoryID > 0 {
			cateIDs[row.CategoryID] = struct{}{}
		}
		if row.CertID != nil && *row.CertID > 0 {
			certIDs[*row.CertID] = struct{}{}
		}
	}
	categories := map[int]map[string]any{}
	if len(cateIDs) > 0 {
		cats, err := s.courses.ListCategoriesByIDs(ctx, keys(cateIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, cat := range cats {
			categories[cat.ID] = map[string]any{"id": cat.ID, "name": cat.Name}
		}
	}
	certs := map[int]map[string]any{}
	if len(certIDs) > 0 {
		list, err := s.courses.ListCertNamesByIDs(ctx, keys(certIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, item := range list {
			certs[item.ID] = map[string]any{"id": item.ID, "name": item.Name}
		}
	}
	chapterCounts := map[int]int{}
	userCounts := map[int]int{}
	if len(ids) > 0 {
		chapters, err := s.courses.ChapterCountsByCourseIDs(ctx, ids)
		if err != nil {
			return page.Body{}, err
		}
		for _, row := range chapters {
			chapterCounts[row.Key] = row.Total
		}
		users, err := s.courses.UserCountsByCourseIDs(ctx, ids)
		if err != nil {
			return page.Body{}, err
		}
		for _, row := range users {
			userCounts[row.Key] = row.Total
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var cat any
		if found, ok := categories[row.CategoryID]; ok {
			cat = found
		}
		var cert any
		if row.CertID != nil {
			if found, ok := certs[*row.CertID]; ok {
				cert = found
			}
		}
		list = append(list, map[string]any{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
			"is_show": row.IsShow, "category_id": row.CategoryID, "is_rec": row.IsRec,
			"min_need_credit": row.MinNeedCredit, "credit": row.Credit,
			"cert_id": row.CertID, "type": row.Type,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"category": cat, "cert": cert,
			"chapters_count": chapterCounts[row.ID], "users_count": userCounts[row.ID],
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetCourseDetail(ctx context.Context, id int) (map[string]any, error) {
	if id <= 0 {
		return nil, biz("课程不存在")
	}
	var row courseRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, biz("课程不存在")
	}
	if err != nil {
		return nil, err
	}

	var category any
	if row.CategoryID > 0 {
		if cat, err := s.courses.FindCategory(ctx, row.CategoryID); err == nil {
			category = map[string]any{"id": cat.ID, "name": cat.Name}
		}
	}
	var creator any
	if row.CreatedID > 0 {
		if admin, err := s.courses.FindAdmin(ctx, row.CreatedID); err == nil {
			creator = map[string]any{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	var paper any
	if row.PaperID != nil && *row.PaperID > 0 {
		if p, err := s.courses.FindPaperScore(ctx, *row.PaperID); err == nil {
			paper = map[string]any{"id": p.ID, "score": p.Score, "pass_score": p.PassScore}
		}
	}
	var cert any
	if row.CertID != nil && *row.CertID > 0 {
		if certRow, err := s.courses.FindCertDetail(ctx, *row.CertID); err == nil {
			cert = map[string]any{
				"id": certRow.ID, "name": certRow.Name,
				"preview_image": certRow.PreviewImage, "template_image": certRow.TemplateImage,
			}
		}
	}

	chaptersCount, _ := s.courses.CountChapters(ctx, id)
	commentsCount, _ := s.courses.CountComments(ctx, id)
	usersCount, _ := s.courses.CountUsers(ctx, id)
	endCount, _ := s.courses.CountEndedUsers(ctx, id)
	likeCount, _ := s.courses.CountLikes(ctx, id)

	pending := 0
	if row.CertID != nil && *row.CertID > 0 {
		pending, err = s.courses.CountPendingCertUsers(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		"id": row.ID, "title": row.Title, "type": row.Type, "thumb": row.Thumb, "desc": row.Desc,
		"content": row.Content, "is_show": row.IsShow, "category_id": row.CategoryID,
		"is_rec": row.IsRec, "min_need_credit": row.MinNeedCredit,
		"credit": row.Credit, "created_id": row.CreatedID, "paper_id": row.PaperID, "cert_id": row.CertID,
		"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
		"category": category, "creator": creator, "paper": paper, "cert": cert,
		"chapters_count": chaptersCount, "comments_count": commentsCount, "users_count": usersCount,
		"user_courses_end_count": endCount, "like_users_count": likeCount,
		"pending_cert_users": pending,
	}, nil
}

func (s *Service) CourseSave(ctx context.Context, adminID int, raw map[string]any) (CertSaveMeta, string, error) {
	if raw == nil {
		return CertSaveMeta{}, "", biz("参数错误")
	}
	title := strings.TrimSpace(asString(raw["title"]))
	if title == "" {
		return CertSaveMeta{}, "", biz("请填写课程标题")
	}
	isShow := 0
	if asInt(raw["is_show"]) == 1 {
		isShow = 1
	}
	var certID any
	if asInt(raw["cert_id"]) <= 0 {
		certID = nil
	} else {
		certID = asInt(raw["cert_id"])
	}
	id := asInt(raw["id"])
	now := time.Now()
	if id > 0 {
		var row courseRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			return CertSaveMeta{}, "", biz("课程不存在")
		}
		if err != nil {
			return CertSaveMeta{}, "", err
		}
		if isShow == 1 {
			ok, err := s.courses.ChaptersExist(ctx, id)
			if err != nil {
				return CertSaveMeta{}, "", err
			}
			if !ok {
				return CertSaveMeta{}, "", biz("请先添加课节后再上架")
			}
		}
		oldCertID := 0
		if row.CertID != nil {
			oldCertID = *row.CertID
		}
		data := map[string]any{
			"title": title, "thumb": asString(raw["thumb"]), "type": asInt(raw["type"]),
			"desc": asString(raw["desc"]), "content": asString(raw["content"]),
			"is_show": isShow, "category_id": asInt(raw["category_id"]),
			"is_rec": asInt(raw["is_rec"]),
			"min_need_credit": asInt(raw["min_need_credit"]), "credit": asInt(raw["credit"]),
			"cert_id": certID, "updated_at": now,
		}
		err = s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", id).
			Select("title", "thumb", "type", "desc", "content", "is_show", "category_id",
				"is_rec", "min_need_credit", "credit", "cert_id", "updated_at").
			Updates(data).Error
		if err != nil {
			return CertSaveMeta{}, "", err
		}
		var fresh courseRow
		_ = s.db.WithContext(ctx).Where("id = ?", id).Take(&fresh).Error
		newCertID := 0
		if fresh.CertID != nil {
			newCertID = *fresh.CertID
		}
		meta, err := s.certSaveMeta(ctx, id, oldCertID, newCertID)
		if err != nil {
			return CertSaveMeta{}, "", err
		}
		return meta, "更新成功", nil
	}
	row := courseRow{
		Title: title, Thumb: asString(raw["thumb"]), Type: asInt(raw["type"]),
		Desc: asString(raw["desc"]), IsShow: 0, CategoryID: asInt(raw["category_id"]),
		IsRec: asInt(raw["is_rec"]),
		MinNeedCredit: asInt(raw["min_need_credit"]), Credit: asInt(raw["credit"]),
		CreatedID: adminID, CreatedAt: &now, UpdatedAt: &now,
	}
	if content := asString(raw["content"]); content != "" {
		row.Content = &content
	}
	if cid := asInt(raw["cert_id"]); cid > 0 {
		row.CertID = &cid
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return CertSaveMeta{}, "", err
	}
	newCertID := 0
	if row.CertID != nil {
		newCertID = *row.CertID
	}
	meta, err := s.certSaveMeta(ctx, row.ID, 0, newCertID)
	if err != nil {
		return CertSaveMeta{}, "", err
	}
	return meta, "课程创建成功", nil
}

func (s *Service) certSaveMeta(ctx context.Context, courseID, oldCertID, newCertID int) (CertSaveMeta, error) {
	changed := newCertID > 0 && oldCertID != newCertID
	pending, issued := 0, 0
	var err error
	if changed {
		pending, err = s.courses.CountPendingCertUsers(ctx, courseID)
		if err != nil {
			return CertSaveMeta{}, err
		}
		issued, err = s.courses.CountIssuedCertUsers(ctx, courseID)
		if err != nil {
			return CertSaveMeta{}, err
		}
	}
	return CertSaveMeta{
		ID: courseID, CertChanged: changed,
		CertSwapped:      oldCertID > 0 && newCertID > 0 && oldCertID != newCertID,
		PendingCertUsers: pending, IssuedCertUsers: issued,
	}, nil
}

func (s *Service) CourseDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("课程不存在")
	}
	var row courseRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return biz("课程不存在")
	}
	if err != nil {
		return err
	}

	chapterIDs, err := s.courses.PluckChapterIDs(ctx, id)
	if err != nil {
		return err
	}
	chapterPaperIDs, err := s.courses.PluckChapterPaperIDs(ctx, id)
	if err != nil {
		return err
	}
	return s.courses.DeleteCourseCascade(ctx, id, row.PaperID, chapterIDs, chapterPaperIDs)
}

func (s *Service) GetCourseAttachList(ctx context.Context, courseID int) ([]map[string]any, error) {
	if courseID <= 0 {
		return nil, biz("参数错误")
	}
	var rows []attachRow
	err := s.db.WithContext(ctx).
		Select("id, name, path, size, download_times, created_at").
		Where("course_id = ? AND deleted_at IS NULL", courseID).Order("id ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		list = append(list, map[string]any{
			"id": row.ID, "name": row.Name, "path": row.Path, "size": row.Size,
			"download_times": row.DownloadTimes, "created_at": model.FormatTime(row.CreatedAt),
		})
	}
	return list, nil
}

func (s *Service) CourseAttachSave(ctx context.Context, raw map[string]any) error {
	if raw == nil {
		return biz("参数错误")
	}
	courseID := asInt(raw["course_id"])
	name := strings.TrimSpace(asString(raw["name"]))
	path := strings.TrimSpace(asString(raw["path"]))
	if courseID <= 0 || name == "" || path == "" {
		return biz("请填写课件名称并上传文件")
	}
	now := time.Now()
	row := attachRow{
		CourseID: courseID, Name: name, Path: path, Size: asInt(raw["size"]),
		CreatedAt: &now, UpdatedAt: &now,
	}
	return s.db.WithContext(ctx).Create(&row).Error
}

func (s *Service) CourseAttachDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("课件不存在")
	}
	var row attachRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return biz("课件不存在")
	}
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Delete(&row).Error
}

func keys(set map[int]struct{}) []int {
	out := make([]int, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func uniquePositive(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
