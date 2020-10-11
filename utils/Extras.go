package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
)

type MapFunc func(index, count int, item interface{}) interface{}

type FindFunc func(index int, item interface{}) bool

func IfNil(a, b interface{}) interface{} {
	if a == nil {
		return b
	}

	return a
}

func IfEmpty(a, b string) string {
	if a == "" {
		return b
	}

	return a
}

func GetEnv(a, b string) string {
	return IfEmpty(os.Getenv(a), b)
}

func GetIntFromEnv(a string, b int) int {
	s := GetEnv(a, "")

	if s == "" {
		return b
	}

	i, err := strconv.Atoi(s)

	CheckPanic(&err)

	return i
}

func ContainsAny(a string, b []string) bool {
	for _, c := range b {
		if strings.Contains(a, c) {
			return true
		}
	}

	return false
}

func CheckStringField(f ...string) bool {
	for i := 0; i < len(f); i++ {
		if f[i] == "" {
			return true
		}
	}

	return false
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func OneLineIf(contion bool, t, f interface{}) interface{} {
	if contion {
		return t
	}

	return f
}

func FindItemOnList(source interface{}, h FindFunc) interface{} {
	s := (source).([]interface{})

	for i := 0; i < len(s); i++ {
		if h(i, s[i]) {
			return s[i]
		}
	}

	return nil
}
