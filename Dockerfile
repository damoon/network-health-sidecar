# build environment
FROM golang:1.17.0@sha256:33ef0040801bb4deabe1db381ee92de1afc81b869ce27d52fb52d24cf37ff543 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

# client
FROM build-env AS client
COPY ./cmd/network-health-client ./cmd/network-health-client
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./cmd/network-health-client

# server
FROM build-env AS server
COPY ./cmd/network-health-server ./cmd/network-health-server
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./cmd/network-health-server

# final
FROM alpine:3.14.2@sha256:e1c082e3d3c45cccac829840a25941e679c25d438cc8412c2fa221cf1a824e6a
RUN apk add --no-cache ca-certificates

COPY --from=client /go/bin/network-health-client /bin/network-health-client
COPY --from=server /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
