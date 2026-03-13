FROM golang:1.25.5-alpine

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 安装 goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main cmd/server/main.go

# 复制启动脚本并设置执行权限
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./entrypoint.sh"]
