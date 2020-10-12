package objets

import (
	"encoding/json"
	. "github.com/MrNinso/MyGoToolBox/lang/ifs"
)

const (
	USER_COL_NAME           = "USERS"
	USER_HEADER_KEY         = "USER"
	UPDATE_EMAIL_HEADER_KEY = "UPDATE_EMAIL"
)

type User struct {
	Email    string
	Password string
	Name     string
	IsAdmin  bool
	Secret   string
}

func UserFromJson(j []byte) (*User, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)

	if err != nil {
		return nil, err
	}

	return UserFromMap(m), nil
}

func UserFromMap(m map[string]interface{}) *User {
	return &User{
		Email:    IfNil(m["Email"], "").(string),
		Password: IfNil(m["Password"], "").(string),
		Name:     IfNil(m["Name"], "").(string),
		IsAdmin:  IfNil(m["IsAdmin"], false).(bool),
		Secret:   IfNil(m["Secret"], "").(string),
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
