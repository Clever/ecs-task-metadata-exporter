package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/docker/engine/api/types"
)

// Source is an abstraction for a provider of task metadata + stats.
// In practice, it'll be the ECS task metadata endpoint, but it can also be mocked out.
type Source interface {
	// Stats returns a map of DockerIDs to stats of the form returned by the Docker daemon's stats endpoint
	Stats() (map[string]types.StatsJSON, error)
	// Metadata retrieves the metadata for the task
	Metadata() (TaskMetadata, error)
}

// TaskMetadata describes the response of `GET MetadataURI/task`
// See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v3.html#task-metadata-endpoint-v3-response
// and v4 is the same: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html#task-metadata-endpoint-v4-response
type TaskMetadata struct {
	Cluster            string
	TaskARN            string
	Family             string
	Revision           string
	DesiredStatus      string
	KnownStatus        string
	Limits             *Limits // Omitted if there are no limits
	PullStartedAt      time.Time
	PullStoppedAt      time.Time
	AvailabilityZone   *string // Only available on Fargate platform version 1.4.0
	ExecutionStoppedAt time.Time
	Containers         []ContainerMetadata
}

// ContainerMetadata describes the response of `GET MetadataURI` which also appears as part of `GET MetadataURI/task`
// See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v3.html#task-metadata-endpoint-v3-response
// and v4 is the same: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html#task-metadata-endpoint-v4-response
type ContainerMetadata struct {
	DockerID      string
	Name          string
	DockerName    string
	Image         string
	ImageID       string
	Labels        map[string]string
	DesiredStatus string
	KnownStatus   string
	Limits        *Limits // Omitted if there are no limits
	CreatedAt     time.Time
	StartedAt     time.Time
	Type          string
	// TODO Networks
}

// Limits is the limits of a container or the whole task
// Only the limits that have been set are non-nil
type Limits struct {
	CPU    *float64
	Memory *uint64
}

// NewMetadataEndpointSource constructs a Source from the base URI to use as if it is the ECS task metadata URI
func NewMetadataEndpointSource(endpointURI string) Source {
	return &metadataEndpointSource{
		Endpoint: endpointURI,
	}
}

type metadataEndpointSource struct {
	Endpoint string
}

func (m metadataEndpointSource) Metadata() (TaskMetadata, error) {
	var ret TaskMetadata
	taskEndpoint := m.Endpoint + "/task"
	resp, err := http.Get(taskEndpoint)
	if err != nil {
		return ret, fmt.Errorf("GET %s: %v", taskEndpoint, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, fmt.Errorf("reading task metadata response body: %v", err)
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("got non-success status code %d from task metadata endpoint %s with response body: %s", resp.StatusCode, taskEndpoint, string(body))
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, fmt.Errorf("unmarshaling task metadata response json: %v", err)
	}
	return ret, nil
}

func (m metadataEndpointSource) Stats() (map[string]types.StatsJSON, error) {
	var ret map[string]types.StatsJSON
	taskEndpoint := m.Endpoint + "/task/stats"
	resp, err := http.Get(taskEndpoint)
	if err != nil {
		return ret, fmt.Errorf("GET %s: %v", taskEndpoint, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, fmt.Errorf("reading task stats response body: %v", err)
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("got non-success status code %d from task stats endpoint %s with response body: %s", resp.StatusCode, taskEndpoint, string(body))
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, fmt.Errorf("unmarshaling task stats response json: %v", err)
	}
	return ret, nil
}

func constantHandler(body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}
}

// ConstantMetadataEndpointHandler creates an http.Handler that can be used as a mock of the ECS task metadata service.
// It will always respond with the same provided task metadata and the same task stats.
func ConstantMetadataEndpointHandler(taskMetadataResponse, taskStatsResponse []byte) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/task", constantHandler(taskMetadataResponse))
	mux.Handle("/task/stats", constantHandler(taskStatsResponse))
	mux.Handle("/", http.NotFoundHandler())
	return mux
}
