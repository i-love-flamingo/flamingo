FROM golang:1.10

RUN mkdir /pact
ADD https://github.com/pact-foundation/pact-go/releases/download/v0.0.11/pact-go_linux_amd64.tar.gz /pact
RUN cd /pact && tar xf pact-go_linux_amd64.tar.gz  && ln -s /pact/pact-go /pact/pact-go_linux_amd64

RUN go get -u github.com/golang/dep/cmd/dep
