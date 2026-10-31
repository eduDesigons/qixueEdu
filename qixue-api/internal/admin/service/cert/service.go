package admincert

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	adminrepo "qixue-api/internal/admin/repo"
	certsvc "qixue-api/internal/service/cert"
	"qixue-api/internal/service/scope"
)

var (
	ErrMissingID      = errors.New("缺少 id")
	ErrNotFound       = errors.New("证书不存在")
	ErrEmptyName      = errors.New("请填写证书名称")
	ErrEmptyTemplate  = errors.New("请上传证书背景图")
	ErrMissingTarget  = errors.New("缺少学员或系列/课程")
	ErrMissingScope   = errors.New("缺少课程或系列")
	ErrCourseNoIssue  = errors.New("该课程未绑定证书模板、已下架或学员未结业")
	ErrPathNoIssue    = errors.New("该系列未绑定证书模板或已下架")
	ErrIssuedNotFound = errors.New("证书记录不存在")
)

type Service struct {
	certs  adminrepo.AdminCertRepository
	db     *gorm.DB
	rdb    *redis.Client
	issuer *certsvc.Issuer
}

func New(certs adminrepo.AdminCertRepository, db *gorm.DB, rdb *redis.Client) *Service {
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

type ListView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	Params        any     `json:"params"`
	IsShow        int     `json:"is_show"`
	CreatedAt     *string `json:"created_at"`
}

type DetailView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	Params        any     `json:"params"`
	IsShow        int     `json:"is_show"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

type OptionView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TemplateImage string  `json:"template_image"`
	PreviewImage  *string `json:"preview_image"`
	IsShow        int     `json:"is_show"`
}

type UserBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type TitleBrief struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CertNameBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IssuedView struct {
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
	User      *UserBrief     `json:"user"`
	Path      *TitleBrief    `json:"path"`
	Course    *TitleBrief    `json:"course"`
	Cert      *CertNameBrief `json:"cert"`
}

type UserCertView struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	CertID    int     `json:"cert_id"`
	PathID    int     `json:"path_id"`
	CourseID  int     `json:"course_id"`
	CertNo    string  `json:"cert_no"`
	Title     string  `json:"title"`
	Image     string  `json:"image"`
	IsValid   int     `json:"is_valid"`
	IssuedAt  *string `json:"issued_at"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type ListQuery struct {
	Name   string
	IsShow *int
	Page   int
	Size   int
}

type SaveInput struct {
	ID            int
	Name          string
	TemplateImage string
	PreviewImage  string
	Params        json.RawMessage
	IsShow        *int
}

type IssueInput struct {
	UserID   int
	PathID   int
	CourseID int
	Force    bool
}

type IssuedListQuery struct {
	Name     string
	CourseID int
	PathID   int
	Page     int
	Size     int
}

type DeleteBlockedError struct {
	Reasons []string
}

func (e *DeleteBlockedError) Error() string {
	return "无法删除：" + strings.Join(e.Reasons, "、") + "，请先解除关联"
}

