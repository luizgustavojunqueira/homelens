// Package internal
package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func readCPUInfo() ([]CPUTime, error) {
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
		fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d",
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

func getCPUUsage(oldSamples []CPUTime, newSamples []CPUTime) []float64 {
	var results []float64
	for i, sample := range newSamples {

		totalDelta := sample.Total() - oldSamples[i].Total()
		idleDelta := sample.Idle - oldSamples[i].Idle

		if totalDelta > 0 {

			usage := float64(totalDelta-idleDelta) / float64(totalDelta) * 100
			results = append(results, usage)
		}
	}

	return results
}

func getCPUAvg(usages []float64) float64 {
	if len(usages) == 0 {
		return 0
	}

	var sum float64
	for _, usage := range usages {
		sum += usage
	}
	return sum / float64(len(usages))
}
