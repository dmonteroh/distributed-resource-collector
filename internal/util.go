package internal

import (
	"fmt"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}

func UniqueString(slice []string) (unique []string) {
	for _, v := range slice {
		skip := false
		for _, u := range unique {
			if v == u {
				skip = true
				break
			}
		}
		if !skip {
			unique = append(unique, v)
		}
	}
	return unique
}
