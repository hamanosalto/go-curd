package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"os"
)

/*
  =========================
  エントリーポイント
  =========================
*/
func main() {

	/* -------------------------
			起動時処理
	   ------------------------- */

	// JSONファイルから Todo を読み込む
	LoadTodos()

	/* -------------------------
			API ルーティング
	   ------------------------- */

	// Todo 一覧取得・追加
	http.HandleFunc("/todos", todosHandler)

	// Todo 個別操作（取得・更新・削除）
	http.HandleFunc("/todos/", todoHandler)

	/* -------------------------
			フロントエンド配信
	   ------------------------- */

	// 実行ディレクトリ取得
	cwd, _ := os.Getwd()

	// frontend ディレクトリを公開
	frontendDir := filepath.Join(cwd, "frontend")
	fs := http.FileServer(http.Dir(frontendDir))

	// ルートアクセス時は index.html を返す
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		}
		// CSS / JS などの静的ファイル配信
		fs.ServeHTTP(w, r)
	})

	/* -------------------------
			サーバー起動
	   ------------------------- */

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
  =========================
  /todos ハンドラ
  =========================
  ・GET  : Todo 一覧取得
  ・POST : Todo 追加
*/
func todosHandler(w http.ResponseWriter, r *http.Request) {

	// JSON を返すことを明示
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// Todo 一覧取得
	case http.MethodGet:
		json.NewEncoder(w).Encode(GetTodos())

	// Todo 追加
	case http.MethodPost:
		var todo Todo

		// リクエストボディを構造体に変換
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Todo 追加処理
		newTodo := AddTodo(todo)

		// 追加後の Todo を返す
		json.NewEncoder(w).Encode(newTodo)

	// 許可されていないメソッド
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/*
  =========================
  /todos/{id} ハンドラ
  =========================
  ・GET    : 個別 Todo 取得
  ・PUT    : 更新
  ・DELETE : 削除
*/
func todoHandler(w http.ResponseWriter, r *http.Request) {

	// URL から ID 部分を取得
	idStr := r.URL.Path[len("/todos/"):]

	// ID を元に Todo を検索
	todo, id, found := GetTodoByID(idStr)

	// POST 以外で存在しない場合は 404
	if !found && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {

	// 個別 Todo 取得
	case http.MethodGet:
		json.NewEncoder(w).Encode(todo)

	// Todo 更新
	case http.MethodPut:
		var updated Todo

		// リクエストボディを構造体に変換
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// 更新処理
		if t, ok := UpdateTodo(id, updated); ok {
			json.NewEncoder(w).Encode(t)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	// Todo 削除
	case http.MethodDelete:
		if DeleteTodo(id) {
			w.WriteHeader(http.StatusNoContent) // 成功（レスポンスなし）
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	// 許可されていないメソッド
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
