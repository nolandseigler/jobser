# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21-bookworm AS build

WORKDIR /app

RUN apt-get update \
&& apt-get install -y ca-certificates \
&& update-ca-certificates

COPY . .

RUN CGO_ENABLED=0 go build -o /app/wordserweb -mod=vendor ./cmd/wordserweb/main.go 


FROM scratch AS final

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/wordserweb /wordserweb

EXPOSE 8080

ENTRYPOINT ["/wordserweb"]