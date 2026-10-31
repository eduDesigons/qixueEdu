package model

import "time"

type Administrator struct {
	ID        int        `gorm:"column:id;primaryKey"`
	CityID    *int       `gorm:"column:city_id"`
	Username  string     `gorm:"column:username"`
	Password  string     `gorm:"column:password"`
	Enable    int        `gorm:"column:enable"`
	Gender    *int       `gorm:"column:gender"`
	Avatar    string     `gorm:"column:avatar"`
	Address   *string    `gorm:"column:address"`
	Email     *string    `gorm:"column:email"`
	NickName  *string    `gorm:"column:nickName"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	IsSuper   int        `gorm:"column:is_super"`
}

func (Administrator) TableName() string { return "administrator" }

type AdminMenu struct {
	ID          int        `gorm:"column:id;primaryKey"`
	Name        string     `gorm:"column:name"`
	Code        string     `gorm:"column:code"`
	Type        string     `gorm:"column:type"`
	ParentID    *int       `gorm:"column:parentId"`
	Path        *string    `gorm:"column:path"`
	Redirect    *string    `gorm:"column:redirect"`
	Icon        *string    `gorm:"column:icon"`
	Component   *string    `gorm:"column:component"`
	Layout      *string    `gorm:"column:layout"`
	KeepAlive   *int       `gorm:"column:keepAlive"`
	Method      *string    `gorm:"column:method"`
	Description *string    `gorm:"column:description"`
	Show        int        `gorm:"column:show"`
	Enable      int        `gorm:"column:enable"`
	Order       *int       `gorm:"column:order"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (AdminMenu) TableName() string { return "administrator_menus" }
