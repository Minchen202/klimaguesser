FROM golang:1.24.1-alpine3.20 AS builder

RUN apk add build-base gcc

WORKDIR /app

COPY . .

EXPOSE ${PORT:-8081}

CMD cd backend && go run .