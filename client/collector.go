package client

import (
	"context"
	"log"
	"time"

	"homelens/shared"
)

func Collect(ctx context.Context, interval time.Duration, out chan<- shared.SystemInfo) error {
	procCollector := newProcessCollector()

	prevCPUTime, err := readCPUTime()
	if err != nil {
		return err
	}
	prevDiskIO, err := readDiskIO()
	if err != nil {
		return err
	}
	prevNetInfo, err := readNetInfo()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			currentCPUTime, err := readCPUTime()
			if err != nil {
				log.Printf("warning: failed to read CPU: %v", err)
				continue
			}

			currentDiskIO, err := readDiskIO()
			if err != nil {
				log.Printf("warning: failed to read disk IO: %v", err)
				continue
			}

			currentNetInfo, err := readNetInfo()
			if err != nil {
				log.Printf("warning: failed to read network info: %v", err)
				continue
			}

			sysInfo := shared.SystemInfo{}

			sysInfo.CPU = getCPU(prevCPUTime, currentCPUTime)

			diskSpace, err := readDiskSpace("/")
			if err != nil {
				log.Printf("warning: failed to read disk space: %v", err)
				continue
			}
			sysInfo.Disk = shared.Disk{
				DiskIOUsage: calcDiskIOUsage(prevDiskIO, currentDiskIO, interval),
				DiskSpace:   diskSpace,
			}

			sysInfo.Network = calcNetUsage(prevNetInfo, currentNetInfo, interval)

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				log.Printf("warning: failed to read memory: %v", err)
				continue
			}

			sysInfo.Temperature = readTempInfo()

			sysInfo.Containers = readDockerContainers()

			sysInfo.Processes = procCollector.readTopProcesses()

			prevCPUTime = currentCPUTime
			prevDiskIO = currentDiskIO
			prevNetInfo = currentNetInfo

			out <- sysInfo
		}
	}
}
