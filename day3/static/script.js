document.addEventListener("DOMContentLoaded", () =>{
    const todoForm = document.getElementById("todoForm");
    const todoInput = document.getElementById("todoInput");
    const todoList = document.getElementById("todoList");
    const exportBtn = document.getElementById("exportBtn");

    // サーバからTODOリストを取得
    const fetchTodos = async () => {
        const response = await fetch("/api/todos");
        const todos = await response.json();
        renderTodos(todos);
    }

    // TODOリストを表示
    const renderTodos = (todos) => {
        todoList.innerHTML = "";
        todos.forEach((todo) => {
            const li = document.createElement("li");
            li.textContent = todo.text;

            const deleteButton = document.createElement("button");
            deleteButton.textContent = "削除";
            deleteButton.onclick = async () => {
                await fetch("/api/todos", {
                    method: "DELETE",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({ id: todo.id})

                });
                fetchTodos();
            };
            li.appendChild(deleteButton);
            todoList.appendChild(li);
        });
    };
    // 新しいTODOを追加
    todoForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const text = todoInput.value.trim();
        if(!text) return;

        await fetch("/api/todos", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({ text }),
        });
        todoInput.value = "";
        fetchTodos();

    });

    exportBtn.addEventListener("click", () =>{
        window.location.href = "/api/export";
    });

    fetchTodos();
})
