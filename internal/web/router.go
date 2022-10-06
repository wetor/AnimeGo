package web

import (
	"AnimeGo/internal/store"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Run(ctx context.Context) {
	store.WG.Add(1)
	go func() {
		defer store.WG.Done()
		r := gin.New()

		r.Use(Cors())             // 跨域中间件
		r.Use(GinLogger(zap.S())) // 日志中间件
		r.Use(GinRecovery(zap.S(), true, func(c *gin.Context, recovered interface{}) {
			if err, ok := recovered.(error); ok {
				c.JSON(ErrSvr(err.Error()))
			} else {
				c.JSON(ErrSvr(recovered.(string)))
			}
		})) // 错误处理中间件
		gin.SetMode("debug")

		r.GET("/ping", Pong)
		r.POST("/api", Download)
		s := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", store.Config.WebApi.Host, store.Config.WebApi.Port),
			Handler: r,
		}
		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				zap.S().Debug(err)
				zap.S().Warn("启动web服务失败")
			}
		}()
		select {
		case <-ctx.Done():
			if err := s.Close(); err != nil {
				zap.S().Debug(err)
				zap.S().Warn("关闭web服务失败")
			}
			zap.S().Info("正常退出")
		}
	}()
}
