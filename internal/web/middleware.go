package web

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

func CheckPath() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Query("path")
		if len(path) == 0 {
			path = ctx.PostForm("path")
		}
		if len(path) == 0 {
			var request models.PathRequest
			err := ctx.ShouldBindJSON(&request)
			if err != nil {
				ctx.JSON(models.ErrIpt("参数解析错误"))
				ctx.Abort()
				return
			}
			path = request.Path
		}
		if len(path) == 0 {
			ctx.JSON(models.ErrIpt("缺少参数: path"))
			ctx.Abort()
			return
		}
		p, err := utils.CheckPath(path)
		if err != nil {
			log.DebugErr(err)
			ctx.JSON(models.ErrIpt("路径参数错误"))
			ctx.Abort()
			return
		}
		ctx.Set("path", p)
		ctx.Next()
	}
}

// KeyAuth 鉴权
func KeyAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenRaw := ctx.Request.FormValue("access_key") // query/form 查找 access_key
		if len(tokenRaw) == 0 {
			tokenRaw = ctx.Request.Header.Get("Access-Key") // header 查找 access_key
			if len(tokenRaw) == 0 {
				ctx.JSON(models.ErrJwt("未发现access_key"))
				ctx.Abort()
				return
			}
		}
		localKey := utils.Sha256(AccessKey)
		if tokenRaw != localKey {
			ctx.JSON(models.ErrJwt("Access key错误"))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Key, Authorization, Token, session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		c.Next()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqPath := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		ext := xpath.Ext(reqPath)
		if ext != "" && ext != ".html" {
			return
		}
		cost := time.Since(start)
		logger.Infof("%s %s {query %s}, %v, %v, 响应: %d, 耗时: %dms",
			c.Request.Method,
			reqPath,
			query,
			zap.String("ip", c.ClientIP()),
			// zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			c.Writer.Status(),
			cost.Milliseconds(),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.SugaredLogger, stack bool, response func(*gin.Context, any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					response(c, err)
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				response(c, err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