func (s *Service) List(ctx context.Context, q ListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	query := s.db.WithContext(ctx).Model(&certRow{}).
		Select("id, name, template_image, preview_image, params, is_show, created_at")
	if name := strings.TrimSpace(q.Name); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if q.IsShow != nil {
		query = query.Where("is_show = ?", *q.IsShow)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []certRow
	if err := query.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list := make([]ListView, 0, len(rows))
	for _, r := range rows {
		list = append(list, ListView{
			ID: r.ID, Name: r.Name, TemplateImage: r.TemplateImage, PreviewImage: r.PreviewImage,
			Params: decodeJSON(r.Params), IsShow: r.IsShow, CreatedAt: model.FormatTime(r.CreatedAt),
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) Detail(ctx context.Context, id int) (DetailView, error) {
	if id <= 0 {
		return DetailView{}, ErrMissingID
	}
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DetailView{}, ErrNotFound
	}
	if err != nil {
		return DetailView{}, err
	}
	return toDetail(row), nil
}

func (s *Service) Options(ctx context.Context) ([]OptionView, error) {
	rows, err := s.certs.ListOptions(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]OptionView, 0, len(rows))
	for _, r := range rows {
		list = append(list, OptionView{
			ID: r.ID, Name: r.Name, TemplateImage: r.TemplateImage,
			PreviewImage: r.PreviewImage, IsShow: r.IsShow,
		})
	}
	return list, nil
}

func (s *Service) SetShow(ctx context.Context, id, isShowRaw int) (DetailView, string, error) {
	if id <= 0 {
		return DetailView{}, "", ErrMissingID
	}
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DetailView{}, "", ErrNotFound
	}
	if err != nil {
		return DetailView{}, "", err
	}
	isShow := 0
	if isShowRaw != 0 {
		isShow = 1
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&certRow{}).Where("id = ?", row.ID).
		Updates(map[string]any{"is_show": isShow, "updated_at": now}).Error; err != nil {
		return DetailView{}, "", err
	}
	if err := s.db.WithContext(ctx).Where("id = ?", row.ID).Take(&row).Error; err != nil {
		return DetailView{}, "", err
	}
	msg := "已下架"
	if isShow == 1 {
		msg = "已上架"
	}
	return toDetail(row), msg, nil
}

func (s *Service) Save(ctx context.Context, in SaveInput) (DetailView, string, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return DetailView{}, "", ErrEmptyName
	}
	templateImage := strings.TrimSpace(in.TemplateImage)
	if templateImage == "" {
		return DetailView{}, "", ErrEmptyTemplate
	}

	paramsRaw := normalizeParams(in.Params)
	isShow := 1
	if in.IsShow != nil && *in.IsShow == 0 {
		isShow = 0
	}
	now := time.Now()
	data := map[string]any{
		"name": name, "template_image": templateImage, "params": paramsRaw,
		"is_show": isShow, "updated_at": now,
	}
	if preview := strings.TrimSpace(in.PreviewImage); preview != "" {
		data["preview_image"] = preview
	}

	if in.ID > 0 {
		var row certRow
		err := s.db.WithContext(ctx).Where("id = ?", in.ID).Take(&row).Error
		if err == nil {
			if err := s.db.WithContext(ctx).Model(&certRow{}).Where("id = ?", in.ID).Updates(data).Error; err != nil {
				return DetailView{}, "", err
			}
			if err := s.db.WithContext(ctx).Where("id = ?", in.ID).Take(&row).Error; err != nil {
				return DetailView{}, "", err
			}
			return toDetail(row), "更新成功", nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return DetailView{}, "", err
		}
	}

	row := certRow{
		Name: name, TemplateImage: templateImage, Params: paramsRaw, IsShow: isShow,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if preview := strings.TrimSpace(in.PreviewImage); preview != "" {
		row.PreviewImage = &preview
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return DetailView{}, "", err
	}
	return toDetail(row), "保存成功", nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrMissingID
	}
	var row certRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	reasons := []string{}
	n, err := s.certs.CountPathsByCert(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		reasons = append(reasons, "已绑定系列")
	}
	n, err = s.certs.CountCoursesByCert(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		reasons = append(reasons, "已绑定课程")
	}
	n, err = s.certs.CountUserCertsByCert(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		reasons = append(reasons, "已有学员证书记录")
	}
	if len(reasons) > 0 {
		return &DeleteBlockedError{Reasons: reasons}
	}
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&certRow{}).Error
}

func (s *Service) Issue(ctx context.Context, in IssueInput) (UserCertView, string, error) {
	if in.UserID <= 0 || (in.PathID <= 0 && in.CourseID <= 0) {
		return UserCertView{}, "", ErrMissingTarget
	}
	var (
		record *certsvc.UserCert
		err    error
	)
	if in.CourseID > 0 {
		record, err = s.issuer.IssueForCourse(ctx, in.UserID, in.CourseID, in.Force)
	} else {
		record, err = s.issuer.IssueForPath(ctx, in.UserID, in.PathID, in.Force)
	}
	if err != nil {
		return UserCertView{}, "", err
	}
	if record == nil {
		if in.CourseID > 0 {
			return UserCertView{}, "", ErrCourseNoIssue
		}
		return UserCertView{}, "", ErrPathNoIssue
	}
	msg := "已提交颁发，证书图生成中"
	if in.Force {
		msg = "已提交重发，证书图生成中"
	}
	return formatUserCert(record), msg, nil
}

func (s *Service) Backfill(ctx context.Context, courseID, pathID int) (int, error) {
	if courseID <= 0 && pathID <= 0 {
		return 0, ErrMissingScope
	}
	if courseID > 0 {
		return s.issuer.BackfillForCourse(ctx, courseID, 500)
	}
	return s.issuer.BackfillForPath(ctx, pathID, 500)
}

