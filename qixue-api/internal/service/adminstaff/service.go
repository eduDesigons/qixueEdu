package adminstaff

import (
	"strconv"
	"strings"
	"time"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/scope"
)

type Service struct {
	staff repo.AdminStaffRepository
	db    *gorm.DB
	rdb   *redis.Client
}

func New(staff repo.AdminStaffRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{staff: staff, db: db, rdb: rdb}
}

type roleView struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Enable    int     `json:"enable"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type cityView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type view struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Enable    int        `json:"enable"`
	NickName  *string    `json:"nickName"`
	Avatar    string     `json:"avatar"`
	CityID    *int       `json:"city_id"`
	IsSuper   int        `json:"is_super"`
	CreatedAt *string    `json:"created_at"`
	UpdatedAt *string    `json:"updated_at"`
	Roles     []roleView `json:"roles"`
	City      *cityView  `json:"city"`
}

func (s *Service) List(c *gin.Context) {
	var body struct {
		Username string   `json:"username"`
		Enable   bind.Int `json:"enable"`
		CityID   bind.Int `json:"city_id"`
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
	query := s.db.WithContext(ctx).Model(&model.Administrator{})
	if body.Username != "" {
		query = query.Where("username LIKE ?", "%"+body.Username+"%")
	}
	if body.Enable.Int() != 0 {
		query = query.Where("enable = ?", body.Enable.Int())
	}
	if limited {
		query = query.Where("city_id = ?", cityID)
	} else if body.CityID.Int() > 0 {
		query = query.Where("city_id = ?", body.CityID.Int())
	}
	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []model.Administrator
	err = query.Order("created_at DESC, id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	list, err := s.decorate(c, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "管理员列表", page.New(current, size, total, list))
}

func (s *Service) AddRole(c *gin.Context) {
	s.changeRole(c, true)
}

func (s *Service) RemoveRole(c *gin.Context) {
	s.changeRole(c, false)
}

func (s *Service) changeRole(c *gin.Context, attach bool) {
	var body struct {
		UserIDs []bind.Int `json:"userIds"`
		RoleID  bind.Int   `json:"roleId"`
	}
	_ = c.ShouldBindJSON(&body)
	count, err := s.staff.CountRole(c.Request.Context(), body.RoleID.Int())
	if response.WriteErr(c, err) {
		return
	}
	if count == 0 {
		response.Fail(c, 1, "角色不存在")
		return
	}
	err = s.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		for _, uid := range bind.Ints(body.UserIDs) {
			if uid <= 0 {
				continue
			}
			if attach {
				if err := tx.Exec("INSERT IGNORE INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", uid, body.RoleID.Int()).Error; err != nil {
					return err
				}
				continue
			}
			if err := tx.Exec("DELETE FROM administrator_role_relation WHERE administratorId = ? AND roleId = ?", uid, body.RoleID.Int()).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if response.WriteErr(c, err) {
		return
	}
	if attach {
		response.OK(c, "添加成功", []any{})
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) Add(c *gin.Context) {
	var body struct {
		Username string     `json:"username"`
		Password string     `json:"password"`
		Enable   bind.Int   `json:"enable"`
		RoleIDs  []bind.Int `json:"roleIds"`
		CityID   any        `json:"city_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Username) == "" || body.Password == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("username = ?", body.Username).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count > 0 {
		response.Fail(c, 1, "用户已存在")
		return
	}
	ownCity, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	var cityID *int
	if limited {
		cityID = &ownCity
	} else if bind.Present(body.CityID) {
		if cid := bind.AsInt(body.CityID); cid > 0 {
			if err := s.validCity(c, cid); err != nil {
				response.WriteErr(c, err)
				return
			}
			cityID = &cid
		}
	}
	hash, err := auth.HashPassword(body.Password)
	if response.WriteErr(c, err) {
		return
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		item := model.Administrator{
			Username: body.Username, Password: hash, Enable: body.Enable.Int(), CityID: cityID,
			CreatedAt: &now, UpdatedAt: &now,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		for _, roleID := range bind.Ints(body.RoleIDs) {
			if roleID <= 0 {
				continue
			}
			if err := tx.Exec("INSERT INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", item.ID, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "添加成功", []any{})
}

func (s *Service) Edit(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	id := asInt(raw["id"])
	if id <= 0 {
		response.Fail(c, 1, "用户不存在")
		return
	}
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("id = ?", id).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "用户不存在")
		return
	}
	updates := map[string]any{}
	if name, ok := raw["username"].(string); ok && name != "" {
		if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("username = ? AND id <> ?", name, id).Count(&count).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if count > 0 {
			response.Fail(c, 1, "用户已存在")
			return
		}
		updates["username"] = name
	}
	if _, ok := raw["enable"]; ok && raw["enable"] != nil && raw["enable"] != "" {
		updates["enable"] = asInt(raw["enable"])
	}
	if _, ok := raw["city_id"]; ok {
		own, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
		if response.WriteErr(c, err) {
			return
		}
		if limited {
			_ = own
			response.Fail(c, 1, "无权分配大区")
			return
		}
		if raw["city_id"] == nil || raw["city_id"] == "" {
			updates["city_id"] = nil
		} else {
			cityID := asInt(raw["city_id"])
			if err := s.validCity(c, cityID); response.WriteErr(c, err) {
				return
			}
			updates["city_id"] = cityID
		}
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, ok := raw["roleIds"]; ok {
			if err := tx.Exec("DELETE FROM administrator_role_relation WHERE administratorId = ?", id).Error; err != nil {
				return err
			}
			for _, roleID := range asInts(raw["roleIds"]) {
				if roleID <= 0 {
					continue
				}
				if err := tx.Exec("INSERT INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", id, roleID).Error; err != nil {
					return err
				}
			}
		}
		if len(updates) == 0 {
			return nil
		}
		return tx.Model(&model.Administrator{}).Where("id = ?", id).Updates(updates).Error
	})
	if response.WriteErr(c, err) {
		return
	}
	cache.Drop(ctx, s.rdb, "admin_"+strconv.Itoa(id))
	response.OK(c, "编辑成功", []any{})
}

