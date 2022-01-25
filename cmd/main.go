package main

import (
	"fmt"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func main() {
	fmt.Println(internal.GetDiskUsage())
}
