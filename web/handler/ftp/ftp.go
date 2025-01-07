package ftp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"oneinstack/core"
	"oneinstack/utils"
	"os"
	"path/filepath"
)

// 列出目录内容
func ListDirectory(c *gin.Context) {
	var input struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	absPath := filepath.Join(filepath.Clean(input.Path))
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	var fileInfos []gin.H
	for _, file := range files {
		info, _ := file.Info()
		fileInfos = append(fileInfos, gin.H{
			"name":  file.Name(),
			"isDir": file.IsDir(),
			"size":  utils.FormatBytes(info.Size()),
		})
	}
	core.HandleSuccess(c, gin.H{"files": fileInfos})
}

// 创建文件或目录
func CreateFileOrDir(c *gin.Context) {
	var input struct {
		Path string `json:"path" binding:"required"`
		Type string `json:"type" binding:"required"` // "file" 或 "dir"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	absPath := filepath.Join(filepath.Clean(input.Path))
	switch input.Type {
	case "file":
		f, err := os.Create(absPath)
		if err != nil {
			core.HandleError(c, http.StatusInternalServerError, err, nil)
			return
		}
		defer f.Close()
	case "dir":
		if err := os.MkdirAll(absPath, 0755); err != nil {
			core.HandleError(c, http.StatusInternalServerError, err, nil)
			return
		}
	default:
		core.HandleError(c, http.StatusInternalServerError, fmt.Errorf("无效类型"), nil)
		return
	}

	core.HandleSuccess(c, "创建成功")

}

// 上传文件
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	path := c.PostForm("path")
	if path == "" {
		path = "/"
	}

	absPath := filepath.Join(filepath.Clean(path), file.Filename)
	if err := c.SaveUploadedFile(file, absPath); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, "上传成功")
}

// 下载文件
func DownloadFile(c *gin.Context) {
	var input struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	filePath := filepath.Join(filepath.Clean(input.Path))
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filepath.Base(filePath))))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	io.Copy(c.Writer, file)
}

// 删除文件或目录
func DeleteFileOrDir(c *gin.Context) {
	var input struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	absPath := filepath.Join(filepath.Clean(input.Path))
	if err := os.RemoveAll(absPath); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, "删除成功")
}
