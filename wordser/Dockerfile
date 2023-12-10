# syntax=docker/dockerfile:1
ARG RUST_VERSION=1.73.0
ARG APP_NAME=wordser

FROM rust:${RUST_VERSION}-slim-bookworm AS build

WORKDIR /app

# if things start breaking the rust crate is why we need pytorch.
# zero effort has been put into keeping them in sync.
# https://github.com/guillaume-be/rust-bert#manual-installation-recommended
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y ca-certificates \
    libssl-dev \
    pkg-config \
    python3-pip \
    && update-ca-certificates

# this is what is busted so isolate it
RUN pip3 install torch==2.0.0 --index-url https://download.pytorch.org/whl/cpu --break-system-packages

ENV LIBTORCH_USE_PYTORCH=1 \
    LIBTORCH=/usr/local/lib/python3.11/dist-packages/torch \
    LD_LIBRARY_PATH=/usr/local/lib/python3.11/dist-packages/torch/lib:$LD_LIBRARY_PATH

COPY . .

RUN cargo build --locked --release

EXPOSE 8080

ENTRYPOINT ["/app/target/release/wordser"]