func (s *Service) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	var item model.Administrator
	err := s.db.WithContext(c.Request.Context()).Where("id = ?", body.ID.Int()).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "用户不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if item.IsSuper != 0 {
		response.Fail(c, 1, "超级管理员不能删除")
		return
	}
	if response.WriteErr(c, s.db.WithContext(c.Request.Context()).Delete(&model.Administrator{}, body.ID.Int()).Error) {
		return
	}
	cache.Drop(c.Request.Context(), s.rdb, "admin_"+strconv.Itoa(body.ID.Int()))
	response.OK(c, "删除成功", []any{})
}

func (s *Service) ResetPassword(c *gin.Context) {
	var body struct {
		ID       bind.Int `json:"id"`
		Password string   `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	hash, err := auth.HashPassword(body.Password)
	if response.WriteErr(c, err) {
		return
	}
	result := s.db.WithContext(c.Request.Context()).Model(&model.Administrator{}).Where("id = ?", body.ID.Int()).Update("password", hash)
	if response.WriteErr(c, result.Error) {
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, 1, "用户不存在")
		return
	}
	response.OK(c, "重置成功", []any{})
}

func (s *Service) validCity(c *gin.Context, id int) error {
	count, err := s.staff.CountRootArea(c.Request.Context(), id)
	if err != nil {
		return err
	}
	if count == 0 {
		return response.BizErr("请选择有效大区")
	}
	return nil
}

func (s *Service) decorate(c *gin.Context, rows []model.Administrator) ([]view, error) {
	if len(rows) == 0 {
		return []view{}, nil
	}
	ids := make([]int, 0, len(rows))
	cityIDs := map[int]struct{}{}
	for _, item := range rows {
		ids = append(ids, item.ID)
		if item.CityID != nil && *item.CityID > 0 {
			cityIDs[*item.CityID] = struct{}{}
		}
	}
	type link struct {
		AdminID   int        `gorm:"column:administratorId"`
		ID        int        `gorm:"column:id"`
		Code      string     `gorm:"column:code"`
		Name      string     `gorm:"column:name"`
		Enable    int        `gorm:"column:enable"`
		CreatedAt *time.Time `gorm:"column:created_at"`
		UpdatedAt *time.Time `gorm:"column:updated_at"`
	}
	var links []link
	err := s.db.WithContext(c.Request.Context()).Raw(`
		SELECT rel.administratorId, role.id, role.code, role.name, role.enable, role.created_at, role.updated_at
		FROM administrator_role_relation rel
		JOIN administrator_roles role ON role.id = rel.roleId
		WHERE rel.administratorId IN ?`, ids).Scan(&links).Error
	if err != nil {
		return nil, err
	}
	roles := map[int][]roleView{}
	for _, item := range links {
		roles[item.AdminID] = append(roles[item.AdminID], roleView{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
		})
	}
	cities := map[int]cityView{}
	if len(cityIDs) > 0 {
		idList := make([]int, 0, len(cityIDs))
		for id := range cityIDs {
			idList = append(idList, id)
		}
		found, err := s.staff.ListAreasByIDs(c.Request.Context(), idList)
		if err != nil {
			return nil, err
		}
		for _, item := range found {
			cities[item.ID] = cityView{ID: item.ID, Name: item.Name}
		}
	}
	list := make([]view, 0, len(rows))
	for _, item := range rows {
		itemRoles := roles[item.ID]
		if itemRoles == nil {
			itemRoles = []roleView{}
		}
		var city *cityView
		if item.CityID != nil {
			if found, ok := cities[*item.CityID]; ok {
				copy := found
				city = &copy
			}
		}
		list = append(list, view{
			ID: item.ID, Username: item.Username, Enable: item.Enable, NickName: item.NickName,
			Avatar: item.Avatar, CityID: item.CityID, IsSuper: item.IsSuper,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
			Roles: itemRoles, City: city,
		})
	}
	return list, nil
}

func asInt(value any) int { return bind.AsInt(value) }

func asInts(value any) []int {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]int, 0, len(items))
	for _, item := range items {
		out = append(out, asInt(item))
	}
	return out
}
