# build environment ###########################################
FROM golang:1.20.2-alpine@sha256:576da1aa73f8ffa3dc6a7577b6032bd834aa84f2e1714d3e7e96b06b49f4e177 AS build-env

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
FROM alpine:3.17.2@sha256:ff6bdca1701f3a8a67e328815ff2346b0e4067d32ec36b7992c1fdc001dc8517 AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-client /bin/network-health-client
COPY --from=build-env /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
