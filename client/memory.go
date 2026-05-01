package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"homelens/shared"
)

func readMemoryUsage() (shared.MemoryUsage, error) {
	memInfo := shared.MemoryUsage{}

	stat, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return memInfo, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stat)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d kB", &memInfo.Total)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %d kB", &memInfo.Available)
		}

		if memInfo.Total != 0 && memInfo.Available != 0 {
			memInfo.Used = memInfo.Total - memInfo.Available
		}
	}

	return memInfo, nil
}

func ConvertKBToGB(kb uint64) float64 {
	return float64(kb) / (1024 * 1024)
}
