# build environment ###########################################
FROM golang:1.19.2-alpine@sha256:845f16d6c1c1501505a9f35978494bcd77a03f4f0cfeef56e3d8788325bea4a3 AS build-env

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
FROM alpine:3.16.2@sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-client /bin/network-health-client
COPY --from=build-env /go/bin/network-health-server /bin/network-health-server

ENTRYPOINT ["network-health-server"]
