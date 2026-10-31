package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func npmBuild(projectDir, script, label string) error {
	if !fileExists(filepath.Join(projectDir, "package.json")) {
		return fmt.Errorf("[%s] 找不到 %s/package.json", label, projectDir)
	}

	// Windows: leftover dart-sass / esbuild keep files open; npm ci then EBUSY on rmdir.
	killNodeModuleLocks(projectDir)

	// Prefer install over ci on Windows — ci deletes node_modules and often hits EBUSY.
	installArgs := []string{"install", "--no-fund", "--no-audit"}
	if runtime.GOOS != "windows" && fileExists(filepath.Join(projectDir, "package-lock.json")) {
		installArgs = []string{"ci", "--no-fund", "--no-audit"}
	}

	fmt.Printf("==> [%s] npm %s && npm run %s\n", label, strings.Join(installArgs, " "), script)

	env := dropEnvKeys(filteredEnv(), "NODE_OPTIONS")
	// admin vite build needs more than 1.5G on Windows; keep install lighter if needed
	heapMB := "4096"
	if v := strings.TrimSpace(os.Getenv("GO_API_NODE_HEAP_MB")); v != "" {
		heapMB = v
	}
	env = append(env, "NODE_OPTIONS=--max-old-space-size="+heapMB)

	if err := runNPM(projectDir, env, installArgs...); err != nil {
		return buildFail(label, err)
	}
	if err := runNPM(projectDir, env, "run", script); err != nil {
		return buildFail(label, err)
	}
	return nil
}

// killNodeModuleLocks stops orphan dart/esbuild under projectDir/node_modules (Windows EBUSY).
func killNodeModuleLocks(projectDir string) {
	if runtime.GOOS != "windows" {
		return
	}
	root, err := filepath.Abs(filepath.Join(projectDir, "node_modules"))
	if err != nil || root == "" {
		return
	}
	rootLower := strings.ToLower(root)
	ps := fmt.Sprintf(
		`Get-Process -Name dart,esbuild -ErrorAction SilentlyContinue | Where-Object { $_.Path -and $_.Path.ToLower().StartsWith('%s') } | Stop-Process -Force -ErrorAction SilentlyContinue`,
		strings.ReplaceAll(rootLower, `'`, `''`),
	)
	_ = exec.Command("powershell", "-NoProfile", "-Command", ps).Run()
}

func runNPM(dir string, env []string, args ...string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// npm is typically a .cmd shim on Windows.
		cmd = exec.Command("cmd", append([]string{"/c", "npm"}, args...)...)
	} else {
		cmd = exec.Command("npm", args...)
	}
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func filteredEnv() []string {
	out := make([]string, 0, 64)
	for _, e := range os.Environ() {
		key, _, _ := strings.Cut(e, "=")
		if strings.EqualFold(key, "NODE_ENV") {
			continue
		}
		out = append(out, e)
	}
	return out
}

func buildFail(label string, err error) error {
	code := exitCode(err)
	hint := ""
	if code == 137 || code == 143 {
		hint = fmt.Sprintf("（多半是内存不足被系统杀掉。本机构建后执行 go run ./cmd/publish %s --skip-build）", label)
	}
	if code > 0 {
		return fmt.Errorf("[%s] 构建失败，退出码 %d%s", label, code, hint)
	}
	return fmt.Errorf("[%s] 构建失败: %v%s", label, err, hint)
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}
