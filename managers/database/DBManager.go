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

func (database DB) CreateUser(user *objets.User) error {
	if database.FindUserIdByEmail(user.Email) == -1 {
		c := database.Use(objets.USER_COL_NAME)
		_, err := c.Insert(*user.ToMap())
		return err
	} else {
		return errors.New("Usuario já existe")
	}
}

func (database DB) CreateImage(img *objets.Image) error {
	err, c := database.DeleteImageIfExist(img.UId, img.PodmanId)

	if err != nil {
		return err
	}

	_, err = c.Insert(*img.ToMap())

	return err
}

func (database DB) DeleteImageIfExist(uId, podmanID string) (error, *db.Col) {
	dbId := -1

	c := database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		if uId == i.UId || podmanID == i.PodmanId {
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

func (database DB) DeleteUser(id int) error {
	return database.Use(objets.IMAGE_COL_NAME).Delete(id)
}

func (database DB) UpdateUserByEmail(email string, user *objets.User) error {
	id := database.FindUserIdByEmail(email)

	if id == -1 {
		return errors.New("Usuario Não encontrado")
	}

	c := database.Use(objets.USER_COL_NAME)

	return c.Update(id, *user.ToMap())
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

func (database DB) FindImagePodmanIdByUId(UId string) string {
	Id := ""

	database.ForEachImage(func(id int, i *objets.Image) (moveOn bool) {
		if i.UId == UId {
			Id = i.PodmanId
			return false
		}
		return true
	})

	return Id
}

func (database DB) ForEachUser(f func(id int, u *objets.User) (moveOn bool)) *db.Col {
	c := database.Use(objets.USER_COL_NAME)

	c.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		u := objets.UserFromJson(doc)

		return f(id, u)
	})

	return c
}

func (database DB) ForEachImage(f func(id int, i *objets.Image) (moveOn bool)) *db.Col {
	c := database.Use(objets.IMAGE_COL_NAME)

	c.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		i := objets.ImageFromJson(doc)

		return f(id, i)
	})

	return c
}

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

	err = database.CreateUser(&objets.User{
		Email:    "admin@admin.com",
		Password: "$2a$12$k4zPsZOfqhiBXcPS2XPfEOVFUiQST0LmVuqwutEkM0IIUTDxJzM5G", //password: admin
		Name:     "Admin User",
		IsAdmin:  true,
	})

	utils.CheckPanic(&err)
}
