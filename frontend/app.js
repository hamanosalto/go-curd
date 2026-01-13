/* =========================
  API設定
   ========================= */

/* Todo API のエンドポイント */
const apiUrl = '/todos';

/* =========================
  日付フォーマット
   ========================= */

/*
  ISO形式の日時文字列を
  「YYYY/MM/DD HH:mm」形式に変換する
*/
function formatDate(iso) {
  if (!iso) return '';           // null / undefined 対策

  const d = new Date(iso);
  if (isNaN(d)) return '';       // 不正な日付対策

  return d.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
}

/* =========================
    Todo 一覧取得・描画
   ========================= */

/*
  サーバーから Todo 一覧を取得し
  テーブルに描画する
*/
async function fetchTodos() {
  const res = await fetch(apiUrl);
  const todos = await res.json();

  const list = document.getElementById('todo-list');
  list.innerHTML = ''; // 既存行を全削除

  todos.forEach(todo => {
    const tr = document.createElement('tr');

    /* 行ID（削除・更新時に使用） */
    tr.id = 'todo-' + todo.id;

    /* 完了状態ならスタイルを変更 */
    tr.className = todo.completed ? 'completed' : '';

    /* 行の中身を生成 */
    tr.innerHTML = `
      <!-- 追加日 -->
      <td class="date-column">
        ${formatDate(todo.created_at)}
      </td>

      <!-- 完了日 -->
      <td class="done-date-column">
        ${formatDate(todo.completed_at)}
      </td>

      <!-- 削除チェック -->
      <td class="delete-column">
        <input type="checkbox" class="delete-checkbox">
      </td>

      <!-- 完了チェック -->
      <td class="complete-column">
        <input type="checkbox"
              class="complete-checkbox"
              ${todo.completed ? 'checked' : ''}
              onchange="toggleComplete(${todo.id}, this.checked)">
      </td>

      <!-- タイトル（インライン編集） -->
      <td class="title-column">
        <span contenteditable="true"
              onkeydown="preventEnter(event)"
              onblur="editTodo(${todo.id}, this.textContent)">
          ${todo.title}
        </span>
      </td>
    `;

    list.appendChild(tr);
  });
}

/* =========================
    Todo 追加
   ========================= */

/*
  入力欄の内容をサーバーに送信し
  新しい Todo を追加する
*/
async function addTodo() {
  const input = document.getElementById('new-todo');
  if (!input.value.trim()) return; // 空文字防止

  await fetch(apiUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: input.value,
      completed: false
    })
  });

  input.value = ''; // 入力欄クリア
  fetchTodos();     // 再描画
}

/* =========================
    タイトル編集
   ========================= */

/*
  タイトル編集確定時に
  サーバーへ更新リクエストを送る
*/
async function editTodo(id, title) {
  const tr = document.getElementById('todo-' + id);

  /* 現在の完了状態を保持 */
  const completed = tr.querySelector('.complete-checkbox').checked;

  await fetch(`${apiUrl}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, completed })
  });

  fetchTodos();
}

/* =========================
    完了状態の切り替え
   ========================= */

/*
  完了チェックボックスのON/OFFを
  サーバーに反映する
*/
async function toggleComplete(id, completed) {
  const tr = document.getElementById('todo-' + id);
  const title = tr.querySelector('span').textContent;

  await fetch(`${apiUrl}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, completed })
  });

  fetchTodos();
}

/* =========================
    選択削除
   ========================= */

/*
  削除チェックが入っている Todo を
  まとめて削除する
*/
async function deleteSelected() {
  const checkboxes = document.querySelectorAll('.delete-checkbox:checked');

  for (const cb of checkboxes) {
    const id = cb.closest('tr').id.replace('todo-', '');
    await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
  }

  fetchTodos();
}

/* =========================
    編集時の Enter 無効化
   ========================= */

/*
  contenteditable の改行を防ぎ
  Enterキーで編集確定にする
*/
function preventEnter(e) {
  if (e.key === 'Enter') {
    e.preventDefault();
    e.target.blur();
  }
}

/* =========================
    初期表示
   ========================= */

/* ページ読み込み時に Todo を取得 */
fetchTodos();
