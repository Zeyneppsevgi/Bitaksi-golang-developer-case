#!/bin/sh
set -eu

export DRIVER_PORT="${DRIVER_PORT:-8080}"
export MATCHING_PORT="${MATCHING_PORT:-8081}"
export SEED_FILE="${SEED_FILE:-/app/driver/data/drivers.csv}"

(
  cd /app/driver
  /usr/local/bin/driver-location
) &
DRIVER_PID=$!

# Wait for driver-location to be healthy before starting matching.
for i in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:${DRIVER_PORT}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

(
  cd /app/matching
  /usr/local/bin/matching
) &
MATCHING_PID=$!

# Optional seed through HTTP endpoint (service path), not direct Mongo import.
if [ "${AUTO_SEED:-true}" = "true" ] && [ -f "${SEED_FILE}" ]; then
  curl -fsS -X POST "http://127.0.0.1:${DRIVER_PORT}/v1/driver-locations/import" \
    -H "X-Internal-Api-Key: ${INTERNAL_API_KEY}" \
    -F "file=@${SEED_FILE}" >/dev/null || true
fi

term() {
  kill "$DRIVER_PID" "$MATCHING_PID" 2>/dev/null || true
  wait "$DRIVER_PID" "$MATCHING_PID" 2>/dev/null || true
}
trap term INT TERM

wait -n "$DRIVER_PID" "$MATCHING_PID"
term
