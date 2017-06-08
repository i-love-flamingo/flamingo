FROM golang:1.8

ADD https://github.com/pact-foundation/pact-go/releases/download/v0.0.8/linux-amd64.tar.gz /pact/linux-amd64.tar.gz

RUN useradd -u 112 jenkins \
    && cd /pact && tar xzf linux-amd64.tar.gz \
    && go get -u github.com/fzipp/gocyclo \
    && go get -u github.com/golang/lint/golint \
    && go get -u github.com/tebeka/go2xunit \
    && go get -u github.com/Masterminds/glide \
    && go get -u github.com/axw/gocov/... \
    && go get -u github.com/AlekSi/gocov-xml \
    && chown -R 112 /go/
