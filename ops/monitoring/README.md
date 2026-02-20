# Monitoring Bootstrap

## Metrics Targets
- backend `/healthz`
- backend request latency and error rate (planned: Prometheus endpoint)
- postgres availability and replication lag (if replicas are added)

## Logs
- backend uses structured Zap logs.
- recommended fields: request_id, route, duration_ms, error.
