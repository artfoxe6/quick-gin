package scaffold

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const templateModulePath = "github.com/artfoxe6/quick-gin"

// Options controls how the scaffold process runs.
type Options struct {
	ModulePath string
	Force      bool
}

// Run creates a new project at the given target path using the quick-gin template.
func Run(target string, opts Options) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return errors.New("project name is required")
	}

	destPath, err := resolveTargetPath(target)
	if err != nil {
		return err
	}

	if err := prepareDestination(destPath, opts.Force); err != nil {
		return err
	}

	sourceDir, err := sourceRoot()
	if err != nil {
		return fmt.Errorf("determine template directory: %w", err)
	}

	if err := copyTemplate(sourceDir, destPath); err != nil {
		return fmt.Errorf("copy template files: %w", err)
	}

	moduleName := strings.TrimSpace(opts.ModulePath)
	if moduleName == "" {
		moduleName = defaultModuleName(destPath)
	}

	if err := updateModule(destPath, moduleName); err != nil {
		return fmt.Errorf("update module path: %w", err)
	}

	if err := rewriteImports(destPath, moduleName); err != nil {
		return fmt.Errorf("rewrite imports: %w", err)
	}

	if err := runGoCommand(destPath, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	fmt.Printf("✅ Project created at %s\n", destPath)
	fmt.Printf("   Module: %s\n", moduleName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", destPath)
	fmt.Println("  2. go run ./cmd/app")
	return nil
}

func resolveTargetPath(target string) (string, error) {
	if filepath.IsAbs(target) {
		return target, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return filepath.Join(wd, target), nil
}

func prepareDestination(path string, force bool) error {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("target %s exists and is not a directory", path)
		}
		if !force {
			return fmt.Errorf("target directory %s already exists (use --force to overwrite)", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove existing directory: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check target directory: %w", err)
	}
	return os.MkdirAll(path, 0o755)
}

func sourceRoot() (string, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok && buildInfo.Main.Path != "" && buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		if dir, err := moduleCacheDir(buildInfo.Main.Path, buildInfo.Main.Version); err == nil {
			return dir, nil
		}
	}

	if dir, err := goEnvPath("GOMOD"); err == nil && dir != "" && dir != os.DevNull {
		return filepath.Dir(dir), nil
	}
	return "", errors.New("unable to locate template source directory")
}

func moduleCacheDir(modulePath, version string) (string, error) {
	modCache, err := goEnvPath("GOMODCACHE")
	if err != nil {
		return "", err
	}
	if modCache == "" {
		return "", errors.New("GOMODCACHE not found")
	}

	escapedPath := escapeModulePath(modulePath)
	escapedVersion := escapeModuleVersion(version)
	dir := filepath.Join(modCache, escapedPath+"@"+escapedVersion)

	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir, nil
	}
	return "", fmt.Errorf("module not found in cache: %s", dir)
}

func goEnvPath(name string) (string, error) {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value, nil
	}
	out, err := exec.Command("go", "env", name).Output()
	if err != nil {
		return "", fmt.Errorf("go env %s: %w", name, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func escapeModulePath(path string) string {
	return escapeUpperAndBang(path)
}

func escapeModuleVersion(version string) string {
	return escapeUpperAndBang(version)
}

func escapeUpperAndBang(input string) string {
	var b strings.Builder
	for _, r := range input {
		switch {
		case r == '!':
			b.WriteString("!!")
		case 'A' <= r && r <= 'Z':
			b.WriteByte('!')
			b.WriteRune(r + ('a' - 'A'))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

var skipDirs = map[string]struct{}{
	".git":              {},
	".github":           {},
	".gocache":          {},
	".idea":             {},
	".vscode":           {},
	"internal/scaffold": {},
}

var skipFiles = map[string]struct{}{
	"main.go": {},
}

func copyTemplate(src, dst string) error {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	absDstWithSep := absDst + string(os.PathSeparator)

	return filepath.WalkDir(absSrc, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == absDst || strings.HasPrefix(path, absDstWithSep) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(absSrc, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if _, ok := skipDirs[rel]; ok {
				return fs.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}

		if _, ok := skipFiles[rel]; ok {
			return nil
		}

		return copyFile(path, filepath.Join(dst, rel), d)
	})
}

func copyFile(src, dst string, entry fs.DirEntry) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	mode := fs.FileMode(0o644)
	if info, err := entry.Info(); err == nil {
		mode = info.Mode()
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func defaultModuleName(dest string) string {
	name := filepath.Base(dest)
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "-")
	if name == "" {
		return templateModulePath
	}
	return name
}

func updateModule(dest, moduleName string) error {
	if moduleName == "" {
		return errors.New("module name is empty")
	}
	if moduleName == templateModulePath {
		return nil
	}
	return runGoCommand(dest, "mod", "edit", "-module", moduleName)
}

func rewriteImports(dest, moduleName string) error {
	if moduleName == templateModulePath {
		return nil
	}
	return filepath.WalkDir(dest, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dest, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if _, skip := skipDirs[rel]; skip || rel == "vendor" {
				return fs.SkipDir
			}
			return nil
		}

		if filepath.Ext(rel) != ".go" && rel != "go.mod" {
			return nil
		}
		return replaceInFile(path, templateModulePath, moduleName)
	})
}

func replaceInFile(path, old, new string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	updated := strings.ReplaceAll(string(data), old, new)
	if updated == string(data) {
		return nil
	}
	return os.WriteFile(path, []byte(updated), 0o644)
}

func runGoCommand(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
