package data

import (
	"net/http"
	"testing"
	"time"

	"github.com/docker/engine/api/types"
	"github.com/go-openapi/swag"
	"github.com/google/go-cmp/cmp"
)

func TestMetadataEndpointSource(t *testing.T) {
	handler := ConstantMetadataEndpointHandler(
		SampleTaskMetadata, SampleTaskStats,
	)
	server := http.Server{
		Addr:    "localhost:8912",
		Handler: handler,
	}
	go server.ListenAndServe()
	defer server.Close()

	m := NewMetadataEndpointSource("http://localhost:8912")

	if meta, err := m.Metadata(); err != nil {
		t.Fatalf("got error from Metadata(): %v", err)
	} else {
		expectedMetadata := TaskMetadata{
			Cluster:            "default",
			TaskARN:            "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
			Family:             "nginx",
			Revision:           "5",
			DesiredStatus:      "RUNNING",
			KnownStatus:        "RUNNING",
			Limits:             nil,
			PullStartedAt:      mustParseTime(time.RFC3339, "2018-02-01T20:55:09.372495529Z"),
			PullStoppedAt:      mustParseTime(time.RFC3339, "2018-02-01T20:55:10.552018345Z"),
			AvailabilityZone:   swag.String("us-east-2b"),
			ExecutionStoppedAt: time.Time{},
			Containers: []ContainerMetadata{
				ContainerMetadata{
					DockerID:   "731a0d6a3b4210e2448339bc7015aaa79bfe4fa256384f4102db86ef94cbbc4c",
					Name:       "~internal~ecs~pause",
					DockerName: "ecs-nginx-5-internalecspause-acc699c0cbf2d6d11700",
					Image:      "amazon/amazon-ecs-pause:0.1.0",
					ImageID:    "",
					Labels: map[string]string{"com.amazonaws.ecs.cluster": "default",
						"com.amazonaws.ecs.container-name":          "~internal~ecs~pause",
						"com.amazonaws.ecs.task-arn":                "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
						"com.amazonaws.ecs.task-definition-family":  "nginx",
						"com.amazonaws.ecs.task-definition-version": "5"},
					DesiredStatus: "RESOURCES_PROVISIONED",
					KnownStatus:   "RESOURCES_PROVISIONED",
					Limits:        &Limits{CPU: swag.Float64(0.0), Memory: swag.Uint64(0)},
					CreatedAt:     mustParseTime(time.RFC3339, "2018-02-01T20:55:08.366329616Z"),
					StartedAt:     mustParseTime(time.RFC3339, "2018-02-01T20:55:09.058354915Z"),
					Type:          "CNI_PAUSE",
				},
				ContainerMetadata{
					DockerID:   "43481a6ce4842eec8fe72fc28500c6b52edcc0917f105b83379f88cac1ff3946",
					Name:       "nginx-curl",
					DockerName: "ecs-nginx-5-nginx-curl-ccccb9f49db0dfe0d901",
					Image:      "nrdlngr/nginx-curl",
					ImageID:    "sha256:2e00ae64383cfc865ba0a2ba37f61b50a120d2d9378559dcd458dc0de47bc165",
					Labels: map[string]string{
						"com.amazonaws.ecs.cluster":                 "default",
						"com.amazonaws.ecs.container-name":          "nginx-curl",
						"com.amazonaws.ecs.task-arn":                "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
						"com.amazonaws.ecs.task-definition-family":  "nginx",
						"com.amazonaws.ecs.task-definition-version": "5",
					},
					DesiredStatus: "RUNNING",
					KnownStatus:   "RUNNING",
					Limits:        &Limits{CPU: swag.Float64(512), Memory: swag.Uint64(512)},
					CreatedAt:     mustParseTime(time.RFC3339, "2018-02-01T20:55:10.554941919Z"),
					StartedAt:     mustParseTime(time.RFC3339, "2018-02-01T20:55:11.064236631Z"),
					Type:          "NORMAL",
				},
			},
		}
		if diff := cmp.Diff(expectedMetadata, meta); diff != "" {
			t.Fatalf("metadata mismatch (-want +got):\n%s", diff)
		}
	}

	if stats, err := m.Stats(); err != nil {
		t.Fatalf("got error from stats(): %v", err)
	} else {
		expectedStats := map[string]types.StatsJSON{
			"43481a6ce4842eec8fe72fc28500c6b52edcc0917f105b83379f88cac1ff3946": types.StatsJSON{
				Stats: types.Stats{
					Read:    mustParseTime(time.RFC3339, "2020-04-06T16:12:01.090148907Z"),
					PreRead: mustParseTime(time.RFC3339, "2020-04-06T16:11:56.083890951Z"),
					BlkioStats: types.BlkioStats{
						IoServiceBytesRecursive: []types.BlkioStatEntry{
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Read", Value: 3452928},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Write", Value: 0},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Sync", Value: 3452928},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Async", Value: 0},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Total", Value: 3452928}},
						IoServicedRecursive: []types.BlkioStatEntry{
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Read", Value: 118},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Write", Value: 0},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Sync", Value: 118},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Async", Value: 0},
							types.BlkioStatEntry{Major: 202, Minor: 26368, Op: "Total", Value: 118}},
						IoQueuedRecursive:      []types.BlkioStatEntry{},
						IoServiceTimeRecursive: []types.BlkioStatEntry{},
						IoWaitTimeRecursive:    []types.BlkioStatEntry{},
						IoMergedRecursive:      []types.BlkioStatEntry{},
						IoTimeRecursive:        []types.BlkioStatEntry{},
						SectorsRecursive:       []types.BlkioStatEntry{}},
					NumProcs:     0,
					StorageStats: types.StorageStats{ReadCountNormalized: 0, ReadSizeBytes: 0, WriteCountNormalized: 0, WriteSizeBytes: 0x0},
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage:        410557100,
							PercpuUsage:       []uint64{410557100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0},
							UsageInKernelmode: 10000000,
							UsageInUsermode:   250000000,
						},
						SystemUsage:    0,
						ThrottlingData: types.ThrottlingData{Periods: 0, ThrottledPeriods: 0, ThrottledTime: 0}},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage:        0,
							PercpuUsage:       []uint64(nil),
							UsageInKernelmode: 0,
							UsageInUsermode:   0,
						},
						SystemUsage:    0,
						ThrottlingData: types.ThrottlingData{Periods: 0, ThrottledPeriods: 0, ThrottledTime: 0}},
					MemoryStats: types.MemoryStats{
						Usage:    4390912,
						MaxUsage: 6488064,
						Stats: map[string]uint64{
							"active_anon":               278528,
							"active_file":               344064,
							"cache":                     3452928,
							"dirty":                     0,
							"hierarchical_memory_limit": 536870912,
							"hierarchical_memsw_limit":  9223372036854772000,
							"inactive_anon":             0,
							"inactive_file":             3108864,
							"mapped_file":               2412544,
							"pgfault":                   2800,
							"pgmajfault":                28,
							"pgpgin":                    3144,
							"pgpgout":                   2233,
							"rss":                       278528,
							"rss_huge":                  0,
							"total_active_anon":         278528,
							"total_active_file":         344064,
							"total_cache":               3452928,
							"total_dirty":               0,
							"total_inactive_anon":       0,
							"total_inactive_file":       3108864,
							"total_mapped_file":         2412544,
							"total_pgfault":             2800,
							"total_pgmajfault":          28,
							"total_pgpgin":              3144,
							"total_pgpgout":             2233,
							"total_rss":                 278528,
							"total_rss_huge":            0,
							"total_unevictable":         0,
							"total_writeback":           0,
							"unevictable":               0,
							"writeback":                 0,
						},
						Failcnt:           0,
						Limit:             9223372036854772000,
						Commit:            0,
						CommitPeak:        0,
						PrivateWorkingSet: 0x0},
				},
				Name: "query-metadata",
				ID:   "1823e1f6-7248-43c3-bed6-eea1fa7501a5query-metadata",
				Networks: map[string]types.NetworkStats{
					"eth1": types.NetworkStats{
						RxBytes:   564655295,
						RxPackets: 384960,
						RxErrors:  0,
						RxDropped: 0,
						TxBytes:   3043269,
						TxPackets: 54355,
						TxErrors:  0,
						TxDropped: 0,
					},
				},
			},
		}
		if diff := cmp.Diff(expectedStats, stats); diff != "" {
			t.Fatalf("stats mismatch (-want +got):\n%s", diff)
		}
	}
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
