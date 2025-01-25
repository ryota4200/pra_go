package main

import (
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
}
