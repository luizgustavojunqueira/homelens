package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"homelens/shared"
)

type CPUTime struct {
	Name      string
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

func (c CPUTime) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal
}

func readCPUTime() ([]CPUTime, error) {
	stat, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	var cpus []CPUTime
	scanner := bufio.NewScanner(strings.NewReader(string(stat)))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu") {
			continue
		}

		var info CPUTime
		fmt.Sscanf(
			line, "%s %d %d %d %d %d %d %d %d %d %d",
			&info.Name,
			&info.User,
			&info.Nice,
			&info.System,
			&info.Idle,
			&info.IOWait,
			&info.IRQ,
			&info.SoftIRQ,
			&info.Steal,
			&info.Guest,
			&info.GuestNice,
		)
		cpus = append(cpus, info)
	}

	return cpus, nil
}

func getCPU(oldSamples []CPUTime, newSamples []CPUTime) []shared.CPU {
	var cpuInfos []shared.CPU
	for i, sample := range newSamples {

		prev := oldSamples[i]

		idle := sample.Idle - prev.Idle
		total := sample.Total() - prev.Total()

		cpuInfos = append(cpuInfos, shared.CPU{
			Name:         sample.Name,
			UsagePercent: (1.0 - float64(idle)/float64(total)) * 100,
		})
	}

	sum := 0.0

	for _, info := range cpuInfos {
		sum += info.UsagePercent
	}

	return cpuInfos
}
