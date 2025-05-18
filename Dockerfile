# ============================ Build Stage ============================
FROM --platform=$BUILDPLATFORM golang:1.24.2-alpine3.19 as build

LABEL org.opencontainers.image.source="https://github.com/bandprotocol/falcon"

RUN apk add --update --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev

ARG TARGETARCH
ARG BUILDARCH

RUN if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then \
    wget -c https://musl.cc/aarch64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr; \
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then \
    wget -c https://musl.cc/x86_64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr; \
    fi

# Set working directory inside the container
WORKDIR /app

ADD . .

RUN if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then \
    export CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++;\
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then \
    export CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++; \
    fi; \
    GOOS=linux GOARCH=$TARGETARCH CGO_ENABLED=1 LDFLAGS='-linkmode external -extldflags "-static"' make install

RUN if [ -d "/go/bin/linux_${TARGETARCH}" ]; then mv /go/bin/linux_${TARGETARCH}/* /go/bin/; fi

# ============================ Final Stage ============================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

# Set working directory inside the container
WORKDIR /app

# Create non-root user for security
RUN addgroup -S falcon && adduser -S falcon -G falcon

# Copy over binaries from the build
COPY --from=build /go/bin/falcon /usr/bin/falcon

# Set ownership for non-root user
RUN chown -R falcon:falcon /app

# Switch to non-root user
USER falcon

ENTRYPOINT ["falcon", "start"]
