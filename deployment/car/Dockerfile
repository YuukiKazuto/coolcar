# 启动编译环境
FROM golang:1.18-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中
COPY . /go/src/coolcar/server

# 编译
WORKDIR /go/src/coolcar/server
RUN go install ./car/...

# 安装grpc-health-probe
RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.11

FROM alpine:3.15
COPY --from=builder /go/bin/car /bin/car
COPY --from=builder /go/bin/grpc-health-probe /bin/grpc-health-probe
ENV ADDR=:8081
ENV WS_ADDR=:8080

# 申明暴露的端口
EXPOSE 8080
EXPOSE 8081

# 设置服务入口
ENTRYPOINT [ "/bin/car" ]
