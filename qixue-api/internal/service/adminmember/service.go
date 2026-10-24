package adminmember

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/appconfig"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/scope"
)

type Service struct {
	members repo.AdminMemberRepository
	db      *gorm.DB
	rdb     *redis.Client
}

func New(members repo.AdminMemberRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{members: members, db: db, rdb: rdb}
}

var (
	hasLetter = regexp.MustCompile(`[a-zA-Z]`)
	hasDigit  = regexp.MustCompile(`\d`)
)

type userRow struct {
	ID         int        `gorm:"column:id;primaryKey"`
	Avatar     string     `gorm:"column:avatar"`
	NickName   string     `gorm:"column:nick_name"`
	Account    *string    `gorm:"column:account"`
	Mobile     string     `gorm:"column:mobile"`
	Email      *string    `gorm:"column:email"`
	Password   string     `gorm:"column:password"`
	PasswordBF string     `gorm:"column:passwordbf"`
	IsActive   int        `gorm:"column:is_active"`
	IsLock     int        `gorm:"column:is_lock"`
	AreaID     int        `gorm:"column:area_id"`
	CityID     *int       `gorm:"column:city_id"`
	CityCode   *string    `gorm:"column:city_code"`
	Level      int        `gorm:"column:level"`
	Credit     int        `gorm:"column:credit"`
	Coins      int        `gorm:"column:coins"`
	CreatedAt  *time.Time `gorm:"column:created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
}

func (userRow) TableName() string { return "users" }

type areaView struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ParentName string `json:"parent_name,omitempty"`
}

type listView struct {
	ID           int       `json:"id"`
	Avatar       string    `json:"avatar"`
	NickName     string    `json:"nick_name"`
	Account      *string   `json:"account"`
	Mobile       string    `json:"mobile"`
	Email        *string   `json:"email"`
	IsActive     int       `json:"is_active"`
	IsLock       int       `json:"is_lock"`
	CreatedAt    *string   `json:"created_at"`
	AreaID       int       `json:"area_id"`
	Level        int       `json:"level"`
	Credit       int       `json:"credit"`
	Coins        int       `json:"coins"`
	Area         *areaView `json:"area"`
	CoursesCount int       `json:"courses_count"`
	PathsCount   int       `json:"paths_count"`
	LikesCount   int       `json:"likes_count"`
	CertsCount   int       `json:"certs_count"`
}

type detailView struct {
	ID           int       `json:"id"`
	Avatar       string    `json:"avatar"`
	NickName     string    `json:"nick_name"`
	Mobile       string    `json:"mobile"`
	IsActive     int       `json:"is_active"`
	IsLock       int       `json:"is_lock"`
	CreatedAt    *string   `json:"created_at"`
	UpdatedAt    *string   `json:"updated_at"`
	AreaID       int       `json:"area_id"`
	Level        int       `json:"level"`
	Credit       int       `json:"credit"`
	Coins        int       `json:"coins"`
	Account      *string   `json:"account"`
	Email        *string   `json:"email"`
	Area         *areaView `json:"area"`
	CoursesCount int       `json:"courses_count"`
	CertsCount   int       `json:"certs_count"`
}

func (s *Service) GetUserList(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Account  string   `json:"account"`
		Email    string   `json:"email"`
		AreaID   bind.Int `json:"area_id"`
		IsActive any      `json:"is_active"`
		IsLock   any      `json:"is_lock"`
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
	query := s.db.WithContext(ctx).Model(&userRow{})
	if limited {
		query = query.Where("city_id = ?", cityID)
	}
	if body.NickName != "" {
		query = query.Where("nick_name LIKE ?", "%"+body.NickName+"%")
	}
	if body.Account != "" {
		query = query.Where("account LIKE ?", "%"+body.Account+"%")
	}
	if body.Email != "" {
		query = query.Where("email LIKE ?", "%"+body.Email+"%")
	}
	if body.AreaID.Int() > 0 {
		query = query.Where("area_id = ?", body.AreaID.Int())
	}
	if bind.Present(body.IsActive) {
		if v := bind.AsInt(body.IsActive); v != 0 {
			query = query.Where("is_active = ?", v)
		}
	}
	if bind.Present(body.IsLock) {
		if bind.AsInt(body.IsLock) == 1 {
			query = query.Where("is_lock = ?", 1)
		} else {
			query = query.Where("(is_lock != ? OR is_lock IS NULL)", 1)
		}
	}
	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []userRow
	err = query.Select("id, avatar, nick_name, account, mobile, is_active, is_lock, created_at, area_id, level, credit, coins, email").
		Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	list, err := s.decorateList(ctx, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserDetail(c *gin.Context) {
	id := atoi(c.Query("id"))
	if id <= 0 {
		response.Fail(c, 1, "学员不存在")
		return
	}
	ctx := c.Request.Context()
	var row userRow
	err := s.db.WithContext(ctx).
		Select("id, avatar, nick_name, mobile, is_active, is_lock, created_at, updated_at, area_id, level, credit, coins, account, email").
		Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "学员不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	area, err := s.loadArea(ctx, row.AreaID, false)
	if response.WriteErr(c, err) {
		return
	}
	courses, err := s.countUserCourses(ctx, []int{id})
	if response.WriteErr(c, err) {
		return
	}
	certs, err := s.countUserCerts(ctx, []int{id})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", detailView{
		ID: row.ID, Avatar: row.Avatar, NickName: row.NickName, Mobile: row.Mobile,
		IsActive: row.IsActive, IsLock: row.IsLock,
		CreatedAt: model.FormatTime(row.CreatedAt), UpdatedAt: model.FormatTime(row.UpdatedAt),
		AreaID: row.AreaID, Level: row.Level, Credit: row.Credit, Coins: row.Coins,
		Account: row.Account, Email: row.Email, Area: area,
		CoursesCount: courses[id], CertsCount: certs[id],
	})
}

func (s *Service) CreateUser(c *gin.Context) {
	var body struct {
		NickName   string   `json:"nick_name"`
		Account    string   `json:"account"`
		Mobile     string   `json:"mobile"`
		AreaID     bind.Int `json:"area_id"`
		CategoryID bind.Int `json:"category_id"`
		Credit     bind.Int `json:"credit"`
		Coins      bind.Int `json:"coins"`
		Password   string   `json:"password"`
		Email      string   `json:"email"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.Password == "" {
		response.Fail(c, 1, "密码不能为空")
		return
	}
	ctx := c.Request.Context()
	closed, err := userCreateClosed(ctx, s.db, s.rdb)
	if response.WriteErr(c, err) {
		return
	}
	if closed {
		response.Fail(c, 1, "当前不在添加学员的时间范围内")
		return
	}
	if body.Mobile != "" {
		var count int64
		if err := s.db.WithContext(ctx).Model(&userRow{}).Where("mobile = ?", body.Mobile).Count(&count).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if count > 0 {
			response.Fail(c, 1, "手机号已存在")
			return
		}
	}
	hash, err := auth.HashPassword(body.Password)
	if response.WriteErr(c, err) {
		return
	}
	credit, ok, err := s.members.TopCategoryBaseCredit(ctx, body.CategoryID.Int())
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.Fail(c, 1, "请选择学员等级")
		return
	}
	cityID, cityCode, err := s.members.FillCityFromArea(ctx, body.AreaID.Int())
	if response.WriteErr(c, err) {
		return
	}
	now := time.Now()
	row := userRow{
		Avatar: "", NickName: body.NickName, Mobile: body.Mobile, Password: hash, PasswordBF: "",
		IsActive: 1, IsLock: -1, AreaID: body.AreaID.Int(), Level: 1, Credit: credit, Coins: body.Coins.Int(),
		CreatedAt: &now, UpdatedAt: &now,
	}
	if body.Email != "" {
		row.Email = &body.Email
	}
	if cityID != nil {
		row.CityID = cityID
	}
	if cityCode != nil {
		row.CityCode = cityCode
	}
	const maxAttempts = 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		account, err := randomAccount()
		if response.WriteErr(c, err) {
			return
		}
		var exists int64
		if err := s.db.WithContext(ctx).Model(&userRow{}).Where("account = ?", account).Count(&exists).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if exists > 0 {
			continue
		}
		row.Account = &account
		err = s.db.WithContext(ctx).Create(&row).Error
		if err == nil {
			response.OK(c, "新增成功", []any{})
			return
		}
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			continue
		}
		response.WriteErr(c, err)
		return
	}
	response.Fail(c, 1, "賬號生成失敗，請稍後重試")
}

