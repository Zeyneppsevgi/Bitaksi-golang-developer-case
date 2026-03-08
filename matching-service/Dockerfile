FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN adduser -D appuser && apk add --no-cache wget
COPY --from=builder /bin/server /usr/local/bin/server
COPY openapi /app/openapi
USER appuser
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/server"]
