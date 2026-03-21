package main

import (
	"fmt"

	"homelens_client/internal"
)

func main() {
	usage, err := internal.GetCPUUsage()
	if err != nil {
		fmt.Printf("Error getting CPU usage: %v\n", err)
		return
	}

	for _, u := range usage {
		fmt.Println(u)
	}
}
