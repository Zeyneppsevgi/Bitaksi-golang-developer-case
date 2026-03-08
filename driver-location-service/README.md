# driver-location-service

Standalone Driver Location API (Go 1.22+, Fiber v3, Mongo geospatial, internal API key auth).

## Run
```bash
docker compose up --build
```

## Docs
- OpenAPI YAML: `http://localhost:8080/openapi.yaml`
- OpenAPI JSON: `http://localhost:8080/openapi.json`
- Swagger UI: `http://localhost:8080/docs/index.html`

## Endpoints
- `GET /healthz` public
- `GET /readyz` public (mongo ping)
- `POST /v1/driver-locations/batch` protected (`X-Internal-Api-Key`)
- `POST /v1/driver-locations/import` protected (`X-Internal-Api-Key`)
- `GET /v1/driver-locations/search` protected (`X-Internal-Api-Key`)

## CSV format
Standard format:
```csv
driverId,longitude,latitude
driver-1,29.01234,41.01234
```
> [!TIP]
> **Robustness:** Column headers are processed case-insensitively (`driverId`, `driverID`, or `DRIVERID` are all accepted). The system also auto-detects semicolon (`;`) delimiters.

## Examples
Batch upsert:
```bash
curl -X POST "http://localhost:8080/v1/driver-locations/batch" \
  -H "Content-Type: application/json" \
  -H "X-Internal-Api-Key: some-secret" \
  -d '{"items":[{"driverId":"driver-1","location":{"type":"Point","coordinates":[29.0,41.0]}}]}'
```

CSV import:

Example CSV file:
https://raw.githubusercontent.com/Zeyneppsevgi/Bitaksi-golang-developer-case/main/driver-location-service/data/sample_drivers.csv

```bash
curl -X POST "http://localhost:8080/v1/driver-locations/import" \
  -H "X-Internal-Api-Key: some-secret" \
  -F "file=@data/drivers.csv"
```

Search nearest:
```bash
curl "http://localhost:8080/v1/driver-locations/search?lon=29.0&lat=41.0&radius_m=3000" \
  -H "X-Internal-Api-Key: some-secret"
```

## Migrations / Seed
- Index migration runs automatically on startup (`2dsphere + unique(driverId)`).
- In operational scenarios, dummy data can be imported via the API.
- Optional seed data for standalone service:
  - `SEED_ON_START=true`
  - `SEED_FILE=data/drivers.csv`

## Tests
```bash
go test ./...
go test -tags=integration ./test/integration/...
```

## Env
- `DRIVER_PORT` (fallback `PORT`) default `8080`
- `MONGO_URI` default `mongodb://mongo:27017`
- `MONGO_DB` default `driver_location`
- `MONGO_COLLECTION` default `driver_locations`
- `INTERNAL_API_KEY` default `some-secret`
- `SEED_ON_START` default `false`
- `SEED_FILE` default `data/drivers.csv`
