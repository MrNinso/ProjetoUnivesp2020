package auth

import (
	"ProjetoUnivesp2020/managers/config"
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/objets"
	"crypto/sha1"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

const (
	charSet        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*"
	COOKIE_EMAIL   = "3ic7k5irhh2az9hkig1oy3"
	COOKIE_TOKEN   = "97b31ae2cd1a382f19a7b95f5ef98016"
	COOKIE_ISADMIN = "9c5dc968ba3e168aac5e7b3e6182e785"

	EMAIL_HEADER_KEY = "Email"
	TOKEN_HEADER_KEY = "Token"
)

//Token = sha256(email_&_md5(senha))
func Login(token, email string) (string, bool) {
	_, User := database.Conn.FindUserByEmail(email)

	if User == nil {
		return "", false
	}

	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(token))

	if err != nil {
		return "", false
	}

	s := CreateSecret()

	hash := sha1.New()
	hash.Write([]byte(token + s))

	bs, _ := bcrypt.GenerateFromPassword(hash.Sum(nil), config.Configs.BcryptCost)

	User.Secret = string(bs)

	if err = database.Conn.UpdateUser(email, User); err != nil {
		return "", false
	}

	return User.Secret, User.IsAdmin
}

func Register(user *objets.User) error {
	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), config.Configs.BcryptCost)

	if err != nil {
		return err
	}

	user.Password = string(b)

	return database.Conn.CreateUser(user)
}

func CheckSecretToken(email, secretToken string) (bool, bool) {
	_, User := database.Conn.FindUserByEmail(email)
	if User == nil {
		return false, false
	}

	return User.Secret == secretToken, User.IsAdmin
}

func CreateSecret() string {
	var output strings.Builder

	for i := 0; i < config.Configs.BcryptCost; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
