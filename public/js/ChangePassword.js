const oldPassword = document.getElementById("oldPassword")
const newPassword = document.getElementById("newPassword")
const confirmPassword = document.getElementById("confimPassword")

const button = document.getElementById("ChangeButton")

button.addEventListener("click", (event) => {
    event.preventDefault()

    const formObject = {
        Password: newPassword.value,
        OldPassword: oldPassword.value
    }

    const ConfirmPassword = confirmPassword.value
    const New = newPassword.value

    console.log("новый пароль", parseInt(ConfirmPassword), "старый пароль", parseInt(New))
    if (parseInt(ConfirmPassword) !== parseInt(New)) {
        alert("Пароли не совпадают")
        return
    }


    fetch("/ChangePassword", {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(formObject)
    })
        .then(async response => {
            // Если ответ не OK, пробуем получить JSON-ошибку
            if (!response.ok) {
                let errorMessage = "Ошибка при смене пароля";
                try {
                    const errorData = await response.json(); // Попробуем прочитать JSON
                    errorMessage = errorData.error || errorMessage;
                } catch (e) {
                    console.warn("Сервер вернул не-JSON ошибку");
                }
                throw new Error(errorMessage);
            }

            // Если ответ пустой, не пытаемся его парсить
            if (response.status === 204) return {};

            return response.json();
        })
        .then(data => {
            alert("Успешная смена пароля");

            if (data.redirect) {
                window.location.href = data.redirect;  // Перенаправляем пользователя
            }
        })
        .catch(error => alert(error.message));

})