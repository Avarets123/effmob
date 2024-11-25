FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o effapp cmd/main.go
RUN go build -o migrator migrations/main.go

FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

USER 1000
COPY --from=build --chown=1000:1000 /app/effapp .
COPY --from=build --chown=1000:1000 /app/migrations ./migrations
COPY --from=build --chown=1000:1000 /app/migrator .


CMD [ "sh", "-c", "/app/migrator && /app/effapp"  ]