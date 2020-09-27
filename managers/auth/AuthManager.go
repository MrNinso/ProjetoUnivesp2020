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

const (
	charSet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*"
	COOKIE_EMAIL = "3ic7k5irhh2az9hkig1oy3"
	COOKIE_TOKEN = "97b31ae2cd1a382f19a7b95f5ef98016"

	EMAIL_HEADER_KEY = "Email"
	TOKEN_HEADER_KEY = "Token"
)

//Token = sha256(email_&_md5(senha))
func Login(token, email string) string {
	User := database.Conn.FindUserByEmail(email)

	if User == nil {
		return ""
	}

	b, _ := bcrypt.GenerateFromPassword([]byte(token), managers.Configs.BcryptCost)
	println(string(b))

	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(token))

	if err != nil {
		return ""
	}

	s := CreateSecret()

	hash := sha1.New()
	hash.Write([]byte(token + s))

	bs, _ := bcrypt.GenerateFromPassword(hash.Sum(nil), managers.Configs.BcryptCost)

	User.Secret = string(bs)

	if err = database.Conn.UpdateUserByEmail(email, User); err != nil {
		return ""
	}

	return User.Secret
}

func Register(user *objets.User) error {
	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), managers.Configs.BcryptCost)

	if err != nil {
		return err
	}

	user.Password = string(b)

	return database.Conn.CreateUser(user)
}

func CheckSecretToken(email, secretToken string) (bool, bool) {
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
