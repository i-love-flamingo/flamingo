FROM golang:1.12

RUN go get -v -u github.com/golang/dep/cmd/dep && mkdir /pact
ADD https://github.com/pact-foundation/pact-go/releases/download/v0.0.13/pact-go_linux_amd64.tar.gz /pact
RUN cd /pact && tar xf pact-go_linux_amd64.tar.gz  && ln -s /pact/pact-go /pact/pact-go_linux_amd64
COPY ./patch /pact/
RUN go get -u golang.org/x/lint/golint && go get -d github.com/go-critic/go-critic/checkers/testdata/_integration/check_main_only

