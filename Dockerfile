
ARG GO_VERSION=1.12
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git gcc libc-dev && \
        mkdir -p /build/etc/ssl/certs && \
        cp /etc/ssl/certs/ca-certificates.crt /build/etc/ssl/certs/

WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /build/app cmd/server/server.go

################################################################################

# alpine required because the built app depends on some dynamicly linked lib,
# thanks to not being able to build without CGO which the sqlite3 pkg depends on
FROM alpine:3.9 AS final
#FROM scratch AS final

COPY --from=builder /build /

ENV HOME /data
WORKDIR $HOME
VOLUME $HOME
EXPOSE 8082

CMD ["/app", "-addr=:8082", "-path=/data/servers.db"]
