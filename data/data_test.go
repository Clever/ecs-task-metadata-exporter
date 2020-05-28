package data

import (
	"net/http"
	"testing"
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
		t.Logf("Got metadata: %+v", meta)
	}

	if stats, err := m.Stats(); err != nil {
		t.Fatalf("got error from stats(): %v", err)
	} else {
		t.Logf("Got stats: %+v", stats)
	}
}
