package client

import (
	"fmt"
	"os"
	"strings"

	"homelens/shared"
)

func readTempInfo() ([]shared.Temperature, error) {
	entries, err := os.ReadDir("/sys/class/thermal")
	if err != nil {
		return []shared.Temperature{}, err
	}

	var temps []shared.Temperature

	avgTemp := 0.0

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
		temps = append(temps, shared.Temperature{
			Zone: e.Name(),
			Temp: float64(raw) / 1000.0,
		})
		avgTemp += float64(raw) / 1000.0
	}

	if len(temps) > 0 {
		avgTemp /= float64(len(temps))
	}

	return temps, nil
}
