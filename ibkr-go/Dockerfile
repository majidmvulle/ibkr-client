FROM golang:1.25 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-s -w" -o server ./cmd/server

FROM scratch AS app

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 50051 8080

CMD ["./server"]
