# ops

In the default compose setup, services run in separate containers:
- `mongo`
- `driver-location-service`
- `matching-service`

## Run the Full Stack with Docker Compose
```bash
cd ops
docker compose up --build
```

## Endpoints
- Token: `http://localhost:8080/v1/token`
- Driver health: `http://localhost:8080/healthz`
- Driver docs: `http://localhost:8080/docs/index.html`
- Matching health: `http://localhost:8081/healthz`
- Matching docs: `http://localhost:8081/docs/index.html`
