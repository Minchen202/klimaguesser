FROM golang:1.24.1-alpine3.20 AS builder
RUN apk add --no-cache build-base gcc git ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/klimaguessr ./backend

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/klimaguessr /app/klimaguessr

EXPOSE ${PORT:-8081}
ENV PORT=${PORT:-8081}
CMD ["/app/klimaguessr"]