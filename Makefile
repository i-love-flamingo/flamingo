.PHONY: test integrationtest

test:
	go test -race -v ./...
	gofmt -l -e -d .
	golint ./...

integrationtest:
	go test -v ./core/cache/... -tags=docker
