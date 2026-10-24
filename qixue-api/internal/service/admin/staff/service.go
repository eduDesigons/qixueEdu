package adminstaff

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/scope"
)

var (
	ErrBadParams      = errors.New("参数错误")
	ErrRoleNotFound   = errors.New("角色不存在")
	ErrUserExists     = errors.New("用户已存在")
	ErrUserNotFound   = errors.New("用户不存在")
	ErrNoCityAssign   = errors.New("无权分配大区")
	ErrInvalidCity    = errors.New("请选择有效大区")
	ErrSuperProtected = errors.New("超级管理员不能删除")
)

type Service struct {
	staff repo.AdminStaffRepository
	db    *gorm.DB
	rdb   *redis.Client
}

func New(staff repo.AdminStaffRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{staff: staff, db: db, rdb: rdb}
}

type RoleView struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Enable    int     `json:"enable"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type CityView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type View struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Enable    int        `json:"enable"`
	NickName  *string    `json:"nickName"`
	Avatar    string     `json:"avatar"`
	CityID    *int       `json:"city_id"`
	IsSuper   int        `json:"is_super"`
	CreatedAt *string    `json:"created_at"`
	UpdatedAt *string    `json:"updated_at"`
	Roles     []RoleView `json:"roles"`
	City      *CityView  `json:"city"`
}

type ListQuery struct {
	Username string
	Enable   int
	CityID   int
	Page     int
	Size     int
}

type AddInput struct {
	Username string
	Password string
	Enable   int
	RoleIDs  []int
	CityID   *int
	HasCity  bool
}

type EditInput struct {
	ID              int
	Username        *string
	Enable          *int
	CityProvided    bool
	ClearCity       bool
	CityID          int
	RoleIDsProvided bool
	RoleIDs         []int
}

func (s *Service) List(ctx context.Context, adminID int, q ListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	query := s.db.WithContext(ctx).Model(&model.Administrator{})
	if q.Username != "" {
		query = query.Where("username LIKE ?", "%"+q.Username+"%")
	}
	if q.Enable != 0 {
		query = query.Where("enable = ?", q.Enable)
	}
	if limited {
		query = query.Where("city_id = ?", cityID)
	} else if q.CityID > 0 {
		query = query.Where("city_id = ?", q.CityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []model.Administrator
	if err := query.Order("created_at DESC, id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list, err := s.decorate(ctx, rows)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) ChangeRole(ctx context.Context, userIDs []int, roleID int, attach bool) error {
	count, err := s.staff.CountRole(ctx, roleID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrRoleNotFound
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, uid := range userIDs {
			if uid <= 0 {
				continue
			}
			if attach {
				if err := tx.Exec("INSERT IGNORE INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", uid, roleID).Error; err != nil {
					return err
				}
				continue
			}
			if err := tx.Exec("DELETE FROM administrator_role_relation WHERE administratorId = ? AND roleId = ?", uid, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) Add(ctx context.Context, adminID int, in AddInput) error {
	if strings.TrimSpace(in.Username) == "" || in.Password == "" {
		return ErrBadParams
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("username = ?", in.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExists
	}
	ownCity, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return err
	}
	var cityID *int
	if limited {
		cityID = &ownCity
	} else if in.HasCity && in.CityID != nil && *in.CityID > 0 {
		if err := s.validCity(ctx, *in.CityID); err != nil {
			return err
		}
		cityID = in.CityID
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		item := model.Administrator{
			Username: in.Username, Password: hash, Enable: in.Enable, CityID: cityID,
			CreatedAt: &now, UpdatedAt: &now,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		for _, roleID := range in.RoleIDs {
			if roleID <= 0 {
				continue
			}
			if err := tx.Exec("INSERT INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", item.ID, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) Edit(ctx context.Context, adminID int, in EditInput) error {
	if in.ID <= 0 {
		return ErrUserNotFound
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("id = ?", in.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotFound
	}
	updates := map[string]any{}
	if in.Username != nil && *in.Username != "" {
		if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("username = ? AND id <> ?", *in.Username, in.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrUserExists
		}
		updates["username"] = *in.Username
	}
	if in.Enable != nil {
		updates["enable"] = *in.Enable
	}
	if in.CityProvided {
		_, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
		if err != nil {
			return err
		}
		if limited {
			return ErrNoCityAssign
		}
		if in.ClearCity {
			updates["city_id"] = nil
		} else {
			if err := s.validCity(ctx, in.CityID); err != nil {
				return err
			}
			updates["city_id"] = in.CityID
		}
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if in.RoleIDsProvided {
			if err := tx.Exec("DELETE FROM administrator_role_relation WHERE administratorId = ?", in.ID).Error; err != nil {
				return err
			}
			for _, roleID := range in.RoleIDs {
				if roleID <= 0 {
					continue
				}
				if err := tx.Exec("INSERT INTO administrator_role_relation (administratorId, roleId) VALUES (?, ?)", in.ID, roleID).Error; err != nil {
					return err
				}
			}
		}
		if len(updates) == 0 {
			return nil
		}
		return tx.Model(&model.Administrator{}).Where("id = ?", in.ID).Updates(updates).Error
	})
	if err != nil {
		return err
	}
	cache.Drop(ctx, s.rdb, "admin_"+strconv.Itoa(in.ID))
	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	var item model.Administrator
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if item.IsSuper != 0 {
		return ErrSuperProtected
	}
	if err := s.db.WithContext(ctx).Delete(&model.Administrator{}, id).Error; err != nil {
		return err
	}
	cache.Drop(ctx, s.rdb, "admin_"+strconv.Itoa(id))
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, id int, password string) error {
	if password == "" {
		return ErrBadParams
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	result := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("id = ?", id).Update("password", hash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *Service) validCity(ctx context.Context, id int) error {
	count, err := s.staff.CountRootArea(ctx, id)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidCity
	}
	return nil
}

func (s *Service) decorate(ctx context.Context, rows []model.Administrator) ([]View, error) {
	if len(rows) == 0 {
		return []View{}, nil
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
	err := s.db.WithContext(ctx).Raw(`
		SELECT rel.administratorId, role.id, role.code, role.name, role.enable, role.created_at, role.updated_at
		FROM administrator_role_relation rel
		JOIN administrator_roles role ON role.id = rel.roleId
		WHERE rel.administratorId IN ?`, ids).Scan(&links).Error
	if err != nil {
		return nil, err
	}
	roles := map[int][]RoleView{}
	for _, item := range links {
		roles[item.AdminID] = append(roles[item.AdminID], RoleView{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
		})
	}
	cities := map[int]CityView{}
	if len(cityIDs) > 0 {
		idList := make([]int, 0, len(cityIDs))
		for id := range cityIDs {
			idList = append(idList, id)
		}
		found, err := s.staff.ListAreasByIDs(ctx, idList)
		if err != nil {
			return nil, err
		}
		for _, item := range found {
			cities[item.ID] = CityView{ID: item.ID, Name: item.Name}
		}
	}
	list := make([]View, 0, len(rows))
	for _, item := range rows {
		itemRoles := roles[item.ID]
		if itemRoles == nil {
			itemRoles = []RoleView{}
		}
		var city *CityView
		if item.CityID != nil {
			if found, ok := cities[*item.CityID]; ok {
				cp := found
				city = &cp
			}
		}
		list = append(list, View{
			ID: item.ID, Username: item.Username, Enable: item.Enable, NickName: item.NickName,
			Avatar: item.Avatar, CityID: item.CityID, IsSuper: item.IsSuper,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
			Roles: itemRoles, City: city,
		})
	}
	return list, nil
}
