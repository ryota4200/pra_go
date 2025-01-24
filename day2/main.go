package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate", calculateHandler)

	log.Println("サーバーを起動します http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッド", http.StatusMethodNotAllowed)
		return
	}

	num1str := r.FormValue("num1")
	num2str := r.FormValue("num2")
	operation := r.FormValue("operation")

	num1, err1 := strconv.Atoi(num1str)
	num2, err2 := strconv.Atoi(num2str)

	if err1 != nil || err2 != nil {
		http.Error(w, "無効な入力です", http.StatusBadRequest)
		return
	}

	var result int
	switch operation {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		result = num1 / num2
	default:
		http.Error(w, "無効な演算子です", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
        <!DOCTYPE html>
        <html>
        <head><title>計算結果</title></head>
        <body>
            <h1>計算結果</h1>
            <p>` + strconv.Itoa(num1) + ` ` + operation + ` ` + strconv.Itoa(num2) + ` = ` + strconv.Itoa(result) + `</p>
            <a href="/">戻る</a>
        </body>
        </html>
    `))

}
