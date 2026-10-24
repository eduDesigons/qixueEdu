package adminmedia

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/service/bunny"
)

var (
	ErrNoFile      = errors.New("文件不存在或无效")
	ErrBadParams   = errors.New("参数错误")
	ErrSaveParams  = errors.New("落库失败: 参数错误")
	ErrNotFound    = errors.New("资源不存在")
	ErrNeedFileMeta = errors.New("请提供有效的文件名与大小")
	ErrVideoTooBig = errors.New("视频文件过大")
)

// ClientError carries a user-facing message (e.g. from storage/uploader).
type ClientError struct{ Message string }

func (e *ClientError) Error() string { return e.Message }

func clientErr(msg string) error { return &ClientError{Message: msg} }

type Service struct {
	db       *gorm.DB
	uploader *bunny.Uploader
	s3       *bunny.S3
	cfg      bunny.Config
}

func New(db *gorm.DB) *Service {
	cfg := bunny.LoadConfig()
	return &Service{
		db:       db,
		uploader: bunny.NewUploader(cfg),
		s3:       bunny.NewS3(cfg),
		cfg:      cfg,
	}
}

type mediaRow struct {
	ID              int64      `gorm:"column:id;primaryKey"`
	Title           string     `gorm:"column:title"`
	Thumb           string     `gorm:"column:thumb"`
	Duration        int        `gorm:"column:duration"`
	Size            int64      `gorm:"column:size"`
	Type            string     `gorm:"column:type"`
	CateID          int        `gorm:"column:cate_id"`
	URL             string     `gorm:"column:url"`
	Vid             string     `gorm:"column:vid"`
	DownloadURL     string     `gorm:"column:download_url"`
	TranscodeStatus int        `gorm:"column:transcode_status"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
}

func (mediaRow) TableName() string { return "media_file" }

type ListItem struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Thumb     string  `json:"thumb"`
	Duration  int     `json:"duration"`
	Size      int64   `json:"size"`
	Type      string  `json:"type"`
	CateID    int     `json:"cate_id"`
	URL       string  `json:"url"`
	CreatedAt *string `json:"created_at"`
}

type ListQuery struct {
	Type   string
	CateID *int
	Title  string
	Page   int
	Size   int
}

type UploadResult struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	Duration        int    `json:"duration"`
	Size            int64  `json:"size"`
	TranscodeStatus int    `json:"transcode_status"`
	CateID          int    `json:"cate_id"`
	ID              int64  `json:"id"`
	Vid             string `json:"vid,omitempty"`
	Thumb           string `json:"thumb,omitempty"`
	DownloadURL     string `json:"download_url,omitempty"`
}

type SaveUploadedInput struct {
	URL      string
	Title    string
	Type     string
	CateID   int
	Thumb    string
	Duration int
	Size     int64
	Vid      string
}

type SaveUploadedResult struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	CateID      int    `json:"cate_id"`
	Thumb       string `json:"thumb"`
	Duration    int    `json:"duration"`
	Size        int64  `json:"size"`
	Vid         string `json:"vid"`
	DownloadURL string `json:"download_url"`
	ID          int64  `json:"id"`
}

type InitVideoInput struct {
	FileName string
	MimeType string
	FileSize int64
	Duration int
	Thumb    *multipart.FileHeader
}

type InitVideoResult struct {
	Mode        string                 `json:"mode"`
	UploadID    string                 `json:"uploadId"`
	ObjectKey   string                 `json:"objectKey"`
	Title       string                 `json:"title"`
	ContentType string                 `json:"contentType"`
	VideoID     string                 `json:"videoId"`
	PartSize    int                    `json:"partSize"`
	PartCount   int                    `json:"partCount"`
	Concurrency int                    `json:"concurrency"`
	Parts       []bunny.MultipartPartURL `json:"parts"`
	ThumbURL    string                 `json:"thumbUrl"`
	Duration    int                    `json:"duration"`
}

func (s *Service) List(ctx context.Context, q ListQuery) (page.Body, error) {
	pageSize := q.Size
	if pageSize <= 0 {
		pageSize = 12
	}
	current, size := page.Clamp(q.Page, pageSize)

	query := s.db.WithContext(ctx).Model(&mediaRow{})
	if t := strings.TrimSpace(q.Type); t != "" {
		query = query.Where("type = ?", t)
	}
	if q.CateID != nil && *q.CateID > 0 {
		query = query.Where("cate_id = ?", *q.CateID)
	}
	if title := strings.TrimSpace(q.Title); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []mediaRow
	if err := query.Select("id, title, thumb, duration, size, type, cate_id, url, created_at").
		Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list := make([]ListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, ListItem{
			ID: r.ID, Title: r.Title, Thumb: r.Thumb, Duration: r.Duration,
			Size: r.Size, Type: r.Type, CateID: r.CateID, URL: r.URL,
			CreatedAt: model.FormatTime(r.CreatedAt),
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) Upload(ctx context.Context, fh *multipart.FileHeader, cateID, maxEdge int) (UploadResult, error) {
	if fh == nil {
		return UploadResult{}, ErrNoFile
	}
	result, err := s.uploader.HandleFormFile(fh, "", maxEdge)
	if err != nil {
		return UploadResult{}, clientErr(err.Error())
	}
	now := time.Now()
	row := mediaRow{
		Title: result.Title, Thumb: result.Thumb, Duration: result.Duration, Size: result.Size,
		Type: result.Type, CateID: cateID, URL: result.URL, Vid: result.Vid,
		DownloadURL: result.DownloadURL, TranscodeStatus: result.TranscodeStatus,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if row.DownloadURL == "" {
		row.DownloadURL = row.URL
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return UploadResult{}, err
	}
	out := UploadResult{
		Type: result.Type, Title: result.Title, URL: result.URL,
		Duration: result.Duration, Size: result.Size, TranscodeStatus: result.TranscodeStatus,
		CateID: cateID, ID: row.ID,
	}
	if result.Vid != "" {
		out.Vid = result.Vid
		out.Thumb = result.Thumb
		out.DownloadURL = result.DownloadURL
	}
	return out, nil
}

func (s *Service) Destroy(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrBadParams
	}
	var row mediaRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	for _, key := range []string{row.URL, row.Thumb} {
		key = strings.TrimLeft(strings.TrimSpace(key), "/")
		if key == "" {
			continue
		}
		_ = s.uploader.Storage().Delete(key)
	}
	return s.db.WithContext(ctx).Delete(&row).Error
}

func (s *Service) SaveUploaded(ctx context.Context, in SaveUploadedInput) (SaveUploadedResult, error) {
	now := time.Now()
	row := mediaRow{
		Title: in.Title, Thumb: in.Thumb, Duration: in.Duration, Size: in.Size,
		Type: in.Type, CateID: in.CateID, URL: in.URL, Vid: in.Vid,
		DownloadURL: in.URL, CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return SaveUploadedResult{}, err
	}
	return SaveUploadedResult{
		URL: in.URL, Title: in.Title, Type: in.Type, CateID: in.CateID,
		Thumb: in.Thumb, Duration: in.Duration, Size: in.Size, Vid: in.Vid,
		DownloadURL: in.URL, ID: row.ID,
	}, nil
}

func (s *Service) InitVideoDirectUpload(in InitVideoInput) (InitVideoResult, error) {
	if in.FileName == "" || in.FileSize <= 0 {
		return InitVideoResult{}, ErrNeedFileMeta
	}
	if in.FileSize > 4*1024*1024*1024 {
		return InitVideoResult{}, ErrVideoTooBig
	}
	videoID := uniqID()
	thumbURL := ""
	if in.Thumb != nil {
		key := "video/" + videoID + "/thumb.jpeg"
		result, err := s.uploader.HandleFormFile(in.Thumb, key, 0)
		if err == nil && result != nil {
			thumbURL = result.URL
		}
	}
	payload, err := s.s3.BuildBrowserDirectVideoUpload(in.FileName, in.MimeType, videoID, in.FileSize)
	if err != nil {
		return InitVideoResult{}, clientErr(err.Error())
	}
	return InitVideoResult{
		Mode: payload.Mode, UploadID: payload.UploadID, ObjectKey: payload.ObjectKey,
		Title: payload.Title, ContentType: payload.ContentType, VideoID: payload.VideoID,
		PartSize: payload.PartSize, PartCount: payload.PartCount, Concurrency: payload.Concurrency,
		Parts: payload.Parts, ThumbURL: thumbURL, Duration: in.Duration,
	}, nil
}

func (s *Service) CompleteVideoDirectUpload(objectKey, uploadID string) error {
	if err := s.s3.CompleteBrowserDirectVideoUpload(objectKey, uploadID); err != nil {
		return clientErr(err.Error())
	}
	return nil
}

func (s *Service) AbortVideoDirectUpload(objectKey, uploadID string) {
	s.s3.AbortBrowserDirectVideoUpload(objectKey, uploadID)
}

func uniqID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}
