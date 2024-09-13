# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Go Modulesをキャッシュするために事前にコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN go build -o line-notify main.go

# Final stage
FROM alpine:latest
WORKDIR /root/

# 必要なファイルをビルドイメージからコピー
COPY --from=builder /app/line-notify .

# 実行時に必要な設定ファイルなどをコピー (.envなど)
COPY .env .

# テンプレートディレクトリをコピー
COPY templates/ ./templates/

# ポートが80でアプリケーションが起動する場合の指定
EXPOSE 80

# アプリケーション実行
CMD ["./line-notify"]
