package bunny

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	multipartPartSize    = 8 * 1024 * 1024
	multipartConcurrency = 6
)

type S3 struct {
	cfg    Config
	client *http.Client
}

func NewS3(cfg Config) *S3 {
	return &S3{
		cfg: cfg,
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

type MultipartPartURL struct {
	PartNumber int    `json:"partNumber"`
	PutURL     string `json:"putUrl"`
}

type MultipartDirectUpload struct {
	Mode        string             `json:"mode"`
	UploadID    string             `json:"uploadId"`
	ObjectKey   string             `json:"objectKey"`
	Title       string             `json:"title"`
	ContentType string             `json:"contentType"`
	VideoID     string             `json:"videoId"`
	PartSize    int                `json:"partSize"`
	PartCount   int                `json:"partCount"`
	Concurrency int                `json:"concurrency"`
	Parts       []MultipartPartURL `json:"parts"`
}

type s3Part struct {
	ETag       string
	PartNumber int
}

func (s *S3) BuildBrowserDirectVideoUpload(fileName, mimeType, videoID string, fileSize int64) (*MultipartDirectUpload, error) {
	videoID = sanitizeID(videoID)
	if videoID == "" {
		return nil, fmt.Errorf("视频目录 ID 无效")
	}
	mimeType = strings.TrimSpace(mimeType)
	if mimeType == "" {
		mimeType = "video/mp4"
	}
	if !isVideoMIME(mimeType) {
		return nil, fmt.Errorf("不支持的视频格式")
	}
	if !s.cfg.S3Ready() {
		return nil, fmt.Errorf("服务端未配置 Bunny 凭证")
	}
	safeName := filepathBase(fileName)
	if safeName == "" || safeName == "." || safeName == ".." {
		return nil, fmt.Errorf("文件名无效")
	}
	ext := strings.ToLower(strings.TrimPrefix(filepathExt(safeName), "."))
	if ext == "" {
		ext = "mp4"
	}
	key := fmt.Sprintf("video/%s/default.%s", videoID, ext)
	expires := 2 * time.Hour
	if fileSize > 1<<30 {
		expires = 6 * time.Hour
	}
	return s.buildMultipart(key, mimeType, safeName, videoID, fileSize, expires)
}

func (s *S3) buildMultipart(key, mimeType, safeName, videoID string, fileSize int64, expires time.Duration) (*MultipartDirectUpload, error) {
	partCount := int((fileSize + multipartPartSize - 1) / multipartPartSize)
	if partCount < 1 {
		partCount = 1
	}
	if partCount > 10000 {
		return nil, fmt.Errorf("分片数量超限")
	}
	uploadID, err := s.createMultipartUpload(key, mimeType)
	if err != nil {
		return nil, err
	}
	ttl := int(expires.Seconds())
	if ttl < 60 {
		ttl = 60
	}
	parts := make([]MultipartPartURL, 0, partCount)
	for n := 1; n <= partCount; n++ {
		parts = append(parts, MultipartPartURL{
			PartNumber: n,
			PutURL:     s.presignUploadPart(key, uploadID, n, ttl),
		})
	}
	return &MultipartDirectUpload{
		Mode:        "multipart",
		UploadID:    uploadID,
		ObjectKey:   key,
		Title:       safeName,
		ContentType: mimeType,
		VideoID:     videoID,
		PartSize:    multipartPartSize,
		PartCount:   partCount,
		Concurrency: multipartConcurrency,
		Parts:       parts,
	}, nil
}

func (s *S3) CompleteBrowserDirectVideoUpload(objectKey, uploadID string) error {
	if !validVideoObjectKey(objectKey) {
		return fmt.Errorf("对象路径无效")
	}
	uploadID = strings.TrimSpace(uploadID)
	if uploadID == "" || len(uploadID) > 2048 {
		return fmt.Errorf("uploadId 无效")
	}
	parts, err := s.listAllParts(objectKey, uploadID)
	if err != nil {
		return fmt.Errorf("合并分片失败: %v", err)
	}
	if len(parts) == 0 {
		return fmt.Errorf("未找到已上传分片")
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })
	if err := s.completeMultipartUpload(objectKey, uploadID, parts); err != nil {
		return fmt.Errorf("合并分片失败: %v", err)
	}
	return nil
}

func (s *S3) AbortBrowserDirectVideoUpload(objectKey, uploadID string) {
	if !validVideoObjectKey(objectKey) {
		return
	}
	uploadID = strings.TrimSpace(uploadID)
	if uploadID == "" {
		return
	}
	_ = s.abortMultipartUpload(objectKey, uploadID)
}

func (s *S3) createMultipartUpload(key, contentType string) (string, error) {
	raw, err := s.request("POST", key, map[string]string{"uploads": ""}, "", map[string]string{
		"content-type": contentType,
	})
	if err != nil {
		return "", err
	}
	var parsed struct {
		UploadID string `xml:"UploadId"`
	}
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("未能创建分片上传会话")
	}
	id := strings.TrimSpace(parsed.UploadID)
	if id == "" {
		return "", fmt.Errorf("未能创建分片上传会话")
	}
	return id, nil
}

