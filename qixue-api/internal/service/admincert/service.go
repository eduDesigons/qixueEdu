package admincert

import (
	"context"
	"encoding/json"
	"fmt"
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
	certsvc "qixue-api/internal/service/cert"
	"qixue-api/internal/service/scope"
)

type Service struct {
	certs  repo.AdminCertRepository
	db     *gorm.DB
	rdb    *redis.Client
	issuer *certsvc.Issuer
}

func New(certs repo.AdminCertRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{certs: certs, db: db, rdb: rdb, issuer: certsvc.NewIssuer(repo.NewCert(db))}
}

type certRow struct {
	ID            int             `gorm:"column:id;primaryKey"`
	Name          string          `gorm:"column:name"`
	TemplateImage string          `gorm:"column:template_image"`
	PreviewImage  *string         `gorm:"column:preview_image"`
	Params        json.RawMessage `gorm:"column:params"`
	IsShow        int             `gorm:"column:is_show"`
	CreatedAt     *time.Time      `gorm:"column:created_at"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at"`
}

func (certRow) TableName() string { return "certs" }

type certListView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	Params        any     `json:"params"`
	IsShow        int     `json:"is_show"`
	CreatedAt     *string `json:"created_at"`
}

type certDetailView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	Params        any     `json:"params"`
	IsShow        int     `json:"is_show"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

type optionView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	IsShow        int     `json:"is_show"`
}

type userBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type titleBrief struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type certNameBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type issuedView struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	CertID    int            `json:"cert_id"`
	PathID    *int           `json:"path_id"`
	CourseID  int            `json:"course_id"`
	CertNo    string         `json:"cert_no"`
	Title     string         `json:"title"`
	Image     string         `json:"image"`
	IsValid   int            `json:"is_valid"`
	IssuedAt  *string        `json:"issued_at"`
	CreatedAt *string        `json:"created_at"`
	User      *userBrief     `json:"user"`
	Path      *titleBrief    `json:"path"`
	Course    *titleBrief    `json:"course"`
	Cert      *certNameBrief `json:"cert"`
}

func (s *Service) List(c *gin.Context) {
	var body struct {
		Name   string   `json:"name"`
		IsShow any      `json:"is_show"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()

	query := s.db.WithContext(ctx).Model(&certRow{}).
		Select("id, name, template_image, preview_image, params, is_show, created_at")
	if name := strings.TrimSpace(body.Name); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if isShow, ok := asOptionalInt(body.IsShow); ok {
		query = query.Where("is_show = ?", isShow)
	}

	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []certRow
	if response.WriteErr(c, query.Order("id DESC").Offset((current-1)*size).Limit(size).Find(&rows).Error) {
		return
	}
	list := make([]certListView, 0, len(rows))
	for _, r := range rows {
		list = append(list, certListView{
			ID: r.ID, Name: r.Name, TemplateImage: r.TemplateImage, PreviewImage: r.PreviewImage,
			Params: decodeJSON(r.Params), IsShow: r.IsShow, CreatedAt: model.FormatTime(r.CreatedAt),
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) Detail(c *gin.Context) {
	id := asInt(c.Query("id"))
	if id <= 0 {
		response.Fail(c, 1, "缺少 id")
		return
	}
	ctx := c.Request.Context()
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "证书不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "ok", certDetailView{
		ID: row.ID, Name: row.Name, TemplateImage: row.TemplateImage, PreviewImage: row.PreviewImage,
		Params: decodeJSON(row.Params), IsShow: row.IsShow,
		CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
	})
}

func (s *Service) Options(c *gin.Context) {
	ctx := c.Request.Context()
	rows, err := s.certs.ListOptions(ctx)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]optionView, 0, len(rows))
	for _, r := range rows {
		list = append(list, optionView{
			ID: r.ID, Name: r.Name, TemplateImage: r.TemplateImage,
			PreviewImage: r.PreviewImage, IsShow: r.IsShow,
		})
	}
	response.OK(c, "ok", list)
}

