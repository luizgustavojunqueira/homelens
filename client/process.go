package client

import (
	"sort"
	"time"

	"homelens/shared"

	"github.com/shirou/gopsutil/v4/process"
)

const topProcessCount = 50

type processCollector struct {
	cache    map[int32]*process.Process
	lastCPU  map[int32]float64
	lastTime time.Time
}

func newProcessCollector() *processCollector {
	return &processCollector{
		cache:   make(map[int32]*process.Process),
		lastCPU: make(map[int32]float64),
	}
}

func (pc *processCollector) readTopProcesses() []shared.Process {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	now := time.Now()
	currentPids := make(map[int32]bool, len(procs))

	type procEntry struct {
		proc *process.Process
		cpu  float64
	}

	entries := make([]procEntry, 0, len(procs))
	for _, p := range procs {
		currentPids[p.Pid] = true

		cached, exists := pc.cache[p.Pid]
		if !exists {
			cached = p
			pc.cache[p.Pid] = cached
		}

		times, err := cached.Times()
		if err != nil {
			continue
		}

		totalCPU := times.User + times.System
		elapsed := now.Sub(pc.lastTime).Seconds()

		var cpuPct float64
		if elapsed > 0 {
			prev := pc.lastCPU[p.Pid]
			cpuPct = (totalCPU - prev) / elapsed * 100
		}
		pc.lastCPU[p.Pid] = totalCPU

		entries = append(entries, procEntry{proc: cached, cpu: cpuPct})
	}

	pc.lastTime = now

	for pid := range pc.cache {
		if !currentPids[pid] {
			delete(pc.cache, pid)
			delete(pc.lastCPU, pid)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].cpu > entries[j].cpu
	})
	if len(entries) > topProcessCount {
		entries = entries[:topProcessCount]
	}

	now64 := now.UnixMilli()
	result := make([]shared.Process, 0, len(entries))
	for _, e := range entries {
		p := e.proc

		var rss uint64
		if memInfo, err := p.MemoryInfo(); err == nil {
			rss = memInfo.RSS
		}

		var memPct float32
		if m, err := p.MemoryPercent(); err == nil {
			memPct = m
		}

		name, _ := p.Name()
		user, _ := p.Username()
		createTime, _ := p.CreateTime()
		cmdline, _ := p.Cmdline()

		result = append(result, shared.Process{
			PID:     int(p.Pid),
			User:    user,
			CPU:     e.cpu,
			Memory:  float64(memPct),
			RSS:     rss,
			Uptime:  int((now64 - createTime) / 1000),
			Name:    name,
			Cmdline: cmdline,
		})
	}

	return result
}
