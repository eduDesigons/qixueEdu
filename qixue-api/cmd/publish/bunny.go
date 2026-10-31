package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type storageItem struct {
	ObjectName  string `json:"ObjectName"`
	Path        string `json:"Path"`
	IsDirectory bool   `json:"IsDirectory"`
}

func storageBaseURL(region string) string {
	region = strings.ToLower(strings.TrimSpace(region))
	if region == "" || region == "de" {
		return "https://storage.bunnycdn.com"
	}
	return "https://" + region + ".storage.bunnycdn.com"
}

func newStorageHTTP() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Minute,
		Transport: &http.Transport{
			MaxIdleConns:        32,
			MaxIdleConnsPerHost: 32,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

func storageRequest(client *http.Client, method, password, region, zone, path string, body io.Reader, contentType string) (*http.Response, error) {
	path = strings.TrimLeft(strings.ReplaceAll(path, "\\", "/"), "/")
	url := strings.TrimRight(storageBaseURL(region), "/") + "/" + strings.TrimRight(zone, "/") + "/"
	if path != "" {
		url += path
	}
	if strings.HasSuffix(path, "/") && !strings.HasSuffix(url, "/") {
		url += "/"
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("AccessKey", password)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return client.Do(req)
}

func storageDelete(client *http.Client, password, region, zone, path string) (int, error) {
	resp, err := storageRequest(client, http.MethodDelete, password, region, zone, path, nil, "")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func storageObjectPath(zone string, item storageItem) string {
	name := strings.Trim(item.ObjectName, "/")
	if name == "" {
		return ""
	}
	dir := strings.Trim(strings.ReplaceAll(item.Path, "\\", "/"), "/")
	parts := []string{}
	if dir != "" {
		parts = strings.Split(dir, "/")
	}
	if len(parts) > 0 && strings.EqualFold(parts[0], zone) {
		parts = parts[1:]
	}
	dirRel := strings.Join(parts, "/")
	base := name
	if dirRel != "" {
		base = dirRel + "/" + name
	}
	if item.IsDirectory {
		return base + "/"
	}
	return base
}

func listStorageLevel(client *http.Client, password, region, zone, prefix string) ([]storageItem, error) {
	prefix = strings.Trim(strings.ReplaceAll(prefix, "\\", "/"), "/")
	listPath := ""
	if prefix != "" {
		listPath = prefix + "/"
	}
	resp, err := storageRequest(client, http.MethodGet, password, region, zone, listPath, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出 %s/%s 失败 HTTP %d", zone, listPath, resp.StatusCode)
	}
	var items []storageItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("列出 %s/%s 返回非 JSON", zone, listPath)
	}
	return items, nil
}

func countStorageLevel(client *http.Client, password, region, zone, prefix string) (int, error) {
	items, err := listStorageLevel(client, password, region, zone, prefix)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, item := range items {
		if storageObjectPath(zone, item) != "" {
			n++
		}
	}
	return n, nil
}

func deleteStorageLevel(client *http.Client, password, region, zone, prefix string) (int, error) {
	prefix = strings.Trim(strings.ReplaceAll(prefix, "\\", "/"), "/")
	items, err := listStorageLevel(client, password, region, zone, prefix)
	if err != nil {
		return 0, err
	}
	var targets []string
	for _, item := range items {
		remote := storageObjectPath(zone, item)
		if remote == "" {
			continue
		}
		if prefix != "" && remote != prefix && !strings.HasPrefix(remote, prefix+"/") {
			continue
		}
		if item.IsDirectory {
			remote = strings.TrimRight(remote, "/") + "/"
		}
		targets = append(targets, remote)
	}
	sort.Slice(targets, func(i, j int) bool {
		return len(targets[i]) > len(targets[j])
	})

	count := 0
	for _, target := range targets {
		fmt.Printf("    DELETE %s\n", target)
		code, err := storageDelete(client, password, region, zone, target)
		if err != nil {
			return count, err
		}
		if code == http.StatusOK {
			count++
		} else if code != http.StatusNotFound {
			fmt.Printf("    删除失败 %s HTTP %d\n", target, code)
		}
	}
	return count, nil
}

// clearStoragePrefix deletes the zone (or prefix) and waits until listing is empty.
func clearStoragePrefix(password, zone, region, prefix string) (int, error) {
	prefix = strings.Trim(strings.ReplaceAll(prefix, "\\", "/"), "/")
	client := newStorageHTTP()
	deleted := 0
	deadline := time.Now().Add(180 * time.Second)
	round := 0

	if prefix != "" {
		fmt.Printf("    DELETE %s/\n", prefix)
		code, err := storageDelete(client, password, region, zone, prefix+"/")
		if err != nil {
			return 0, err
		}
		if code == http.StatusOK {
			deleted++
		} else if code != http.StatusNotFound {
			fmt.Printf("    目录删除 HTTP %d，改为删顶层子项\n", code)
			n, err := deleteStorageLevel(client, password, region, zone, prefix)
			if err != nil {
				return deleted, err
			}
			deleted += n
		}
	} else {
		n, err := deleteStorageLevel(client, password, region, zone, "")
		if err != nil {
			return 0, err
		}
		deleted += n
	}

	for {
		left, err := countStorageLevel(client, password, region, zone, prefix)
		if err != nil {
			return deleted, err
		}
		if left == 0 {
			if round > 0 {
				fmt.Println("    远端已清空")
			}
			return deleted, nil
		}
		if time.Now().After(deadline) {
			return deleted, fmt.Errorf("清空远端超时，仍有 %d 项未删完，已中止上传", left)
		}
		round++
		fmt.Printf("    仍有 %d 项，等待删除完成\n", left)
		time.Sleep(2 * time.Second)
		if round%5 == 0 {
			n, err := deleteStorageLevel(client, password, region, zone, prefix)
			if err != nil {
				return deleted, err
			}
			deleted += n
		}
	}
}

func uploadFilesConcurrent(password, zone, region string, files map[string]string, remoteBase string, concurrency int) error {
	if len(files) == 0 {
		return nil
	}
	if concurrency < 1 {
		concurrency = 1
	}
	client := newStorageHTTP()
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var done atomic.Int64
	total := int64(len(files))
	errCh := make(chan string, len(files))

	for _, rel := range sortedRels(files) {
		abs := files[rel]
		remote := rel
		if remoteBase != "" {
			remote = remoteBase + "/" + rel
		}
		remote = strings.ReplaceAll(remote, "\\", "/")

		wg.Add(1)
		sem <- struct{}{}
		go func(remote, abs string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := uploadOne(client, password, region, zone, remote, abs); err != nil {
				n := done.Add(1)
				msg := err.Error()
				fmt.Printf("    [%d/%d] FAIL %s: %s\n", n, total, remote, msg)
				errCh <- remote
				return
			}
			n := done.Add(1)
			fmt.Printf("    [%d/%d] ok %s\n", n, total, remote)
		}(remote, abs)
	}
	wg.Wait()
	close(errCh)
	var failed []string
	for r := range errCh {
		failed = append(failed, r)
	}
	if len(failed) > 0 {
		return fmt.Errorf("上传失败 %d 个文件，例如: %s", len(failed), failed[0])
	}
	return nil
}

func uploadOne(client *http.Client, password, region, zone, remote, abs string) error {
	f, err := os.Open(abs)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := storageRequest(client, http.MethodPut, password, region, zone, remote, f, "application/octet-stream")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

func purgePullZoneCaches(apiKey string, pullZoneIDs []int, label string) error {
	if apiKey == "" {
		fmt.Printf("==> [%s] 跳过 CDN 缓存清理（未配置 BUNNY_ACCOUNT_API_KEY）\n", label)
		return nil
	}
	if len(pullZoneIDs) == 0 {
		fmt.Printf("==> [%s] skip CDN purge (no pull zone id)\n", label)
		return nil
	}

	client := &http.Client{Timeout: 60 * time.Second}
	fmt.Printf("==> [%s] purge CDN cache (%d zones)\n", label, len(pullZoneIDs))
	for _, id := range pullZoneIDs {
		url := fmt.Sprintf("https://api.bunny.net/pullzone/%d/purgeCache", id)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("{}")))
		if err != nil {
			return err
		}
		req.Header.Set("AccessKey", apiKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("[%s] CDN purge zone #%d failed: %v", label, id, err)
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		code := resp.StatusCode
		if code != http.StatusNoContent && code != http.StatusOK {
			msg := strings.TrimSpace(string(body))
			if msg != "" {
				return fmt.Errorf("[%s] CDN purge zone #%d HTTP %d: %s", label, id, code, msg)
			}
			return fmt.Errorf("[%s] CDN purge zone #%d HTTP %d", label, id, code)
		}
		fmt.Printf("    #%d purge 已触发（HTTP %d）\n", id, code)
	}
	return nil
}
