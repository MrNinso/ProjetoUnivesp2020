package database

import (
	"ProjetoUnivesp2020/objets"
	"ProjetoUnivesp2020/utils"
	"errors"
	"github.com/HouzuoGuo/tiedot/db"
)

var Conn = InitDataBase("./temp")

type DB struct {
	*db.DB
}

// Users
func (database DB) CreateUser(user *objets.User) error {
	if database.FindUserIdByEmail(user.Email) == -1 {
		c := database.Use(objets.USER_COL_NAME)
		_, err := c.Insert(*user.ToMap())
		return err
	} else {
		return errors.New("Usuario já existe")
	}
}

func (database DB) UpdateUser(email string, user *objets.User) error {
	id := database.FindUserIdByEmail(email)

	if id == -1 {
		return errors.New("Usuario Não encontrado")
	}

	c := database.Use(objets.USER_COL_NAME)

	return c.Update(id, *user.ToMap())
}

func (database DB) ListAllUsers() *[]objets.User {
	us := make([]objets.User, 0)

	_ = database.ForEachUser(func(id int, u *objets.User) (moveOn bool) {
		us = append(us, objets.User{
			Name:    u.Name,
			Email:   u.Email,
			IsAdmin: u.IsAdmin,
		})

		return true
	})

	return &us
}

func (database DB) DeleteUser(user *objets.User) error {
	id := database.FindUserIdByEmail(user.Email)

	if id == -1 {
		return errors.New("Usuario não encontrado")
	}

	return database.DeleteUserByID(id)
}

func (database DB) DeleteUserByID(id int) error {
	return database.Use(objets.IMAGE_COL_NAME).Delete(id)
}

func (database DB) ForEachUser(f func(id int, u *objets.User) (moveOn bool)) *db.Col {
	c := database.Use(objets.USER_COL_NAME)

	c.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		u, _ := objets.UserFromJson(doc)

		return f(id, u)
	})

	return c
}

func (database DB) FindUserIdByEmail(email string) int {
	uId := -1

	database.ForEachUser(func(id int, u *objets.User) (moveOn bool) {
		if u.Email == email {
			uId = id
			return false
		}

		return true
	})

	return uId
}

func (database DB) FindUserByEmail(email string) *objets.User {
	var U *objets.User

	database.ForEachUser(func(id int, u *objets.User) (moveOn bool) {
		if u.Email == email {
			U = u
			return false
		}

		return true
	})

	return U
}

// Images
func (database DB) CreateImage(img *objets.Image) error {
	c, id := database.FindImageUIDByDockerName(img.Name)

	if id != "" {
		return errors.New("Já existe uma image")
	}

	_, err := c.Insert(*img.ToMap())

	return err
}

func (database DB) UpdateImage(UID string, img *objets.Image) error {
	id := database.FindImageIDByUID(UID)

	if id == -1 {
		return errors.New("Imagem não encontrada")
	}

	c := database.Use(objets.IMAGE_COL_NAME)

	img.UId = UID

	return c.Update(id, *img.ToMap())
}

func (database DB) ListAllImages() *[]objets.Image {
	imgs := make([]objets.Image, 0)

	_ = database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		imgs = append(imgs, *i)
		return true
	})

	return &imgs
}

func (database DB) FindImageDockerNameByUID(UID string) string {
	name := ""

	_ = database.ForEachImage(func(i int, img *objets.Image) (moveOn bool) {
		if img.UId == UID {
			name = img.Name
			return false
		}

		return true
	})

	return name
}

func (database DB) FindImageIDByUID(UID string) int {
	Id := -1

	_ = database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		if i.UId == UID {
			Id = id
			return false
		}

		return true
	})

	return Id
}

func (database DB) FindImageUIDByDockerName(dockerName string) (*db.Col, string) {
	Id := ""

	c := database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		if i.Name == dockerName {
			Id = i.UId
			return false
		}

		return true
	})

	return c, Id
}

func (database DB) DeleteImageIfExist(uId string) (error, *db.Col) {
	dbId := -1

	c := database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		if uId == i.UId {
			dbId = id
			return false
		}
		return true
	})

	if dbId != -1 {
		return c.Delete(dbId), c
	}

	return nil, c
}

func (database DB) ForEachImage(f func(id int, i *objets.Image) (moveOn bool)) *db.Col {
	c := database.Use(objets.IMAGE_COL_NAME)

	c.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		i, _ := objets.ImageFromJson(doc)

		return f(id, i)
	})

	return c
}

//Rooms
func (database DB) CreateRoom(r *objets.Room) error {
	if database.FindRoomIdByUId(r.GetUID()) == -1 {
		c := database.Use(objets.ROOM_COL_NAME)
		_, err := c.Insert(*r.ToMap())
		return err
	} else {
		return errors.New("Sala já existe")
	}
}

func (database DB) ListAllRooms() *[]objets.Room {
	rooms := make([]objets.Room, 0)

	database.ForEachRoom(func(id int, r *objets.Room) (moveOn bool) {
		rooms = append(rooms, *r)
		return true
	})

	return &rooms
}

func (database DB) UpdateRoom(UId string, r *objets.Room) error {
	id := database.FindRoomIdByUId(UId)

	if id == -1 {
		return errors.New("Usuario Não encontrado")
	}

	c := database.Use(objets.USER_COL_NAME)

	return c.Update(id, *r.ToMap())
}

func (database DB) FindRoomByUID(UId string) *objets.Room {
	var R *objets.Room

	R = nil

	database.ForEachRoom(func(id int, r *objets.Room) (moveOn bool) {
		if r.GetUID() == UId {
			R = r
			return false
		}

		return true
	})

	return R
}

func (database DB) FindRoomIdByUId(UId string) int {
	Id := -1

	database.ForEachRoom(func(id int, r *objets.Room) (moveOn bool) {
		if r.GetUID() == UId {
			Id = id
			return false
		}

		return true
	})

	return Id
}

func (database DB) ForEachRoom(f func(id int, r *objets.Room) (moveOn bool)) *db.Col {
	c := database.Use(objets.ROOM_COL_NAME)

	c.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		r := objets.RoomFromJson(doc)

		return f(id, r)
	})

	return c
}

// Database setup
func InitDataBase(path string) *DB {
	d, err := db.OpenDB(path)

	utils.CheckPanic(&err)

	database := &DB{d}

	if database.ColExists(objets.USER_COL_NAME) {
		return database
	}

	createDatabase(database)

	return database
}

func createDatabase(database *DB) {
	err := database.Create(objets.USER_COL_NAME)

	utils.CheckPanic(&err)

	err = database.Create(objets.IMAGE_COL_NAME)

	utils.CheckPanic(&err)

	err = database.Create(objets.ROOM_COL_NAME)

	utils.CheckPanic(&err)

	err = database.CreateUser(&objets.User{
		Email:    "admin@admin.com",
		Password: "$2a$12$k4zPsZOfqhiBXcPS2XPfEOVFUiQST0LmVuqwutEkM0IIUTDxJzM5G", //password: admin
		Name:     "Admin User",
		IsAdmin:  true,
	})

	utils.CheckPanic(&err)
}
