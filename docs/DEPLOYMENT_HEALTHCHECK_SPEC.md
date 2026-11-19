# Production Deployment & Healthcheck Guidelines

## Container Healthcheck Probes
- **Go Gateway:** `HTTP GET /api/health` every 10s (timeout: 2s, retries: 3)
- **Next.js Web:** `HTTP GET /` every 15s
- **MinIO S3:** `HTTP GET /minio/health/live` every 30s

Zero-downtime rolling updates supported via Docker Compose and Kubernetes deployments.
