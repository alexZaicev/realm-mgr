# ------------------------------------------------------------------------------
# Build image
# ------------------------------------------------------------------------------
FROM golang:1.18.8-alpine3.16 AS build

# install build tools
RUN set -eux; \
	apk add -U --no-cache \
		curl \
		git  \
		make \
	;

# add new user
ARG USER=default
RUN addgroup -g 1000 ${USER} \
        && adduser -h /build -D -u 1000 -G ${USER} ${USER} \
    ;

USER ${USER}

WORKDIR /build

# copy vendored dependencies
COPY --chown=${USER}:${USER} go.mod go.sum ./
COPY --chown=${USER}:${USER} ./vendor ./vendor
COPY --chown=${USER}:${USER} ./internal ./internal
COPY --chown=${USER}:${USER} ./cmd ./cmd
COPY --chown=${USER}:${USER} ./proto ./proto
COPY --chown=${USER}:${USER} ./Makefile ./Makefile
COPY --chown=${USER}:${USER} ./tests ./tests

ARG VERSION
RUN set -eux; \
    CGO_ENABLED=0 make build VERSION=${VERSION};

# ------------------------------------------------------------------------------
# Additional build steps for grpc image
# ------------------------------------------------------------------------------

FROM alpine:3.16 AS grpc-build

WORKDIR /build

ARG GRPC_HEALTH_PROBE_VERSION=0.4.4
ARG GRPC_HEALTH_PROBE_CHECKSUM=0cc53c34fb392e3c7e035b598cd17f82d05a822a10655c0bbbd2b41e9136fab8
ARG HTTPS_PROXY
RUN set -eux; \
    # The default mirror seems to be a networking ticking timebomb for DIND builds
    # see <https://github.com/gliderlabs/docker-alpine/issues/307#issuecomment-649642293>
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.16/main" > /etc/apk/repositories; \
    apk add -U --no-cache \
        ca-certificates \
        wget \
    ; \
    # install grpc health probe
    export https_proxy=${HTTPS_PROXY}; \
    wget -O grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64; \
    echo "${GRPC_HEALTH_PROBE_CHECKSUM}  grpc_health_probe" | sha256sum -c; \
    chmod 755 grpc_health_probe

# ------------------------------------------------------------------------------
# Build runtime image
# ------------------------------------------------------------------------------
FROM alpine:3.16 AS base-runtime

ARG HTTPS_PROXY
RUN set -eux; \
    # The default mirror seems to be a networking ticking timebomb for DIND builds
    # see <https://github.com/gliderlabs/docker-alpine/issues/307#issuecomment-649642293>
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.16/main" > /etc/apk/repositories; \
    apk add -U --no-cache \
        ca-certificates \
        curl \
    ;

# add new user
ARG USER=default
RUN set -eux; \
    addgroup ${USER}; \
    adduser -D -G ${USER} ${USER};

USER ${USER}
WORKDIR /home/${USER}

# ------------------------------------------------------------------------------
# realm-mgr-grpc runtime image
# ------------------------------------------------------------------------------

FROM base-runtime AS realm-mgr-grpc

COPY --from=build /build/dist/realm-mgr-grpc /usr/local/bin/realm-mgr-grpc
COPY --from=grpc-build /build/grpc_health_probe /usr/local/bin/grpc_health_probe
COPY build/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/usr/local/bin/realm-mgr-grpc"]
