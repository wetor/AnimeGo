package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wetor/AnimeGo/docs"
	"github.com/wetor/AnimeGo/internal/store"
	"os"
)

// @termsOfService https://github.com/wetor/AnimeGo
// @license.name MIT
// @license.url https://www.mit-license.org/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Access-Key

func InitSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "AnimeGo"
	docs.SwaggerInfo.Description = "Golang开发的自动追番与下载工具"
	docs.SwaggerInfo.Version = os.Getenv("ANIMEGO_VERSION")
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", store.Config.Setting.WebApi.Host, store.Config.Setting.WebApi.Port)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
