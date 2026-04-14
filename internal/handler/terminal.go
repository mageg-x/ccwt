package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ccwt/ccwt/internal/middleware"
	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type WsMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
	Rows uint16          `json:"rows,omitempty"`
	Cols uint16          `json:"cols,omitempty"`
}

type ResizeData struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}

// TerminalWS WebSocket 终端连接
func TerminalWS(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr, _ = c.Cookie("token")
	}
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	claims, err := middleware.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: user=%s err=%v", claims.Username, err)
		return
	}
	defer conn.Close()

	sessID := c.Query("session_id")
	rows := uint16(24)
	cols := uint16(80)

	// 获取或创建 PTY 会话
	var sess *service.PtySession
	if sessID != "" {
		sess = service.Pty.Get(sessID)
	}
	if sess == nil {
		sessID = uuid.New().String()
		sess, err = service.Pty.Create(sessID, claims.Username, rows, cols)
		if err != nil {
			log.Printf("PTY 创建失败: user=%s err=%v", claims.Username, err)
			conn.WriteJSON(gin.H{"type": "error", "data": "终端创建失败"})
			return
		}
	}

	// 发送会话 ID
	conn.WriteJSON(gin.H{"type": "session", "data": sessID})

	// 发送回滚缓冲区数据（断线重连恢复）
	if buf := sess.Buf.Bytes(); len(buf) > 0 {
		conn.WriteMessage(websocket.BinaryMessage, buf)
	}

	// 订阅终端输出
	outCh := sess.Subscribe()
	defer sess.Unsubscribe(outCh)

	// PTY 输出 → WebSocket
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case data, ok := <-outCh:
				if !ok {
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					return
				}
			case <-sess.Done():
				conn.WriteJSON(gin.H{"type": "exit"})
				return
			}
		}
	}()

	// WebSocket → PTY 输入
	conn.SetReadDeadline(time.Time{})
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			var wsMsg WsMsg
			if json.Unmarshal(msg, &wsMsg) != nil {
				// 非 JSON 控制消息，按普通输入透传到 PTY。
				sess.Write(msg)
				continue
			}
			switch wsMsg.Type {
			case "resize":
				if wsMsg.Rows > 0 && wsMsg.Cols > 0 {
					service.Pty.Resize(sessID, wsMsg.Rows, wsMsg.Cols)
				} else {
					var rd ResizeData
					json.Unmarshal(wsMsg.Data, &rd)
					service.Pty.Resize(sessID, rd.Rows, rd.Cols)
				}
			case "input":
				var input string
				if err := json.Unmarshal(wsMsg.Data, &input); err != nil {
					// 兼容 data 直接传入字符串字节的情况
					input = string(wsMsg.Data)
				}
				sess.Write([]byte(input))
			case "ping":
				conn.WriteJSON(gin.H{"type": "pong"})
			default:
				// 不识别的 JSON 文本，按原始文本透传
				sess.Write(msg)
			}
		} else if msgType == websocket.BinaryMessage {
			sess.Write(msg)
		}
	}

	<-done
}

// ListTerminals 列出用户终端
func ListTerminals(c *gin.Context) {
	username, _ := c.Get("username")
	sessions := service.Pty.List(username.(string))
	var out []gin.H
	for _, s := range sessions {
		out = append(out, gin.H{"id": s.ID, "created_at": s.CreateAt})
	}
	if out == nil {
		out = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"sessions": out})
}

// CloseTerminal 关闭终端
func CloseTerminal(c *gin.Context) {
	id := c.Param("id")
	username, _ := c.Get("username")
	sess := service.Pty.Get(id)
	if sess == nil || sess.UserName != username.(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}
	service.Pty.Close(id)
	c.JSON(http.StatusOK, gin.H{"message": "已关闭"})
}
