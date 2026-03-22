package internal

import "syscall"

type DiskSpace struct {
	Path      string
	Total     uint64
	Available uint64
}

func (d DiskSpace) UsagePercent() float64 {
	used := d.Total - d.Available

	return float64(used) / float64(d.Total) * 100
}

func readDiskSpace(path string) (DiskSpace, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskSpace{}, err
	}

	return DiskSpace{
		Path:      path,
		Total:     stat.Blocks * uint64(stat.Bsize),
		Available: stat.Bavail * uint64(stat.Bsize),
	}, nil
}

func ConvertBytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
