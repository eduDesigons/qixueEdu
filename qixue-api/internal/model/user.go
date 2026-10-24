package model

import "time"

type User struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	Avatar       string     `gorm:"column:avatar"`
	NickName     string     `gorm:"column:nick_name"`
	Mobile       string     `gorm:"column:mobile"`
	Password     string     `gorm:"column:password"`
	Passwordbf   string     `gorm:"column:passwordbf"`
	IsActive     int        `gorm:"column:is_active"`
	IsLock       int        `gorm:"column:is_lock"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	InviteUserID int        `gorm:"column:invite_user_id"`
	RegisterIP   string     `gorm:"column:register_ip"`
	AreaID       int        `gorm:"column:area_id"`
	CityID       *int       `gorm:"column:city_id"`
	CityCode     *string    `gorm:"column:city_code"`
	LastLoginID  int        `gorm:"column:last_login_id"`
	Level        int        `gorm:"column:level"`
	Credit       int        `gorm:"column:credit"`
	Coins        int        `gorm:"column:coins"`
	Account      *string    `gorm:"column:account"`
	Gender       string     `gorm:"column:gender"`
	Age          int        `gorm:"column:age"`
	Degree       string     `gorm:"column:degree"`
	Profession   string     `gorm:"column:profession"`
	Job          string     `gorm:"column:job"`
	Marital      *string    `gorm:"column:marital"`
	Target       *string    `gorm:"column:target"`
	Email        *string    `gorm:"column:email"`
}

func (User) TableName() string { return "users" }
