const login = document.getElementById("loginAuth")
const email = document.getElementById("emailAuth")
const password = document.getElementById("passwordAuth")
const roleSelect = document.getElementById('Role');
const button = document.getElementById("CreateButtonUsers")

button.addEventListener("click", (event) => {
    event.preventDefault()

    const formObject = {
        Login: login.value,
        Email: emailAuth.value,
        Password: passwordAuth.value,
        Role: roleSelect.value
    }
    console.log(formObject)
    if (formObject.Role === "") {
        alert("Выберите роль")
        return
    }

    console.log(formObject)

    fetch("/CreateAuth", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(formObject)
    })
        .then(async response => {
            if (response.ok) {
                location.reload();  // Перезагружаем страницу, чтобы обновить таблицу
            }

            let errorMessage = "Произошла ошибка";
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorMessage;
            } catch (e) {
                console.warn("Сервер вернул не-JSON ошибку", e);
            }
            throw new Error(errorMessage);

            if (response.status === 204) return {}

            return response.json()
        })
        .then(data => {
            alert(data)
        })
        .catch(error => alert(error.message));
});
