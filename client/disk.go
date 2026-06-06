package client

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	"homelens/shared"
)

const sectorToMB = 512.0 / (1024.0 * 1024.0)

type DiskIO struct {
	Name           string `json:"name"`
	SectorsRead    uint64 `json:"sectors_read"`
	SectorsWritten uint64 `json:"sectors_written"`
	IOMs           uint64 `json:"io_ms"`
}

func readDiskSpace(path string) (shared.DiskSpace, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return shared.DiskSpace{}, err
	}

	return shared.DiskSpace{
		Path:         path,
		Total:        stat.Blocks * uint64(stat.Bsize),
		Available:    stat.Bavail * uint64(stat.Bsize),
		Used:         (stat.Blocks - stat.Bavail) * uint64(stat.Bsize),
		UsagePercent: float64(stat.Blocks-stat.Bavail) / float64(stat.Blocks) * 100,
	}, nil
}

var wholeDiskRe = regexp.MustCompile(`^(nvme\d+n\d+|sd[a-z]+)$`)

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
		if !wholeDiskRe.MatchString(name) {
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

func calcDiskIOUsage(prev, current []DiskIO, interval time.Duration) []shared.DiskIOUsage {
	secs := interval.Seconds()
	var results []shared.DiskIOUsage

	for i, c := range current {
		p := prev[i]

		results = append(results, shared.DiskIOUsage{
			Name:      c.Name,
			ReadMBps:  float64(c.SectorsRead-p.SectorsRead) * sectorToMB / secs,
			WriteMBps: float64(c.SectorsWritten-p.SectorsWritten) * sectorToMB / secs,
			IOPercent: float64(c.IOMs-p.IOMs) / (secs * 1000) * 100,
		})
	}

	return results
}

func ConvertBytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
