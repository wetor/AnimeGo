package web

import (
	"embed"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/utils"
)

//go:embed static
var staticFiles embed.FS

func InitStatic(router gin.IRouter) {
	InitAssets(router)
}

func InitAssets(router gin.IRouter) {
	// 将静态文件夹嵌入到二进制文件中，并设置静态文件路由
	router.GET("/*filepath", func(c *gin.Context) {
		// 从嵌入式文件系统中获取文件
		filePath := c.Param("filepath")
		if filePath == "/" || filePath == "" {
			filePath = "index.html"
		}
		var err error
		var content []byte
		localPath := path.Join(constant.WebPath, filePath)
		if utils.IsExist(localPath) {
			content, err = os.ReadFile(localPath)
		} else {
			content, err = staticFiles.ReadFile(path.Join("static", filePath))
		}
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
