#build stage
ARG BASE=golang:1.18-alpine3.16
FROM ${BASE} AS builder

ARG ALPINE_PKG_BASE="make git gcc libc-dev libsodium-dev zeromq-dev"
ARG ALPINE_PKG_EXTRA=""

LABEL license='SPDX-License-Identifier: Apache-2.0' \
    copyright='Copyright (c) 2022: Intel'
RUN sed -e 's/dl-cdn[.]alpinelinux.org/dl-4.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}
WORKDIR /simple-filter-xml

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
ARG MAKE="make build"
RUN $MAKE

# Docker Client
FROM alpine:3.16

# REST Communications
RUN apk add --update --no-cache ca-certificates zeromq dumb-init curl libc6-compat

RUN apk add --no-cache tzdata
ENV TZ=Asia/Ho_Chi_Minh

# # Copy files
# COPY --from=builder /simple-filter-xml/Attribution.txt /Attribution.txt
# COPY --from=builder /simple-filter-xml/LICENSE /LICENSE

COPY --from=builder /simple-filter-xml/app-service-demo /app-service-demo

COPY --from=builder /simple-filter-xml/res /res

EXPOSE 59708

CMD ["./app-service-demo"]
# ENTRYPOINT ["/simple-filter-xml"] 
