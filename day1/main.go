package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logFile, err := os.OpenFile("./access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ログファイルの生成に失敗しました", err)
		return
	}
	defer logFile.Close()
	// ログ出力をファイルに設定
	log.SetOutput(logFile)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate", caluculateHandler)

	fmt.Println("Stand by server: http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>電卓アプリ</title>
		<script>
		// ボタンを押したらフォームを送信
		function calculate(op) {
			document.getElementById("operation").value = op;
			document.getElementById("calcForm").submit();
		}
		</script>
	</head>
	<body>
		<h1>電卓アプリ</h1>
		<form id="calcForm" action="/calculate" method="post">
			<label for="num1">数字1:</label>
			<input type="text" id="num1" name="num1" required>
			<br>
			<label for="num2">数字2:</label>
			<input type="text" id="num2" name="num2" required>
			<br><br>
			<input type="hidden" id="operation" name="operation">
			<button type="button" onclick="calculate('+')">+</button>
			<button type="button" onclick="calculate('-')">-</button>
		</form>
	</body>
	</html>
	`

	fmt.Fprint(w, tmpl)
}

func caluculateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("caluculateHandlerが呼び出されました")

	if r.Method != http.MethodPost {
		log.Println("無効なリクセストメソッド:", r.Method)
		http.Error(w, "無効なリクエストメソッド", http.StatusMethodNotAllowed)
		return
	}

	// フォームデータの取得
	num1Str := r.FormValue("num1")
	num2Str := r.FormValue("num2")
	operation := r.FormValue("operation")

	// 入力されたデータをログに記録
	log.Printf("入力された数値： num1=%s, num2=%s, operation=%s\n", num1Str, num2Str, operation)

	// 文字列を数値に変換
	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	if err1 != nil || err2 != nil {
		http.Error(w, "無効な入力です。通知を入力してください", http.StatusBadRequest)
		log.Println("無効な入力がありました:", err1, err2)
		return
	}

	// 計算
	var result int
	switch operation {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	default:
		log.Println("無効な演算子です", operation)
		http.Error(w, "無効な演算子です", http.StatusBadRequest)
	}

	// 計算結果をログに出力
	log.Printf("計算結果: %d %s %d = %d\n", num1, operation, num2, result)

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>計算結果</title>
	</head>
	<body>
		<h1>計算結果</h1>
		<p>{{.Num1}} {{.Operation}} {{.Num2}} = {{.Result}}</p>
		<a href="/">戻る</a>
	</body>
	</html>
	`

	data := struct {
		Num1      int
		Num2      int
		Result    int
		Operation string
	}{
		Num1:      num1,
		Num2:      num2,
		Result:    result,
		Operation: operation,
	}

	t, err := template.New("result").Parse(tmpl)
	if err != nil {
		log.Println("テンプレート解析エラー", err)
		http.Error(w, "内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "テンプレートの実行エラー", http.StatusInternalServerError)
		log.Println("テンプレート実行エラー", err)
	}

}
