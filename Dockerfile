# build environment
FROM golang:1.17.0@sha256:4f5b9100c3660dd36da84ae865de6746234627e8456d04f594cf7e3c140cd079 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .

# final
FROM alpine:3.14.1@sha256:eb3e4e175ba6d212ba1d6e04fc0782916c08e1c9d7b45892e9796141b1d379ae

RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-sidecar /bin/network-health-sidecar

ENTRYPOINT ["network-health-sidecar"]
