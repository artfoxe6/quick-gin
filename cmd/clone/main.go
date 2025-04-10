package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
)

func main() {
	var path, pkgName string
	flag.StringVar(&path, "path", "", "project path")
	flag.StringVar(&pkgName, "package", "", "package name")
	flag.Parse()

	if path == "" || pkgName == "" {
		fmt.Println("project path is required, eg: --path=/opt/app/test")
		return
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	err := run(getDirName(), path, pkgName)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func run(srcDir, destDir, packageName string) error {
	cmd := exec.Command("cp", "-r", srcDir+"/.", destDir)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	cmd = exec.Command("rm", "-rf", destDir+"/.git")
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	files := getAllFiles(destDir)
	for _, file := range files {
		err := replaceFileContent(file, getPackage(), packageName)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func replaceFileContent(filePath, oldStr, newStr string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	newContent := strings.ReplaceAll(string(content), oldStr, newStr)
	if _, err = file.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return err
	}
	if err = file.Truncate(0); err != nil {
		fmt.Println("Error truncating file:", err)
		return err
	}
	_, err = file.WriteAt([]byte(newContent), 0)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}
func getAllFiles(dir string) []string {
	var files []string
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			files = append(files, getAllFiles(dir+"/"+fileInfo.Name())...)
		} else {
			files = append(files, dir+"/"+fileInfo.Name())
		}
	}
	return files
}

func getDirName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func getPackage() string {
	info, _ := debug.ReadBuildInfo()
	if info.Main.Path != "" {
		return info.Main.Path
	}
	out, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(out), "\n", "")
}
