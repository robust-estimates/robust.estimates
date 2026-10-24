FROM golang:1.27.3-bookworm AS build
ENV CGO_ENABLED=0 GOBIN=/out
RUN go install github.com/minio/minio@v0.0.0-20260212201848-7aac2a2c5b7c
COPY deploy/minio/healthcheck.go /healthcheck.go
RUN go build -trimpath -o /out/healthcheck /healthcheck.go
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/minio /usr/local/bin/minio
COPY --from=build /out/healthcheck /usr/local/bin/healthcheck
WORKDIR /data
EXPOSE 9000 9001
ENTRYPOINT ["/usr/local/bin/minio"]
