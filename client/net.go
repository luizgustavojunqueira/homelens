package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"homelens/shared"
)

type NetInfo struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

func readNetInfo() ([]NetInfo, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	var nets []NetInfo

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.Split(line, ":")
		name := strings.TrimSpace(parts[0])
		if name == "lo" {
			continue
		}

		if name[0] != 'e' && name[0] != 'w' {
			continue
		}

		fields := strings.Fields(parts[1])
		var info NetInfo
		info.Name = name
		fmt.Sscanf(fields[0], "%d", &info.RxBytes)
		fmt.Sscanf(fields[8], "%d", &info.TxBytes)
		nets = append(nets, info)
	}

	return nets, nil
}

func calcNetUsage(prev, curr []NetInfo, interval time.Duration) []shared.Network {
	secs := interval.Seconds()

	var netUsages []shared.Network

	for i, c := range curr {
		p := prev[i]

		netUsages = append(netUsages, shared.Network{
			Name:  c.Name,
			RxBps: float64(c.RxBytes-p.RxBytes) / secs,
			TxBps: float64(c.TxBytes-p.TxBytes) / secs,
		})
	}

	return netUsages
}
