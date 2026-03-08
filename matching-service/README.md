# matching-service

Standalone matching API (Go 1.25+, Fiber v3). Validates User JWT (`authenticated=true`) and calls the driver-location service.

## Run standalone
```bash
docker compose up --build
```

## Docs
- OpenAPI YAML: `http://localhost:8081/openapi.yaml`
- OpenAPI JSON: `http://localhost:8081/openapi.json`
- Swagger UI: `http://localhost:8081/docs/index.html`

## Endpoints
- `GET /v1/token`
- `GET /v1/match?lon=...&lat=...&radius_m=...` (Bearer JWT required)

## Generate JWT
Use the `/v1/token` endpoint to generate a JWT.


## Example request
```bash

TOKEN=$(curl -s http://localhost:8081/v1/token | jq -r .data.token)


curl "http://localhost:8081/v1/match?lon=29.0&lat=41.0&radius_m=3000" \
  -H "Authorization: Bearer ${TOKEN}"
```


## Full stack run
- To run both services in the same container: `../ops/docker-compose.yml`
- Alternative: `docker-compose.full.yml`

## Tests
```bash
go test ./...
go test -tags=integration ./test/integration/...
```

## Env
- `MATCHING_PORT` (fallback `PORT`) default `8081`
- `DRIVER_LOCATION_BASE_URL` default `http://driver-location-service:8080`
- `INTERNAL_API_KEY` default `some-secret`
- `USER_JWT_SECRET` default `user-secret`
