let action = "send"

const actionSend = "send"
const actionSignIn = "signin"

const msgSendOTP = "Please provide the email address to which we will send the password."
const msgEnterPassword = "Enter the password that came to your email."
const msgError = "Error. Try again later."

const formPassword = document.querySelector("#form_password")
const formMessage = document.querySelector("#form_message")
const formAnotherEmail = document.querySelector("#form_another_email")
const formError = document.querySelector("#form_error")
const submit = document.querySelector("#submit")

const email = document.querySelector("#email")
const password = document.querySelector("#password")

function loader(show) {
    document.querySelector("#loader").style.display = show == true ? "block" : "none"
    submit.disabled = show == true
}


async function sendEmail(e) {
    loader(true)
    formError.innerText = ""
    res = await fetch("/api/v0/auth/otp", {
        method: "POST",
        body: JSON.stringify({
            "email": email.value
        })
    })
    if (res.status == 200) {
        json = await res.json()
        if (json.status == true) {
            formPassword.style.display = "block"
            formAnotherEmail.style.display = "block"
            formMessage.innerText = msgEnterPassword
            action = actionSignIn
        }
    } else {
        formError.innerText = msgError
    }
    loader(false)
}

async function sendOTP() {
    loader(true)
    formError.innerText = ""
    res = await fetch("/api/v0/auth/signIn", {
        method: "POST",
        body: JSON.stringify({
            "email": email.value,
            "password": password.value
        })
    })
    if (res.status == 200) {
        json = await res.json()
        if (json.status == true) {
            document.location.href = "/"
        }
    }else {
        formError.innerText = msgError
    }
    loader(false)
}

function startForm() {
    formPassword.style.display = "none"
    formAnotherEmail.style.display = "none"
    formMessage.innerText = msgSendOTP
    action = actionSend
}

formAnotherEmail.addEventListener("click", startForm)


submit.addEventListener("click", function(e) {
    e.preventDefault()
    if (action == actionSend) {
        sendEmail(e)
        return
    }
    if (action == actionSignIn) {
        sendOTP(e)
        return
    }
})