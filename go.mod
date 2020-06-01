module github.com/Clever/ecs-task-metadata-exporter

go 1.13

require (
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200514230353-811a247d06e8+incompatible // indirect; <- actually pinned to commit for v19.03.9 , but they removed their go.mod so things go weird
	github.com/docker/engine v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-openapi/swag v0.19.9
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/go-cmp v0.4.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.6.0
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/Clever/kayvee-go.v6 v6.23.0
	gotest.tools v2.2.0+incompatible // indirect
)
