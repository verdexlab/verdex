# Build
FROM golang:1.21-alpine AS build
RUN apk add build-base
WORKDIR /pkg
COPY . /pkg
RUN go mod download
RUN go build -ldflags "-X github.com/verdexlab/verdex/verdex/core.releaseEnvironment=release-docker" -o verdex-binary

# Binary
FROM alpine:3.20.3
RUN apk upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates
COPY --from=build /pkg/verdex-binary /usr/local/bin/verdex

ENTRYPOINT ["verdex"]
