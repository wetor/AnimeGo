package web

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wetor/AnimeGo/internal/web/docs"
)

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag fmt -g swagger.go
//go:generate swag init -g swagger.go

//	@termsOfService	https://github.com/wetor/AnimeGo
//	@license.name	MIT
//	@license.url	https://www.mit-license.org/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Access-Key

func InitSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "AnimeGo"
	docs.SwaggerInfo.Description = "Golang开发的自动追番与下载工具"
	docs.SwaggerInfo.Version = os.Getenv("ANIMEGO_VERSION")
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", Host, Port)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
