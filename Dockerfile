# Build webchaind in a stock Go builder container
FROM golang:1.12-stretch as builder

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y gcc make git

ADD . /go/src/github.com/webchain-network/webchaind

WORKDIR /go/src/github.com/webchain-network/webchaind

RUN make cmd/webchaind

# Pull webchaind into a second stage deploy ubuntu container
FROM ubuntu:18.04

RUN apt-get update && apt-get install -y openssh-server iputils-ping iperf3 && apt-get clean
COPY --from=builder /go/src/github.com/webchain-network/webchaind/bin/webchaind /usr/local/bin/

EXPOSE 39573 31440 31440/udp
ENTRYPOINT ["webchaind"]
