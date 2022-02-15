package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		//fmt.Println(value)
		return value
	}
	return fallback
}

func DateFormatID(d int64) string {
	t := time.Unix(d, 0)
	layout := "2006-01-02T15:04:05"
	return t.Format(layout)
}

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

// Adds every key and value in map to the gin context as middleware. Allows access to these variables from inside the handlers
func EnviromentMiddleware(variables map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range variables {
			if key != "" && value != "" {
				c.Set(key, value)
				c.Next()
			}
		}
	}
}

func RecoverEndpoint(c *gin.Context) {
	if err := recover(); err != nil {
		msg := "Error: [Recovered] "
		switch errType := err.(type) {
		case string:
			msg += err.(string)
		case error:
			msg += errType.Error()
		default:
		}
		fmt.Println(msg)
		c.JSON(400, gin.H{"error": msg})
	}
}
