package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/Clever/kayvee-go.v6/logger"

	"github.com/Clever/ecs-task-metadata-exporter/data"
)

// We can try to detect an ECS metadata endpoint from env vars starting with the newest supported version and descending down.
// First, look for V4 (Fargate Platform Version >= 1.4.0 or ECS container agent >= 1.39.0).
// Failing that, look for V3 (Fargate Platform Version >= 1.3.0 or ECS container agent >= 1.21.0).
// There is a V2, which is at 169.254.170.2/v2 for Fargte >= 1.10 or (ECS container agent >= 1.17.0 AND using awsvpc network mode).
// But we don't support it at present (for example, it uses /metadata instead of / and /stats instead of /task/stats).
const (
	ECSMetadataURIV4Var = "ECS_CONTAINER_METADATA_URI_V4"
	ECSMetadataURIV3Var = "ECS_CONTAINER_METADATA_URI"
)

const defaultPort = "9659"

var mainLogger = logger.New("ecs-task-metadata-exporter")

func main() {
	port := defaultPort
	if portStr, ok := os.LookupEnv("PORT"); ok {
		port = portStr
	}
	if logFields, ok := os.LookupEnv("ADDITIONAL_LOG_FIELDS"); ok {
		var fields map[string]string
		if err := json.Unmarshal([]byte(logFields), &fields); err != nil {
			mainLogger.WarnD("bad-additional-log-fields", logger.M{
				"error":                 fmt.Sprintf("decoding ADDITIONAL_LOG_FIELDS: %v", err),
				"ADDITIONAL_LOG_FIELDS": logFields,
			})
		} else {
			for k, v := range fields {
				mainLogger.AddContext(k, v)
			}
		}
	}

	var endpoint string
	if os.Getenv("IS_LOCAL") != "" {
		handler := data.ConstantMetadataEndpointHandler(
			data.SampleTaskMetadata, data.SampleTaskStats,
		)
		server := http.Server{
			Addr:    "localhost:8912",
			Handler: handler,
		}
		go server.ListenAndServe()
		defer server.Close()

		endpoint = "http://localhost:8912"
		mainLogger.InfoD("using-source", logger.M{
			"source": "localhost",
			"uri":    endpoint,
		})
	} else {
		endpoint = mustGetECSMetadataURI()
	}
	s := data.NewMetadataEndpointSource(endpoint)

	c := NewCollector(s, mainLogger)
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	promServerOpts := promhttp.HandlerOpts{
		// Log errors from the http server to our main logger with title promhttp-error
		ErrorLog: kayveePrintlnLogger{l: mainLogger, title: "promhttp-error"},
	}
	http.Handle("/metrics", promhttp.HandlerFor(reg, promServerOpts))
	if os.Getenv("EXPOSE_RAW_DATA") != "" {
		setupDebugRoutes(s)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

	log.Println("ecs-task-metadata-exporter exited without error")
}

// kayveePrintlnLogger implements the prometheus.Logger interface using a kayvee logger.Logger
type kayveePrintlnLogger struct {
	title string
	l     logger.KayveeLogger
}

func (k kayveePrintlnLogger) Println(v ...interface{}) {
	k.l.InfoD(k.title, logger.M{
		"line": fmt.Sprintln(v...),
	})
}

func mustGetECSMetadataURI() string {
	if uri, ok := os.LookupEnv(ECSMetadataURIV4Var); ok {
		mainLogger.InfoD("using-source", logger.M{
			"source": "ECSMetadataURIV4",
			"uri":    uri,
		})
		return uri
	}
	if uri, ok := os.LookupEnv(ECSMetadataURIV3Var); ok {
		mainLogger.InfoD("using-source", logger.M{
			"source": "ECSMetadataURIV3",
			"uri":    uri,
		})
		return uri
	}
	panic(fmt.Errorf("couldn't detect ECS metadata endpoint (tried env vars %s and %s)", ECSMetadataURIV4Var, ECSMetadataURIV3Var))
}

func setupDebugRoutes(source data.Source) {
	http.Handle("/_debug/metadata", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		metadata, err := source.Metadata()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error": "getting metadata from source: %s"}`, err)))
		}
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error": "marshalling metadata: %s"}`, err)))
		}
		res.Write(metadataBytes)
	}))
	http.Handle("/_debug/stats", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		stats, err := source.Stats()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error": "getting stats from source: %s"}`, err)))
		}
		statsBytes, err := json.Marshal(stats)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error": "marshalling stats: %s"}`, err)))
		}
		res.Write(statsBytes)
	}))
}
