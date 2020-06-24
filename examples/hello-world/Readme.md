# Flamingo Hello World

## Run it:

```
go run . serve
```

Open: http://localhost:3322

## TLS / HTTPS

Generate certificates:

```
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost:3322
```

Run it:
```
go run . serve -c cert.pem -k key.pem
```

Open: https://localhost:3322
