package main

import (
	"ProjetoUnivesp2020/managers"
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/managers/podman"
	"ProjetoUnivesp2020/managers/routes/api"
	"ProjetoUnivesp2020/managers/routes/terminalSocket"
	"ProjetoUnivesp2020/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	err := podman.KillAllTerminals()
	configs := managers.LoadConfigs()

	_ = database.InitDataBase("./temp")

	utils.CheckPanic(&err)

	router := gin.Default()

	//router.Static("/", "./public/") TODO REDIRECT /home
	router.Static("/res", "./public/res")
	router.Static("/md", "./public/temp")

	router.GET("/terminal", terminalSocket.HandleTerminalSocket)
	router.GET("/api/:ACTION/:ARG", api.Handles.Run)

	router.POST("/api/:ACTION/:ARG", api.Handles.Run)

	_ = router.RunTLS(
		configs.Bind,
		configs.SSL.CertPath,
		configs.SSL.KeyPath,
	)
}
