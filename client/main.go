package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"homelens_client/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	out := make(chan internal.SystemInfo)

	go func() {
		defer close(out)
		if _, err := internal.Collect(ctx, 1*time.Second, out); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Collecting CPU usage data. Press Ctrl+C to stop.")
	for sysInfo := range out {
		data, err := json.Marshal(sysInfo)
		if err != nil {
			fmt.Printf("Error marshaling system info: %v\n", err)
			continue
		}

		fmt.Printf("System Info: %s\n", string(data))

	}

	cancel()
}
