FROM golang:1.27.3-bookworm AS build
WORKDIR /src
COPY contracts ./contracts
COPY platform ./platform
COPY gateway ./gateway
WORKDIR /src/gateway
ENV GOWORK=off CGO_ENABLED=0
RUN go mod download && go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY gateway/migrations /app/migrations
COPY api/public/openapi.yaml /app/api/openapi.yaml
COPY docs/swagger.html /app/docs/swagger.html
ENV OPENAPI_FILE=/app/api/openapi.yaml SWAGGER_FILE=/app/docs/swagger.html
USER 65532:65532
ENTRYPOINT ["/app/server"]
