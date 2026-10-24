FROM golang:1.27.3-bookworm AS build
WORKDIR /src
COPY contracts ./contracts
COPY platform ./platform
COPY company-service ./company-service
WORKDIR /src/company-service
ENV GOWORK=off CGO_ENABLED=0
RUN go mod download && go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY company-service/migrations /app/migrations
USER 65532:65532
ENTRYPOINT ["/app/server"]