func (s *Service) IssuedList(ctx context.Context, adminID int, q IssuedListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}

	rows, total, err := s.certs.ListIssued(ctx, adminrepo.AdminCertIssuedFilter{
		CourseID: q.CourseID, PathID: q.PathID,
		Name: strings.TrimSpace(q.Name), CityID: cityID, Limited: limited,
	}, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
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
	if err != nil {
		return page.Body{}, err
	}
	paths, err := s.loadTitles(ctx, "paths", uniquePositive(pathIDs))
	if err != nil {
		return page.Body{}, err
	}
	courses, err := s.loadTitles(ctx, "courses", uniquePositive(courseIDs))
	if err != nil {
		return page.Body{}, err
	}
	certs, err := s.loadCertNames(ctx, uniquePositive(certIDs))
	if err != nil {
		return page.Body{}, err
	}

	list := make([]IssuedView, 0, len(rows))
	for _, r := range rows {
		var user *UserBrief
		if u, ok := users[r.UserID]; ok {
			cp := u
			user = &cp
		}
		var path *TitleBrief
		if r.PathID != nil {
			if p, ok := paths[*r.PathID]; ok {
				cp := p
				path = &cp
			}
		}
		var course *TitleBrief
		if co, ok := courses[r.CourseID]; ok {
			cp := co
			course = &cp
		}
		var cert *CertNameBrief
		if ce, ok := certs[r.CertID]; ok {
			cp := ce
			cert = &cp
		}
		list = append(list, IssuedView{
			ID: r.ID, UserID: r.UserID, CertID: r.CertID, PathID: r.PathID, CourseID: r.CourseID,
			CertNo: r.CertNo, Title: r.Title, Image: r.Image, IsValid: r.IsValid,
			IssuedAt: model.FormatTime(r.IssuedAt), CreatedAt: model.FormatTime(r.CreatedAt),
			User: user, Path: path, Course: course, Cert: cert,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) Revoke(ctx context.Context, id int) (UserCertView, string, error) {
	if id <= 0 {
		return UserCertView{}, "", ErrMissingID
	}
	var row certsvc.UserCert
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserCertView{}, "", ErrIssuedNotFound
	}
	if err != nil {
		return UserCertView{}, "", err
	}
	if row.IsValid != 1 {
		return formatUserCert(&row), "证书已是撤回状态", nil
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&certsvc.UserCert{}).Where("id = ?", row.ID).
		Updates(map[string]any{"is_valid": 0, "updated_at": now}).Error; err != nil {
		return UserCertView{}, "", err
	}
	if err := s.db.WithContext(ctx).Where("id = ?", row.ID).Take(&row).Error; err != nil {
		return UserCertView{}, "", err
	}
	return formatUserCert(&row), "已撤回", nil
}

func toDetail(row certRow) DetailView {
	return DetailView{
		ID: row.ID, Name: row.Name, TemplateImage: row.TemplateImage, PreviewImage: row.PreviewImage,
		Params: decodeJSON(row.Params), IsShow: row.IsShow,
		CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
	}
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]UserBrief, error) {
	out := map[int]UserBrief{}
	rows, err := s.certs.ListUsersBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = UserBrief{ID: r.ID, NickName: r.NickName, Avatar: r.Avatar}
	}
	return out, nil
}

func (s *Service) loadTitles(ctx context.Context, table string, ids []int) (map[int]TitleBrief, error) {
	out := map[int]TitleBrief{}
	rows, err := s.certs.ListTitles(ctx, table, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = TitleBrief{ID: r.ID, Title: r.Title}
	}
	return out, nil
}

func (s *Service) loadCertNames(ctx context.Context, ids []int) (map[int]CertNameBrief, error) {
	out := map[int]CertNameBrief{}
	rows, err := s.certs.ListCertNames(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = CertNameBrief{ID: r.ID, Name: r.Name}
	}
	return out, nil
}

func formatUserCert(row *certsvc.UserCert) UserCertView {
	if row == nil {
		return UserCertView{}
	}
	return UserCertView{
		ID: row.ID, UserID: row.UserID, CertID: row.CertID,
		PathID: row.PathID, CourseID: row.CourseID, CertNo: row.CertNo,
		Title: row.Title, Image: row.Image, IsValid: row.IsValid,
		IssuedAt: model.FormatTime(row.IssuedAt), CreatedAt: model.FormatTime(row.CreatedAt),
		UpdatedAt: model.FormatTime(row.UpdatedAt),
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
