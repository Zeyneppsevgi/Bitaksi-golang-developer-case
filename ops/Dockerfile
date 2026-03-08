FROM golang:1.25-alpine AS build-driver
WORKDIR /src/driver-location-service
COPY driver-location-service/go.mod driver-location-service/go.sum ./
RUN go mod download
COPY driver-location-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/driver-location ./cmd/server

FROM golang:1.25-alpine AS build-matching
WORKDIR /src/matching-service
COPY matching-service/go.mod matching-service/go.sum ./
RUN go mod download
COPY matching-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/matching ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN adduser -D appuser && apk add --no-cache curl
COPY --from=build-driver /bin/driver-location /usr/local/bin/driver-location
COPY --from=build-matching /bin/matching /usr/local/bin/matching
COPY driver-location-service/openapi /app/driver/openapi
COPY driver-location-service/data /app/driver/data
COPY matching-service/openapi /app/matching/openapi
COPY ops/start.sh /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/start.sh && chown -R appuser:appuser /app
USER appuser
EXPOSE 8080 8081
ENTRYPOINT ["/usr/local/bin/start.sh"]
