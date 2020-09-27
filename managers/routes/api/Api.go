package api

import (
	"ProjetoUnivesp2020/managers"
	"ProjetoUnivesp2020/managers/auth"
	"github.com/gin-gonic/gin"
	"os/exec"
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
							c.String(403, "Forbidden")
						}
					} else {
						h.r(c, c.Param("ARG"))
					}
				} else {
					c.String(403, "Forbidden")
				}
			} else {
				h.r(c, c.Param("ARG"))
			}
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

var RoomManager = managers.RenderRooms() //TODO

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
	}, true, true},
	{"ListRooms", func(c *gin.Context, args string) {
		if args == "json" {
			c.JSON(200, RoomManager)
		} else if args == "csv" {
			c.String(200, RoomManager.ToCSV())
		} else {
			c.JSON(400, gin.H{"error": "Bad Request"})
		}

	}, false, true},

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

	}, true, true},
	{"Login", func(c *gin.Context, args string) {
		token, err := c.Cookie(auth.COOKIE_TOKEN)

		//login with cookie
		if err == nil {
			email, err := c.Cookie(auth.COOKIE_EMAIL)

			if err == nil {
				login, admin := auth.CheckSecretToken(email, token)
				if login {
					c.JSON(200, gin.H{
						"isAdmin": admin,
					})
				} else {
					c.String(403, "auto login failed")
				}
				return
			}
		}

		//login with email and password
		email := c.PostForm(auth.EMAIL_HEADER_KEY)
		token = c.PostForm(auth.TOKEN_HEADER_KEY)

		if email == "" || token == "" {
			c.String(400, "Bad Request")
			return
		}

		secret := auth.Login(token, email)

		if secret == "" {
			c.String(401, "Login Fail")
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
}
