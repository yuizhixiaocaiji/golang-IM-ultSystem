FROM golang:1.18 AS build

ADD . /usr/local/go/src/golang-im-ulisystem

WORKDIR /usr/local/go/src/golang-im-ulisystem

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server

FROM alpine:3.12

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on

RUN echo "https://mirrors.aliyun.com/alpine/v3.12/main/" > /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates && \
    echo "host: files dns" > /etc/nsswitch.conf && \
    mkdir -p /www/config

WORKDIR /www

COPY --from=build /usr/local/go/src/golang-im-ulisystem/api_server /usr/bin/api_server
ADD ./config /www/config

RUN chmod +x /usr/bin/api_server

ENTRYPOINT ["api_server"]

