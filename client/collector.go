package client

import (
	"context"
	"time"

	"homelens/shared"
)

func Collect(ctx context.Context, interval time.Duration, out chan<- shared.SystemInfo) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUTime []CPUTime
	var prevDiskIO []DiskIO
	var prevNetInfo []NetInfo

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			currentCPUTime, err := readCPUTime()
			if err != nil {
				return err
			}

			currentDiskIO, err := readDiskIO()
			if err != nil {
				return err
			}

			currentNetInfo, err := readNetInfo()
			if err != nil {
				return err
			}

			if prevCPUTime == nil {
				prevCPUTime = currentCPUTime
				continue
			}

			if prevDiskIO == nil {
				prevDiskIO = currentDiskIO
				continue
			}

			if prevNetInfo == nil {
				prevNetInfo = currentNetInfo
				continue
			}

			sysInfo := shared.SystemInfo{}

			sysInfo.CPU = getCPU(prevCPUTime, currentCPUTime)

			diskSpace, err := readDiskSpace("/")
			if err != nil {
				return err
			}
			sysInfo.Disk = shared.Disk{
				DiskIOUsage: calcDiskIOUsage(prevDiskIO, currentDiskIO, interval),
				DiskSpace:   diskSpace,
			}

			sysInfo.Network = calcNetUsage(prevNetInfo, currentNetInfo, interval)

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				return err
			}

			sysInfo.Temperature = readTempInfo()

			sysInfo.Containers = readDockerContainers()

			sysInfo.Processes = readTopProcesses()

			prevCPUTime = currentCPUTime
			prevDiskIO = currentDiskIO
			prevNetInfo = currentNetInfo

			out <- sysInfo
		}
	}
}
