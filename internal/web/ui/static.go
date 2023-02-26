package ui

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
)

//go:embed static
var staticFiles embed.FS

func RegisterStatic(router gin.IRouter) {
	// 将静态文件夹嵌入到二进制文件中，并设置静态文件路由
	router.GET("/static/*filepath", func(c *gin.Context) {
		// 从嵌入式文件系统中获取文件
		filePath := c.Param("filepath")
		content, err := staticFiles.ReadFile(path.Join("static", filePath))
		if err != nil {
			// 如果文件不存在，则返回 404 错误
			c.Status(http.StatusNotFound)
			return
		}

		// 设置 MIME 类型并返回文件内容
		c.Data(http.StatusOK, getContentType(filePath), content)
	})
}

// 根据文件后缀名获取 MIME 类型
func getContentType(filePath string) string {
	if strings.HasSuffix(filePath, ".css") {
		return "text/css"
	} else if strings.HasSuffix(filePath, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(filePath, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filePath, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		return "image/jpeg"
	} else {
		return "text/plain"
	}
}
