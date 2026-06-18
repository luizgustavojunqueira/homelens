package client

import (
	"sort"
	"time"

	"homelens/shared"

	"github.com/shirou/gopsutil/v4/process"
)

var processCache = make(map[int32]*process.Process)

func readTopProcesses() []shared.Process {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	var processList []shared.Process
	now := time.Now().UnixMilli()

	currentPids := make(map[int32]bool)

	for _, p := range procs {
		currentPids[p.Pid] = true

		cachedProc, exists := processCache[p.Pid]
		if !exists {
			cachedProc = p
			processCache[p.Pid] = cachedProc
		}

		cpu, err := cachedProc.Percent(0)
		if err != nil {
			continue
		}

		mem, err := cachedProc.MemoryPercent()
		if err != nil {
			continue
		}

		memInfo, err := cachedProc.MemoryInfo()
		var rss uint64
		if err == nil {
			rss = memInfo.RSS
		}

		name, _ := cachedProc.Name()
		user, _ := cachedProc.Username()
		createTime, _ := cachedProc.CreateTime()
		cmdline, _ := cachedProc.Cmdline()

		uptimeSeconds := int((now - createTime) / 1000)

		processList = append(processList, shared.Process{
			PID:     int(cachedProc.Pid),
			User:    user,
			CPU:     cpu,
			Memory:  float64(mem),
			RSS:     rss,
			Uptime:  uptimeSeconds,
			Name:    name,
			Cmdline: cmdline,
		})
	}

	for pid := range processCache {
		if !currentPids[pid] {
			delete(processCache, pid)
		}
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].CPU > processList[j].CPU
	})

	if len(processList) > 10 {
		processList = processList[:10]
	}

	return processList
}
