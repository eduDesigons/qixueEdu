package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// packageAPI cross-compiles the API (linux/amd64 by default) into deploy/release/.
func packageAPI(repoRoot, apiRoot string, skipBuild bool) error {
	goos := envOr("GO_API_GOOS", "linux")
	goarch := envOr("GO_API_GOARCH", "amd64")
	outDir := filepath.Join(repoRoot, "deploy", "release", "qixue-api")
	binName := "server"
	binPath := filepath.Join(outDir, binName)
	zipPath := filepath.Join(repoRoot, "deploy", "release", "qixue-api.zip")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("[api] 无法创建 %s: %v", outDir, err)
	}

	if skipBuild && fileExists(binPath) {
		fmt.Printf("==> [api] 跳过构建，沿用 %s\n", binPath)
	} else {
		fmt.Printf("==> [api] go build (%s/%s) → %s\n", goos, goarch, binPath)
		if err := buildAPIBinary(apiRoot, binPath, goos, goarch); err != nil {
			return err
		}
	}

	if err := syncFonts(apiRoot, outDir); err != nil {
		return err
	}
	if err := copyIfExists(
		filepath.Join(apiRoot, ".env.example"),
		filepath.Join(outDir, ".env.example"),
	); err != nil {
		return err
	}
	svcSrc := filepath.Join(repoRoot, "deploy", "security", "qixue-api.service")
	if err := copyIfExists(svcSrc, filepath.Join(outDir, "qixue-api.service")); err != nil {
		return err
	}

	fmt.Printf("==> [api] 打包 zip → %s\n", zipPath)
	if err := zipDir(outDir, zipPath, "qixue-api"); err != nil {
		return fmt.Errorf("[api] 打 zip 失败: %v", err)
	}

	st, _ := os.Stat(binPath)
	zst, _ := os.Stat(zipPath)
	fmt.Printf("==> [api] 完成\n")
	fmt.Printf("    目录 %s\n", outDir)
	if st != nil {
		fmt.Printf("    二进制 %s (%d bytes, %s/%s)\n", binPath, st.Size(), goos, goarch)
	}
	if zst != nil {
		fmt.Printf("    压缩包 %s (%d bytes)\n", zipPath, zst.Size())
	}
	return nil
}

func buildAPIBinary(apiRoot, outPath, goos, goarch string) error {
	_ = os.Remove(outPath)
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", outPath, "./cmd/server")
	cmd.Dir = apiRoot
	env := dropEnvKeys(filteredEnv(), "GOOS", "GOARCH", "CGO_ENABLED", "GOPROXY", "GOSUMDB")
	proxy := envOr("GO_API_GOPROXY", envOr("GOPROXY", "https://goproxy.cn,direct"))
	sumdb := envOr("GO_API_GOSUMDB", envOr("GOSUMDB", "off"))
	env = append(env,
		"CGO_ENABLED=0",
		"GOOS="+goos,
		"GOARCH="+goarch,
		"GOPROXY="+proxy,
		"GOSUMDB="+sumdb,
	)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return buildFail("api", err)
	}
	if !fileExists(outPath) {
		return fmt.Errorf("[api] 构建后找不到 %s", outPath)
	}
	// linux binary on Windows has no +x bit that survives; chmod on unix build host
	if runtime.GOOS != "windows" {
		_ = os.Chmod(outPath, 0o755)
	}
	return nil
}

func syncFonts(apiRoot, outDir string) error {
	src := filepath.Join(apiRoot, "storage", "fonts")
	dst := filepath.Join(outDir, "storage", "fonts")
	if st, err := os.Stat(src); err != nil || !st.IsDir() {
		fmt.Printf("==> [api] 警告：找不到 %s，证书合成可能缺字体\n", src)
		return nil
	}
	_ = os.RemoveAll(dst)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("[api] 无法创建字体目录: %v", err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if err := copyFile(filepath.Join(src, name), filepath.Join(dst, name)); err != nil {
			return fmt.Errorf("[api] 复制字体失败 %s: %v", name, err)
		}
		n++
	}
	fmt.Printf("==> [api] 复制字体 %d 个 → storage/fonts\n", n)
	return nil
}

func copyIfExists(src, dst string) error {
	if !fileExists(src) {
		fmt.Printf("==> [api] 跳过（不存在）%s\n", src)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("[api] 复制失败 %s → %s: %v", src, dst, err)
	}
	return nil
}

func zipDir(srcDir, zipPath, rootPrefix string) error {
	_ = os.Remove(zipPath)
	if err := os.MkdirAll(filepath.Dir(zipPath), 0o755); err != nil {
		return err
	}
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()
	zw := zip.NewWriter(out)
	defer zw.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		name := filepath.ToSlash(filepath.Join(rootPrefix, rel))
		if info.IsDir() {
			_, err := zw.Create(name + "/")
			return err
		}
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		return err
	}
	return zw.Close()
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func dropEnvKeys(env []string, keys ...string) []string {
	deny := map[string]struct{}{}
	for _, k := range keys {
		deny[strings.ToUpper(k)] = struct{}{}
	}
	out := make([]string, 0, len(env))
	for _, e := range env {
		key, _, _ := strings.Cut(e, "=")
		if _, bad := deny[strings.ToUpper(key)]; bad {
			continue
		}
		out = append(out, e)
	}
	return out
}
