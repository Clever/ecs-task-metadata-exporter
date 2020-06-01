package metrics

import (
	"testing"

	"github.com/docker/engine/api/types"
)

func Test_CpuUsage(t *testing.T) {
	stats := types.StatsJSON{Stats: types.Stats{
		CPUStats: types.CPUStats{
			SystemUsage: 100,
			CPUUsage: types.CPUUsage{
				TotalUsage: 30,
			},
		},
		PreCPUStats: types.CPUStats{
			SystemUsage: 80,
			CPUUsage: types.CPUUsage{
				TotalUsage: 20,
			},
		},
	}}
	usage := cpuUsage(stats)
	// out of the 20 total system CPU seconds between previous and current, 10 were used by the container; thus we expect 10.0/20.0 = 0.5
	if usage != 0.5 {
		t.Fatalf("Got CPU usage %f; expecting %f", usage, 0.5)
	}

}