func (s *Service) UpdateUser(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	id := bind.AsInt(raw["id"])
	if id <= 0 {
		response.Fail(c, 1, "学员不存在")
		return
	}
	ctx := c.Request.Context()
	var row userRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "学员不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	mobile, _ := raw["mobile"].(string)
	if mobile != "" {
		var count int64
		if err := s.db.WithContext(ctx).Model(&userRow{}).Where("mobile = ? AND id <> ?", mobile, id).Count(&count).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if count > 0 {
			response.Fail(c, 1, "手机号已存在")
			return
		}
	}
	updates := map[string]any{
		"is_active":  1,
		"updated_at": time.Now(),
	}
	if v, ok := raw["nick_name"]; ok {
		if s, ok := v.(string); ok {
			updates["nick_name"] = s
		}
	}
	if _, ok := raw["mobile"]; ok {
		updates["mobile"] = mobile
	}
	if v, ok := raw["email"]; ok {
		if s, ok := v.(string); ok {
			updates["email"] = s
		} else if v == nil {
			updates["email"] = ""
		}
	}
	if _, ok := raw["area_id"]; ok {
		areaID := bind.AsInt(raw["area_id"])
		updates["area_id"] = areaID
		cityID, cityCode, err := s.members.FillCityFromArea(ctx, areaID)
		if response.WriteErr(c, err) {
			return
		}
		if cityID != nil {
			updates["city_id"] = *cityID
		}
		if cityCode != nil {
			updates["city_code"] = *cityCode
		}
	}
	if catID := bind.AsInt(raw["category_id"]); catID > 0 {
		next, ok, err := s.members.TopCategoryBaseCredit(ctx, catID)
		if response.WriteErr(c, err) {
			return
		}
		if !ok {
			response.Fail(c, 1, "请选择学员等级")
			return
		}
		// 当前学分已高于该等级基础学分时保留，避免重选等级把学分压低
		if row.Credit <= next {
			updates["credit"] = next
		}
	} else if _, ok := raw["credit"]; ok {
		updates["credit"] = bind.AsInt(raw["credit"])
	}
	if _, ok := raw["coins"]; ok {
		updates["coins"] = bind.AsInt(raw["coins"])
	}
	if password, _ := raw["password"].(string); password != "" {
		if len(password) < 8 {
			response.Fail(c, 1, "密码长度不能少于8位")
			return
		}
		if !hasLetter.MatchString(password) {
			response.Fail(c, 1, "密码必须包含至少一个字母")
			return
		}
		if !hasDigit.MatchString(password) {
			response.Fail(c, 1, "密码必须包含至少一个数字")
			return
		}
		hash, err := auth.HashPassword(password)
		if response.WriteErr(c, err) {
			return
		}
		updates["password"] = hash
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Model(&userRow{}).Where("id = ?", id).Updates(updates).Error) {
		return
	}
	response.OK(c, "编辑成功", []any{})
}

