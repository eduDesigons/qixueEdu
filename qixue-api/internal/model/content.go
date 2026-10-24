package model

import "time"

type Banner struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name"`
	Sort      int        `gorm:"column:sort"`
	ImgURL    string     `gorm:"column:img_url"`
	LinkURL   *string    `gorm:"column:link_url"`
	Type      int        `gorm:"column:type"`
	Status    int        `gorm:"column:status"`
	Content   *string    `gorm:"column:content"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (Banner) TableName() string { return "banners" }

type Notice struct {
	ID          uint       `gorm:"column:id;primaryKey"`
	Title       string     `gorm:"column:title"`
	Status      int        `gorm:"column:status"`
	PublishTime time.Time  `gorm:"column:publish_time"`
	AdminID     int        `gorm:"column:admin_id"`
	Content     string     `gorm:"column:content"`
	ViewTimes   int        `gorm:"column:view_times"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (Notice) TableName() string { return "notices" }

type MediaFile struct {
	ID              uint       `gorm:"column:id;primaryKey"`
	Title           string     `gorm:"column:title"`
	Thumb           string     `gorm:"column:thumb"`
	Duration        int        `gorm:"column:duration"`
	Size            int64      `gorm:"column:size"`
	StorageFileID   string     `gorm:"column:storage_file_id"`
	TranscodeStatus int        `gorm:"column:transcode_status"`
	RefCount        int        `gorm:"column:ref_count"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
	IsOpen          int        `gorm:"column:is_open"`
	Type            string     `gorm:"column:type"`
	CateID          int        `gorm:"column:cate_id"`
	URL             *string    `gorm:"column:url"`
	Vid             string     `gorm:"column:vid"`
	DownloadURL     *string    `gorm:"column:download_url"`
}

func (MediaFile) TableName() string { return "media_file" }

type MediaCate struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name"`
	Type      *string    `gorm:"column:type"`
	IsInit    int        `gorm:"column:is_init"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	Sort      int        `gorm:"column:sort"`
}

func (MediaCate) TableName() string { return "media_cate" }

type Area struct {
	ID         uint       `gorm:"column:id;primaryKey"`
	Name       string     `gorm:"column:name"`
	ParentID   int        `gorm:"column:parent_id"`
	ParentName string     `gorm:"column:parent_name"`
	Type       int        `gorm:"column:type"`
	IsShow     int        `gorm:"column:is_show"`
	Sort       int        `gorm:"column:sort"`
	CreatedAt  *time.Time `gorm:"column:created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
	Code       *string    `gorm:"column:code"`
}

func (Area) TableName() string { return "areas" }

type Verse struct {
	ID      int     `gorm:"column:id;primaryKey"`
	Content *string `gorm:"column:content"`
}

func (Verse) TableName() string { return "verse" }

type AppConfig struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	Key       string     `gorm:"column:key"`
	Value     *string    `gorm:"column:value"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (AppConfig) TableName() string { return "app_config" }

type RankSnapshot struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	SiteCode     string     `gorm:"column:site_code"`
	Type         string     `gorm:"column:type"`
	Period       string     `gorm:"column:period"`
	PeriodKey    string     `gorm:"column:period_key"`
	PeriodLabel  string     `gorm:"column:period_label"`
	SnapshotDate time.Time  `gorm:"column:snapshot_date"`
	Payload      string     `gorm:"column:payload"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (RankSnapshot) TableName() string { return "rank_snapshots" }
