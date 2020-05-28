package metrics

import (
	"encoding/json"
	"testing"

	"github.com/docker/engine/api/types"
)

// testStats pulled from https://docs.docker.com/engine/api/v1.30/#operation/ContainerExport
const testStats string = `{
  "read": "2015-01-08T22:57:31.547920715Z",
  "pids_stats": {
    "current": 3
  },
  "networks": {
    "eth0": {
      "rx_bytes": 5338,
      "rx_dropped": 0,
      "rx_errors": 0,
      "rx_packets": 36,
      "tx_bytes": 648,
      "tx_dropped": 0,
      "tx_errors": 0,
      "tx_packets": 8
    },
    "eth5": {
      "rx_bytes": 4641,
      "rx_dropped": 0,
      "rx_errors": 0,
      "rx_packets": 26,
      "tx_bytes": 690,
      "tx_dropped": 0,
      "tx_errors": 0,
      "tx_packets": 9
    }
  },
  "memory_stats": {
    "stats": {
      "total_pgmajfault": 0,
      "cache": 0,
      "mapped_file": 0,
      "total_inactive_file": 0,
      "pgpgout": 414,
      "rss": 6537216,
      "total_mapped_file": 0,
      "writeback": 0,
      "unevictable": 0,
      "pgpgin": 477,
      "total_unevictable": 0,
      "pgmajfault": 0,
      "total_rss": 6537216,
      "total_rss_huge": 6291456,
      "total_writeback": 0,
      "total_inactive_anon": 0,
      "rss_huge": 6291456,
      "hierarchical_memory_limit": 67108864,
      "total_pgfault": 964,
      "total_active_file": 0,
      "active_anon": 6537216,
      "total_active_anon": 6537216,
      "total_pgpgout": 414,
      "total_cache": 0,
      "inactive_anon": 0,
      "active_file": 0,
      "pgfault": 964,
      "inactive_file": 0,
      "total_pgpgin": 477
    },
    "max_usage": 6651904,
    "usage": 6537216,
    "failcnt": 0,
    "limit": 67108864
  },
  "blkio_stats": {},
  "cpu_stats": {
    "cpu_usage": {
      "percpu_usage": [
        10000000,
        20000000,
        30000000,
        40000000
      ],
      "usage_in_usermode": 80000000,
      "total_usage": 100000000,
      "usage_in_kernelmode": 10000000
    },
    "system_cpu_usage": 1000000000,
    "online_cpus": 4,
    "throttling_data": {
      "periods": 0,
      "throttled_periods": 0,
      "throttled_time": 0
    }
  },
  "precpu_stats": {
    "cpu_usage": {
      "percpu_usage": [
        9000000,
        19000000,
        29000000,
        39000000
      ],
      "usage_in_usermode": 78000000,
      "total_usage": 96000000,
      "usage_in_kernelmode": 9500000
    },
    "system_cpu_usage": 950000000,
    "online_cpus": 4,
    "throttling_data": {
      "periods": 0,
      "throttled_periods": 0,
      "throttled_time": 0
    }
  }
}`

func Test_Unmarshal(t *testing.T) {
	var stats types.StatsJSON
	err := json.Unmarshal([]byte(testStats), &stats)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

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
