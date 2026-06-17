package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"homelens/shared"
)

type TempStrategy interface {
	Read() []shared.Temperature
}

func readTempInfo() []shared.Temperature {
	strategies := []TempStrategy{
		thermalZoneStrategy{},
		hwmonStrategy{},
	}

	var allTemps []shared.Temperature
	for _, strategy := range strategies {
		allTemps = append(allTemps, strategy.Read()...)
	}
	return allTemps
}

type thermalZoneStrategy struct{}

func (t thermalZoneStrategy) Read() []shared.Temperature {
	entries, err := os.ReadDir("/sys/class/thermal")
	if err != nil {
		return nil
	}

	var temps []shared.Temperature

	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "thermal_zone") {
			continue
		}

		basePath := filepath.Join("/sys/class/thermal", e.Name())
		zoneName := e.Name()

		if typeData, err := os.ReadFile(filepath.Join(basePath, "type")); err == nil {
			zoneName = strings.TrimSpace(string(typeData))
		}

		data, err := os.ReadFile(filepath.Join(basePath, "temp"))
		if err != nil {
			continue
		}

		var raw uint64
		fmt.Sscanf(string(data), "%d", &raw)
		temps = append(temps, shared.Temperature{
			Zone: zoneName,
			Temp: float64(raw) / 1000.0,
		})
	}

	return temps
}

type hwmonStrategy struct{}

func (h hwmonStrategy) Read() []shared.Temperature {
	entries, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return nil
	}

	var temps []shared.Temperature

	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "hwmon") {
			continue
		}

		componentName := e.Name()

		basePath := filepath.Join("/sys/class/hwmon", componentName)

		if nameData, err := os.ReadFile(filepath.Join(basePath, "name")); err == nil {
			componentName = strings.TrimSpace(string(nameData))
		}

		inputs, err := filepath.Glob(filepath.Join(basePath, "temp*_input"))
		if err != nil || len(inputs) == 0 {
			continue
		}

		for _, inputPath := range inputs {

			baseName := filepath.Base(inputPath)
			prefix := strings.TrimSuffix(baseName, "_input")
			label := componentName + "_" + prefix
			if labelData, err := os.ReadFile(filepath.Join(basePath, prefix+"label")); err == nil {
				label = componentName + "_" + strings.TrimSpace(string(labelData))
			}

			data, err := os.ReadFile(inputPath)
			if err != nil {
				continue
			}

			var raw uint64
			fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &raw)
			if err != nil {
				continue
			}
			temps = append(temps, shared.Temperature{
				Zone: label,
				Temp: float64(raw) / 1000.0,
			})

		}
	}
	return temps
}
