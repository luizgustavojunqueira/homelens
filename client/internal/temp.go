package internal

import (
	"fmt"
	"os"
	"strings"
)

type TempInfo struct {
	Zone string  `json:"zone"`
	Temp float64 `json:"temp_c"`
}

func readTempInfo() ([]TempInfo, error) {
	entries, err := os.ReadDir("/sys/class/thermal")
	if err != nil {
		return nil, err
	}

	var temps []TempInfo

	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "thermal_zone") {
			continue
		}

		data, err := os.ReadFile("/sys/class/thermal/" + e.Name() + "/temp")
		if err != nil {
			continue
		}

		var raw uint64
		fmt.Sscanf(string(data), "%d", &raw)
		temps = append(temps, TempInfo{
			Zone: e.Name(),
			Temp: float64(raw) / 1000.0,
		})

	}
	return temps, nil
}