func (s *Service) SetShow(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id"`
		IsShow any      `json:"is_show"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 {
		response.Fail(c, 1, "缺少 id")
		return
	}
	ctx := c.Request.Context()
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", body.ID.Int()).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "证书不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	isShow := 0
	if asInt(body.IsShow) != 0 {
		isShow = 1
	}
	now := time.Now()
	if response.WriteErr(c, s.db.WithContext(ctx).Model(&certRow{}).Where("id = ?", row.ID).
		Updates(map[string]any{"is_show": isShow, "updated_at": now}).Error) {
		return
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Where("id = ?", row.ID).Take(&row).Error) {
		return
	}
	msg := "已下架"
	if isShow == 1 {
		msg = "已上架"
	}
	response.OK(c, msg, certDetailView{
		ID: row.ID, Name: row.Name, TemplateImage: row.TemplateImage, PreviewImage: row.PreviewImage,
		Params: decodeJSON(row.Params), IsShow: row.IsShow,
		CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
	})
}

func (s *Service) Save(c *gin.Context) {
	var body struct {
		ID            bind.Int        `json:"id"`
		Name          string          `json:"name"`
		TemplateImage string          `json:"template_image"`
		PreviewImage  string          `json:"preview_image"`
		Params        json.RawMessage `json:"params"`
		IsShow        any             `json:"is_show"`
	}
	_ = c.ShouldBindJSON(&body)
	name := strings.TrimSpace(body.Name)
	if name == "" {
		response.Fail(c, 1, "请填写证书名称")
		return
	}
	templateImage := strings.TrimSpace(body.TemplateImage)
	if templateImage == "" {
		response.Fail(c, 1, "请上传证书背景图")
		return
	}

	paramsRaw := normalizeParams(body.Params)
	isShow := 1
	if body.IsShow != nil {
		if asInt(body.IsShow) == 0 {
			isShow = 0
		}
	}
	now := time.Now()
	data := map[string]any{
		"name": name, "template_image": templateImage, "params": paramsRaw,
		"is_show": isShow, "updated_at": now,
	}
	if preview := strings.TrimSpace(body.PreviewImage); preview != "" {
		data["preview_image"] = preview
	}

	ctx := c.Request.Context()
	id := body.ID.Int()
	if id > 0 {
		var row certRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == nil {
			if response.WriteErr(c, s.db.WithContext(ctx).Model(&certRow{}).Where("id = ?", id).Updates(data).Error) {
				return
			}
			if response.WriteErr(c, s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error) {
				return
			}
			response.OK(c, "更新成功", certDetailView{
				ID: row.ID, Name: row.Name, TemplateImage: row.TemplateImage, PreviewImage: row.PreviewImage,
				Params: decodeJSON(row.Params), IsShow: row.IsShow,
				CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
			})
			return
		}
		if err != gorm.ErrRecordNotFound && response.WriteErr(c, err) {
			return
		}
	}

	data["created_at"] = now
	row := certRow{
		Name: name, TemplateImage: templateImage, Params: paramsRaw, IsShow: isShow,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if preview := strings.TrimSpace(body.PreviewImage); preview != "" {
		row.PreviewImage = &preview
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
		return
	}
	response.OK(c, "保存成功", certDetailView{
		ID: row.ID, Name: row.Name, TemplateImage: row.TemplateImage, PreviewImage: row.PreviewImage,
		Params: decodeJSON(row.Params), IsShow: row.IsShow,
		CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
	})
}

func (s *Service) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	id := body.ID.Int()
	if id <= 0 {
		response.Fail(c, 1, "缺少 id")
		return
	}
	ctx := c.Request.Context()
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "证书不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}

	reasons := []string{}
	n, err := s.certs.CountPathsByCert(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if n > 0 {
		reasons = append(reasons, "已绑定系列")
	}
	n, err = s.certs.CountCoursesByCert(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if n > 0 {
		reasons = append(reasons, "已绑定课程")
	}
	n, err = s.certs.CountUserCertsByCert(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if n > 0 {
		reasons = append(reasons, "已有学员证书记录")
	}
	if len(reasons) > 0 {
		response.Fail(c, 1, "无法删除："+strings.Join(reasons, "、")+"，请先解除关联")
		return
	}
	// certs table has no SoftDeletes
	if response.WriteErr(c, s.db.WithContext(ctx).Where("id = ?", id).Delete(&certRow{}).Error) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func (s *Service) Issue(c *gin.Context) {
	var body struct {
		UserID   bind.Int `json:"user_id"`
		PathID   bind.Int `json:"path_id"`
		CourseID bind.Int `json:"course_id"`
		Force    bool     `json:"force"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.UserID.Int() <= 0 || (body.PathID.Int() <= 0 && body.CourseID.Int() <= 0) {
		response.Fail(c, 1, "缺少学员或系列/课程")
		return
	}
	ctx := c.Request.Context()
	var (
		record *certsvc.UserCert
		err    error
	)
	if body.CourseID.Int() > 0 {
		record, err = s.issuer.IssueForCourse(ctx, body.UserID.Int(), body.CourseID.Int(), body.Force)
	} else {
		record, err = s.issuer.IssueForPath(ctx, body.UserID.Int(), body.PathID.Int(), body.Force)
	}
	if response.WriteErr(c, err) {
		return
	}
	if record == nil {
		if body.CourseID.Int() > 0 {
			response.Fail(c, 1, "该课程未绑定证书模板、已下架或学员未结业")
		} else {
			response.Fail(c, 1, "该系列未绑定证书模板或已下架")
		}
		return
	}
	msg := "已提交颁发，证书图生成中"
	if body.Force {
		msg = "已提交重发，证书图生成中"
	}
	response.OK(c, msg, formatUserCert(record))
}

