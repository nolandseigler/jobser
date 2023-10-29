# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21-bookworm AS build

# https://github.com/chemidy/smallest-secured-golang-docker-image/blob/master/docker/scratch.Dockerfile
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

RUN apt-get update \
&& apt-get install -y ca-certificates \
&& update-ca-certificates

COPY . .

RUN CGO_ENABLED=0 go build -o /app/wordserweb -mod=vendor ./cmd/wordserweb/main.go


FROM scratch AS final

USER 10001

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/wordserweb /wordserweb

EXPOSE 8080

ENTRYPOINT ["/wordserweb"]