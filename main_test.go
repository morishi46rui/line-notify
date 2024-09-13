package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// LINE Notify APIをモックするためにhttptest.Serverを使用
func mockLineNotifyServer(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドの確認
		if r.Method != http.MethodPost {
			t.Errorf("期待されるメソッドはPOSTですが、実際には%sが使用されました", r.Method)
		}

		// ヘッダーの確認
		if auth := r.Header.Get("Authorization"); auth != "Bearer test_token" {
			t.Errorf("期待されるAuthorizationヘッダーは'Bearer test_token'ですが、実際には%sが使用されました", auth)
		}
		if contentType := r.Header.Get("Content-Type"); contentType != "application/x-www-form-urlencoded" {
			t.Errorf("期待されるContent-Typeは'application/x-www-form-urlencoded'ですが、実際には%sが使用されました", contentType)
		}

		// フォームデータの確認
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		expectedBody := "message=Test+message"
		if bodyString != expectedBody {
			t.Errorf("期待されるボディは'%s'ですが、実際には'%s'が送信されました", expectedBody, bodyString)
		}

		// 成功レスポンスを返す
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":200,"message":"ok"}`))
	})
	return httptest.NewServer(handler)
}

// sendHandler関数の正常動作時のテスト
func TestSendHandlerSuccess(t *testing.T) {
	// LINE Notify APIのモックを作成
	ts := mockLineNotifyServer(t)
	defer ts.Close()

	// モックサーバーのURLにlineNotifyAPIを上書き
	originalLineNotifyAPI := lineNotifyAPI
	lineNotifyAPI = ts.URL
	defer func() { lineNotifyAPI = originalLineNotifyAPI }()

	// アクセストークン用の環境変数を設定
	os.Setenv("LINE_NOTIFY_ACCESS_TOKEN", "test_token")

	// フォームデータの準備
	form := url.Values{}
	form.Add("message", "Test message")

	// POSTリクエストの作成
	req, err := http.NewRequest("POST", "/send", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// レスポンスの記録
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sendHandler)

	// ハンドラを実行
	handler.ServeHTTP(rr, req)

	// ステータスコードの確認
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("期待されるステータスコードは%dですが、実際には%dが返されました", http.StatusSeeOther, status)
	}

	// リダイレクトの確認
	if location := rr.Header().Get("Location"); location != "/" {
		t.Errorf("期待されるリダイレクト先は'/'ですが、実際には'%s'にリダイレクトされました", location)
	}
}

// メッセージが無い場合のsendHandler関数のテスト
func TestSendHandlerMissingMessage(t *testing.T) {
	// メッセージ無しのPOSTリクエストを作成
	req, err := http.NewRequest("POST", "/send", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// レスポンスの記録
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sendHandler)

	// ハンドラを実行
	handler.ServeHTTP(rr, req)

	// ステータスコードの確認
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("期待されるステータスコードは%dですが、実際には%dが返されました", http.StatusBadRequest, status)
	}
}

// sendLineNotify関数のテスト
func TestSendLineNotify(t *testing.T) {
	// LINE Notify APIのモックを作成
	ts := mockLineNotifyServer(t)
	defer ts.Close()

	// モックサーバーのURLにlineNotifyAPIを上書き
	originalLineNotifyAPI := lineNotifyAPI
	lineNotifyAPI = ts.URL
	defer func() { lineNotifyAPI = originalLineNotifyAPI }()

	// 関数の実行
	err := sendLineNotify("Test message", "test_token")
	if err != nil {
		t.Errorf("エラーが発生しないことを期待しましたが、%vが発生しました", err)
	}
}

// 無効なトークンを使用したsendLineNotify関数のテスト
func TestSendLineNotifyInvalidToken(t *testing.T) {
	// エラーレスポンスを返すモックサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 認証エラーを返す
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":401,"message":"Invalid access token"}`))
	}))
	defer ts.Close()

	// モックサーバーのURLにlineNotifyAPIを上書き
	originalLineNotifyAPI := lineNotifyAPI
	lineNotifyAPI = ts.URL
	defer func() { lineNotifyAPI = originalLineNotifyAPI }()

	// 関数の実行
	err := sendLineNotify("Test message", "invalid_token")
	if err == nil {
		t.Errorf("無効なトークンによるエラーが発生することを期待しましたが、エラーが発生しませんでした")
	}
}
