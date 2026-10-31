package bunny

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Storage struct {
	cfg    Config
	client *http.Client
}

func NewStorage(cfg Config) *Storage {
	return &Storage{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}
}

func (s *Storage) UploadFile(localPath, key, contentType string) (string, error) {
	key = objectKey(key)
	if key == "" {
		return "", fmt.Errorf("对象路径无效")
	}
	if s.cfg.Local {
		return s.saveLocal(localPath, key)
	}
	if !s.cfg.StorageReady() {
		return "", fmt.Errorf("服务端未配置 Bunny 凭证")
	}
	f, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("文件不存在")
	}
	defer f.Close()

	url := s.cfg.storageBaseURL() + "/" + s.cfg.StorageZone + "/" + key
	req, err := http.NewRequest(http.MethodPut, url, f)
	if err != nil {
		return "", err
	}
	req.Header.Set("AccessKey", s.cfg.APIKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("文件上传失败: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("文件上传失败: %s", msg)
	}
	return key, nil
}

func (s *Storage) UploadBytes(data []byte, key, contentType string) (string, error) {
	tmp, err := os.CreateTemp("", "bunny-up-*")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	defer os.Remove(path)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()
	return s.UploadFile(path, key, contentType)
}

// GetContents returns object bytes by key (remote storage or local public/).
func (s *Storage) GetContents(key string) ([]byte, error) {
	key = objectKey(key)
	if key == "" || strings.Contains(key, "..") {
		return nil, fmt.Errorf("对象 key 为空")
	}
	if s.cfg.Local {
		path := s.localPath(key)
		if path == "" {
			return nil, fmt.Errorf("路径无效")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("文件不存在")
			}
			return nil, fmt.Errorf("读取文件失败: %v", err)
		}
		return data, nil
	}
	if !s.cfg.StorageReady() {
		return nil, fmt.Errorf("服务端未配置 Bunny 凭证")
	}
	url := s.cfg.storageBaseURL() + "/" + s.cfg.StorageZone + "/" + key
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("AccessKey", s.cfg.APIKey)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(body))
		if len(msg) > 200 {
			msg = msg[:200]
		}
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("下载文件失败: %s", msg)
	}
	return body, nil
}

func (s *Storage) Delete(key string) error {
	key = objectKey(key)
	if key == "" {
		return nil
	}
	if s.cfg.Local {
		path := s.localPath(key)
		if path == "" {
			return nil
		}
		_ = os.Remove(path)
		return nil
	}
	if !s.cfg.StorageReady() {
		return fmt.Errorf("服务端未配置 Bunny 凭证")
	}
	url := s.cfg.storageBaseURL() + "/" + s.cfg.StorageZone + "/" + key
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("AccessKey", s.cfg.APIKey)
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("删除文件失败: %s", msg)
	}
	return nil
}

func (s *Storage) saveLocal(localPath, key string) (string, error) {
	dest := s.localPath(key)
	if dest == "" {
		return "", fmt.Errorf("路径无效")
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("无法创建上传目录")
	}
	in, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("文件不存在")
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("保存文件失败")
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return "", fmt.Errorf("保存文件失败")
	}
	return key, nil
}

func (s *Storage) localPath(key string) string {
	key = objectKey(key)
	if key == "" || strings.Contains(key, "..") {
		return ""
	}
	root := filepath.Join(".", "public")
	return filepath.Join(append([]string{root}, strings.Split(key, "/")...)...)
}

func objectKey(key string) string {
	return strings.TrimLeft(strings.TrimSpace(key), "/")
}
