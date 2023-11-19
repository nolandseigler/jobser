# syntax=docker/dockerfile:1
ARG RUST_VERSION=1.73.0
ARG APP_NAME=wordser

FROM rust:${RUST_VERSION}-slim-bookworm AS build

WORKDIR /app

RUN apt-get update \
&& apt-get install -y ca-certificates \
libssl-dev \
pkg-config \
&& update-ca-certificates

COPY . .

RUN cargo build --locked --release

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/target/release/wordser"]