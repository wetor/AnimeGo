package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/errors"
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
				zap.S().Debugf("服务器错误，err: %v", errors.NewAniErrorD(err))
				c.JSON(ErrSvr("服务器错误"))
			} else {
				zap.S().Debug(recovered.(string))
				c.JSON(ErrSvr(recovered.(string)))
			}
		})) // 错误处理中间件
		if Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		r.GET("/ping", Ping)
		apiRoot := r.Group("/api")
		apiRoot.Use(KeyAuth())
		apiRoot.POST("/rss", Rss)
		apiRoot.POST("/plugin/config", PluginConfigPost)
		apiRoot.GET("/plugin/config", PluginConfigGet)
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
		zap.S().Infof("github.com/wetor/AnimeGo Web服务已启动: http://%s", s.Addr)
		select {
		case <-ctx.Done():
			if err := s.Close(); err != nil {
				zap.S().Debug(err)
				zap.S().Warn("关闭web服务失败")
			}
			zap.S().Debug("正常退出 web")
		}
	}()
}
