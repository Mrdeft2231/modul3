let auth = document.getElementById("loginAuth")
let emailAuth = document.getElementById("emailAuth")
let passwordAuth = document.getElementById("passwordAuth")

function checkInputs() {
    const auth = document.getElementById("loginAuth").value;
    const emailAuth = document.getElementById("emailAuth").value;
    const passwordAuth = document.getElementById("passwordAuth").value;

    // Если все поля не пустые, активируем кнопку
    if (auth && emailAuth && passwordAuth) {
        document.getElementById("submitButton").disabled = false;
    } else {
        document.getElementById("submitButton").disabled = true;
    }
}


auth.addEventListener("input", checkInputs);
emailAuth.addEventListener("input", checkInputs);
passwordAuth.addEventListener("input", checkInputs);


const form = document.getElementById("AuthForm")
const button = document.getElementById("submitButton")

button.addEventListener("click", (event) => {
    event.preventDefault()

    data = new FormData()

    data.append("login", auth.value)
    data.append("email", emailAuth.value)
    data.append("password", passwordAuth.value)

    // const formObject = Object.fromEntries(data.entries());
    // console.log(formObject);

    fetch("http://127.0.0.1:8080/createAuth", {
        method: "POST",

        body: data
    })
        .then(response => response.json())
        .then(data => console.log("Ответ сервера:", data))
        .catch(error => console.error("Ошибка:", error));
});



document.addEventListener("DOMContentLoaded", function () {
    checkInputs();
});

