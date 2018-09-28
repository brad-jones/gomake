# Upstream image sources
FROM docker:18.09.0 as docker
FROM node:11.2.0-alpine AS node
FROM golang:1.11.2-alpine3.8 AS go

FROM alpine:3.8 AS final

# This is used by both the docker and golang install
# see: https://github.com/docker-library/docker/blob/master/18.09/Dockerfile
# see: https://github.com/docker-library/golang/blob/master/1.11/alpine3.8/Dockerfile
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

# Copy docker
COPY --from=docker /usr/local/. /usr/local/

# Copy golang
COPY --from=go /usr/local/. /usr/local/
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Copy node & yarn
COPY --from=node /usr/lib/. /usr/lib/
COPY --from=node /usr/local/. /usr/local/
COPY --from=node /opt/. /opt/

# Install what we can directly from the Alpine packages.
# Mandatory tooling for circleci images is also included.
# see: https://circleci.com/docs/2.0/custom-images/#required-tools-for-primary-containers
RUN apk --no-cache add \
    bash \
    ca-certificates \
    curl \
    gcc \
    git \
    gzip \
    make \
    openssh-client \
    tar \
    xz
