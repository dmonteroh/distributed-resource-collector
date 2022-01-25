package internal

import (
	"fmt"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}
