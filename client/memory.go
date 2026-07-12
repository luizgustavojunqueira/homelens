package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"homelens/shared"
)

func readMemoryUsage() (shared.Memory, error) {
	stat, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return shared.Memory{}, err
	}

	var memInfo shared.Memory

	scanner := bufio.NewScanner(strings.NewReader(string(stat)))
	for scanner.Scan() {
		line := scanner.Text()
		var val uint64
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			if n, err := fmt.Sscanf(line, "MemTotal: %d kB", &val); n == 1 && err == nil {
				memInfo.Total = val * 1024
			}
		case strings.HasPrefix(line, "MemAvailable:"):
			if n, err := fmt.Sscanf(line, "MemAvailable: %d kB", &val); n == 1 && err == nil {
				memInfo.Available = val * 1024
			}
		}
	}

	if memInfo.Total > 0 && memInfo.Available <= memInfo.Total {
		memInfo.Used = memInfo.Total - memInfo.Available
	}

	return memInfo, nil
}
