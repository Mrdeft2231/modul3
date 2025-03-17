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

  fetch("/Auth", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(formObject)
  })
    .then(async response => {
      if (!response.ok) {
        let errorMessage = "Неправильный логин или пароль"
        try {
          const errorData = await response.json();
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.warn("Сервер вернул не-JSON ошибку")
        }
        throw new Error(errorMessage)
      }
      if (response.status === 204) return {}

      return response.json()
    })
      .then(data => {
        if (data.redirect) {
          window.location.href = data.redirect;  // Перенаправляем пользователя
        }
      })
    .catch(error => alert(error.message));
});



document.addEventListener("DOMContentLoaded", function () {
  checkInputs();
});

