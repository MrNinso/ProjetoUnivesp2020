package main

import (
	"ProjetoUnivesp2020/managers"
	"ProjetoUnivesp2020/managers/routes/api"
	"ProjetoUnivesp2020/managers/routes/terminalSocket"
	"ProjetoUnivesp2020/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

var FrontPages = []string{"login", "room", "rooms", "admin"}

func init() {
	rand.Seed(time.Now().Unix())
	//err := docker.KillAllTerminals()
	//utils.CheckPanic(&err)
}

func main() {
	router := gin.Default()

	router.Static("/res", "./public/res")
	router.Static("/md", "./public/temp")

	router.GET("/terminal/:ID", terminalSocket.HandleTerminalSocket)

	router.GET("/api/:ACTION/*ARG", api.Handles.Run)
	router.POST("/api/:ACTION/*ARG", api.Handles.Run)

	router.GET("/app/*D", func(c *gin.Context) {
		d := c.Param("D")

		if utils.ContainsAny(d, FrontPages) {
			c.File("./public/site/build/index.html")
			return
		}

		if strings.Contains(d, "..") {
			c.String(450, "Nice try")
			return
		}

		if d == "/" {
			c.Redirect(301, "/app/login")
			return
		}

		c.File("./public/site/build/" + d)
	})

	_ = router.RunTLS(
		managers.Configs.Bind,
		managers.Configs.SSL.CertPath,
		managers.Configs.SSL.KeyPath,
	)
}