func (s *S3) listAllParts(key, uploadID string) ([]s3Part, error) {
	var parts []s3Part
	var marker *int
	for {
		query := map[string]string{"uploadId": uploadID}
		if marker != nil {
			query["part-number-marker"] = strconv.Itoa(*marker)
		}
		raw, err := s.request("GET", key, query, "", nil)
		if err != nil {
			return nil, err
		}
		var parsed struct {
			IsTruncated          string `xml:"IsTruncated"`
			NextPartNumberMarker int    `xml:"NextPartNumberMarker"`
			Parts                []struct {
				PartNumber int    `xml:"PartNumber"`
				ETag       string `xml:"ETag"`
			} `xml:"Part"`
		}
		if err := xml.Unmarshal(raw, &parsed); err != nil {
			return nil, fmt.Errorf("Bunny S3 返回不是 XML")
		}
		for _, p := range parsed.Parts {
			parts = append(parts, s3Part{ETag: p.ETag, PartNumber: p.PartNumber})
		}
		truncated := strings.EqualFold(parsed.IsTruncated, "true")
		if !truncated {
			break
		}
		next := parsed.NextPartNumberMarker
		if next == 0 && len(parts) > 0 {
			next = parts[len(parts)-1].PartNumber
		}
		marker = &next
	}
	return parts, nil
}

func (s *S3) completeMultipartUpload(key, uploadID string, parts []s3Part) error {
	var b strings.Builder
	b.WriteString("<CompleteMultipartUpload>")
	for _, p := range parts {
		etag := xmlEscape(p.ETag)
		b.WriteString("<Part><PartNumber>")
		b.WriteString(strconv.Itoa(p.PartNumber))
		b.WriteString("</PartNumber><ETag>")
		b.WriteString(etag)
		b.WriteString("</ETag></Part>")
	}
	b.WriteString("</CompleteMultipartUpload>")
	_, err := s.request("POST", key, map[string]string{"uploadId": uploadID}, b.String(), map[string]string{
		"content-type": "application/xml",
	})
	return err
}

func (s *S3) abortMultipartUpload(key, uploadID string) error {
	_, err := s.request("DELETE", key, map[string]string{"uploadId": uploadID}, "", nil)
	return err
}

func (s *S3) presignUploadPart(key, uploadID string, partNumber, expiresSeconds int) string {
	if expiresSeconds < 60 {
		expiresSeconds = 60
	}
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	date := amzDate[:8]
	scope := date + "/" + s.cfg.Region + "/s3/aws4_request"
	query := map[string]string{
		"partNumber":          strconv.Itoa(partNumber),
		"uploadId":            uploadID,
		"X-Amz-Algorithm":     "AWS4-HMAC-SHA256",
		"X-Amz-Credential":    s.cfg.s3Bucket() + "/" + scope,
		"X-Amz-Date":          amzDate,
		"X-Amz-Expires":       strconv.Itoa(expiresSeconds),
		"X-Amz-SignedHeaders": "host",
	}
	cq := encodeCanonicalQuery(query)
	uri := s.canonicalURI(key)
	host := s.host()
	canonical := strings.Join([]string{
		"PUT", uri, cq, "host:" + host + "\n", "host", "UNSIGNED-PAYLOAD",
	}, "\n")
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256", amzDate, scope, sha256Hex([]byte(canonical)),
	}, "\n")
	query["X-Amz-Signature"] = hmacSHA256Hex(s.signingKey(date), stringToSign)
	return s.cfg.s3Endpoint() + uri + "?" + encodeCanonicalQuery(query)
}

