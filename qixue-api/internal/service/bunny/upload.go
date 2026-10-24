package bunny

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

var dangerousMIME = map[string]struct{}{
	"application/x-php": {}, "application/x-sh": {}, "application/x-python": {},
	"application/x-perl": {}, "application/x-ruby": {}, "application/x-java": {},
	"application/x-c": {}, "application/x-c++": {}, "application/x-java-source": {},
}

type UploadResult struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	DownloadURL     string `json:"download_url,omitempty"`
	Thumb           string `json:"thumb,omitempty"`
	Duration        int    `json:"duration"`
	Size            int64  `json:"size"`
	Vid             string `json:"vid,omitempty"`
	TranscodeStatus int    `json:"transcode_status"`
	CateID          int    `json:"cate_id,omitempty"`
	ID              int64  `json:"id,omitempty"`
}

type Uploader struct {
	storage *Storage
}

func NewUploader(cfg Config) *Uploader {
	return &Uploader{storage: NewStorage(cfg)}
}

func (u *Uploader) Storage() *Storage { return u.storage }

func (u *Uploader) HandleFormFile(fileHeader *multipart.FileHeader, forcedKey string, maxEdge int) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, fmt.Errorf("文件不存在或无效")
	}
	if maxEdge < 320 || maxEdge > 1920 {
		maxEdge = 0
	}
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("文件不存在或无效")
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "qx-upload-*")
	if err != nil {
		return nil, err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return nil, fmt.Errorf("读取上传文件失败")
	}
	tmp.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if _, bad := dangerousMIME[mimeType]; bad {
		return nil, fmt.Errorf("文件类型不允许上传")
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileHeader.Filename), "."))
	fileType := detectFileType(mimeType, ext)

	if fileType == "video" {
		return u.handleVideo(tmpPath, fileHeader.Filename, mimeType, fileHeader.Size)
	}
	return u.handleNonVideo(tmpPath, fileHeader.Filename, mimeType, fileType, forcedKey, maxEdge, fileHeader.Size)
}

func (u *Uploader) UploadLocal(path, key, contentType string) (string, error) {
	return u.storage.UploadFile(path, key, contentType)
}

func (u *Uploader) handleVideo(path, title, mimeType string, size int64) (*UploadResult, error) {
	if !isVideoMIME(mimeType) {
		return nil, fmt.Errorf("不支持的视频格式")
	}
	videoID := uniqID()
	videoKey := fmt.Sprintf("video/%s/default.mp4", videoID)
	objectKey, err := u.storage.UploadFile(path, videoKey, "video/mp4")
	if err != nil {
		return nil, fmt.Errorf("视频处理失败: %v", err)
	}
	return &UploadResult{
		Type:            "video",
		Vid:             videoID,
		Thumb:           "",
		Duration:        0,
		Size:            size,
		Title:           title,
		URL:             objectKey,
		DownloadURL:     objectKey,
		TranscodeStatus: 0,
	}, nil
}

func (u *Uploader) handleNonVideo(path, title, mimeType, fileType, forcedKey string, maxEdge int, originalSize int64) (*UploadResult, error) {
	uploadPath := path
	uploadMIME := mimeType
	cleanup := ""
	defer func() {
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()

	finalKey := forcedKey
	if finalKey == "" {
		finalKey = fileType + "/" + time.Now().Format("20060102") + "/" + uniqID() + filepath.Ext(title)
	}

	if fileType == "image" {
		compressed, newMIME, ok := compressImage(path, mimeType, maxEdge)
		if ok && compressed != path {
			uploadPath = compressed
			cleanup = compressed
			uploadMIME = newMIME
			if newMIME == "image/jpeg" {
				finalKey = replaceExt(finalKey, ".jpg")
			}
		}
	}

	info, err := os.Stat(uploadPath)
	if err != nil {
		return nil, fmt.Errorf("文件不存在")
	}
	objectKey, err := u.storage.UploadFile(uploadPath, finalKey, uploadMIME)
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size <= 0 {
		size = originalSize
	}
	return &UploadResult{
		Duration:        0,
		Size:            size,
		TranscodeStatus: 0,
		URL:             objectKey,
		Type:            fileType,
		Title:           title,
	}, nil
}

func detectFileType(mimeType, extension string) string {
	if extension == "apk" || extension == "ipa" {
		return "app"
	}
	switch mimeType {
	case "video/mp4", "video/quicktime", "video/x-msvideo", "video/x-flv", "video/webm", "video/mpeg":
		return "video"
	case "audio/mpeg", "audio/wav", "audio/ogg":
		return "audio"
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return "image"
	case "application/vnd.android.package-archive":
		return "app"
	default:
		return "document"
	}
}

func compressImage(path, mimeType string, maxEdge int) (string, string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return path, mimeType, false
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return path, mimeType, false
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	stat, _ := os.Stat(path)
	origSize := int64(0)
	if stat != nil {
		origSize = stat.Size()
	}

	limitEdge := maxEdge >= 320
	saveJPEG := limitEdge && (mimeType == "image/png" || mimeType == "image/gif" || mimeType == "image/webp")
	if limitEdge && !saveJPEG && w <= maxEdge && h <= maxEdge {
		return path, mimeType, false
	}
	if !limitEdge && origSize < 500*1024 {
		return path, mimeType, false
	}

	maxW, maxH := 1920, 1920
	if limitEdge {
		maxW, maxH = maxEdge, maxEdge
	}
	nw, nh := w, h
	if maxW > 0 && nw > maxW {
		nh = nh * maxW / nw
		nw = maxW
	}
	if maxH > 0 && nh > maxH {
		nw = nw * maxH / nh
		nh = maxH
	}
	quality := 85
	if origSize >= 2*1024*1024 {
		quality = 70
	} else if origSize >= 1024*1024 {
		quality = 75
	} else if origSize >= 500*1024 {
		quality = 80
	}
	if quality < 60 {
		quality = 60
	}

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	if saveJPEG {
		draw.Draw(dst, dst.Bounds(), image.White, image.Point{}, draw.Src)
	}
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	outMIME := mimeType
	ext := filepath.Ext(path)
	if saveJPEG {
		outMIME = "image/jpeg"
		ext = ".jpg"
	}
	tmp, err := os.CreateTemp("", "qx-img-*"+ext)
	if err != nil {
		return path, mimeType, false
	}
	defer tmp.Close()

	switch outMIME {
	case "image/jpeg", "image/jpg":
		if err := jpeg.Encode(tmp, dst, &jpeg.Options{Quality: quality}); err != nil {
			os.Remove(tmp.Name())
			return path, mimeType, false
		}
	default:
		// Keep a JPEG fallback for unsupported encoders after resize.
		if err := jpeg.Encode(tmp, dst, &jpeg.Options{Quality: quality}); err != nil {
			os.Remove(tmp.Name())
			return path, mimeType, false
		}
		outMIME = "image/jpeg"
	}
	info, err := tmp.Stat()
	if err != nil || info.Size() >= origSize && !limitEdge {
		os.Remove(tmp.Name())
		return path, mimeType, false
	}
	return tmp.Name(), outMIME, true
}

func uniqID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func replaceExt(key, ext string) string {
	lower := strings.ToLower(key)
	for _, old := range []string{".png", ".gif", ".webp"} {
		if strings.HasSuffix(lower, old) {
			return key[:len(key)-len(old)] + ext
		}
	}
	return key
}

// ReadFormFile is a tiny helper used by tests/callers that already have bytes.
func ReadFormFile(fh *multipart.FileHeader) ([]byte, error) {
	if fh == nil {
		return nil, fmt.Errorf("empty")
	}
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
