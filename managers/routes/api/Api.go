package api

import (
	"ProjetoUnivesp2020/managers/auth"
	"ProjetoUnivesp2020/managers/config"
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/managers/docker"
	"ProjetoUnivesp2020/managers/log"
	"ProjetoUnivesp2020/managers/room"
	"ProjetoUnivesp2020/objets"
	"fmt"
	"github.com/MrNinso/MyGoToolBox/lang/crypto"
	"github.com/MrNinso/MyGoToolBox/lang/ifs"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"strconv"
	"strings"
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
					apiError(c, 403, "Forbidden")
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

var roomManager = room.RenderRooms()

// API Handles
var Handles = &handles{
	{"CreateRoom", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(room.RoomJsonHeaderKey)

		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		r, err := objets.NewRoomFromJson(uuid.New().String(), []byte(jsonString))

		if err != nil || r == nil {
			apiError(c, 400, err)
			return
		}

		if ifs.IfAnyStringEmpty(r.GetTitle(), r.GetContetMd(), r.GetImageUID()) {
			apiError(c, 400, "Bad Request")
			return
		}

		if err := database.Conn.CreateRoom(r); err != nil {
			apiError(c, 500, err)
			return
		}

		if err := room.AddRoomToManager(roomManager, r); err != nil {
			apiError(c, 500, err)
			return
		}

		c.JSON(200, r)
	}, true, true},
	{"GetRoomByID", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		room := roomManager.GetRoomByID(strings.Replace(args, "/", "", 1))

		if room == nil {
			apiError(c, 404, "Room not found")
			return
		}

		c.JSON(200, gin.H{
			"title":    room.GetTitle(),
			"imageUId": room.GetImageUID(),
		})
	}, false, true},
	{"ListRooms", func(c *gin.Context, args string) {
		json, err := roomManager.ToJson()

		if err != nil {
			apiError(c, 500, err.Error())
		}

		c.String(200, json)
	}, false, true},
	{"UpdateRoom", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(room.RoomJsonHeaderKey)

		if jsonString == "" || args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		r, err := objets.NewRoomFromJson(args, []byte(jsonString))

		if err != nil || r == nil {
			apiError(c, 400, err)
			return
		}

		if ifs.IfAnyStringEmpty(r.GetTitle(), r.GetContetMd(), r.GetImageUID()) {
			apiError(c, 400, "Bad Request")
			return
		}

		if err := database.Conn.UpdateRoom(args, r); err != nil {
			apiError(c, 500, err)
			return
		}

		room.UpdateRoomFromManager(roomManager, r)

		c.JSON(200, r)
	}, true, true},
	{"DeleteRoom", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		pos := roomManager.GetRoomPosByID(args)

		if pos == -1 {
			apiError(c, 404, "Room não encotrada")
			return
		}

		room.RemoveRoomFromManager(roomManager, pos)

		c.String(200, "")
	}, true, true},
	{"RenderRoomByUID", func(c *gin.Context, args string) {
		if args == "" {
			c.JSON(400, gin.H{"error": "Bad Request"})
			return
		}

		if r := roomManager.GetRoomByID(args); r != nil {
			err := roomManager.RenderRoom(r)
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
		roomManager = room.RenderRooms()
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

					go func() {
						l := log.Manager.GetLogLevel(ifs.IfReturn(admin, "infoAdmin", "infoUser").(string))
						l.AppendLog(fmt.Sprintf("%s - entrou", email))
					}()

					c.JSON(200, gin.H{"isAdmin": admin})
				} else {
					apiError(c, 403, "auto login failed")
				}
				return
			}
		}

		//login with email and password
		email := c.GetHeader(auth.EMAIL_HEADER_KEY)
		token = c.GetHeader(auth.TOKEN_HEADER_KEY)

		if ifs.IfAnyStringEmpty(email, token) {
			apiError(c, 400, "Bad Request")
			return
		}

		secret, admin := auth.Login(token, email)

		if secret == "" {
			apiError(c, 401, "Login Fail")
			return
		}

		go func() {
			l := log.Manager.GetLogLevel(ifs.IfReturn(admin, "infoAdmin", "infoUser").(string))
			l.AppendLog(fmt.Sprintf("%s - entrou", email))
		}()

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

		c.SetCookie(
			auth.COOKIE_ISADMIN, ifs.IfReturn(admin, "1", "0").(string), 0,
			"/", c.Request.Host,
			true, false,
		)

		c.JSON(200, gin.H{})

	}, false, false},
	{"CreateUser", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.UserHeaderKey)
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
		jsonString := c.Request.Header.Get(objets.UserHeaderKey)

		if jsonString == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		user, err := objets.UserFromJson([]byte(jsonString))

		if err != nil || user == nil {
			apiError(c, 400, err)
			return
		}

		oldEmail := c.Request.Header.Get(objets.UpdateEmailHeaderKey)

		if oldEmail == "" {
			apiError(c, 400, "Bad Request")
		}

		_, oldUser := database.Conn.FindUserByEmail(oldEmail)

		if oldUser == nil {
			apiError(c, 404, "User not found")
			return
		}

		if user.Password != "" {
			b, err := bcrypt.GenerateFromPassword([]byte(user.Password), config.Configs.BcryptCost)

			if err != nil {
				apiError(c, 500, err)
			}

			user.Password = string(b)
		}

		if oldEmail != user.Email {
			if id, _ := database.Conn.FindUserByEmail(user.Email); id != -1 {
				apiError(c, 400, "Email já cadastrado")
				return
			}
		}

		err = database.Conn.UpdateUser(oldUser.Email, user)

		if err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")
	}, true, true},
	{"DeleteUser", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.UserHeaderKey)
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

		if ifs.IfAnyStringEmpty(img.Name, img.Dockerfile) {
			apiError(c, 400, "Bad Request")
			return
		}

		img.UId = uuid.New().String()
		img.DockerImageName = crypto.ToMD5Hash(img.UId)

		if err = docker.BuildImage(img.DockerImageName, img.Dockerfile); err != nil {
			apiError(c, 500, err)
			return
		}

		img.Created = strconv.FormatInt(time.Now().Unix(), 10)

		if err := database.Conn.CreateImage(img); err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")
	}, true, true},
	{"UpdateImage", func(c *gin.Context, args string) {
		jsonString := c.Request.Header.Get(objets.UserHeaderKey)

		if ifs.IfAnyStringEmpty(jsonString, args) {
			apiError(c, 400, "Bad Request")
			return
		}

		img, err := objets.ImageFromJson([]byte(jsonString))

		if err != nil || img == nil {
			apiError(c, 400, err)
			return
		}

		if ifs.IfAnyStringEmpty(img.Name, img.Dockerfile, img.UId) {
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
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		_, image := database.Conn.FindImageByUID(args)

		if image == nil {
			apiError(c, 404, "Imagem não encontrada")
			return
		}

		if err := docker.RemoveImage(image.Name); err != nil {
			apiError(c, 500, err)
			return
		}

		c.String(200, "")

	}, true, true},
	{"ListImages", func(c *gin.Context, args string) {
		c.JSON(200, database.Conn.ListAllImages())
	}, true, true},
	{"GetImageByUId", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}
		_, img := database.Conn.FindImageByUID(args)

		if img == nil {
			apiError(c, 404, "Not Found")
			return
		}

		c.JSON(200, img)
	}, true, true},
	{"GetImageByUIdName", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}
		_, img := database.Conn.FindImageByUID(args)

		if img == nil {
			apiError(c, 404, "Not Found")
			return
		}

		c.String(200, img.Name)
	}, true, true},

	{"LimparLog", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}
		l := log.Manager.GetLogLevel(args)

		if l == nil {
			apiError(c, 404, "Not Found")
			return
		}

		l.ClearLog()
		c.String(200, "")
	}, true, true},
	{"ListarLogs", func(c *gin.Context, args string) {
		c.JSON(200, log.Manager.GetAllLogsLevels())
	}, true, true},
	{"ShowLogs", func(c *gin.Context, args string) {
		if args == "" {
			apiError(c, 400, "Bad Request")
			return
		}

		l := log.Manager.GetLogLevel(args)

		if l == nil {
			apiError(c, 404, "Not Found")
			return
		}

		c.JSON(200, l.ReadLog())
	}, true, true},
}
