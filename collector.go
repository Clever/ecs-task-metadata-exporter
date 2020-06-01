package main

import (
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/Clever/kayvee-go.v6/logger"

	"github.com/Clever/ecs-task-metadata-exporter/data"
	"github.com/Clever/ecs-task-metadata-exporter/metrics"
)

type collector struct {
	Source data.Source
	Logger logger.KayveeLogger
}

// NewCollector returns a prometheus.Collector configured to collect Docker metrics
func NewCollector(source data.Source, l logger.KayveeLogger) prometheus.Collector {
	// Set a logger with discarded output instead of nil, so we can call methods on log without panicing/checking for nil every time.
	if l == nil {
		l = logger.New("")
		l.SetOutput(ioutil.Discard)
	}
	return collector{
		Source: source,
		Logger: l,
	}
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	// By construction, no metrics will ever change, but at describe time, we haven't inspected the metadata yet, so we don't have the full list
	prometheus.DescribeByCollect(c, ch)
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	meta, err := c.Source.Metadata()
	if err != nil {
		c.Logger.ErrorD("retrieving-metadata", logger.M{
			"error": err.Error(),
		})
		return
	}

	commonLabels := map[string]string{
		"Cluster":                meta.Cluster,
		"TaskARN":                meta.TaskARN,
		"TaskDefinitionFamily":   meta.Family,
		"TaskDefinitionRevision": meta.Revision,
	}
	if meta.AvailabilityZone != nil {
		commonLabels["AvailabilityZone"] = *meta.AvailabilityZone
	}
	statusDesc := prometheus.NewDesc(metrics.Prefix+"exporter_up", "1 if no issues were encountered during the scrape, 0 if errors occured", nil, commonLabels)

	stats, err := c.Source.Stats()
	if err != nil {
		c.Logger.ErrorD("retrieving-stats", logger.M{
			"error": err.Error(),
		})
		status, err := prometheus.NewConstMetric(statusDesc, prometheus.GaugeValue, 0.0)
		if err != nil {
			c.Logger.ErrorD("reporting-exporter-up-metric", logger.M{
				"up": 0,
			})
		} else {
			ch <- status
		}
		return
	}
	exporterIsUp := 1.0
	for _, container := range meta.Containers {
		// container.Type is used by ECS to distinguish containers internal to ECS from ones that are part of the task
		if container.Type != "NORMAL" {
			continue
		}
		containerID := container.DockerID
		labels := map[string]string{}
		for k, v := range commonLabels {
			labels[k] = v
		}
		labels["ContainerName"] = container.Name
		containerStats, ok := stats[containerID]
		if !ok {
			containersInStats := []string{}
			for k := range stats {
				containersInStats = append(containersInStats, k)
			}
			c.Logger.ErrorD("missing-container", logger.M{
				"missing":                containerID,
				"containers-in-metadata": meta.Containers,
				"containers-in-stats":    containersInStats,
			})
			exporterIsUp = 0.0
			continue
		}
		containerMetrics, err := metrics.StatsToMetrics(containerStats, metrics.DefaultMetrics, labels)
		if err != nil {
			c.Logger.ErrorD("converting-stats", logger.M{
				"error": err.Error(),
			})
			exporterIsUp = 0.0
			continue
		}
		for _, m := range containerMetrics {
			ch <- m
		}
	}
	status, err := prometheus.NewConstMetric(statusDesc, prometheus.GaugeValue, exporterIsUp)
	if err != nil {
		c.Logger.ErrorD("reporting-exporter-up-metric", logger.M{
			"up": exporterIsUp,
		})
	} else {
		ch <- status
	}

}