func (s *Service) Backfill(c *gin.Context) {
	var body struct {
		CourseID bind.Int `json:"course_id"`
		PathID   bind.Int `json:"path_id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.CourseID.Int() <= 0 && body.PathID.Int() <= 0 {
		response.Fail(c, 1, "缺少课程或系列")
		return
	}
	ctx := c.Request.Context()
	var (
		count int
		err   error
	)
	if body.CourseID.Int() > 0 {
		count, err = s.issuer.BackfillForCourse(ctx, body.CourseID.Int(), 500)
	} else {
		count, err = s.issuer.BackfillForPath(ctx, body.PathID.Int(), 500)
	}
	if response.WriteErr(c, err) {
		return
	}
	if count <= 0 {
		response.OK(c, "暂无需要补发的学员", gin.H{"count": 0})
		return
	}
	response.OK(c, fmt.Sprintf("已为 %d 名学员补发证书", count), gin.H{"count": count})
}

func (s *Service) IssuedList(c *gin.Context) {
	var body struct {
		Name     string   `json:"name"`
		CourseID bind.Int `json:"course_id"`
		PathID   bind.Int `json:"path_id"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}

	rows, total, err := s.certs.ListIssued(ctx, repo.AdminCertIssuedFilter{
		CourseID: body.CourseID.Int(), PathID: body.PathID.Int(),
		Name: strings.TrimSpace(body.Name), CityID: cityID, Limited: limited,
	}, (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}

	userIDs, pathIDs, courseIDs, certIDs := []int{}, []int{}, []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		if r.PathID != nil && *r.PathID > 0 {
			pathIDs = append(pathIDs, *r.PathID)
		}
		if r.CourseID > 0 {
			courseIDs = append(courseIDs, r.CourseID)
		}
		if r.CertID > 0 {
			certIDs = append(certIDs, r.CertID)
		}
	}
	users, err := s.loadUsers(ctx, uniquePositive(userIDs))
	if response.WriteErr(c, err) {
		return
	}
	paths, err := s.loadTitles(ctx, "paths", uniquePositive(pathIDs))
	if response.WriteErr(c, err) {
		return
	}
	courses, err := s.loadTitles(ctx, "courses", uniquePositive(courseIDs))
	if response.WriteErr(c, err) {
		return
	}
	certs, err := s.loadCertNames(ctx, uniquePositive(certIDs))
	if response.WriteErr(c, err) {
		return
	}

	list := make([]issuedView, 0, len(rows))
	for _, r := range rows {
		var user *userBrief
		if u, ok := users[r.UserID]; ok {
			cp := u
			user = &cp
		}
		var path *titleBrief
		if r.PathID != nil {
			if p, ok := paths[*r.PathID]; ok {
				cp := p
				path = &cp
			}
		}
		var course *titleBrief
		if co, ok := courses[r.CourseID]; ok {
			cp := co
			course = &cp
		}
		var cert *certNameBrief
		if ce, ok := certs[r.CertID]; ok {
			cp := ce
			cert = &cp
		}
		list = append(list, issuedView{
			ID: r.ID, UserID: r.UserID, CertID: r.CertID, PathID: r.PathID, CourseID: r.CourseID,
			CertNo: r.CertNo, Title: r.Title, Image: r.Image, IsValid: r.IsValid,
			IssuedAt: model.FormatTime(r.IssuedAt), CreatedAt: model.FormatTime(r.CreatedAt),
			User: user, Path: path, Course: course, Cert: cert,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) Revoke(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 {
		response.Fail(c, 1, "缺少 id")
		return
	}
	ctx := c.Request.Context()
	var row certsvc.UserCert
	err := s.db.WithContext(ctx).Where("id = ?", body.ID.Int()).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "证书记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if row.IsValid != 1 {
		response.OK(c, "证书已是撤回状态", formatUserCert(&row))
		return
	}
	now := time.Now()
	if response.WriteErr(c, s.db.WithContext(ctx).Model(&certsvc.UserCert{}).Where("id = ?", row.ID).
		Updates(map[string]any{"is_valid": 0, "updated_at": now}).Error) {
		return
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Where("id = ?", row.ID).Take(&row).Error) {
		return
	}
	response.OK(c, "已撤回", formatUserCert(&row))
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]userBrief, error) {
	out := map[int]userBrief{}
	rows, err := s.certs.ListUsersBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = userBrief{ID: r.ID, NickName: r.NickName, Avatar: r.Avatar}
	}
	return out, nil
}

