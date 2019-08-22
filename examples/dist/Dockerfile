# Example Dockerfile for Flamingo/Go based Projects

# Builder
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY . .
RUN go build -o app .

# Final image
FROM scratch

# add timezone data and ssl root certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# add artifacts
ADD config /config

# add binary
ADD app /app

ENTRYPOINT ["/app"]
CMD ["serve"]

