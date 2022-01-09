###############################################################################
# Build container
###############################################################################
FROM golang:1.17-alpine3.15 as build

# Install build tools
RUN apk add binutils

# Add source code
ADD . /go/src/github.com/deesel/wol/
WORKDIR /go/src/github.com/deesel/wol/cmd/wol

# Build and install module
RUN go install && strip /go/bin/wol


###############################################################################
# Final container
###############################################################################
FROM alpine:3.15 as final

# Install application
COPY --from=build /go/bin/wol /bin/wol

ENTRYPOINT /bin/wol

EXPOSE 8001
