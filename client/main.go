package main

import (
	"context"
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
		fmt.Println("CPU Average: ", sysInfo.CPUAvg)
		for _, cpuUsage := range sysInfo.CPUUsages {
			fmt.Println("CPU Usage:", cpuUsage)
		}
	}

	cancel()
}
