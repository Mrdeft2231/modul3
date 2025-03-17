document.addEventListener("DOMContentLoaded", () => {
    // Запрос к серверу для получения списка пользователей
    fetch("http://localhost:8080/GetUser")  // Замените на свой эндпоинт
        .then(response => response.json())  // Преобразуем ответ в JSON
        .then(data => {
            console.log(data)
            const login = data.login;  // Получаем массив пользователей из ответа

            // Получаем тело таблицы
            const tableBody = document.getElementById("user");
            tableBody.innerHTML = "";  // Очищаем таблицу перед заполнением

                // Создаём строку таблицы
                const row = document.createElement("h1");
                row.innerHTML = `
                 ${login}
                `;

                // Добавляем строку в таблицу
                tableBody.appendChild(row);
        })
        .catch(error => console.error("Ошибка загрузки пользователей:", error));
});

function deleteCookie() {
    console.log("работает куки")
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; domain=localhost";
}

document.getElementById("LeaveButton").addEventListener("click", () => {
    console.log("работает кнопка")
    deleteCookie()
})