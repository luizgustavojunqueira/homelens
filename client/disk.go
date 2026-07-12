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
	Name           string
	SectorsRead    uint64
	SectorsWritten uint64
	IOMs           uint64
}

func readDiskSpace(path string) (shared.DiskSpace, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return shared.DiskSpace{}, err
	}

	var usagePercent float64
	if stat.Blocks > 0 {
		usagePercent = float64(stat.Blocks-stat.Bavail) / float64(stat.Blocks) * 100
	}

	return shared.DiskSpace{
		Path:         path,
		Total:        stat.Blocks * uint64(stat.Bsize),
		Available:    stat.Bavail * uint64(stat.Bsize),
		Used:         (stat.Blocks - stat.Bavail) * uint64(stat.Bsize),
		UsagePercent: usagePercent,
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
	if secs == 0 {
		return nil
	}

	prevByName := make(map[string]DiskIO, len(prev))
	for _, p := range prev {
		prevByName[p.Name] = p
	}

	var results []shared.DiskIOUsage
	for _, c := range current {
		p, ok := prevByName[c.Name]
		if !ok {
			continue
		}

		results = append(results, shared.DiskIOUsage{
			Name:      c.Name,
			ReadMBps:  float64(c.SectorsRead-p.SectorsRead) * sectorToMB / secs,
			WriteMBps: float64(c.SectorsWritten-p.SectorsWritten) * sectorToMB / secs,
			IOPercent: float64(c.IOMs-p.IOMs) / (secs * 1000) * 100,
		})
	}

	return results
}
