# build environment
FROM golang:1.16.7@sha256:47ccd2936048419069f459360885cf71a7ce5896ac3a4263a1c050e160b7d936 AS build-env

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
