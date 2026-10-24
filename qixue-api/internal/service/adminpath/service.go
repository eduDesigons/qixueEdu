package adminpath

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

type Service struct {
	paths repo.AdminPathRepository
	db    *gorm.DB
	rdb   *redis.Client
}

func New(paths repo.AdminPathRepository, db *gorm.DB, rdb *redis.Client) *Service {
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
	IsOpen        int            `gorm:"column:is_open"`
	MinNeedCredit int            `gorm:"column:min_need_credit"`
	CertID        *int           `gorm:"column:cert_id"`
	CreatedID     int            `gorm:"column:created_id"`
	CreatedAt     *time.Time     `gorm:"column:created_at"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (pathRow) TableName() string { return "paths" }

func (s *Service) GetPathCategories(c *gin.Context) {
	rows, err := s.paths.ListPathCategories(c.Request.Context())
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{"id": row.ID, "name": row.Name})
	}
	response.OK(c, "", list)
}

func (s *Service) GetPathList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
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

	ctx := c.Request.Context()
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
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []pathRow
	err := query.Order(sort + " " + order).Offset((current - 1) * size).Limit(size).Find(&rows).Error
	if response.WriteErr(c, err) {
		return
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

	categories := map[int]gin.H{}
	if len(cateIDs) > 0 {
		cats, err := s.paths.ListCategoriesByIDs(ctx, keys(cateIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, cat := range cats {
			categories[cat.ID] = gin.H{"id": cat.ID, "name": cat.Name}
		}
	}
	creators := map[int]gin.H{}
	if len(creatorIDs) > 0 {
		admins, err := s.paths.ListAdminsByIDs(ctx, keys(creatorIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, admin := range admins {
			creators[admin.ID] = gin.H{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	certs := map[int]gin.H{}
	if len(certIDs) > 0 {
		list, err := s.paths.ListCertNamesByIDs(ctx, keys(certIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, item := range list {
			certs[item.ID] = gin.H{"id": item.ID, "name": item.Name}
		}
	}
	courseCounts := map[int]int{}
	userCounts := map[int]int{}
	if len(ids) > 0 {
		courses, err := s.paths.CourseCountsByPathIDs(ctx, ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range courses {
			courseCounts[row.Key] = row.Total
		}
		users, err := s.paths.UserCountsByPathIDs(ctx, ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range users {
			userCounts[row.Key] = row.Total
		}
	}

	list := make([]gin.H, 0, len(rows))
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
		list = append(list, gin.H{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
			"is_show": row.IsShow, "category_id": row.CategoryID, "is_rec": row.IsRec,
			"is_open": row.IsOpen, "min_need_credit": row.MinNeedCredit, "cert_id": row.CertID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"created_id": row.CreatedID, "category": cat, "creator": creator, "cert": cert,
			"courses_count": courseCounts[row.ID], "users_count": userCounts[row.ID],
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) GetPathDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	ctx := c.Request.Context()
	var row pathRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "系列不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}

	var category any
	if row.CategoryID > 0 {
		if cat, err := s.paths.FindCategory(ctx, row.CategoryID); err == nil {
			category = gin.H{"id": cat.ID, "name": cat.Name}
		}
	}
	var creator any
	if row.CreatedID > 0 {
		if admin, err := s.paths.FindAdmin(ctx, row.CreatedID); err == nil {
			creator = gin.H{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	var cert any
	if row.CertID != nil && *row.CertID > 0 {
		if certRow, err := s.paths.FindCertDetail(ctx, *row.CertID); err == nil {
			cert = gin.H{
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
		if response.WriteErr(c, err) {
			return
		}
	}

	response.OK(c, "", gin.H{
		"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
		"is_show": row.IsShow, "content": row.Content, "category_id": row.CategoryID,
		"is_rec": row.IsRec, "is_open": row.IsOpen, "min_need_credit": row.MinNeedCredit,
		"cert_id": row.CertID, "created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt), "created_id": row.CreatedID,
		"category": category, "creator": creator, "cert": cert,
		"courses_count": coursesCount, "users_count": usersCount,
		"user_path_end_count": endCount, "applys_count": applyCount,
		"pending_cert_users": pending,
	})
}

func (s *Service) PathSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	title := strings.TrimSpace(asString(raw["title"]))
	if title == "" {
		response.Fail(c, 1, "请填写系列标题")
		return
	}
	var certID any
	if asInt(raw["cert_id"]) <= 0 {
		certID = nil
	} else {
		certID = asInt(raw["cert_id"])
	}
	ctx := c.Request.Context()
	id := asInt(raw["id"])
	now := time.Now()
	if id > 0 {
		var row pathRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "系列不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		oldCertID := 0
		if row.CertID != nil {
			oldCertID = *row.CertID
		}
		data := map[string]any{
			"title": title, "thumb": asString(raw["thumb"]), "desc": asString(raw["desc"]),
			"content": asString(raw["content"]), "category_id": asInt(raw["category_id"]),
			"is_rec": asInt(raw["is_rec"]), "is_open": asInt(raw["is_open"]),
			"min_need_credit": asInt(raw["min_need_credit"]), "cert_id": certID,
			"is_show": asInt(raw["is_show"]), "updated_at": now,
		}
		err = s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).
			Select("title", "thumb", "desc", "content", "category_id", "is_rec", "is_open",
				"min_need_credit", "cert_id", "is_show", "updated_at").
			Updates(data).Error
		if response.WriteErr(c, err) {
			return
		}
		var fresh pathRow
		_ = s.db.WithContext(ctx).Where("id = ?", id).Take(&fresh).Error
		newCertID := 0
		if fresh.CertID != nil {
			newCertID = *fresh.CertID
		}
		meta, err := s.certSaveMeta(ctx, id, oldCertID, newCertID)
		if response.WriteErr(c, err) {
			return
		}
		response.OK(c, "更新成功", meta)
		return
	}

	row := pathRow{
		Title: title, Thumb: asString(raw["thumb"]), Desc: asString(raw["desc"]),
		CategoryID: asInt(raw["category_id"]), IsRec: asInt(raw["is_rec"]),
		IsOpen: asInt(raw["is_open"]), MinNeedCredit: asInt(raw["min_need_credit"]),
		IsShow: asInt(raw["is_show"]), CreatedID: middleware.UID(c),
		CreatedAt: &now, UpdatedAt: &now,
	}
	if content := asString(raw["content"]); content != "" {
		row.Content = &content
	}
	if cid := asInt(raw["cert_id"]); cid > 0 {
		row.CertID = &cid
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
		return
	}
	newCertID := 0
	if row.CertID != nil {
		newCertID = *row.CertID
	}
	meta, err := s.certSaveMeta(ctx, row.ID, 0, newCertID)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "系列创建成功", meta)
}

func (s *Service) certSaveMeta(ctx context.Context, pathID, oldCertID, newCertID int) (gin.H, error) {
	changed := newCertID > 0 && oldCertID != newCertID
	pending, issued := 0, 0
	var err error
	if changed {
		pending, err = s.paths.CountPendingCertUsers(ctx, pathID)
		if err != nil {
			return nil, err
		}
		issued, err = s.paths.CountIssuedCertUsers(ctx, pathID)
		if err != nil {
			return nil, err
		}
	}
	return gin.H{
		"id": pathID, "cert_changed": changed,
		"cert_swapped":       oldCertID > 0 && newCertID > 0 && oldCertID != newCertID,
		"pending_cert_users": pending, "issued_cert_users": issued,
	}, nil
}

func (s *Service) PathDelete(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	ctx := c.Request.Context()
	var row pathRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "系列不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	courseCount, err := s.paths.CountPathCourses(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if courseCount > 0 {
		response.Fail(c, 1, "当前系列下存在课程无法删除")
		return
	}
	userCount, err := s.paths.CountPathUsers(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if userCount > 0 {
		response.Fail(c, 1, "学员已报名，无法删除")
		return
	}
	if response.WriteErr(c, s.paths.DeletePathCascade(ctx, id)) {
		return
	}
	response.OK(c, "删除成功", []any{})
}
