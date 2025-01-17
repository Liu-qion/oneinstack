package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
)

// DecompressTarGz 跨平台解压 tar.gz 文件
func DecompressTarGz(src string, dest string) error {
	// 打开 tar.gz 文件
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解压 gzip
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("无法创建 gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// 解压 tar 文件
	tarReader := tar.NewReader(gzipReader)

	// 逐个解包文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 文件读取完毕
		}
		if err != nil {
			return fmt.Errorf("解包文件时出错: %v", err)
		}

		// 构建解压路径
		targetPath := filepath.Join(dest, header.Name)

		// 如果是 Windows 系统，需要做路径兼容
		if runtime.GOOS == "windows" {
			targetPath = filepath.ToSlash(targetPath)
		}

		// 创建目标文件或目录
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("无法创建目录: %v", err)
			}
		case tar.TypeReg:
			// 创建文件
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("无法创建文件: %v", err)
			}
			defer outFile.Close()

			// 将文件内容写入目标文件
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return fmt.Errorf("写入文件时出错: %v", err)
			}
		default:
			return fmt.Errorf("不支持的文件类型: %v", header.Typeflag)
		}
	}

	return nil
}

// DownloadFile 下载文件到指定路径
func DownloadFile(url string, destPath string) error {
	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// 创建文件
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// 将响应体复制到文件中
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	fmt.Printf("Downloaded %s successfully.\n", destPath)
	return nil
}

// SetExecPermissions 设置指定目录及其子目录下所有文件的执行权限
func SetExecPermissions(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只为文件设置执行权限，忽略目录
		if !info.IsDir() {
			err = os.Chmod(path, 0755)
			if err != nil {
				return fmt.Errorf("无法设置 %s 的执行权限: %v", path, err)
			}
			fmt.Printf("已设置执行权限: %s\n", path)
		}
		return nil
	})
}

func GetLogContent(logFilePath string) (string, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return "", fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	var content []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
		content = append(content, '\n')
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return string(content), nil
}

func FormatBytes(bytes int64) string {
	if bytes < 0 {
		return fmt.Sprintf("-%s", FormatBytes(-bytes))
	}

	units := []string{"B", "KB", "MB", "GB", "TB"}
	var unitIndex int
	var value float64

	for unitIndex = 0; bytes >= 1024 && unitIndex < len(units)-1; unitIndex++ {
		value = float64(bytes) / 1024
		bytes = int64(value)
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%dB", bytes)
	}

	// 格式化输出，保留两位小数
	return fmt.Sprintf("%.2f%s", value, units[unitIndex])
}

func LookupUser(uid int) string {
	user, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return fmt.Sprintf("无法查找用户: %v", err)
	}
	return user.Name
}

func LookupGroup(gid int) string {
	group, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		return fmt.Sprintf("无法查找组: %v", err)
	}
	return group.Name
}
