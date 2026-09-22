//go:build !linux

package system

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v4/sensors"
)

func getSensorsDetails() ([]*DetailTable, error) {
	t := newDetailTable("Sensors", "Sensor", "Value")
	temps, _ := sensors.SensorsTemperatures()
	sort.Slice(temps, func(i, j int) bool { return temps[i].SensorKey < temps[j].SensorKey })
	for _, s := range temps {
		if s.Temperature <= 0 {
			continue
		}
		t.Add(s.SensorKey, fmt.Sprintf("%.1f °C", s.Temperature))
	}
	return []*DetailTable{t}, nil
}
