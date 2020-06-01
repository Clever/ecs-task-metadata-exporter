FROM alpine:3.10

RUN apk add ca-certificates
RUN update-ca-certificates

COPY kvconfig.yml /bin/kvconfig.yml
COPY bin/ecs-task-metadata-exporter /bin/ecs-task-metadata-exporter

CMD ["/bin/ecs-task-metadata-exporter"]
