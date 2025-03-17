document.addEventListener("DOMContentLoaded", () => {
    // Запрос к серверу для получения списка пользователей
    fetch("http://localhost:8080/GetUsers")  // Замените на свой эндпоинт
        .then(response => response.json())  // Преобразуем ответ в JSON
        .then(data => {
            const users = data.users;  // Получаем массив пользователей из ответа

            // Получаем тело таблицы
            const tableBody = document.getElementById("userTableBody");
            tableBody.innerHTML = "";  // Очищаем таблицу перед заполнением

            // Проходим по каждому пользователю и создаём строку таблицы
            users.forEach(user => {
                const createdAt = new Date(user.CreateUser);  // Преобразуем строку в объект Date
                const currentDate = new Date();

                const endDate = new Date(createdAt);
                endDate.setDate(endDate.getDate() + 30);  // Прибавляем 30 дней

                const timeDiff = endDate - currentDate;  // Разница во времени
                const daysRemaining = Math.ceil(timeDiff / (1000 * 60 * 60 * 24));  // Переводим в дни

                // Создаём строку таблицы
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${user.Login}</td>
                    <td>${user.Email}</td>
                    <td>${user.Password}</td>
                    <td><button class="statusBlock" data-id="${user.Id}">${user.Status === 1 ? "активен ✅" : "неактивен 🚫"}</button></td>
                    <td>${user.Role}</td>
                    <td>${daysRemaining}</td>
                    <td><button class="deleteUser" data-id="${user.Id}">Удалить</button></td>
                `;

                // Добавляем строку в таблицу
                tableBody.appendChild(row);
            });

            // Добавляем обработчик события для удаления пользователей
            document.querySelectorAll(".deleteUser").forEach(button => {
                button.addEventListener("click", (event) => {
                    const userId = event.target.getAttribute("data-id");
                    deleteUser(userId);
                });
            });



            document.querySelectorAll(".statusBlock").forEach(button => {
                button.addEventListener("click", (event) => {
                    const userId = event.target.getAttribute("data-id");
                    StatusUser(userId);
                });
            });

            document.querySelectorAll()
        })
        .catch(error => console.error("Ошибка загрузки пользователей:", error));
});

// Функция для удаления пользователя
function deleteUser(Id) {
    fetch(`http://localhost:8080/DeleteUser/${Id}`, {
        method: "DELETE",
    })
        .then(response => {
            if (response.ok) {
                alert("Пользователь удалён");
                location.reload();  // Перезагружаем страницу, чтобы обновить таблицу
            } else {
                alert("Ошибка удаления");
            }
        })
        .catch(error => console.error("Ошибка удаления:", error));
}

function StatusUser(Id) {
    fetch(`http://localhost:8080/StatusPut/${Id}`, {
        method: "PUT",
    })
        .then(response => {
            if (response.ok) {
                alert("Статус пользоватеоя изменён");
                location.reload();  // Перезагружаем страницу, чтобы обновить таблицу
            } else {
                alert("Ошибка изменения статуса");
            }
        })
        .catch(error => console.error("Ошибка изменения статуса:", error));
}

function deleteCookie() {
    console.log("работает куки")
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; domain=localhost";
}

document.getElementById("LeaveButton").addEventListener("click", () => {
    console.log("работает кнопка")
    deleteCookie()
})