func (s *S3) request(method, key string, query map[string]string, body string, headers map[string]string) ([]byte, error) {
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	date := amzDate[:8]
	payloadHash := sha256Hex([]byte(body))
	host := s.host()
	signed := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	for name, value := range headers {
		signed[strings.ToLower(name)] = strings.TrimSpace(value)
	}
	names := make([]string, 0, len(signed))
	for name := range signed {
		names = append(names, name)
	}
	sort.Strings(names)
	var canonicalHeaders strings.Builder
	for _, name := range names {
		canonicalHeaders.WriteString(name)
		canonicalHeaders.WriteByte(':')
		canonicalHeaders.WriteString(signed[name])
		canonicalHeaders.WriteByte('\n')
	}
	signedHeaderNames := strings.Join(names, ";")
	uri := s.canonicalURI(key)
	canonical := strings.Join([]string{
		method, uri, encodeCanonicalQuery(query), canonicalHeaders.String(), signedHeaderNames, payloadHash,
	}, "\n")
	scope := date + "/" + s.cfg.Region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256", amzDate, scope, sha256Hex([]byte(canonical)),
	}, "\n")
	signature := hmacSHA256Hex(s.signingKey(date), stringToSign)

	httpHeaders := http.Header{}
	for _, name := range names {
		if name == "host" {
			continue
		}
		httpHeaders.Set(name, signed[name])
	}
	httpHeaders.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+s.cfg.s3Bucket()+"/"+scope+
		", SignedHeaders="+signedHeaderNames+", Signature="+signature)

	u := s.cfg.s3Endpoint() + uri
	if len(query) > 0 {
		u += "?" + encodeCanonicalQuery(query)
	}
	req, err := http.NewRequest(method, u, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = httpHeaders
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Bunny S3 %s %d %s", method, resp.StatusCode, s3ErrorText(raw))
	}
	return raw, nil
}

func (s *S3) signingKey(date string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+s.cfg.s3Secret()), date)
	kRegion := hmacSHA256(kDate, s.cfg.Region)
	kService := hmacSHA256(kRegion, "s3")
	return hmacSHA256(kService, "aws4_request")
}

func (s *S3) canonicalURI(key string) string {
	segments := strings.Split(strings.TrimLeft(key, "/"), "/")
	encoded := make([]string, len(segments))
	for i, seg := range segments {
		encoded[i] = rawPathEscape(seg)
	}
	return "/" + rawPathEscape(s.cfg.s3Bucket()) + "/" + strings.Join(encoded, "/")
}

func (s *S3) host() string {
	u, err := url.Parse(s.cfg.s3Endpoint())
	if err != nil || u.Host == "" {
		return strings.TrimPrefix(strings.TrimPrefix(s.cfg.s3Endpoint(), "https://"), "http://")
	}
	return u.Host
}

func encodeCanonicalQuery(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(query[k]))
	}
	return strings.Join(pairs, "&")
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(data))
	return m.Sum(nil)
}

func hmacSHA256Hex(key []byte, data string) string {
	return hex.EncodeToString(hmacSHA256(key, data))
}

func rawPathEscape(s string) string {
	return strings.ReplaceAll(url.PathEscape(s), "+", "%20")
}

func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

func s3ErrorText(raw []byte) string {
	var parsed struct {
		Message string `xml:"Message"`
	}
	if xml.Unmarshal(raw, &parsed) == nil && strings.TrimSpace(parsed.Message) != "" {
		return strings.TrimSpace(parsed.Message)
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return "存储服务错误"
	}
	if len(text) > 300 {
		text = text[:300]
	}
	return text
}

func validVideoObjectKey(key string) bool {
	key = objectKey(key)
	if key == "" {
		return false
	}
	parts := strings.Split(key, "/")
	if len(parts) != 3 || parts[0] != "video" {
		return false
	}
	if sanitizeID(parts[1]) != parts[1] {
		return false
	}
	name := parts[2]
	if !strings.HasPrefix(strings.ToLower(name), "default.") {
		return false
	}
	ext := strings.ToLower(filepathExt(name))
	return ext != "" && ext != "."
}

func sanitizeID(id string) string {
	var b strings.Builder
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		}
	}
	return b.String()
}

func filepathBase(name string) string {
	name = strings.ReplaceAll(name, "\x00", "")
	name = strings.ReplaceAll(name, "\\", "/")
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func filepathExt(name string) string {
	i := strings.LastIndex(name, ".")
	if i < 0 {
		return ""
	}
	return name[i:]
}

func isVideoMIME(mime string) bool {
	switch mime {
	case "video/mp4", "video/quicktime", "video/x-msvideo", "video/x-flv", "video/webm", "video/mpeg":
		return true
	default:
		return false
	}
}
