# build environment ###########################################
FROM golang:1.20.2-alpine@sha256:ff56122f00ff5f5f06af287b26cd0ab6c210361e3fac209c900b0cb13ec4cb43 AS build-env

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
FROM alpine:3.17.3@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126 AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-client /bin/network-health-client
COPY --from=build-env /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
