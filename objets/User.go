package objets

import (
	"ProjetoUnivesp2020/utils"
	"encoding/json"
)

const USER_COL_NAME = "USERS"

type User struct {
	Email    string
	Password string
	Name     string
	IsAdmin  bool
	Secret   string
}

func UserFromJson(j []byte) *User {
	m := make(map[string]interface{})
	_ = json.Unmarshal(j, &m)
	return UserFromMap(m)
}

func UserFromMap(m map[string]interface{}) *User {
	return &User{
		Email:    utils.IfNil(m["Email"], "").(string),
		Password: utils.IfNil(m["Password"], "").(string),
		Name:     utils.IfNil(m["Name"], "").(string),
		IsAdmin:  utils.IfNil(m["IsAdmin"], false).(bool),
		Secret:   utils.IfNil(m["Secret"], "").(string),
	}
}

func (u User) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"Email":    u.Email,
		"Password": u.Password,
		"Name":     u.Name,
		"IsAdmin":  u.IsAdmin,
		"Secret":   u.Secret,
	}
}
