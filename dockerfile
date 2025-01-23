FROM golang:1.22.3-alpine3.19 as build

# Install minimum necessary dependencies
RUN apk add --no-cache ca-certificates build-base git

WORKDIR /app

# Add source files
COPY . .

# install binary
RUN go mod download
RUN make install

############################################################################
FROM alpine:3.16

# Install CA certificates for secure connections
RUN apk add --no-cache ca-certificates

# Copy over binaries from the build
COPY --from=build /go/bin/falcon /usr/bin/falcon
