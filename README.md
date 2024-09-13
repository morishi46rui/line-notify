# Line Notify Webアプリケーション

Webフォームを通じてLINE Notifyにメッセージを送信できるシンプルなWebアプリケーションです。
Go言語で作成されており、Dockerを使用して簡単にセットアップと実行が可能です。

## 機能
- Webフォームを通じてLINE Notifyにメッセージを送信
- 簡単なインターフェースと入力バリデーション
- Dockerでの簡単なセットアップとデプロイ

## 必要なもの
- Docker
- Docker Compose
- LINE Notify アクセストークン

## セットアップ手順

### 1. リポジトリをクローン
```
git clone https://github.com/morishi46rui/line-notify.git
cd line-notify
```

### 2. LINE Notifyのアクセストークンの取得
https://notify-bot.line.me/ja/ にアクセスしログイン。
マイページからトークルームを選択し、アクセストークンを発行。

### 3. .envファイルの設定
プロジェクトのルートディレクトリに .env ファイルを作成し、以下の内容を記述してください。
your-line-notify-access-token は、実際のLINE Notifyのアクセストークンに置き換えてください。
```
LINE_NOTIFY_ACCESS_TOKEN=your-line-notify-access-token
```

### 4. Dockerでアプリケーションをビルドして実行
以下のコマンドを実行して、Dockerイメージをビルドし、アプリケーションを起動します。
docker-compose up --build

### 5. アプリケーションにアクセス
ブラウザを開いて、以下のURLにアクセスしてください。
http://localhost:8080

メッセージを入力して、LINE Notifyに送信するためのシンプルなフォームが表示されます。

## アプリケーションの動作
1. ユーザーがフォームにメッセージを入力します。
2. メッセージがPOSTリクエストで/sendエンドポイントに送信されます。
3. サーバーはリクエストを処理し、.envファイルに設定されたアクセストークンを使用してLINE Notifyにメッセージを送信します。

## ファイル構成
line-notify/
│
├── Dockerfile                 # GoアプリケーションのためのDockerfile
├── docker-compose.yml          # Docker Composeの設定ファイル
├── go.mod                      # Goモジュールの設定ファイル
├── go.sum                      # Goモジュールの依存関係ファイル
├── main.go                     # メインのGoアプリケーションファイル
└── .env                        # 環境変数の設定ファイル (.gitには含まれません)

## 環境変数
このアプリケーションでは、以下の環境変数を使用します。
- LINE_NOTIFY_ACCESS_TOKEN: LINE Notify APIのトークン
