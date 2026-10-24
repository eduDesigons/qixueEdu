package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func publishLocal(opt publishOpt) error {
	label := opt.Label
	distDir := opt.DistDir
	destDir := opt.PublicDir
	prefix := strings.Trim(opt.URLPrefix, "/")
	if prefix != "admin" && prefix != "web" {
		return fmt.Errorf("[%s] 只能发布到 qixue-api/public/%s", label, prefix)
	}
	if filepath.Base(filepath.Dir(destDir)) != "public" || filepath.Base(destDir) != prefix {
		return fmt.Errorf("[%s] 只能发布到 qixue-api/public/%s", label, prefix)
	}

	indexBytes, err := os.ReadFile(filepath.Join(distDir, "index.html"))
	if err != nil {
		return fmt.Errorf("[%s] 读取 index.html 失败: %v", label, err)
	}
	quoted := regexp.QuoteMeta(prefix + "/")
	re := regexp.MustCompile(`(?:src|href)=(["'])/(?!` + quoted + `)`)
	if re.Match(indexBytes) {
		where := "qixue-web 构建环境里的 VITE_PUBLIC_PATH"
		if label == "admin" {
			where = "qixue-admin/.env.production 的 VITE_PUBLIC_PATH"
		}
		return fmt.Errorf("[%s] 打包资源在站点根目录，和 public/%s 对不上。请把 %s 设为 /%s/ 后重新构建。", label, prefix, where, prefix)
	}

	publicDir := filepath.Dir(destDir)
	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		return fmt.Errorf("[%s] 无法创建 %s", label, publicDir)
	}
	if st, err := os.Stat(destDir); err == nil && st.IsDir() {
		if err := removeTree(destDir); err != nil {
			return fmt.Errorf("[%s] 清理旧目录失败: %v", label, err)
		}
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("[%s] 无法创建 %s", label, destDir)
	}

	files, err := listLocalFiles(distDir)
	if err != nil {
		return fmt.Errorf("[%s] 列举 dist 失败: %v", label, err)
	}
	fmt.Printf("==> [%s] 复制 %d 个文件到 public/%s\n", label, len(files), prefix)
	for rel, abs := range files {
		rel = strings.ReplaceAll(rel, "\\", "/")
		if rel == "" || strings.Contains(rel, "..") {
			return fmt.Errorf("[%s] 打包路径无效", label)
		}
		target := filepath.Join(destDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("[%s] 无法创建目录", label)
		}
		if err := copyFile(abs, target); err != nil {
			return fmt.Errorf("[%s] 复制失败 %s", label, rel)
		}
	}
	fmt.Printf("==> [%s] 完成。页面在 /%s/\n", label, prefix)
	return nil
}

func removeTree(dir string) error {
	return os.RemoveAll(dir)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
