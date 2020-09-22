package utils

import (
	"os"
	"strconv"
	"strings"
)

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
