package adminpath

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
	paths adminrepo.AdminPathRepository
	db    *gorm.DB
	rdb   *redis.Client
}

func New(paths adminrepo.AdminPathRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{paths: paths, db: db, rdb: rdb}
}

type pathRow struct {
	ID            int            `gorm:"column:id;primaryKey"`
	Title         string         `gorm:"column:title"`
	Thumb         string         `gorm:"column:thumb"`
	Desc          string         `gorm:"column:desc"`
	Content       *string        `gorm:"column:content"`
	IsShow        int            `gorm:"column:is_show"`
	CategoryID    int            `gorm:"column:category_id"`
	IsRec         int            `gorm:"column:is_rec"`
	MinNeedCredit int            `gorm:"column:min_need_credit"`
	CertID        *int           `gorm:"column:cert_id"`
	CreatedID     int            `gorm:"column:created_id"`
	CreatedAt     *time.Time     `gorm:"column:created_at"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (pathRow) TableName() string { return "paths" }

type CertSaveMeta struct {
	ID               int  `json:"id"`
	CertChanged      bool `json:"cert_changed"`
	CertSwapped      bool `json:"cert_swapped"`
	PendingCertUsers int  `json:"pending_cert_users"`
	IssuedCertUsers  int  `json:"issued_cert_users"`
}

func (s *Service) GetPathCategories(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.paths.ListPathCategories(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		list = append(list, map[string]any{"id": row.ID, "name": row.Name})
	}
	return list, nil
}

func (s *Service) GetPathList(ctx context.Context, raw map[string]any) (page.Body, error) {
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

	query := s.db.WithContext(ctx).Model(&pathRow{})
	if id := asInt(raw["id"]); id > 0 {
		query = query.Where("id = ?", id)
	}
	if title := strings.TrimSpace(asString(raw["title"])); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if cid := asInt(raw["cid"]); cid > 0 {
		query = query.Where("category_id = ?", cid)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []pathRow
	if err := query.Order(sort + " " + order).Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}

	ids := make([]int, 0, len(rows))
	cateIDs := map[int]struct{}{}
	creatorIDs := map[int]struct{}{}
	certIDs := map[int]struct{}{}
	for _, row := range rows {
		ids = append(ids, row.ID)
		if row.CategoryID > 0 {
			cateIDs[row.CategoryID] = struct{}{}
		}
		if row.CreatedID > 0 {
			creatorIDs[row.CreatedID] = struct{}{}
		}
		if row.CertID != nil && *row.CertID > 0 {
			certIDs[*row.CertID] = struct{}{}
		}
	}

	categories := map[int]map[string]any{}
	if len(cateIDs) > 0 {
		cats, err := s.paths.ListCategoriesByIDs(ctx, keys(cateIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, cat := range cats {
			categories[cat.ID] = map[string]any{"id": cat.ID, "name": cat.Name}
		}
	}
	creators := map[int]map[string]any{}
	if len(creatorIDs) > 0 {
		admins, err := s.paths.ListAdminsByIDs(ctx, keys(creatorIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, admin := range admins {
			creators[admin.ID] = map[string]any{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	certs := map[int]map[string]any{}
	if len(certIDs) > 0 {
		list, err := s.paths.ListCertNamesByIDs(ctx, keys(certIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, item := range list {
			certs[item.ID] = map[string]any{"id": item.ID, "name": item.Name}
		}
	}
	courseCounts := map[int]int{}
	userCounts := map[int]int{}
	if len(ids) > 0 {
		courses, err := s.paths.CourseCountsByPathIDs(ctx, ids)
		if err != nil {
			return page.Body{}, err
		}
		for _, row := range courses {
			courseCounts[row.Key] = row.Total
		}
		users, err := s.paths.UserCountsByPathIDs(ctx, ids)
		if err != nil {
			return page.Body{}, err
		}
		for _, row := range users {
			userCounts[row.Key] = row.Total
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var cat, creator, cert any
		if found, ok := categories[row.CategoryID]; ok {
			cat = found
		}
		if found, ok := creators[row.CreatedID]; ok {
			creator = found
		}
		if row.CertID != nil {
			if found, ok := certs[*row.CertID]; ok {
				cert = found
			}
		}
		list = append(list, map[string]any{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
			"is_show": row.IsShow, "category_id": row.CategoryID, "is_rec": row.IsRec,
			"min_need_credit": row.MinNeedCredit, "cert_id": row.CertID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"created_id": row.CreatedID, "category": cat, "creator": creator, "cert": cert,
			"courses_count": courseCounts[row.ID], "users_count": userCounts[row.ID],
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetPathDetail(ctx context.Context, id int) (map[string]any, error) {
	if id <= 0 {
		return nil, biz("系列不存在")
	}
	var row pathRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, biz("系列不存在")
	}
	if err != nil {
		return nil, err
	}

	var category any
	if row.CategoryID > 0 {
		if cat, err := s.paths.FindCategory(ctx, row.CategoryID); err == nil {
			category = map[string]any{"id": cat.ID, "name": cat.Name}
		}
	}
	var creator any
	if row.CreatedID > 0 {
		if admin, err := s.paths.FindAdmin(ctx, row.CreatedID); err == nil {
			creator = map[string]any{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	var cert any
	if row.CertID != nil && *row.CertID > 0 {
		if certRow, err := s.paths.FindCertDetail(ctx, *row.CertID); err == nil {
			cert = map[string]any{
				"id": certRow.ID, "name": certRow.Name,
				"preview_image": certRow.PreviewImage, "template_image": certRow.TemplateImage,
			}
		}
	}

	coursesCount, _ := s.paths.CountPathCourses(ctx, id)
	usersCount, _ := s.paths.CountPathUsers(ctx, id)
	endCount, _ := s.paths.CountPathEndedUsers(ctx, id)
	applyCount, _ := s.paths.CountPathApplies(ctx, id)

	pending := 0
	if row.CertID != nil && *row.CertID > 0 {
		pending, err = s.paths.CountPendingCertUsers(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
		"is_show": row.IsShow, "content": row.Content, "category_id": row.CategoryID,
		"is_rec": row.IsRec, "min_need_credit": row.MinNeedCredit,
		"cert_id": row.CertID, "created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt), "created_id": row.CreatedID,
		"category": category, "creator": creator, "cert": cert,
		"courses_count": coursesCount, "users_count": usersCount,
		"user_path_end_count": endCount, "applys_count": applyCount,
		"pending_cert_users": pending,
	}, nil
}

func (s *Service) PathSave(ctx context.Context, adminID int, raw map[string]any) (CertSaveMeta, string, error) {
	if raw == nil {
		return CertSaveMeta{}, "", biz("参数错误")
	}
	title := strings.TrimSpace(asString(raw["title"]))
	if title == "" {
		return CertSaveMeta{}, "", biz("请填写系列标题")
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
		var row pathRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			return CertSaveMeta{}, "", biz("系列不存在")
		}
		if err != nil {
			return CertSaveMeta{}, "", err
		}
		oldCertID := 0
		if row.CertID != nil {
			oldCertID = *row.CertID
		}
		data := map[string]any{
			"title": title, "thumb": asString(raw["thumb"]), "desc": asString(raw["desc"]),
			"content": asString(raw["content"]), "category_id": asInt(raw["category_id"]),
			"is_rec": asInt(raw["is_rec"]),
			"min_need_credit": asInt(raw["min_need_credit"]), "cert_id": certID,
			"is_show": asInt(raw["is_show"]), "updated_at": now,
		}
		err = s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).
			Select("title", "thumb", "desc", "content", "category_id", "is_rec",
				"min_need_credit", "cert_id", "is_show", "updated_at").
			Updates(data).Error
		if err != nil {
			return CertSaveMeta{}, "", err
		}
		var fresh pathRow
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

	row := pathRow{
		Title: title, Thumb: asString(raw["thumb"]), Desc: asString(raw["desc"]),
		CategoryID: asInt(raw["category_id"]), IsRec: asInt(raw["is_rec"]),
		MinNeedCredit: asInt(raw["min_need_credit"]),
		IsShow: asInt(raw["is_show"]), CreatedID: adminID,
		CreatedAt: &now, UpdatedAt: &now,
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
	return meta, "系列创建成功", nil
}

func (s *Service) certSaveMeta(ctx context.Context, pathID, oldCertID, newCertID int) (CertSaveMeta, error) {
	changed := newCertID > 0 && oldCertID != newCertID
	pending, issued := 0, 0
	var err error
	if changed {
		pending, err = s.paths.CountPendingCertUsers(ctx, pathID)
		if err != nil {
			return CertSaveMeta{}, err
		}
		issued, err = s.paths.CountIssuedCertUsers(ctx, pathID)
		if err != nil {
			return CertSaveMeta{}, err
		}
	}
	return CertSaveMeta{
		ID: pathID, CertChanged: changed,
		CertSwapped:      oldCertID > 0 && newCertID > 0 && oldCertID != newCertID,
		PendingCertUsers: pending, IssuedCertUsers: issued,
	}, nil
}

func (s *Service) PathDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("系列不存在")
	}
	var row pathRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return biz("系列不存在")
	}
	if err != nil {
		return err
	}
	courseCount, err := s.paths.CountPathCourses(ctx, id)
	if err != nil {
		return err
	}
	if courseCount > 0 {
		return biz("当前系列下存在课程无法删除")
	}
	userCount, err := s.paths.CountPathUsers(ctx, id)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return biz("学员已报名，无法删除")
	}
	return s.paths.DeletePathCascade(ctx, id)
}
