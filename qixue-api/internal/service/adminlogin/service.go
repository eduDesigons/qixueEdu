package adminlogin

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/menutree"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/captcha"
)

var (
	ErrBadCaptcha = errors.New("bad captcha")
	ErrBlocked    = errors.New("blocked")
	ErrPassword   = errors.New("bad password")
	ErrDisabled   = errors.New("disabled")
	ErrNotFound   = errors.New("not found")
	ErrEmpty      = errors.New("empty")
	ErrNickName   = errors.New("nickname")
)

type Service struct {
	db    *gorm.DB
	rdb   *redis.Client
	jwt   auth.Issuer
	slide *captcha.Slide
	limit *auth.Limiter
	login repo.AdminLoginRepository
}

func New(db *gorm.DB, rdb *redis.Client, jwt auth.Issuer, slide *captcha.Slide) *Service {
	return &Service{
		db:    db,
		rdb:   rdb,
		jwt:   jwt,
		slide: slide,
		limit: auth.NewLimiter(rdb),
		login: repo.NewAdminLogin(db),
	}
}

type Profile struct {
	ID         int     `json:"id"`
	Gender     *int    `json:"gender"`
	Enable     int     `json:"enable"`
	NickName   *string `json:"nickName"`
	Avatar     string  `json:"avatar"`
	Address    *string `json:"address"`
	Email      *string `json:"email"`
	Username   string  `json:"username"`
	CityID     *int    `json:"city_id"`
	CreateTime *string `json:"createTime"`
	UpdateTime *string `json:"updateTime"`
}

type ProfilePatch struct {
	NickName   *string
	Gender     *int
	Address    *string
	Email      *string
	Avatar     *string
	HasNick    bool
	HasGender  bool
	HasAddress bool
	HasEmail   bool
	HasAvatar  bool
}

func (s *Service) Captcha(ctx context.Context) (captcha.View, error) {
	return s.slide.Issue(ctx)
}

func (s *Service) CheckCaptcha(ctx context.Context, token string, x, y int) (string, error) {
	pass, err := s.slide.Verify(ctx, token, x, y)
	if errors.Is(err, captcha.ErrNotFound) || errors.Is(err, captcha.ErrMismatch) {
		return "", ErrBadCaptcha
	}
	return pass, err
}

func (s *Service) Login(ctx context.Context, username, password, pass string) (string, error) {
	ok, err := s.slide.Consume(ctx, pass)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrBadCaptcha
	}
	if s.limit.Blocked(ctx, "admin", username) {
		return "", ErrBlocked
	}
	admin, err := s.findByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	hash := ""
	if admin != nil {
		hash = admin.Password
	}
	if admin == nil || !auth.PasswordOK(hash, password) {
		s.limit.Record(ctx, "admin", username)
		return "", ErrPassword
	}
	if admin.Enable != 1 {
		return "", ErrDisabled
	}
	s.limit.Clear(ctx, "admin", username)
	token, _, err := s.jwt.Sign(admin.ID, time.Now())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) UserInfo(ctx context.Context, id int) (Profile, error) {
	admin, err := s.findByID(ctx, id)
	if err != nil {
		return Profile{}, err
	}
	if admin == nil {
		return Profile{}, ErrNotFound
	}
	return toProfile(*admin), nil
}

func (s *Service) Permissions(ctx context.Context, id int) ([]menutree.Node, error) {
	super, err := s.isSuper(ctx, id)
	if err != nil {
		return nil, err
	}
	var rows []model.AdminMenu
	query := s.db.WithContext(ctx).Model(&model.AdminMenu{}).Order("`order` ASC, id ASC")
	if !super {
		query = query.Where(
			"id IN (SELECT permissionId FROM administrator_permissions WHERE roleId IN (SELECT roleId FROM administrator_role_relation WHERE administratorId = ?))",
			id,
		)
		roleCount, err := s.login.CountRoles(ctx, id)
		if err != nil {
			return nil, err
		}
		if roleCount == 0 {
			return []menutree.Node{}, nil
		}
	}
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return menutree.Tree(rows), nil
}

func (s *Service) UpdateProfile(ctx context.Context, id int, patch ProfilePatch) error {
	if !patch.HasNick && !patch.HasGender && !patch.HasAddress && !patch.HasEmail && !patch.HasAvatar {
		return ErrEmpty
	}
	admin, err := s.findByID(ctx, id)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrNotFound
	}
	updates := map[string]any{}
	if patch.HasNick {
		name := ""
		if patch.NickName != nil {
			name = strings.TrimSpace(*patch.NickName)
		}
		if utf8.RuneCountInString(name) > 10 {
			return ErrNickName
		}
		updates["nickName"] = name
	}
	if patch.HasGender {
		if patch.Gender == nil {
			updates["gender"] = nil
		} else {
			updates["gender"] = *patch.Gender
		}
	}
	if patch.HasAddress {
		updates["address"] = stringOrNil(patch.Address)
	}
	if patch.HasEmail {
		updates["email"] = stringOrNil(patch.Email)
	}
	if patch.HasAvatar {
		updates["avatar"] = stringOrNil(patch.Avatar)
	}
	return s.db.WithContext(ctx).Model(&model.Administrator{}).Where("id = ?", id).Updates(updates).Error
}

func (s *Service) ChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error {
	admin, err := s.findByID(ctx, id)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrNotFound
	}
	if !auth.PasswordOK(admin.Password, oldPassword) {
		return ErrPassword
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Model(&model.Administrator{}).Where("id = ?", id).Update("password", hash).Error; err != nil {
		return err
	}
	cache.Drop(ctx, s.rdb, "admin_"+strconv.Itoa(id))
	return nil
}

func (s *Service) findByUsername(ctx context.Context, username string) (*model.Administrator, error) {
	var admin model.Administrator
	err := s.db.WithContext(ctx).Where("username = ?", username).Take(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *Service) findByID(ctx context.Context, id int) (*model.Administrator, error) {
	var admin model.Administrator
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *Service) isSuper(ctx context.Context, id int) (bool, error) {
	return s.login.IsSuper(ctx, id)
}

func stringOrNil(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func toProfile(admin model.Administrator) Profile {
	return Profile{
		ID:         admin.ID,
		Gender:     admin.Gender,
		Enable:     admin.Enable,
		NickName:   admin.NickName,
		Avatar:     admin.Avatar,
		Address:    admin.Address,
		Email:      admin.Email,
		Username:   admin.Username,
		CityID:     admin.CityID,
		CreateTime: model.FormatTime(admin.CreatedAt),
		UpdateTime: model.FormatTime(admin.UpdatedAt),
	}
}
