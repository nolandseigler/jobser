# syntax=docker/dockerfile:1
ARG RUST_VERSION=1.73.0
ARG APP_NAME=wordser

FROM rust:${RUST_VERSION}-slim-bookworm AS build


WORKDIR /app

RUN apt-get update \
&& apt-get install -y ca-certificates \
&& update-ca-certificates

COPY . .

RUN cargo build --locked --release


FROM debian:bookworm-slim AS final

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/target/release/$APP_NAME /wordser

EXPOSE 8080

ENTRYPOINT ["/wordser"]