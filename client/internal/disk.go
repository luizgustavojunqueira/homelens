package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type DiskSpace struct {
	Path         string  `json:"path"`
	Total        uint64  `json:"total"`
	Available    uint64  `json:"available"`
	Used         uint64  `json:"used"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskIO struct {
	Name           string `json:"name"`
	SectorsRead    uint64 `json:"sectors_read"`
	SectorsWritten uint64 `json:"sectors_written"`
	IOMs           uint64 `json:"io_ms"`
}

type DiskIOUsage struct {
	Name      string  `json:"name"`
	ReadMBps  float64 `json:"read_mbps"`
	WriteMBps float64 `json:"write_mbps"`
	IOPercent float64 `json:"io_percent"`
}

func readDiskSpace(path string) (DiskSpace, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskSpace{}, err
	}

	return DiskSpace{
		Path:         path,
		Total:        stat.Blocks * uint64(stat.Bsize),
		Available:    stat.Bavail * uint64(stat.Bsize),
		Used:         (stat.Blocks - stat.Bavail) * uint64(stat.Bsize),
		UsagePercent: float64(stat.Blocks-stat.Bavail) / float64(stat.Blocks) * 100,
	}, nil
}

func readDiskIO() ([]DiskIO, error) {
	data, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		return nil, err
	}

	var diskIOs []DiskIO
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}

		name := fields[2]
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
			continue
		}

		var d DiskIO

		d.Name = name
		fmt.Sscanf(fields[5], "%d", &d.SectorsRead)
		fmt.Sscanf(fields[9], "%d", &d.SectorsWritten)
		fmt.Sscanf(fields[12], "%d", &d.IOMs)
		diskIOs = append(diskIOs, d)
	}

	return diskIOs, nil
}

func calcDiskIOUsage(prev, current []DiskIO, interval time.Duration) []DiskIOUsage {
	secs := interval.Seconds()
	var results []DiskIOUsage

	for i, c := range current {
		p := prev[i]

		results = append(results, DiskIOUsage{
			Name:      c.Name,
			ReadMBps:  float64(c.SectorsRead-p.SectorsRead) * 512 / (1024 * 1024) / secs,
			WriteMBps: float64(c.SectorsWritten-p.SectorsWritten) * 512 / (1024 * 1024) / secs,
			IOPercent: float64(c.IOMs-p.IOMs) / (secs * 1000) * 100,
		})
	}

	return results
}

func ConvertBytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