func (s *Service) SetLock(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id"`
		IsLock bind.Int `json:"is_lock"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 || (body.IsLock.Int() != 1 && body.IsLock.Int() != -1) {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	var row userRow
	err := s.db.WithContext(ctx).Where("id = ?", body.ID.Int()).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "学员不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	now := time.Now()
	if response.WriteErr(c, s.db.WithContext(ctx).Model(&row).Updates(map[string]any{
		"is_lock": body.IsLock.Int(), "updated_at": now,
	}).Error) {
		return
	}
	cache.Drop(ctx, s.rdb, "user_is_lock_"+strconv.Itoa(body.ID.Int()))
	msg := "已解除拉黑"
	if body.IsLock.Int() == 1 {
		msg = "已拉黑该学员"
	}
	response.OK(c, msg, []any{})
}

func (s *Service) DeleteUser(c *gin.Context) {
	id := atoi(c.Query("id"))
	if id <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	var row userRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "学员不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if err := s.db.WithContext(ctx).Delete(&row).Error; err != nil {
		response.Fail(c, 1, "删除失败："+err.Error())
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) decorateList(ctx context.Context, rows []userRow) ([]listView, error) {
	list := make([]listView, 0, len(rows))
	if len(rows) == 0 {
		return list, nil
	}
	ids := make([]int, 0, len(rows))
	areaIDs := make([]int, 0, len(rows))
	areaSeen := map[int]bool{}
	for _, row := range rows {
		ids = append(ids, row.ID)
		if row.AreaID > 0 && !areaSeen[row.AreaID] {
			areaSeen[row.AreaID] = true
			areaIDs = append(areaIDs, row.AreaID)
		}
	}
	areas, err := s.loadAreas(ctx, areaIDs, true)
	if err != nil {
		return nil, err
	}
	courses, err := s.countUserCourses(ctx, ids)
	if err != nil {
		return nil, err
	}
	paths, err := s.countUserPaths(ctx, ids)
	if err != nil {
		return nil, err
	}
	likes, err := s.countUserLikes(ctx, ids)
	if err != nil {
		return nil, err
	}
	certs, err := s.countUserCerts(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		var area *areaView
		if a, ok := areas[row.AreaID]; ok {
			cp := a
			area = &cp
		}
		list = append(list, listView{
			ID: row.ID, Avatar: row.Avatar, NickName: row.NickName, Account: row.Account,
			Mobile: row.Mobile, Email: row.Email, IsActive: row.IsActive, IsLock: row.IsLock,
			CreatedAt: model.FormatTime(row.CreatedAt), AreaID: row.AreaID, Level: row.Level,
			Credit: row.Credit, Coins: row.Coins, Area: area,
			CoursesCount: courses[row.ID], PathsCount: paths[row.ID],
			LikesCount: likes[row.ID], CertsCount: certs[row.ID],
		})
	}
	return list, nil
}

func (s *Service) loadArea(ctx context.Context, id int, withParent bool) (*areaView, error) {
	if id <= 0 {
		return nil, nil
	}
	m, err := s.loadAreas(ctx, []int{id}, withParent)
	if err != nil {
		return nil, err
	}
	if a, ok := m[id]; ok {
		cp := a
		return &cp, nil
	}
	return nil, nil
}

func (s *Service) loadAreas(ctx context.Context, ids []int, withParent bool) (map[int]areaView, error) {
	out := map[int]areaView{}
	rows, err := s.members.ListAreasByIDs(ctx, ids, withParent)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		v := areaView{ID: row.ID, Name: row.Name}
		if withParent {
			v.ParentName = row.ParentName
		}
		out[row.ID] = v
	}
	return out, nil
}

func (s *Service) countUserCourses(ctx context.Context, ids []int) (map[int]int, error) {
	return s.countByUser(ctx, `
		SELECT uc.user_id AS uid, COUNT(*) AS total
		FROM user_course uc
		INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
		WHERE uc.user_id IN ?
		GROUP BY uc.user_id
	`, ids)
}

func (s *Service) countUserPaths(ctx context.Context, ids []int) (map[int]int, error) {
	return s.countByUser(ctx, `
		SELECT up.user_id AS uid, COUNT(*) AS total
		FROM user_path up
		INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
		WHERE up.user_id IN ?
		GROUP BY up.user_id
	`, ids)
}

func (s *Service) countUserLikes(ctx context.Context, ids []int) (map[int]int, error) {
	return s.countByUser(ctx, `
		SELECT user_id AS uid, COUNT(*) AS total
		FROM user_like_courses
		WHERE user_id IN ?
		GROUP BY user_id
	`, ids)
}

func (s *Service) countUserCerts(ctx context.Context, ids []int) (map[int]int, error) {
	return s.countByUser(ctx, `
		SELECT user_id AS uid, COUNT(*) AS total
		FROM user_cert
		WHERE user_id IN ?
		GROUP BY user_id
	`, ids)
}

func (s *Service) countByUser(ctx context.Context, sql string, ids []int) (map[int]int, error) {
	out := map[int]int{}
	if len(ids) == 0 {
		return out, nil
	}
	type countRow struct {
		UID   int `gorm:"column:uid"`
		Total int `gorm:"column:total"`
	}
	var rows []countRow
	if err := s.db.WithContext(ctx).Raw(sql, ids).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.UID] = row.Total
	}
	return out, nil
}

func userCreateClosed(ctx context.Context, db *gorm.DB, rdb *redis.Client) (bool, error) {
	raw, err := appconfig.Get(ctx, db, rdb, "memberRule")
	if err != nil {
		return false, err
	}
	start, end := "", ""
	if m, ok := raw.(map[string]any); ok {
		start = strings.TrimSpace(asString(m["register_start"]))
		end = strings.TrimSpace(asString(m["register_end"]))
	}
	now := time.Now()
	if start != "" {
		if t, ok := parseFlexibleTime(start); ok && now.Before(t) {
			return true, nil
		}
	}
	if end != "" {
		if t, ok := parseFlexibleTime(end); ok && now.After(t) {
			return true, nil
		}
	}
	return false, nil
}

func randomAccount() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(999000000))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(n.Int64()+1000000, 10), nil
}

func parseFlexibleTime(s string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func asString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func atoi(s string) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
