package utils

import (
	"fmt"
)

func CheckPanic(err *error) {
	if *err != nil {
		fmt.Println(*err)
		panic(err)
	}
}
