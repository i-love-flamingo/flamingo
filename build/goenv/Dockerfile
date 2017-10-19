FROM golang:1.9

RUN mkdir /pact
ADD https://github.com/pact-foundation/pact-go/releases/download/v0.0.8/linux-amd64.tar.gz /pact
RUN cd /pact && tar xf linux-amd64.tar.gz

RUN go get -u github.com/golang/dep/cmd/dep
