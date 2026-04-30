# Online Shop

Backend API service for an online shop, built with Go and deployed on Railway.

## Tech Stack

- **Go 1.26.2** — API service runtime
- **Gin** — HTTP routing
- **PostgreSQL 18** — primary database
- **Viper** — environment and `.env` configuration loading
- **Zap** — structured logging
- **Docker** — containerized builds and local stack
- **Railway** — production deployment

## Prerequisites

- Go 1.26
- Docker and Docker Compose

## Local Development

Create a local environment file:

```bash
cp .env.example .env
```

Start the local infrastructure:

```bash
make up-infra
```

Run the API:

```bash
make run
```

The API is available at `http://localhost:8080`.

To run the full Docker Compose stack:

```bash
make up
```

Rebuild the API image when starting the stack:

```bash
BUILD=1 make up
```

## Configuration

The application reads configuration from OS environment variables and, for local development, an optional `.env` file. OS environment variables override `.env` values.

Required runtime variables:

| Name           | Purpose                      |
| -------------- | ---------------------------- |
| `PORT`         | HTTP port the API listens on |
| `DATABASE_URL` | PostgreSQL connection string |

Production should also set:

| Name       | Purpose                                                 |
| ---------- | ------------------------------------------------------- |
| `GIN_MODE` | Set to `release` to run Gin without debug-mode warnings |

## Running Tests

```bash
make test
```

Equivalent direct command:

```bash
go test ./...
```

## CI/CD

This project follows a trunk-based development strategy:

- **CI** runs on every pull request to `main` — Go tests, Go build, and Docker build must pass before merge
- **Deploy** triggers automatically when a pull request is merged to `main`
- Production deployment uses Railway through the `RAILWAY_TOKEN` GitHub secret

Railway deploys the `api` service in the `production` environment.

## Health Checks

The health endpoints are intentionally public for container and platform probes.

```http request
GET /healthz
```

```http request
GET /readyz
```
