FROM ubuntu:16.04

ADD https://github.com/pact-foundation/pact-go/releases/download/v0.0.8/linux-amd64.tar.gz /pact
RUN mv /pact/pact-go_linux_amd64 /pact/pact-go

EXPOSE 6666

CMD /pact/pact-go daemon
