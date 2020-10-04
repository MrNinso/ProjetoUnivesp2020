package api

import (
	"ProjetoUnivesp2020/managers"
	"ProjetoUnivesp2020/managers/auth"
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/objets"
	"ProjetoUnivesp2020/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/google/uuid"
	"io/ioutil"
	"os/exec"
	"time"
)

type handle func(c *gin.Context, args string)

type action struct {
	name         string
	r            handle
	requireAdmin bool
	requireUser  bool
}

type handles []action

func (hs handles) Run(c *gin.Context) {
	if a := c.Param("ACTION"); a != "" {
		if h := hs.GetActionByName(a); h != nil {
			if h.requireUser {
				token, errToken := c.Cookie(auth.COOKIE_TOKEN)
				email, errEmail := c.Cookie(auth.COOKIE_EMAIL)

				if errToken != nil || errEmail != nil {
					c.String(400, "Bad Request")
					return
				}

				isValidLogin, isAdmin := auth.CheckSecretToken(email, token)

				if isValidLogin {
					if h.requireAdmin {
						if isAdmin {
							h.r(c, c.Param("ARG"))
						} else {
							apiError(c, 403, "Forbidden")
						}
					} else {
						h.r(c, c.Param("ARG"))
					}
				} else {
					apiError(c, 403, "Forbidden")
				}
			} else {
				h.r(c, c.Param("ARG"))
			}
		} else {
			apiError(c, 404, "Not Found")
		}
	} else {
		apiError(c, 400, "Bad request")
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

func apiError(c *gin.Context, code int, erroMsg interface{}) {
	c.JSON(code, gin.H{"error": erroMsg})
}

var RoomManager = managers.RenderRooms()

var Handles = &handles{
	{"CreateRoom", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(managers.ROOM_JSON_HEADER_KEY)

		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		room, err := objets.NewRoomFromJson(uuid.New().String(), []byte(jsonString))

		if err != nil || room == nil {
			apiError(c, 400, err)
			return
		}

		if utils.CheckStringField(room.GetTitle(), room.GetContetMd(), room.GetImageUID()) {
			apiError(c, 400, "Bad Request")
			return
		}

		if err := database.Conn.CreateRoom(room); err != nil {
			apiError(c, 500, err)
			return
		}

		if err := managers.AddRoomToManager(RoomManager, room); err != nil {
			apiError(c, 500, err)
			return
		}

		c.JSON(200, room)
	}, true, true},
	{"ListRooms", func(c *gin.Context, args string) {
		c.JSON(200, RoomManager)
	}, false, true},
	{"UpdateRoom", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(managers.ROOM_JSON_HEADER_KEY)

		if jsonString == "" || args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		room, err := objets.NewRoomFromJson(args, []byte(jsonString))

		if err != nil || room == nil {
			apiError(c, 400, err)
			return
		}

		if utils.CheckStringField(room.GetTitle(), room.GetContetMd(), room.GetImageUID()) {
			apiError(c, 400, "Bad Request")
			return
		}

		if err := database.Conn.UpdateRoom(args, room); err != nil {
			apiError(c, 500, err)
			return
		}

		managers.UpdateRoomFromManager(RoomManager, room)

		c.JSON(200, room)
	}, true, true},
	{"DeleteRoom", func(c *gin.Context, args string) {
		//TODO
	}, true, true},
	{"RenderRoomByUID", func(c *gin.Context, args string) {
		if args == "" {
			c.JSON(400, gin.H{"error": "Bad Request"})
			return
		}

		if r := RoomManager.GetRoomByID(args); r != nil {
			err := RoomManager.RenderRoom(r)
			if err != nil {
				apiError(c, 500, err)
				return
			}

			c.JSON(200, r)
		} else {
			apiError(c, 404, "Not Found")
		}
	}, true, true},
	{"RenderAllRooms", func(c *gin.Context, args string) {
		RoomManager = managers.RenderRooms()
	}, true, true},
	{"LiveRoomRender", func(c *gin.Context, args string) {
		md, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, string(markdown.ToHTML(md, nil, nil)))

	}, true, true},

	{"Login", func(c *gin.Context, args string) {
		token, err := c.Cookie(auth.COOKIE_TOKEN)

		//login with cookie
		if err == nil {
			email, err := c.Cookie(auth.COOKIE_EMAIL)

			if err == nil {
				login, admin := auth.CheckSecretToken(email, token)
				if login {
					c.JSON(200, gin.H{"isAdmin": admin})
				} else {
					apiError(c, 403, "auto login failed")
				}
				return
			}
		}

		//login with email and password
		email := c.PostForm(auth.EMAIL_HEADER_KEY)
		token = c.PostForm(auth.TOKEN_HEADER_KEY)

		if utils.CheckStringField(email, token) {
			apiError(c, 400, "Bad Request")
			return
		}

		secret := auth.Login(token, email)

		if secret == "" {
			apiError(c, 401, "Login Fail")
			return
		}

		c.SetCookie(
			auth.COOKIE_EMAIL, email, 0,
			"/", c.Request.Host,
			true, false,
		)

		c.SetCookie(
			auth.COOKIE_TOKEN, secret, 0,
			"/", c.Request.Host,
			true, false,
		)

		c.JSON(200, gin.H{})

	}, false, false},
	{"CreateUser", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.USER_HEADER_KEY)
		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		user, err := objets.UserFromJson([]byte(jsonString))

		if err != nil || user == nil {
			apiError(c, 400, err)
			return
		}

		if err := database.Conn.CreateUser(user); err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")
	}, true, true},
	{"UpdateUser", func(c *gin.Context, args string) {
		//TODO
	}, true, true},
	{"DeleteUser", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.USER_HEADER_KEY)
		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		user, err := objets.UserFromJson([]byte(jsonString))

		if err != nil || user == nil {
			apiError(c, 400, err)
			return
		}

		if err := database.Conn.DeleteUser(user); err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")
	}, true, true},
	{"ListUsers", func(c *gin.Context, args string) {
		c.JSON(200, database.Conn.ListAllUsers())
	}, true, true},

	{"CreateImage", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.IMAGE_HEADER_KEY)

		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		img, err := objets.ImageFromJson([]byte(jsonString))

		if err != nil || img == nil {
			apiError(c, 400, err)
			return
		}

		if utils.CheckStringField(img.Name, img.DockerFile) {
			apiError(c, 400, "Bad Request")
			return
		}

		//TODO BUILD DA IMAGEN

		img.UId = uuid.New().String()
		img.Created = time.Now().Unix()

		if err := database.Conn.CreateImage(img); err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")
	}, true, true},
	{"UpdateImage", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.USER_HEADER_KEY)

		if utils.CheckStringField(jsonString, args) {
			apiError(c, 400, "Bad Request")
			return
		}

		img, err := objets.ImageFromJson([]byte(jsonString))

		if err != nil || img == nil {
			apiError(c, 400, err)
			return
		}

		if utils.CheckStringField(img.Name, img.DockerFile, img.UId) {
			apiError(c, 400, "Bad Request")
			return
		}

		if err := database.Conn.UpdateImage(args, img); err != nil {
			apiError(c, 500, err)
			return
		}

		c.JSON(200, "")
	}, true, true},
	{"DeleteImage", func(c *gin.Context, args string) {
		//TODO
	}, true, true},
	{"ListImages", func(c *gin.Context, args string) {
		c.JSON(200, database.Conn.ListAllImages())
	}, true, true},

	{"ListRunningContainers", func(c *gin.Context, args string) {
		cmd := exec.Command("docker", "ps")

		defer func() {
			_ = cmd.Process.Kill()
		}()

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
			apiError(c, 400, "Bad Request")
		}

	}, true, true},
}