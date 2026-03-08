# ops

Varsayılan compose ile servisler ayrı container olarak çalışır:
- `mongo`
- `driver-location-service`
- `matching-service`

## Ayrı container çalıştırma (önerilen)
```bash
cd ops
docker compose up --build
```

## Tek container (opsiyonel)
```bash
cd ops
docker compose -f docker-compose.single.yml up --build
```

## Endpointler
- Driver health: `http://localhost:8080/healthz`
- Driver docs: `http://localhost:8080/docs/index.html`
- Matching health: `http://localhost:8081/healthz`
- Matching docs: `http://localhost:8081/docs/index.html`
