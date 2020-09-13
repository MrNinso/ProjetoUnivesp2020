package api

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"testsGin/rooms"
)

type handle func(c *gin.Context, args string)

type action struct {
	name         string
	r            handle
	requireAdmin bool
}

type handles []action

func (hs handles) Run(c *gin.Context) {
	if a := c.Param("ACTION"); a != "" {
		if h := hs.GetActionByName(a); h != nil {
			h.r(c, c.Param("ARG")) //TODO TEST ADMIN
		}
	} else {
		c.JSON(500, gin.H{"error": "Bad request"})
		return
	}
}

func (hs handles) GetActionByName(name string) *action {
	for i := 0; i < len(hs); i++ {
		if hs[i].name == name {
			return &hs[i]
		}
	}

	return nil
}

var RoomManager = rooms.Render() //TODO

var Handles = &handles{
	{"RenderRoom", func(c *gin.Context, args string) {
		if args == "" {
			c.JSON(400, gin.H{"error": "Bad Request"})
			return
		}

		if r := RoomManager.GetRoomByID(args); r != nil {
			RoomManager.RenderRoom(r)
			c.JSON(200, r)
		} else {
			c.JSON(404, gin.H{"error": "Not Found"})
		}
	}, true},

	{"ListRooms", func(c *gin.Context, args string) {
		if args == "json" {
			c.JSON(200, RoomManager)
		} else if args == "csv" {
			c.String(200, RoomManager.ToCSV())
		} else {
			c.JSON(400, gin.H{"error": "Bad Request"})
		}

	}, false},

	{"ListRunningContainers", func(c *gin.Context, args string) {
		cmd := exec.Command("pomdman", "ps")

		defer cmd.Process.Kill()

		err := cmd.Run()

		if err != nil {
			_ = c.Error(err)
			return
		}

		//TODO

		if args == "json" {
			c.JSON(200, RoomManager)
		} else if args == "csv" {
			c.String(200, RoomManager.ToCSV())
		} else {
			c.JSON(400, gin.H{"error": "Bad Request"})
		}

	}, false},
}
