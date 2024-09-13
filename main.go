package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var lineNotifyAPI = "https://notify-api.line.me/api/notify"

func main() {
	// .envファイルの読み込み
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("環境変数ファイル(.env)の読み込みに失敗しました")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/send", sendHandler)
	fmt.Println("サーバーが http://localhost:8080 で起動しています")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 絶対パスを確認
	path, err := filepath.Abs("templates/index.html")
	if err != nil {
		log.Fatalf("テンプレートファイルへのパスを取得できませんでした: %v", err)
	}

	// パスを出力
	log.Printf("テンプレートファイルの絶対パス: %s", path)

	// テンプレートファイルの読み込み
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "テンプレートファイルが見つかりません", http.StatusInternalServerError)
		return
	}

	// テンプレートの実行
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "テンプレートのレンダリングに失敗しました", http.StatusInternalServerError)
		return
	}
}


func sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "メッセージを入力してください", http.StatusBadRequest)
		return
	}

	accessToken := os.Getenv("LINE_NOTIFY_ACCESS_TOKEN")
	if accessToken == "" {
		http.Error(w, "LINE Notifyのアクセストークンが見つかりません", http.StatusInternalServerError)
		return
	}

	// LINE Notifyにメッセージを送信
	err := sendLineNotify(message, accessToken)
	if err != nil {
		http.Error(w, "メッセージの送信に失敗しました", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func sendLineNotify(message, token string) error {
	data := url.Values{}
	data.Set("message", message)

	// リクエストの作成
	req, err := http.NewRequest("POST", lineNotifyAPI, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	// 必要なヘッダーを設定
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスステータスを表示
	log.Printf("レスポンスステータス: %s", resp.Status)

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("レスポンスの読み込みに失敗しました: %v", err)
	}
	log.Printf("レスポンスボディ: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("通知の送信に失敗しました: %s", resp.Status)
	}

	return nil
}
