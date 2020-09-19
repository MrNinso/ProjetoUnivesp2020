package auth

import (
	"ProjetoUnivesp2020/managers"
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/objets"
	"crypto/sha1"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

var charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*"

//Token = sha256(email_&_md5(senha))
func Login(token, email string) string {
	User := database.Conn.FindUserByEmail(email)

	if User == nil {
		return ""
	}

	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(token))

	if err != nil {
		return ""
	}

	s := CreateSecret()

	hash := sha1.New()
	hash.Write([]byte(token + s))

	User.Secret = string(hash.Sum(nil))

	if err = database.Conn.UpdateUserByEmail(email, User); err != nil {
		return ""
	}

	return s
}

func Register(user *objets.User) error {
	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), managers.Configs.BcryptCost)

	if err != nil {
		return err
	}

	user.Password = string(b)

	return database.Conn.CreateUser(user)
}

func CheckSecret(email, secretToken string) (bool, bool) {
	User := database.Conn.FindUserByEmail(email)
	if User == nil {
		return false, false
	}

	return User.Secret == secretToken, User.IsAdmin
}

func CreateSecret() string {
	var output strings.Builder

	for i := 0; i < managers.Configs.BcryptCost; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
