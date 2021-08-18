# build environment
FROM golang:1.17.0@sha256:634cda4edda00e59167e944cdef546e2d62da71ef1809387093a377ae3404df0 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY sidecar sidecar

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./sidecar/cmd

# final
FROM alpine:3.14.1@sha256:eb3e4e175ba6d212ba1d6e04fc0782916c08e1c9d7b45892e9796141b1d379ae

RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/cmd /bin/sidecar

ENTRYPOINT ["sidecar"]
