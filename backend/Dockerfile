FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS dev
RUN go install github.com/air-verse/air@v1.52.3
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM base AS build
RUN go build -o /easy-arbitra-server ./cmd/server

FROM alpine:3.20 AS prod
WORKDIR /app
COPY --from=build /easy-arbitra-server /usr/local/bin/easy-arbitra-server
COPY config/config.yaml /app/config/config.yaml
EXPOSE 8080
CMD ["easy-arbitra-server"]
