package client

import (
	"fmt"
	"os"
	"strings"

	"homelens/shared"
)

func readTempInfo() (shared.Temp, error) {
	entries, err := os.ReadDir("/sys/class/thermal")
	if err != nil {
		return shared.Temp{}, err
	}

	var temps []shared.TempInfo

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
		temps = append(temps, shared.TempInfo{
			Zone: e.Name(),
			Temp: float64(raw) / 1000.0,
		})
		avgTemp += float64(raw) / 1000.0
	}

	if len(temps) > 0 {
		avgTemp /= float64(len(temps))
	}

	temp := shared.Temp{
		TempAvg:  avgTemp,
		TempInfo: temps,
	}

	return temp, nil
}
