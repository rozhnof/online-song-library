FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/song-service /app/cmd/main.go


FROM alpine:3 AS app

WORKDIR /app

COPY --from=builder /app/song-service /app/song-service
COPY --from=builder /app/config /app/config

CMD ["/app/song-service"]