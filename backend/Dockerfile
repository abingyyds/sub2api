FROM golang:1.24-alpine

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git postgresql-client

# 安装 goose 和 wire
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/google/wire/cmd/wire@latest

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 生成 Ent ORM 代码（包含 PaymentOrder 等实体）
RUN go generate ./ent/...

# 生成 Wire 依赖注入代码
RUN cd cmd/server && wire

# 构建应用
RUN go build -o main cmd/server/main.go

# 复制启动脚本并设置执行权限
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./entrypoint.sh"]
