FROM golang:1.15 as builder
WORKDIR /go/src/github.com/ruixiaoedu/hiot
COPY . .
ENV GOPROXY https://goproxy.cn,direct
RUN go build -ldflags="-w -s" -o hiot ./cmd

FROM debian:bullseye
WORKDIR /
COPY --from=builder /go/src/github.com/ruixiaoedu/hiot/hiot /hiot
EXPOSE 1883

CMD ["/hiot"]