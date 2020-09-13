package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"syscall"
	"testsGin/api"
	"testsGin/podman"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"testsGin/utils"
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

func handleTerminalSocket(w http.ResponseWriter, r *http.Request, imageID string) {
	conn, err := upgrader.Upgrade(w, r, nil)

	utils.CheckPanic(&err)

	terminal, err := podman.StartTerminal("debian") //TODO BANCO DE IMAGENS

	utils.CheckPanic(&err)

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
			//TODO LOG USUARIO DESCONECTADO
			break
		}

		if msgType == websocket.TextMessage {
			fmt.Println("msgType == TextMessage") //TODO MELHORAR LOG
			_ = conn.WriteMessage(websocket.TextMessage, []byte("msgType == TextMessage"))
			continue
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)

		utils.CheckPanic(&err)

		if read != 1 {
			fmt.Println("bytes", read) //TODO MELHORAR LOG
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

func main() {
	err := podman.KillAllTerminals()

	utils.CheckPanic(&err)

	router := gin.Default()

	router.Static("/res", "./public/res")

	router.Static("/room/", "./public/html/room")

	router.Static("/md", "./public/temp")

	router.GET("/terminal", func(c *gin.Context) {
		handleTerminalSocket(c.Writer, c.Request, "debian")
	})

	router.GET("/api/:ACTION/:ARG", api.Handles.Run)

	router.POST("/api/:ACTION/:ARG", api.Handles.Run)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/home")
	})

	_ = router.RunTLS(
		":1443",
		"certs/server.crt",
		"certs/server.key",
	)
}
