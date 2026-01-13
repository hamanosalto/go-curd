package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

/*
  =========================
  Todo 構造体
  =========================
*/
type Todo struct {
	ID          int        `json:"id"`            // 一意なID
	Title       string     `json:"title"`         // タイトル
	Completed   bool       `json:"completed"`     // 完了フラグ
	CreatedAt   time.Time  `json:"created_at"`    // 追加日時
	CompletedAt *time.Time `json:"completed_at"`  // 完了日時（未完了時は null）
}

/*
  =========================
  グローバル変数
  =========================
*/
var todos []Todo          // Todo 一覧
var nextID int = 1        // 次に採番するID

const todoFile = "todos.json" // 永続化用 JSON ファイル

/*
  =========================
  起動時処理
  =========================
*/
func LoadTodos() {

	// JSONファイルが存在しない場合は新規作成
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		todos = []Todo{}
		saveFile()
		log.Println("todos.json がなかったので新規作成しました")
		return
	}

	// JSONファイル読み込み
	file, err := os.ReadFile(todoFile)
	if err != nil {
		log.Fatal(err)
	}

	// JSON → 構造体へ変換
	if err := json.Unmarshal(file, &todos); err != nil {
		log.Fatal(err)
	}

	/*
		nextID の決定と
		既存データの補完処理
	*/
	maxID := 0
	updated := false

	for i, t := range todos {

		// 最大IDを取得
		if t.ID > maxID {
			maxID = t.ID
		}

		// 完了済みなのに CompletedAt が無い場合は補完
		if t.Completed && t.CompletedAt == nil {
			now := time.Now()
			todos[i].CompletedAt = &now
			updated = true
		}
	}

	// 次に使うIDを設定
	nextID = maxID + 1

	// 補完があった場合は保存
	if updated {
		saveFile()
	}
}

/*
  =========================
  JSON保存（内部用）
  =========================
*/
func saveFile() {

	// 構造体 → JSON（整形あり）
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		log.Println("Error marshaling todos:", err)
		return
	}

	// ファイル書き込み
	if err := os.WriteFile(todoFile, data, 0644); err != nil {
		log.Println("Error writing file:", err)
	}
}

/*
  =========================
  外部向け保存関数
  =========================
*/
func SaveTodos() {
	saveFile()
}

/*
  =========================
  CRUD 操作
  =========================
*/

// Todo 一覧取得
func GetTodos() []Todo {
	return todos
}

// Todo 追加
func AddTodo(todo Todo) Todo {

	// ID 採番
	todo.ID = nextID
	nextID++

	// 追加日時を自動設定
	todo.CreatedAt = time.Now()

	todos = append(todos, todo)
	saveFile()

	return todo
}

// Todo 更新
func UpdateTodo(id int, updated Todo) (Todo, bool) {

	for i, t := range todos {
		if t.ID == id {

			// ID と追加日時は保持
			updated.ID = id
			updated.CreatedAt = t.CreatedAt

			/*
				完了状態の変化に応じて
				完了日時を制御
			*/
			if !t.Completed && updated.Completed {
				// 未完了 → 完了
				now := time.Now()
				updated.CompletedAt = &now

			} else if t.Completed && !updated.Completed {
				// 完了 → 未完了
				updated.CompletedAt = nil

			} else {
				// 状態変化なし
				updated.CompletedAt = t.CompletedAt
			}

			todos[i] = updated
			saveFile()

			return updated, true
		}
	}

	return Todo{}, false
}

// Todo 削除
func DeleteTodo(id int) bool {

	for i, t := range todos {
		if t.ID == id {

			// 対象要素を除外
			todos = append(todos[:i], todos[i+1:]...)
			saveFile()

			return true
		}
	}

	return false
}

// ID 文字列から Todo を取得
func GetTodoByID(idStr string) (Todo, int, bool) {

	// ID を数値に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return Todo{}, 0, false
	}

	// Todo 検索
	for _, t := range todos {
		if t.ID == id {
			return t, id, true
		}
	}

	return Todo{}, id, false
}
