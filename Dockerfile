FROM alpine:latest

RUN adduser -u 10001 -D app-runner

MAINTAINER proxy "370474613@qq.com"

## 设置 操作系统时区
RUN rm -rf /etc/localtime
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime


ENV RUN_ENV docker
RUN mkdir -p /home/proxy
RUN mkdir -p /home/proxy-pool/logs
COPY ./proxy /home/proxy

RUN chmod -R 755 /home/proxy
RUN chmod -R 777 /home/proxy-pool/logs

WORKDIR /home/proxy
EXPOSE 80

USER app-runner
CMD  ["./proxy"]