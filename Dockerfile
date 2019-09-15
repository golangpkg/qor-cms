# 使用 as build 方式创建最小镜像部署。
# https://blog.csdn.net/freewebsys/article/details/80177036

FROM golang:alpine AS build-env

ENV APP_DIR /go/src/github.com/golangpkg/qor-cms
RUN mkdir -p $APP_DIR

ADD . $APP_DIR

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
# RUN cd $APP_DIR && go build -ldflags -s -a -installsuffix cgo main.go && ls -lh
RUN cd $APP_DIR && go build -ldflags '-d -w -s' && ls -lh

# 拷贝编译后的文件。
FROM alpine

ENV APP_DIR /go/src/github.com/golangpkg/qor-cms
WORKDIR /app
COPY --from=build-env $APP_DIR/qor-cms /app
COPY --from=build-env $APP_DIR/app /app/app
COPY --from=build-env $APP_DIR/conf/app_docker.conf /app/conf/app.conf
COPY --from=build-env $APP_DIR/static /app/static
COPY --from=build-env $APP_DIR/views /app/views

RUN ls /app

EXPOSE 9000

ENTRYPOINT (cd /app && ./qor-cms)
