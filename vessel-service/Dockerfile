# 我们使用Golang的官方镜像，它包含了所有用于构建Golang应用的依赖和工具
# 请注意`as builder`，这命名了我们这个镜像，以后可以用`builder`来指向它
FROM golang:alpine as builder

# 创建一个文件夹，即 /app，并且这个文件夹就是我们的工作目录
RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

# 将代码复制到工作目录中
COPY . .

# 下载依赖包
RUN go mod download

# 编译我们的二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vessel-service

# 在这里，我们第二次使用了FROM，它告诉Docker，我们从这里要开始第二个构建流程
# 最新版的Linux Alpine。Linux Alpine是一个轻量的Linux发行版，专门为运行容器化的应用而生
FROM alpine:latest

# 确保安全性的包
RUN apk --no-cache add ca-certificates

# 创建一个文件夹，即 /app，并且这个文件夹就是我们的工作目录
RUN mkdir /app
WORKDIR /app
# 在这里，我们不再直接从我们的宿主机器中直接拷贝二进制文件,
# 而是从一个叫做`builder`的容器中获取。它会从我们起先构建的镜像中
# 获得已经编译好的文件并引入到这个容器里。
COPY --from=builder /app/vessel-service .

CMD ["./vessel-service"]