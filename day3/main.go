package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type TodoItem struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	todoList  []TodoItem
	idCounter int
	mutex     sync.Mutex
)

func main() {
	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ルートとAPIエンドポイントの設定
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/todos", todosHandler)
	http.HandleFunc("/api/export", exportHandler)

	log.Print("サーバを起動します http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ホームページのハンドラー
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// TODOアプリを管理するAPIハンドラー
func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODOリストを取得
		mutex.Lock()
		defer mutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todoList)
	case http.MethodPost:
		// TODOリストにアイテム追加
		var newTodo TodoItem
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			http.Error(w, "無効なデータ", http.StatusBadRequest)
			return
		}
		mutex.Lock()
		defer mutex.Unlock()
		idCounter++
		newTodo.ID = idCounter
		todoList = append(todoList, newTodo)
		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete:
		// TODOアイテムの削除
		var toDelete struct {
			ID int `json:"ID"`
		}
		err := json.NewDecoder(r.Body).Decode(&toDelete)
		if err != nil {
			http.Error(w, "無効なデータ", http.StatusBadRequest)
			return
		}
		mutex.Lock()
		defer mutex.Unlock()
		for i, item := range todoList {
			if item.ID == toDelete.ID {
				todoList = append(todoList[:i], todoList[i+1:]...)
				break
			}
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "無効なリクエスト", http.StatusMethodNotAllowed)

	}
}

// TODOリストをエクスポポート
func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "無効なメソッド", http.StatusMethodNotAllowed)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachemnt; filename="todo_list.json"`)
	json.NewEncoder(w).Encode(todoList)

}