func (s *Service) loadTitles(ctx context.Context, table string, ids []int) (map[int]titleBrief, error) {
	out := map[int]titleBrief{}
	rows, err := s.certs.ListTitles(ctx, table, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = titleBrief{ID: r.ID, Title: r.Title}
	}
	return out, nil
}

func (s *Service) loadCertNames(ctx context.Context, ids []int) (map[int]certNameBrief, error) {
	out := map[int]certNameBrief{}
	rows, err := s.certs.ListCertNames(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = certNameBrief{ID: r.ID, Name: r.Name}
	}
	return out, nil
}

func formatUserCert(row *certsvc.UserCert) gin.H {
	if row == nil {
		return nil
	}
	return gin.H{
		"id": row.ID, "user_id": row.UserID, "cert_id": row.CertID,
		"path_id": row.PathID, "course_id": row.CourseID, "cert_no": row.CertNo,
		"title": row.Title, "image": row.Image, "is_valid": row.IsValid,
		"issued_at":  model.FormatTime(row.IssuedAt),
		"created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt),
	}
}

func decodeJSON(raw json.RawMessage) any {
	if len(raw) == 0 || string(raw) == "null" {
		return []any{}
	}
	var v any
	if json.Unmarshal(raw, &v) == nil {
		return v
	}
	return []any{}
}

func normalizeParams(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 || string(raw) == "null" {
		return json.RawMessage("[]")
	}

	var asString string
	if json.Unmarshal(raw, &asString) == nil {
		asString = strings.TrimSpace(asString)
		if asString == "" {
			return json.RawMessage("[]")
		}
		var v any
		if json.Unmarshal([]byte(asString), &v) == nil {
			encoded, err := json.Marshal(v)
			if err == nil {
				return encoded
			}
		}
		return json.RawMessage("[]")
	}
	var v any
	if json.Unmarshal(raw, &v) == nil {
		encoded, err := json.Marshal(v)
		if err == nil {
			return encoded
		}
	}
	return json.RawMessage("[]")
}

func asInt(v any) int { return bind.AsInt(v) }

func asOptionalInt(v any) (int, bool) {
	if !bind.Present(v) {
		return 0, false
	}
	return bind.AsInt(v), true
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
