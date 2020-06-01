package data

// SampleTaskMetadata is an example of the response from ECS_TASK_METADATA_URI/task
// This was taken from https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v3.html
var SampleTaskMetadata = []byte(`{
  "Cluster": "default",
  "TaskARN": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
  "Family": "nginx",
  "Revision": "5",
  "DesiredStatus": "RUNNING",
  "KnownStatus": "RUNNING",
  "Containers": [
    {
      "DockerId": "731a0d6a3b4210e2448339bc7015aaa79bfe4fa256384f4102db86ef94cbbc4c",
      "Name": "~internal~ecs~pause",
      "DockerName": "ecs-nginx-5-internalecspause-acc699c0cbf2d6d11700",
      "Image": "amazon/amazon-ecs-pause:0.1.0",
      "ImageID": "",
      "Labels": {
        "com.amazonaws.ecs.cluster": "default",
        "com.amazonaws.ecs.container-name": "~internal~ecs~pause",
        "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
        "com.amazonaws.ecs.task-definition-family": "nginx",
        "com.amazonaws.ecs.task-definition-version": "5"
      },
      "DesiredStatus": "RESOURCES_PROVISIONED",
      "KnownStatus": "RESOURCES_PROVISIONED",
      "Limits": {
        "CPU": 0,
        "Memory": 0
      },
      "CreatedAt": "2018-02-01T20:55:08.366329616Z",
      "StartedAt": "2018-02-01T20:55:09.058354915Z",
      "Type": "CNI_PAUSE",
      "Networks": [
        {
          "NetworkMode": "awsvpc",
          "IPv4Addresses": [
            "10.0.2.106"
          ]
        }
      ]
    },
    {
      "DockerId": "43481a6ce4842eec8fe72fc28500c6b52edcc0917f105b83379f88cac1ff3946",
      "Name": "nginx-curl",
      "DockerName": "ecs-nginx-5-nginx-curl-ccccb9f49db0dfe0d901",
      "Image": "nrdlngr/nginx-curl",
      "ImageID": "sha256:2e00ae64383cfc865ba0a2ba37f61b50a120d2d9378559dcd458dc0de47bc165",
      "Labels": {
        "com.amazonaws.ecs.cluster": "default",
        "com.amazonaws.ecs.container-name": "nginx-curl",
        "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
        "com.amazonaws.ecs.task-definition-family": "nginx",
        "com.amazonaws.ecs.task-definition-version": "5"
      },
      "DesiredStatus": "RUNNING",
      "KnownStatus": "RUNNING",
      "Limits": {
        "CPU": 512,
        "Memory": 512
      },
      "CreatedAt": "2018-02-01T20:55:10.554941919Z",
      "StartedAt": "2018-02-01T20:55:11.064236631Z",
      "Type": "NORMAL",
      "Networks": [
        {
          "NetworkMode": "awsvpc",
          "IPv4Addresses": [
            "10.0.2.106"
          ]
        }
      ]
    }
  ],
  "PullStartedAt": "2018-02-01T20:55:09.372495529Z",
  "PullStoppedAt": "2018-02-01T20:55:10.552018345Z",
  "AvailabilityZone": "us-east-2b"
}
`)

// SampleTaskStats is an example of the response from ECS_TASK_METADATA_URI/task/stats
// This was modified from https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html.
// The only modification was to replace the container ID with one that matches SampleTaskMetadata
var SampleTaskStats = []byte(`{
    "43481a6ce4842eec8fe72fc28500c6b52edcc0917f105b83379f88cac1ff3946": {
        "read": "2020-04-06T16:12:01.090148907Z",
        "preread": "2020-04-06T16:11:56.083890951Z",
        "pids_stats": {

        },
        "blkio_stats": {
            "io_service_bytes_recursive": [
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Read",
                    "value": 3452928
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Write",
                    "value": 0
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Sync",
                    "value": 3452928
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Async",
                    "value": 0
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Total",
                    "value": 3452928
                }
            ],
            "io_serviced_recursive": [
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Read",
                    "value": 118
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Write",
                    "value": 0
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Sync",
                    "value": 118
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Async",
                    "value": 0
                },
                {
                    "major": 202,
                    "minor": 26368,
                    "op": "Total",
                    "value": 118
                }
            ],
            "io_queue_recursive": [

            ],
            "io_service_time_recursive": [

            ],
            "io_wait_time_recursive": [

            ],
            "io_merged_recursive": [

            ],
            "io_time_recursive": [

            ],
            "sectors_recursive": [

            ]
        },
        "num_procs": 0,
        "storage_stats": {

        },
        "cpu_stats": {
            "cpu_usage": {
                "total_usage": 410557100,
                "percpu_usage": [
                    410557100,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "usage_in_kernelmode": 10000000,
                "usage_in_usermode": 250000000
            },
            "throttling_data": {
                "periods": 0,
                "throttled_periods": 0,
                "throttled_time": 0
            }
        },
        "precpu_stats": {
            "cpu_usage": {
                "total_usage": 0,
                "usage_in_kernelmode": 0,
                "usage_in_usermode": 0
            },
            "throttling_data": {
                "periods": 0,
                "throttled_periods": 0,
                "throttled_time": 0
            }
        },
        "memory_stats": {
            "usage": 4390912,
            "max_usage": 6488064,
            "stats": {
                "active_anon": 278528,
                "active_file": 344064,
                "cache": 3452928,
                "dirty": 0,
                "hierarchical_memory_limit": 536870912,
                "hierarchical_memsw_limit": 9223372036854772000,
                "inactive_anon": 0,
                "inactive_file": 3108864,
                "mapped_file": 2412544,
                "pgfault": 2800,
                "pgmajfault": 28,
                "pgpgin": 3144,
                "pgpgout": 2233,
                "rss": 278528,
                "rss_huge": 0,
                "total_active_anon": 278528,
                "total_active_file": 344064,
                "total_cache": 3452928,
                "total_dirty": 0,
                "total_inactive_anon": 0,
                "total_inactive_file": 3108864,
                "total_mapped_file": 2412544,
                "total_pgfault": 2800,
                "total_pgmajfault": 28,
                "total_pgpgin": 3144,
                "total_pgpgout": 2233,
                "total_rss": 278528,
                "total_rss_huge": 0,
                "total_unevictable": 0,
                "total_writeback": 0,
                "unevictable": 0,
                "writeback": 0
            },
            "limit": 9223372036854772000
        },
        "name": "query-metadata",
        "id": "1823e1f6-7248-43c3-bed6-eea1fa7501a5query-metadata",
        "networks": {
            "eth1": {
                "rx_bytes": 564655295,
                "rx_packets": 384960,
                "rx_errors": 0,
                "rx_dropped": 0,
                "tx_bytes": 3043269,
                "tx_packets": 54355,
                "tx_errors": 0,
                "tx_dropped": 0
            }
        }
    }
}`)
