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
		fmt.Println("CPU Average: ", sysInfo.CPUInfo.CPUAvg)
		for i, cpuUsage := range sysInfo.CPUInfo.UsagePercent {
			fmt.Printf("CPU %s Usage: %.2f%%\n", sysInfo.CPUInfo.Name[i], cpuUsage)
		}
		fmt.Printf("Memory: Total: %f GB, Available: %f GB\n", internal.ConvertKBToGB(sysInfo.Memory.Total), internal.ConvertKBToGB(sysInfo.Memory.Available))

		fmt.Printf("Disk: Path: %s, Total: %f GB, Available: %f GB, Usage: %.2f%%\n", sysInfo.DiskSpace.Path, internal.ConvertBytesToGB(sysInfo.DiskSpace.Total), internal.ConvertBytesToGB(sysInfo.DiskSpace.Available), sysInfo.DiskSpace.UsagePercent())

	}

	cancel()
}
