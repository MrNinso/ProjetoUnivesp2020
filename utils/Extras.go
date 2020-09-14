package utils

import "os"

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