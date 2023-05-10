package websocket

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

type WebSocket struct {
	upgrader   *websocket.Upgrader
	wsConnLock sync.Mutex
	wsConns    []*websocket.Conn
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (w *WebSocket) Start(ctx context.Context) {
	WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 renamer")
					exit = true
					return
				case logData := <-Notify:
					if logger.GetLogNotify() {
						w.wsConnLock.Lock()
						for _, conn := range w.wsConns {
							if err := conn.WriteMessage(websocket.TextMessage, logData); err != nil {
								log.Warnf("Failed to send log message to WebSocket client: %v", err)
							}
						}
						w.wsConnLock.Unlock()
					}
				}
			}()
			if exit {
				return
			}
		}
	}()
}

func (w *WebSocket) addConn(conn *websocket.Conn) {
	w.wsConnLock.Lock()
	w.wsConns = append(w.wsConns, conn)
	w.wsConnLock.Unlock()
}

func (w *WebSocket) deleteConn(conn *websocket.Conn) {
	w.wsConnLock.Lock()
	for i, c := range w.wsConns {
		if c == conn {
			// Remove closed connection from the slice
			w.wsConns = append(w.wsConns[:i], w.wsConns[i+1:]...)
			break
		}
	}
	w.wsConnLock.Unlock()
}

func (w *WebSocket) wsHandler(resp http.ResponseWriter, req *http.Request, before func(), after func()) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Warnf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	if before != nil {
		before()
	}
	w.addConn(conn)
	defer func() {
		w.deleteConn(conn)
		_ = conn.Close()
		if after != nil {
			after()
		}
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
