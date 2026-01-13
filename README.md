# Go Todo App

Go（標準ライブラリのみ）で作成した、  
**シンプルな Todo 管理アプリ（CRUD）** です。

ブラウザから操作でき、Todo は JSON ファイルに永続化されます。

---

## ✨ 機能

- Todo の追加
- Todo の編集（タイトル変更）
- 完了 / 未完了の切り替え
- 完了日時の自動記録
- Todo の削除（複数選択可）
- 追加日・完了日の表示
- JSON ファイルによるデータ保存

---

## 🛠 使用技術

### バックエンド
- Go
- net/http
- encoding/json

### フロントエンド
- HTML
- CSS
- JavaScript（Vanilla JS）

※ フレームワーク未使用

---

## 📂 ディレクトリ構成

go-curd/

├── main.go # サーバー起動・ルーティング

├── todo.go # Todo のデータ管理・CRUD 処理

├── frontend/

│ ├── index.html # 画面（HTML）

│ ├── style.css # スタイル（CSS）

│ └── script.js # フロントのロジック（JavaScript）

├── todos.json # データ保存用（※ Git 管理外）

├── .gitignore

└── README.md

---

🚀 起動方法

1️⃣ リポジトリをクローン

```bash
git clone https://github.com/hamanosalto/go-curd.git
cd go-curd

2️⃣ サーバー起動
go run .

3️⃣ ブラウザでアクセス
http://localhost:8080

🗂️ データ構造（Todo）
{
  "id": 1,
  "title": "やること",
  "completed": false,
  "created_at": "2026-01-13T16:39:15+09:00",
  "completed_at": null
}

| フィールド        | 説明               |
| ------------ | ---------------- |
| id           | Todo の一意な ID     |
| title        | タイトル             |
| completed    | 完了状態             |
| created_at   | 追加日時             |
| completed_at | 完了日時（未完了時は null） |

🛠 使用技術
・Go（標準ライブラリのみ）
・net/http
・HTML / CSS / JavaScript（Vanilla JS）
・JSON ファイルによる簡易永続化