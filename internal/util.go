package internal

import (
	"fmt"
	"os"
	"strings"
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

func CustomContains(str string, subStrings ...string) bool {
	if len(subStrings) == 0 {
		return true
	}

	for _, subString := range subStrings {
		if strings.Contains(str, subString) {
			return true
		}
	}
	return false
}

func InDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
