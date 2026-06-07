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
	}

	return temps, nil
}
