package client

import (
	"context"
	"time"

	"homelens/shared"
)

func Collect(ctx context.Context, interval time.Duration, out chan<- shared.SystemInfo) (shared.SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUTime []CPUTime
	var prevDiskIO []DiskIO
	var prevNetInfo []NetInfo

	for {
		select {
		case <-ctx.Done():
			return shared.SystemInfo{}, ctx.Err()

		case <-ticker.C:
			currentCPUTime, err := readCPUTime()
			if err != nil {
				return shared.SystemInfo{}, err
			}

			currentDiskIO, err := readDiskIO()
			if err != nil {
				return shared.SystemInfo{}, err
			}

			currentNetInfo, err := readNetInfo()
			if err != nil {
				return shared.SystemInfo{}, err
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

			sysInfo.CPUUsage = getCPUUsage(prevCPUTime, currentCPUTime)
			sysInfo.DiskIOUsage = calcDiskIOUsage(prevDiskIO, currentDiskIO, interval)
			sysInfo.NetUsage = calcNetUsage(prevNetInfo, currentNetInfo, interval)

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				return shared.SystemInfo{}, err
			}

			sysInfo.DiskSpace, err = readDiskSpace("/")
			if err != nil {
				return shared.SystemInfo{}, err
			}

			sysInfo.Temperature, err = readTempInfo()
			if err != nil {
				return shared.SystemInfo{}, err
			}

			prevCPUTime = currentCPUTime
			prevDiskIO = currentDiskIO
			prevNetInfo = currentNetInfo

			out <- sysInfo
		}
	}
}
