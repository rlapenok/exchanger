FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020
RUN templ generate ./web/templ/...

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/exchanger ./cmd

FROM alpine:3.21

RUN apk add --no-cache ca-certificates netcat-openbsd curl \
    && curl -fsSL https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/local/bin migrate

WORKDIR /app

COPY migrations ./migrations
COPY --from=builder /out/exchanger /app/exchanger
COPY docker/entrypoint.sh /app/entrypoint.sh

RUN sed -i 's/\r$//' /app/entrypoint.sh \
    && chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
