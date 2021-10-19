# build environment ###########################################
FROM golang:1.17.2-alpine@sha256:5519c8752f6b53fc8818dc46e9fda628c99c4e8fd2d2f1df71e1f184e71f47dc AS build-env

WORKDIR /app

# entrypoint
RUN apk add --no-cache entr
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

# dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# client
COPY ./cmd/network-health-client ./cmd/network-health-client
RUN go install ./cmd/network-health-client

# server
COPY ./cmd/network-health-server ./cmd/network-health-server
COPY ./pkg ./pkg
RUN go install ./cmd/network-health-server

# production image ############################################
FROM alpine:3.14.2@sha256:e1c082e3d3c45cccac829840a25941e679c25d438cc8412c2fa221cf1a824e6a AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-client /bin/network-health-client
COPY --from=build-env /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
