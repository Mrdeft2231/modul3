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

  const formObject = {
    Login: auth.value,
    Email: emailAuth.value,
    Password: passwordAuth.value
  }

  console.log(formObject)

  fetch("/createAuth", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(formObject)
  })
    .then(response => response.json())
    .then(data => console.log("Ответ сервера:", data))
    .catch(error => console.error("Ошибка:", error));
});



document.addEventListener("DOMContentLoaded", function () {
  checkInputs();
});

