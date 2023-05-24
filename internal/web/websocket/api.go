package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/logger"
)

// Log godoc
//
//	@Summary WebSocket日志监听接口
//	@Description 监听日志接口
//	@Tags websocket
//	@Produce json
//	@Success 101 {string} string "Switching Protocols"
//	@Security ApiKeyAuth
//	@Router /websocket/log [get]
func (w *WebSocket) Log(c *gin.Context) {
	if c.Request.Header.Get("Upgrade") != "websocket" {
		c.String(http.StatusOK, "")
		return
	}
	w.wsHandler(c.Writer, c.Request,
		func() {
			logger.EnableLogNotify()
		},
		func() {
			if len(w.wsConns) == 0 {
				logger.DisableLogNotify()
			}
		})
}
