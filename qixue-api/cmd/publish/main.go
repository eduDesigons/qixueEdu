// Command publish builds admin/web static sites and/or packages the API binary.
//
//	go run ./cmd/publish admin
//	go run ./cmd/publish web --skip-build
//	go run ./cmd/publish admin web --no-purge
//	go run ./cmd/publish api
//	go run ./cmd/publish admin web api
//
// STATIC_PUBLISH in .env: bunny | local.
// API output: ../deploy/release/qixue-api/ and qixue-api.zip.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	skipBuild := false
	noPurge := false
	var targets []string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--skip-build":
			skipBuild = true
		case "--no-purge":
			noPurge = true
		default:
			if strings.HasPrefix(arg, "--") {
				fatalf("未知参数: %s", arg)
			}
			targets = append(targets, strings.ToLower(strings.TrimSpace(arg)))
		}
	}
	if len(targets) == 0 {
		targets = []string{"admin"}
	}

	apiRoot, err := findAPIRoot()
	if err != nil {
		fatalf("%v", err)
	}
	repoRoot := filepath.Dir(apiRoot)
	envPath := filepath.Join(apiRoot, ".env")
	cfg, err := loadEnvFile(envPath)
	if err != nil {
		fatalf("%v", err)
	}

	needStatic := false
	for _, t := range targets {
		if t != "api" && t != "go" && t != "server" {
			needStatic = true
			break
		}
	}
	if needStatic {
		if _, err := publishDriver(cfg); err != nil {
			fatalf("%v", err)
		}
	}

	for _, t := range targets {
		switch t {
		case "api", "go", "server":
			if err := packageAPI(repoRoot, apiRoot, skipBuild); err != nil {
				fatalf("%v", err)
			}
		default:
			if err := publishTarget(repoRoot, apiRoot, cfg, t, skipBuild, noPurge); err != nil {
				fatalf("%v", err)
			}
		}
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func findAPIRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := cwd
	for {
		envFile := filepath.Join(dir, ".env")
		modFile := filepath.Join(dir, "go.mod")
		if fileExists(envFile) && fileExists(modFile) {
			data, err := os.ReadFile(modFile)
			if err == nil && strings.Contains(string(data), "module qixue-api") {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback: cwd is qixue-api even without matching go.mod content.
	if fileExists(filepath.Join(cwd, ".env")) {
		return cwd, nil
	}
	return "", fmt.Errorf("找不到 qixue-api/.env（请在 qixue-api 目录执行，或按 .env.example 填写）")
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func publishTarget(repoRoot, apiRoot string, cfg map[string]string, target string, skipBuild, noPurge bool) error {
	driver, err := publishDriver(cfg)
	if err != nil {
		return err
	}
	apiKey := cfgGet(cfg, "BUNNY_ACCOUNT_API_KEY")
	region := cfgGet(cfg, "BUNNY_REGION", "sg")
	publicRoot := filepath.Join(apiRoot, "public")

	var opt publishOpt
	switch target {
	case "admin":
		dir := filepath.Join(repoRoot, "qixue-admin")
		opt = publishOpt{
			Label:       "admin",
			Zone:        cfgGet(cfg, "BUNNY_ADMIN_STORAGE_ZONE"),
			Password:    cfgGet(cfg, "BUNNY_ADMIN_STORAGE_PASSWORD"),
			Region:      cfgGet(cfg, "BUNNY_ADMIN_REGION", region),
			DistDir:     filepath.Join(dir, "dist"),
			ProjectDir:  dir,
			NPMScript:   "build",
			PullZoneIDs: parsePullZoneIDs(cfgGet(cfg, "BUNNY_ADMIN_PULL_ZONE_ID")),
			PurgeAfter:  !noPurge && envYes(cfg, "BUNNY_ADMIN_PURGE_AFTER_PUBLISH", true),
			APIKey:      apiKey,
			Driver:      driver,
			PublicDir:   filepath.Join(publicRoot, "admin"),
			URLPrefix:   "admin",
		}
	case "h5", "web":
		dir := filepath.Join(repoRoot, "qixue-web")
		dist := cfgGet(cfg, "BUNNY_H5_DIST", "qixue-web/dist")
		if !isAbsPath(dist) {
			dist = filepath.Join(repoRoot, filepath.FromSlash(dist))
		}
		pullRaw := cfgGet(cfg, "BUNNY_H5_PULL_ZONE_IDS")
		if pullRaw == "" {
			pullRaw = cfgGet(cfg, "BUNNY_H5_PULL_ZONE_ID")
		}
		opt = publishOpt{
			Label:       "web",
			Zone:        cfgGet(cfg, "BUNNY_H5_STORAGE_ZONE"),
			Password:    cfgGet(cfg, "BUNNY_H5_STORAGE_PASSWORD"),
			Region:      cfgGet(cfg, "BUNNY_H5_REGION", region),
			DistDir:     dist,
			ProjectDir:  dir,
			NPMScript:   "build",
			PullZoneIDs: parsePullZoneIDs(pullRaw),
			PurgeAfter:  !noPurge && envYes(cfg, "BUNNY_H5_PURGE_AFTER_PUBLISH", true),
			APIKey:      apiKey,
			Driver:      driver,
			PublicDir:   filepath.Join(publicRoot, "web"),
			URLPrefix:   "web",
		}
	default:
		return fmt.Errorf("未知 target: %s（admin|web|api）", target)
	}
	return publishStatic(opt, skipBuild)
}

type publishOpt struct {
	Label       string
	Zone        string
	Password    string
	Region      string
	DistDir     string
	ProjectDir  string
	NPMScript   string
	PullZoneIDs []int
	PurgeAfter  bool
	APIKey      string
	Driver      string
	PublicDir   string
	URLPrefix   string
}

func publishDriver(cfg map[string]string) (string, error) {
	driver := strings.ToLower(cfgGet(cfg, "STATIC_PUBLISH", "bunny"))
	if driver != "bunny" && driver != "local" {
		return "", fmt.Errorf("STATIC_PUBLISH 只能是 bunny 或 local")
	}
	return driver, nil
}

func publishStatic(opt publishOpt, skipBuild bool) error {
	label := opt.Label
	if !skipBuild {
		if err := npmBuild(opt.ProjectDir, opt.NPMScript, label); err != nil {
			return err
		}
	} else {
		fmt.Printf("==> [%s] 跳过构建\n", label)
	}

	index := filepath.Join(opt.DistDir, "index.html")
	if st, err := os.Stat(opt.DistDir); err != nil || !st.IsDir() || !fileExists(index) {
		return fmt.Errorf("[%s] 找不到 %s/index.html", label, opt.DistDir)
	}

	if opt.Driver == "local" {
		return publishLocal(opt)
	}

	if opt.Zone == "" || opt.Password == "" {
		return fmt.Errorf("[%s] 未配置存储区名称或密码", label)
	}

	files, err := listLocalFiles(opt.DistDir)
	if err != nil {
		return fmt.Errorf("[%s] 列举 dist 失败: %v", label, err)
	}

	const concurrency = 16
	remoteBase := ""
	fmt.Printf("==> [%s] 清空远端（整区）\n", label)
	cleared, err := clearStoragePrefix(opt.Password, opt.Zone, opt.Region, remoteBase)
	if err != nil {
		return err
	}
	fmt.Printf("    远端已空（删除 %d 项）\n", cleared)

	fmt.Printf("==> [%s] 上传 %d 个文件（并发 %d）\n", label, len(files), concurrency)
	if err := uploadFilesConcurrent(opt.Password, opt.Zone, opt.Region, files, remoteBase, concurrency); err != nil {
		return err
	}

	if opt.PurgeAfter {
		if err := purgePullZoneCaches(opt.APIKey, opt.PullZoneIDs, label); err != nil {
			return err
		}
	} else {
		fmt.Printf("==> [%s] 跳过 CDN 缓存清理\n", label)
	}

	fmt.Printf("==> [%s] 完成\n", label)
	return nil
}

func isAbsPath(p string) bool {
	if filepath.IsAbs(p) {
		return true
	}
	// Unix-style absolute from env on Windows.
	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "\\") {
		return true
	}
	return false
}
