# build environment
FROM golang:1.16.7 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .

# final
FROM alpine:3.14.1

RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-sidecar /bin/network-health-sidecar

ENTRYPOINT ["network-health-sidecar"]
