package internal

import "syscall"

type DiskSpace struct {
	Path         string  `json:"path"`
	Total        uint64  `json:"total"`
	Available    uint64  `json:"available"`
	Used         uint64  `json:"used"`
	UsagePercent float64 `json:"usage_percent"`
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

func ConvertBytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
