package adminmedia

import (
	"strconv"
	"strings"
	"time"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/service/bunny"
)

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

func (s *Service) List(c *gin.Context) {
	var body struct {
		Type   string   `json:"type"`
		CateID any      `json:"cateId"`
		Title  string   `json:"title"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	pageSize := body.Size.Int()
	if pageSize <= 0 {
		pageSize = 12
	}
	current, size := page.Clamp(body.Page.Int(), pageSize)
	ctx := c.Request.Context()

	query := s.db.WithContext(ctx).Model(&mediaRow{})
	if t := strings.TrimSpace(body.Type); t != "" {
		query = query.Where("type = ?", t)
	}
	if cateID, ok := asOptionalInt(body.CateID); ok && cateID > 0 {
		query = query.Where("cate_id = ?", cateID)
	}
	if title := strings.TrimSpace(body.Title); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []mediaRow
	if response.WriteErr(c, query.Select("id, title, thumb, duration, size, type, cate_id, url, created_at").
		Order("id DESC").Offset((current-1)*size).Limit(size).Find(&rows).Error) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		list = append(list, gin.H{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "duration": r.Duration,
			"size": r.Size, "type": r.Type, "cate_id": r.CateID, "url": r.URL,
			"created_at": model.FormatTime(r.CreatedAt),
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		response.Fail(c, 1, "文件不存在或无效")
		return
	}
	cateID := bind.Query(c, "cateId")
	maxEdge := bind.Query(c, "maxEdge")
	result, err := s.uploader.HandleFormFile(fh, "", maxEdge)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
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
	if response.WriteErr(c, s.db.WithContext(c.Request.Context()).Create(&row).Error) {
		return
	}
	out := gin.H{
		"type": result.Type, "title": result.Title, "url": result.URL,
		"duration": result.Duration, "size": result.Size, "transcode_status": result.TranscodeStatus,
		"cate_id": cateID, "id": row.ID,
	}
	if result.Vid != "" {
		out["vid"] = result.Vid
		out["thumb"] = result.Thumb
		out["download_url"] = result.DownloadURL
	}
	response.OK(c, "ok", out)
}

func (s *Service) Destroy(c *gin.Context) {
	var body struct {
		ID any `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID == nil {
		body.ID = c.Query("id")
	}
	id := bind.AsInt64(body.ID)
	if id <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	var row mediaRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "资源不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	for _, key := range []string{row.URL, row.Thumb} {
		key = strings.TrimLeft(strings.TrimSpace(key), "/")
		if key == "" {
			continue
		}
		_ = s.uploader.Storage().Delete(key)
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Delete(&row).Error) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func (s *Service) SaveUploaded(c *gin.Context) {
	var body struct {
		URL      string   `json:"url"`
		Title    string   `json:"title"`
		Type     string   `json:"type"`
		CateID   bind.Int `json:"cate_id"`
		Thumb    string   `json:"thumb"`
		Duration int      `json:"duration"`
		Size     int64    `json:"size"`
		Vid      string   `json:"vid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "落库失败: 参数错误")
		return
	}
	now := time.Now()
	row := mediaRow{
		Title: body.Title, Thumb: body.Thumb, Duration: body.Duration, Size: body.Size,
		Type: body.Type, CateID: body.CateID.Int(), URL: body.URL, Vid: body.Vid,
		DownloadURL: body.URL, CreatedAt: &now, UpdatedAt: &now,
	}
	if response.WriteErr(c, s.db.WithContext(c.Request.Context()).Create(&row).Error) {
		return
	}
	response.OK(c, "ok", gin.H{
		"url": body.URL, "title": body.Title, "type": body.Type, "cate_id": body.CateID.Int(),
		"thumb": body.Thumb, "duration": body.Duration, "size": body.Size, "vid": body.Vid,
		"download_url": body.URL, "id": row.ID,
	})
}

func (s *Service) InitVideoDirectUpload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	mimeType := c.PostForm("mimeType")
	fileSize, _ := strconv.ParseInt(c.PostForm("fileSize"), 10, 64)
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	if fileName == "" || fileSize <= 0 {
		response.Fail(c, 1, "请提供有效的文件名与大小")
		return
	}
	if fileSize > 4*1024*1024*1024 {
		response.Fail(c, 1, "视频文件过大")
		return
	}
	videoID := uniqID()
	thumbURL := ""
	if fh, err := c.FormFile("thumb"); err == nil && fh != nil {
		key := "video/" + videoID + "/thumb.jpeg"
		result, err := s.uploader.HandleFormFile(fh, key, 0)
		if err == nil && result != nil {
			thumbURL = result.URL
		}
	}
	payload, err := s.s3.BuildBrowserDirectVideoUpload(fileName, mimeType, videoID, fileSize)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}
	response.OK(c, "ok", gin.H{
		"mode": payload.Mode, "uploadId": payload.UploadID, "objectKey": payload.ObjectKey,
		"title": payload.Title, "contentType": payload.ContentType, "videoId": payload.VideoID,
		"partSize": payload.PartSize, "partCount": payload.PartCount, "concurrency": payload.Concurrency,
		"parts": payload.Parts, "thumbUrl": thumbURL, "duration": duration,
	})
}

func (s *Service) CompleteVideoDirectUpload(c *gin.Context) {
	var body struct {
		ObjectKey string `json:"objectKey"`
		UploadID  string `json:"uploadId"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ObjectKey == "" {
		body.ObjectKey = c.PostForm("objectKey")
		body.UploadID = c.PostForm("uploadId")
	}
	if err := s.s3.CompleteBrowserDirectVideoUpload(body.ObjectKey, body.UploadID); err != nil {
		response.Fail(c, 1, err.Error())
		return
	}
	response.OK(c, "ok", nil)
}

func (s *Service) AbortVideoDirectUpload(c *gin.Context) {
	var body struct {
		ObjectKey string `json:"objectKey"`
		UploadID  string `json:"uploadId"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ObjectKey == "" {
		body.ObjectKey = c.PostForm("objectKey")
		body.UploadID = c.PostForm("uploadId")
	}
	s.s3.AbortBrowserDirectVideoUpload(body.ObjectKey, body.UploadID)
	response.OK(c, "ok", nil)
}

func asOptionalInt(v any) (int, bool) {
	if !bind.Present(v) {
		return 0, false
	}
	return bind.AsInt(v), true
}

func uniqID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}
