FROM node:12.19.1-slim as node-builder
ARG NPM_REGISTRY=https://registry.npm.taobao.org
ARG GOPROXY="https://mirrors.aliyun.com/goproxy/"
WORKDIR /vue-admin
COPY web/vue-admin .
RUN rm -rf vue.config.js &&\
    npm config set registry "$NPM_REGISTRY" && \
    npm install --unsafe-perm --registry=https://registry.npm.taobao.org --sass_binary_site=https://npm.taobao.org/mirrors/node-sass
RUN npm run build --scripts-prepend-node-path=auto

FROM golang:1.16.4-buster AS go-builder

WORKDIR /alertmanager_notifier
COPY . .
COPY --from=node-builder /vue-admin/dist web/dist
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=1 GOOS=linux go build -a -trimpath \
    -o bin/alertmanager_notifier \
    -ldflags "-X alertmanager_notifier/pkg/version.Version=`cat VERSION` -X alertmanager_notifier/pkg/version.Revision=`git rev-parse HEAD` -X alertmanager_notifier/pkg/version.Branch=`git rev-parse --abbrev-ref HEAD` -X alertmanager_notifier/pkg/version.BuildUser=`whoami` -X alertmanager_notifier/pkg/version.BuildDate=`date +%Y%m%d-%H:%M:%S`"  \
    cmd/alertmanager_notifier/main.go

FROM alpine:3.12

LABEL MAINTAINER=lyle.lai

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /opt/glibc
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.29-r0/glibc-2.29-r0.apk && \
    apk add glibc-2.29-r0.apk

WORKDIR /opt/alertmanager_notifier

COPY --from=go-builder /alertmanager_notifier/bin /opt/alertmanager_notifier/bin
#COPY --from=node-builder /vue-admin/dist/ /usr/share/nginx/html/
#COPY nginx-default.conf /etc/nginx/conf.d/default.conf
RUN mkdir /opt/alertmanager_notifier/logs && \
    mkdir /opt/alertmanager_notifier/data && \
    mkdir /opt/alertmanager_notifier/conf && \
    rm -rf /opt/glibc
COPY conf/settings.yaml conf/settings.yaml

EXPOSE 8080

ENTRYPOINT ["./bin/alertmanager_notifier"]