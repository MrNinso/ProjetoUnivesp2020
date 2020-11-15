package terminalsocket

import (
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/managers/docker"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"syscall"
	"unsafe"
)

type windowSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
	X    uint16
	Y    uint16
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handle Terminal WebSocket
func HandleTerminalSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(500, err)
	}

	_, image := database.Conn.FindImageByUID(c.Param("ID"))

	if image == nil {
		c.String(404, "Image Not Found")
		return
	}

	terminal, err := docker.StartTerminal(image.DockerImageName)

	if err != nil {
		c.JSON(500, err)
	}

	defer func() {
		_ = conn.Close()
		terminal.Kill()
	}()

	go func() {
		for {
			buf := make([]byte, upgrader.WriteBufferSize)
			read, err := terminal.TTY.Read(buf)

			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				fmt.Println(err.Error())
				return
			}

			_ = conn.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		msgType, reader, err := conn.NextReader()

		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("msgType == TextMessage"))
			continue
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)

		if err != nil {
			c.JSON(500, err)
		}

		if read != 1 {
			return
		}

		switch dataTypeBuf[0] {
		case 0:
			_, _ = io.Copy(terminal.TTY, reader)
		case 1:
			decoder := json.NewDecoder(reader)
			resizeMessage := windowSize{}
			err := decoder.Decode(&resizeMessage)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte("Error decoding resize message: "+err.Error()))
				continue
			}

			_, _, errno := syscall.Syscall(
				syscall.SYS_IOCTL,
				terminal.TTY.Fd(),
				syscall.TIOCSWINSZ,
				uintptr(unsafe.Pointer(&resizeMessage)),
			)

			if errno != 0 {
				fmt.Println("Unable to resize terminal")
			}
		}

	}
}
