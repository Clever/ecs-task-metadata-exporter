# ecs-task-metadata-exporter

A Prometheus exporter for monitoring ECS containers using the ECS task metadata endpoint. `ecs-task-metadata-exporter` can be run as a sidecar by including it as an additional container in any task definition for deployment onto ECS. It targets tasks run on EC2 instances running ECS agent version `>= v1.21.0` or tasks running on Fargate with platform version `>= 1.3.0`. Windows services are not supported.

## Metrics

As of now, the set of metrics is fairly minimal. Once the system is validated as functioning, more will be added, including network and disk usage metrics.
For each container, the following are exported, distinguished by the `ContainerName` label.
- `ecs_container_mem_usage_bytes`: The current memory in use.
- `ecs_container_mem_max_usage_bytes`: The maximum memory the container has had in use at one time since creation.
- `ecs_container_mem_limit_bytes`: The maximum memory the container can use, as per the task definition and container runtime.
- `ecs_container_cpu_usage`: A number from 0 to 1 which represents the ratio of CPU time used by this container compared to the whole host system in a short interval before your request.

In addition, the following metrics about the entire task are exported (no `ContainerName` label is applied):
- `ecs_container_exporter_up`: 1.0 if no errors were encountered during the the scrape, and 0.0 otherwise. If it returns 0.0, any metrics that were able to be constructed will still be exported.

## Labels

By default, all metrics are labeled with:
- `Cluster`: Name of the ECS cluster.
- `TaskARN`: Full ARN of the task.
- `TaskDefinitionFamily`: Name of the task definition family this task is a part of
- `TaskDefinitionRevision`: Revision of the family.
- `AvailabilityZone`: AZ this task is running in (subject to availability of this information from the ECS task metadata.)

In addition, metrics specific to the a container, rather than the whole task, get:
- `ContainerName`: The name of the container as specified in the ECS task definition.

In theory, other things could be included, i.e. the Docker labels from the task definition, or shortcuts like just the task ID itself rather than the full ARN. These were chosen to be roughly minimal from which anything else can be deduced by anyone who can look at the task definition.

In some ways, most of these labels are against the spirit of [Target labels, not static scraped labels](https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels). However, they are included on the idea that it is may be harder to determine this information on the scraper side depending on how the scraper find this instance.

## Configuration

Configuration is in the form of environment variables, as they are easy to provide to the container via the task definition when deploying to ECS.

- `PORT`: sets the port on which it will listen for HTTP GET requests to the `/metrics` endpoint. The default is 9659, as listed on https://github.com/prometheus/prometheus/wiki/Default-port-allocations .
- `ADDITIONAL_LOG_FIELDS`: add key:value pairs to the logs emitted. It should be valid JSON with values strings. If it is invalid, it will be ignored with a warning. It can be useful for configuring with information about the container it is being deployed with, for example.

## Developing

If you run locally with `make run`, a mock ECS metadata endpoint will be set to run and listen on port `8912`, and the exporter will be run normally but looking to `http://localhost:8912` instead of looking for a real ECS metadata endpoint. The mock endpoint just returns contant data, but it be can tuned to your use case for testing and developing locally.

There is a small test suite. `make test` runs `go test` and also lints and vets.
test
