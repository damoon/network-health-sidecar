# build environment
FROM golang:1.17.0@sha256:06e92e576fc7a7067a268d47727f3083c0a564331bfcbfdde633157fc91fb17d AS build-env

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
FROM alpine:3.14.1@sha256:eb3e4e175ba6d212ba1d6e04fc0782916c08e1c9d7b45892e9796141b1d379ae
RUN apk add --no-cache ca-certificates

COPY --from=client /go/bin/network-health-client /bin/network-health-client
COPY --from=server /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
