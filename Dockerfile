# build environment
FROM golang:1.17.0@sha256:50b5f2673c0995eb3cda07150611acd3bd0523b0b3da5bd079e9b9d905f8e1cc AS build-env

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
