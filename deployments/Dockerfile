FROM golang:1.24.3 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o app ./cmd/app

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y libc6

COPY --from=build /app/app /usr/local/bin/app

COPY config/config.yml /config/config.yml

EXPOSE 8080

CMD ["/usr/local/bin/app"]
