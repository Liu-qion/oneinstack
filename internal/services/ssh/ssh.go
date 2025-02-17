package ssh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func OpenWebShell(c *gin.Context) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 创建带有伪终端的命令
	cmd := exec.Command("bash")
	ptmx, err := pty.StartWithAttrs(cmd, &pty.Winsize{Rows: 24, Cols: 80}, &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	})
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to start shell: "+err.Error()))
		return
	}
	defer ptmx.Close()
	defer cmd.Process.Kill()

	// 优化输出处理
	go func() {
		buf := make([]byte, 4096) // 增大缓冲区
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
				return
			}
			// 使用二进制消息类型发送原始终端数据
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// 优化输入处理
	go func() {
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			switch messageType {
			case websocket.TextMessage:
				// 先尝试解析窗口大小
				var size struct {
					Rows uint16 `json:"rows"`
					Cols uint16 `json:"cols"`
				}
				if err := json.Unmarshal(data, &size); err == nil {
					pty.Setsize(ptmx, &pty.Winsize{
						Rows: size.Rows,
						Cols: size.Cols,
					})
					continue
				}

				// 处理普通文本输入
				input := append(data, '\r') // 添加回车符
				if _, err := ptmx.Write(input); err != nil {
					return
				}

			case websocket.BinaryMessage:
				// 直接写入二进制数据
				if _, err := ptmx.Write(data); err != nil {
					return
				}
			}
		}
	}()

	// 保持连接
	select {}
}
