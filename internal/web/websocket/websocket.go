package websocket

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const LogNotifyCap = 1000 // 暂停监听状态最多储存的日志数量

type WebSocket struct {
	upgrader   *websocket.Upgrader
	wsConnLock sync.Mutex
	wsConns    []*websocket.Conn
	logList    *list.List
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logList: list.New(),
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
				defer utils.HandleError(func(err error) {
					log.Errorf("%+v", err)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 websocket")
					exit = true
					return
				case logData := <-Notify:
					if logger.GetLogNotify() == logger.NotifyEnabled {
						data := bytes.NewBuffer(nil)
						if w.logList.Len() > 0 {
							data.WriteString(fmt.Sprintf(`{"type":"log","count":%d}`, w.logList.Len()))
							for e := w.logList.Front(); e != nil; e = e.Next() {
								data.WriteString("\n\n")
								data.WriteString(e.Value.(string))
							}
							w.logList.Init()
						} else {
							data.WriteString(`{"type":"log","count":1}`)
						}
						data.WriteString("\n\n")
						data.Write(logData)

						w.wsConnLock.Lock()
						for _, conn := range w.wsConns {
							if err := conn.WriteMessage(websocket.TextMessage, data.Bytes()); err != nil {
								log.DebugErr(err)
								log.Warnf("[WebSocket] 发送消息失败")
							}
						}
						w.wsConnLock.Unlock()
					} else if logger.GetLogNotify() == logger.NotifyPaused {
						w.logList.PushBack(string(logData))
						if w.logList.Len() > LogNotifyCap {
							w.logList.Remove(w.logList.Front())
						}
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
			w.wsConns = append(w.wsConns[:i], w.wsConns[i+1:]...)
			break
		}
	}
	w.wsConnLock.Unlock()
}

func (w *WebSocket) wsHandler(resp http.ResponseWriter, req *http.Request, before func(), after func()) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.DebugErr(err)
		log.Warnf("[WebSocket] 请求升级为WebSocket失败")
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
	cmd := &Command{}
	cmd.Init()
	exit := false
	cmd.SetActionFunc("terminate", func() error {
		exit = true
		return nil
	})
	for {
		messageType, data, err := conn.ReadMessage()
		if messageType == websocket.CloseMessage {
			break
		}
		if err != nil {
			log.DebugErr(err)
			log.Warnf("[WebSocket] 异常结束")
			break
		}
		err = json.Unmarshal(data, cmd)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("[WebSocket] 命令解析失败")
		}
		err = cmd.Execute()
		if err != nil {
			log.DebugErr(err)
			log.Warnf("[WebSocket] 执行命令失败：%s", data)
		} else {
			log.Infof("[WebSocket] 执行命令成功：%s", data)
		}
		if exit {
			break
		}
	}
}
