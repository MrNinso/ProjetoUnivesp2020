package objets

import (
	"encoding/json"
	. "github.com/MrNinso/MyGoToolBox/lang/ifs"
)

const (
	// User Column name in database
	UserColName = "USERS"
	// User json key header in HTTP request
	UserHeaderKey = "USER"
	// User email key header in a update user HTTP request
	UpdateEmailHeaderKey = "UPDATE_EMAIL"
)

// User type
type User struct {
	Email    string
	Password string
	Name     string
	IsAdmin  bool
	Secret   string
}

// Create User from json
func UserFromJson(j []byte) (*User, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)

	if err != nil {
		return nil, err
	}

	return UserFromMap(m), nil
}

// Create User from json
func UserFromMap(m map[string]interface{}) *User {
	return &User{
		Email:    IfNil(m["Email"], "").(string),
		Password: IfNil(m["Password"], "").(string),
		Name:     IfNil(m["Name"], "").(string),
		IsAdmin:  IfNil(m["IsAdmin"], false).(bool),
		Secret:   IfNil(m["Secret"], "").(string),
	}
}

// Convert User to map
func (u User) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"Email":    u.Email,
		"Password": u.Password,
		"Name":     u.Name,
		"IsAdmin":  u.IsAdmin,
		"Secret":   u.Secret,
	}
}
