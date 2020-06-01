package metrics

import (
	"fmt"

	"github.com/docker/engine/api/types"
	"github.com/prometheus/client_golang/prometheus"
)

// Prefix will be prepended to all metrics so that they can be distinguished from similar metrics coming from other sources
const Prefix = "ecs_container_"

// MetricConfig is a specification of a single metric that can be extracted from the docker stats
type MetricConfig struct {
	Name    string
	Help    string
	Type    prometheus.ValueType
	ValueFn func(types.StatsJSON) float64
}

// DefaultMetrics is a slice of default metrics to use
var DefaultMetrics = []MetricConfig{
	{
		Name:    "mem_usage_bytes",
		Help:    "Current memory usage",
		Type:    prometheus.GaugeValue,
		ValueFn: func(s types.StatsJSON) float64 { return float64(s.MemoryStats.Usage) },
	},
	{
		Name:    "mem_max_usage_bytes",
		Help:    "Maximum memory usage",
		Type:    prometheus.GaugeValue,
		ValueFn: func(s types.StatsJSON) float64 { return float64(s.MemoryStats.MaxUsage) },
	},
	{
		Name:    "mem_limit_bytes",
		Help:    "Memory limit",
		Type:    prometheus.GaugeValue,
		ValueFn: func(s types.StatsJSON) float64 { return float64(s.MemoryStats.Limit) },
	},
	{
		Name:    "cpu_usage",
		Help:    "CPU usage from 0 to 1 of the container as a ratio of total CPU usage on the host",
		Type:    prometheus.GaugeValue,
		ValueFn: cpuUsage,
	},
}

// StatsToMetrics converts docker's StatsJSON into constant Prometheus metrics
func StatsToMetrics(stats types.StatsJSON, configs []MetricConfig, labels prometheus.Labels) ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}
	for _, config := range configs {
		m, err := prometheus.NewConstMetric(
			prometheus.NewDesc(Prefix+config.Name, config.Help, nil /* variable labels */, labels),
			config.Type,
			config.ValueFn(stats),
		)
		if err != nil {
			// NewConstMetric can fail if variable labels are the wrong length (not applicable here) or Desc is invalid (shouldn't come up)
			return nil, fmt.Errorf("prometheus.NewConstMetric(%s): %v", config.Name, err)
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// cpuUsage returns the fraction from 0 to 1 of CPU time being used by the container.
func cpuUsage(stats types.StatsJSON) float64 {
	// On linux systems, docker reports CPU usage as nanoseconds of CPU time used since the container started. It also reports total system CPU nanoseconds.
	// Ref: https://github.com/moby/moby/blob/master/api/types/stats.go
	// When asking for stats, it gives two sets of those nanosecond totals, a newer one under stats.CPUStats and an older one under stats.PreCPUStats.
	// (Pre is for previous. Where the previous data comes from I'm not sure.)
	// Thus, we can calculate how much CPU each container is using as a fraction of the total CPU used on the host.
	// All of these numbers are totals over all the CPUs on the system, so the number of CPUs doesn't come into play
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	containerDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)

	if systemDelta > 0.0 && containerDelta > 0.0 {
		return containerDelta / systemDelta
	}
	return 0.0
}
