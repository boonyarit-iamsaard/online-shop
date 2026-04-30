FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot AS final

WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/api"]
