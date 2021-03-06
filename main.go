package main

import (
	"ProjetoUnivesp2020/managers/config"
	"ProjetoUnivesp2020/managers/docker"
	"ProjetoUnivesp2020/managers/routes/api"
	"ProjetoUnivesp2020/managers/routes/terminalsocket"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

var frontPages = []string{"login", "room", "home", "admin", "demo"}

func init() {
	rand.Seed(time.Now().Unix())
	_ = docker.KillAllTerminals()
}

func main() {
	router := gin.Default()

	router.Static("/res", "./public/res")
	router.Static("/md", "./public/res/temp")

	router.GET("/terminal/:ID", terminalsocket.HandleTerminalSocket)

	router.GET("/api/:ACTION/*ARG", api.Handles.Run)
	router.POST("/api/:ACTION/*ARG", api.Handles.Run)

	router.GET("/app/*D", func(c *gin.Context) {
		d := c.Param("D")

		for i := 0; i < len(frontPages); i++ {
			if strings.Contains(d, frontPages[i]) {
				c.File("./public/site/build/index.html")
				return
			}
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

	go func() {
		if gin.Mode() == gin.DebugMode {
			log.Println("Debug mode on 0.0.0.0:6060/debug/pprof")
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}
	}()

	_ = router.RunTLS(
		config.Configs.Bind,
		config.Configs.SSL.CertPath,
		config.Configs.SSL.KeyPath,
	)
}
