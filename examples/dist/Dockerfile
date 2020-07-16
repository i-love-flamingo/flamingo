# Example Dockerfile for Flamingo/Go based Projects

# Builder
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache ca-certificates tzdata git && update-ca-certificates
COPY . /app
ARG APPVERSION=develop
RUN cd /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X flamingo.me/flamingo/v3/framework/flamingo.appVersion=${APPVERSION}" -o app .

# Final image
FROM scratch

# add timezone data and ssl root certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# add artifacts
#ADD config /config

# add binary
COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]
CMD ["serve"]
