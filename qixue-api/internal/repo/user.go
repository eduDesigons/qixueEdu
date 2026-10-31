package repo

import (
	"context"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (*model.User, error)
	FindProfile(ctx context.Context, id int) (map[string]any, error)
	Exists(ctx context.Context, id int) (bool, error)
	Update(ctx context.Context, id int, fields map[string]any) error
	FindPassword(ctx context.Context, id int) (hash string, err error)
	CreditCoins(ctx context.Context, id int) (credit, coins int, err error)
	FindLoginByAccount(ctx context.Context, account string) (*LoginUser, error)
}

// LoginUser is the subset of users needed for webapp login.
type LoginUser struct {
	ID       int     `gorm:"column:id"`
	Account  *string `gorm:"column:account"`
	Password string  `gorm:"column:password"`
	CityCode *string `gorm:"column:city_code"`
	IsActive int     `gorm:"column:is_active"`
	IsLock   int     `gorm:"column:is_lock"`
}

type userRepo struct{ db *gorm.DB }

func NewUser(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id int) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindProfile(ctx context.Context, id int) (map[string]any, error) {
	var row map[string]any
	err := r.db.WithContext(ctx).Table("users").
		Select("id, avatar, nick_name, account, mobile, gender, age, degree, profession, job, marital, target, coins").
		Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *userRepo) Exists(ctx context.Context, id int) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&n).Error
	return n > 0, err
}

func (r *userRepo) Update(ctx context.Context, id int, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	if _, ok := fields["updated_at"]; !ok {
		fields["updated_at"] = time.Now()
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *userRepo) FindPassword(ctx context.Context, id int) (string, error) {
	var u model.User
	err := r.db.WithContext(ctx).Select("id, password").Where("id = ?", id).Take(&u).Error
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

func (r *userRepo) CreditCoins(ctx context.Context, id int) (int, int, error) {
	var u model.User
	err := r.db.WithContext(ctx).Select("credit, coins").Where("id = ?", id).Take(&u).Error
	if err != nil {
		return 0, 0, err
	}
	return u.Credit, u.Coins, nil
}

func (r *userRepo) FindLoginByAccount(ctx context.Context, account string) (*LoginUser, error) {
	var row LoginUser
	err := r.db.WithContext(ctx).Table("users").
		Select("id, account, password, city_code, is_active, is_lock").
		Where("account = ?", account).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}
