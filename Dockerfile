FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api

# Pinned distroless static-debian12 nonroot image by digest.
FROM gcr.io/distroless/static-debian12@sha256:a9329520abc449e3b14d5bc3a6ffae065bdde0f02667fa10880c49b35c109fd1 AS final

WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/api"]
