FROM golang:1.8

RUN useradd -u 112 jenkins \
    && go get -u github.com/fzipp/gocyclo\
    && go get -u github.com/golang/lint/golint\
    && go get -u github.com/tebeka/go2xunit\
    && go get -u github.com/Masterminds/glide
