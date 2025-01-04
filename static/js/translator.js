const API = "/api/v0"

const imgPlay = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-play-circle" viewBox="0 0 16 16">  <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>  <path d="M6.271 5.055a.5.5 0 0 1 .52.038l3.5 2.5a.5.5 0 0 1 0 .814l-3.5 2.5A.5.5 0 0 1 6 10.5v-5a.5.5 0 0 1 .271-.445"/></svg>'
const imgPause = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-pause-circle" viewBox="0 0 16 16">  <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>  <path d="M5 6.25a1.25 1.25 0 1 1 2.5 0v3.5a1.25 1.25 0 1 1-2.5 0zm3.5 0a1.25 1.25 0 1 1 2.5 0v3.5a1.25 1.25 0 1 1-2.5 0z"/></svg>'
const imgRecognize = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chat-dots" viewBox="0 0 16 16">  <path d="M5 8a1 1 0 1 1-2 0 1 1 0 0 1 2 0m4 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0m3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2"/>  <path d="m2.165 15.803.02-.004c1.83-.363 2.948-.842 3.468-1.105A9 9 0 0 0 8 15c4.418 0 8-3.134 8-7s-3.582-7-8-7-8 3.134-8 7c0 1.76.743 3.37 1.97 4.6a10.4 10.4 0 0 1-.524 2.318l-.003.011a11 11 0 0 1-.244.637c-.079.186.074.394.273.362a22 22 0 0 0 .693-.125m.8-3.108a1 1 0 0 0-.287-.801C1.618 10.83 1 9.468 1 8c0-3.192 3.004-6 7-6s7 2.808 7 6-3.004 6-7 6a8 8 0 0 1-2.088-.272 1 1 0 0 0-.711.074c-.387.196-1.24.57-2.634.893a11 11 0 0 0 .398-2"/></svg>'
const imgMicrophon = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-mic" viewBox="0 0 16 16"><path d="M3.5 6.5A.5.5 0 0 1 4 7v1a4 4 0 0 0 8 0V7a.5.5 0 0 1 1 0v1a5 5 0 0 1-4.5 4.975V15h3a.5.5 0 0 1 0 1h-7a.5.5 0 0 1 0-1h3v-2.025A5 5 0 0 1 3 8V7a.5.5 0 0 1 .5-.5"/><path d="M10 8a2 2 0 1 1-4 0V3a2 2 0 1 1 4 0zM8 0a3 3 0 0 0-3 3v5a3 3 0 0 0 6 0V3a3 3 0 0 0-3-3"/></svg>'
const imgRecord = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-record-fill" viewBox="0 0 16 16">  <path fill-rule="evenodd" d="M8 13A5 5 0 1 0 8 3a5 5 0 0 0 0 10"/></svg>'
const imgTranslate = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-translate" viewBox="0 0 16 16">  <path d="M4.545 6.714 4.11 8H3l1.862-5h1.284L8 8H6.833l-.435-1.286zm1.634-.736L5.5 3.956h-.049l-.679 2.022z"/>  <path d="M0 2a2 2 0 0 1 2-2h7a2 2 0 0 1 2 2v3h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-3H2a2 2 0 0 1-2-2zm2-1a1 1 0 0 0-1 1v7a1 1 0 0 0 1 1h7a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zm7.138 9.995q.289.451.63.846c-.748.575-1.673 1.001-2.768 1.292.178.217.451.635.555.867 1.125-.359 2.08-.844 2.886-1.494.777.665 1.739 1.165 2.93 1.472.133-.254.414-.673.629-.89-1.125-.253-2.057-.694-2.82-1.284.681-.747 1.222-1.651 1.621-2.757H14V8h-3v1.047h.765c-.318.844-.74 1.546-1.272 2.13a6 6 0 0 1-.415-.492 2 2 0 0 1-.94.31"/></svg>'
const imgUpload = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-cloud-upload" viewBox="0 0 16 16">  <path fill-rule="evenodd" d="M4.406 1.342A5.53 5.53 0 0 1 8 0c2.69 0 4.923 2 5.166 4.579C14.758 4.804 16 6.137 16 7.773 16 9.569 14.502 11 12.687 11H10a.5.5 0 0 1 0-1h2.688C13.979 10 15 8.988 15 7.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 2.825 10.328 1 8 1a4.53 4.53 0 0 0-2.941 1.1c-.757.652-1.153 1.438-1.153 2.055v.448l-.445.049C2.064 4.805 1 5.952 1 7.318 1 8.785 2.23 10 3.781 10H6a.5.5 0 0 1 0 1H3.781C1.708 11 0 9.366 0 7.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383"/>  <path fill-rule="evenodd" d="M7.646 4.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 5.707V14.5a.5.5 0 0 1-1 0V5.707L5.354 7.854a.5.5 0 1 1-.708-.708z"/></svg>'
const imgClose = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x" viewBox="0 0 16 16">  <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/></svg>'

let longTimeVoiceMessage = 5000 // 5 sec
let languages = {
    "RU": "ru",
    "EN": "en",
}

let timerID = null

let lang1value = "ru-RU"
let lang2value = "en-US"

const chat = document.querySelector("#chat")
const lang1 = document.querySelector("#lang1")
const btnLang1 = document.querySelector("#btnGroupDropLanguage1")
const lang2 = document.querySelector("#lang2")
const btnLang2 = document.querySelector("#btnGroupDropLanguage2")

const record1 = document.querySelector("#record1")
record1.innerHTML = imgMicrophon
record1.addEventListener("click", startRecord(record1, "record1", "click"))
// record1.addEventListener("touchstart", startRecord(record1, "record1", "start"))
// record1.addEventListener("touchend", startRecord(record1, "record1", "end"))

const record2 = document.querySelector("#record2")
record2.innerHTML = imgMicrophon
record2.addEventListener("click", startRecord(record2, "record2", "click"))
// record2.addEventListener("touchstart", startRecord(record2, "record2", "start"))
// record2.addEventListener("touchend", startRecord(record2, "record2", "end"))

const voice_played = document.querySelector("#voice_played")
const loader = document.querySelector("#loader")

function ErrorAlert() {
    let timerID = null
    const er = document.querySelector("#controler__error")
    return function(text) {
        if (timerID) {
            clearTimeout(timerID)
            timerID = null
        }
        timerID = setTimeout(function() {
            er.innerText = ""
            timerID = null
        }, 10000)
        er.innerText = text
    }
}

const setErrorAlert = ErrorAlert()

async function getBlobContent(url) {
    res = await fetch(url)
    content = await res.blob()

    return content
}

function timeoutStopRecord(e, name) {
    return function() {
        startRecord(e, name, "end")()
        setErrorAlert(`Maximum length ${longTimeVoiceMessage / 1000} seconds`)
    }
}

function recognizing(text, src, name, next) {
    return async function() {
        console.log(text, src, name)
        let content = await getBlobContent(src)
        let formData = new FormData()
        formData.append("data", content)
        formData.append("language", name == "record2" ? lang2value : lang1value)
        fetch(API+"/audio/recognize", {
            method: "POST",
            body: formData
        })
        .then(async data => {
            console.log(data)
            if (data.status == 200) {
                body = await data.json()
                text.innerText = body.text
                text.classList.remove("hide")
                if (next != null) {
                    next()
                }
                return
            }
        })
        // text.innerText = "speech to text "+src
    }
}

function translating(msg, text, name, error, next) {
    return async function() {
        if (text.innerText == "") {
            error.innerText = "message is empty"
            return
        }
        fetch(API+"/text/translate", {
            method: "POST",
            body: JSON.stringify({
                sourceLang:  name == "record2" ? lang2value : lang1value,
                targetLang: name == "record2" ? lang1value : lang2value,
                text: text.innerText
            })
        })
        .then(async data => {
            console.log(data)
            if (data.status == 200) {
                body = await data.json()
                msg.innerText = body.text
                msg.classList.remove("hide")
                if (next != null) {
                    next()
                }
                return
            }
            
            error.innerText = "something went wrong, try again later"
        })
        // msg.innerText = "translate text: "+text.innerText
    }
}

function speech(audio, btn, text, name, error, next) {
    return async function() {
        error.innerText = ""
        if (text.innerText == "") {
            error.innerText = "message is empty"
            return
        }
        fetch(API+"/text/speech", {
            method: "POST",
            body: JSON.stringify({
                lang: name == "record2" ? lang1value : lang2value,
                text: text.innerText
            })
        })
        .then(async data => {
            console.log(data)
            if (data.status == 200) {
                body = await data.blob()
                src = URL.createObjectURL(body)
                audio.src = src
                audio.load()
                btn.classList.remove("hide")
                if (next != null) {
                    next()
                }
                return
            }
            
            error.innerText = "something went wrong, try again later"
            if (data.status == 400) {
                data.json().then(j => {
                    error.innerText = j.error
                })
            }
        })
    }
}

function getInfo() {
    fetch(API+"/info", {
        method: "GET",
    })
    .then(data => {
        console.log(data)
        if (data.status == 200) {
            data.json().then(body => {
                if (body.recognize) {
                    if (body.recognize?.languages) {
                        renderLanguage(body.recognize?.languages)
                    }
                    if (body.recognize?.maximumLength && body.recognize?.maximumLength > 0) {
                        longTimeVoiceMessage = body.recognize.maximumLength * 1000
                    }
                }
            })

            return
        }
        
        error.innerText = "something went wrong, try again later"
    })
}

function audioPlay(audio, btn) {
    let status = "pause"
    audio.addEventListener("ended", function() {
        btn.innerHTML = imgPlay
        status = "pause"
    })
    return function() {
        if (status == "pause") {
            audio.play()
            btn.innerHTML = imgPause
            status = "play"
        } else {
            audio.pause()
            btn.innerHTML = imgPlay
            status = "pause"
        }
    }
}

async function addAudio(src, name) {
    let audio = document.createElement("audio")
    audio.src = src

    let error = document.createElement("div")
    error.classList.add("error")

    let recognizedMsg = document.createElement("div")
    recognizedMsg.classList.value = "recognized_message hide"

    let translatedMsg = document.createElement("div")
    translatedMsg.classList.value = "translated_message hide"
    
    let audioTranslated = document.createElement("audio")

    let playTranslated = document.createElement("button")
    playTranslated.innerHTML = imgPlay
    playTranslated.classList.value = "btn btn-primary hide"

    const eventPlaySpeech = audioPlay(audioTranslated, playTranslated)

    playTranslated.addEventListener("click", eventPlaySpeech)

    const eventGetSpeech = speech(audioTranslated, playTranslated, translatedMsg, name, error, eventPlaySpeech)
    const eventGetTranslate = translating(translatedMsg, recognizedMsg, name, error, eventGetSpeech)
    const eventGetRecognize = recognizing(recognizedMsg, src, name, eventGetTranslate)

    let play = document.createElement("button")
    play.innerHTML = imgPlay
    play.classList.value = "btn btn-primary"
    play.addEventListener("click", audioPlay(audio, play))

    let recognize = document.createElement("button")
    recognize.innerHTML = imgRecognize
    recognize.classList.value = "btn btn-primary"
    recognize.addEventListener("click", eventGetRecognize)

    let translate = document.createElement("button")
    translate.innerHTML = imgTranslate
    translate.classList.value = "btn btn-primary"
    translate.addEventListener("click", eventGetTranslate)


    let tts = document.createElement("button")
    tts.innerHTML = imgUpload
    tts.classList.value = "btn btn-primary"
    tts.addEventListener("click", eventGetSpeech)

    let btnRemove = document.createElement("button")
    btnRemove.innerHTML = imgClose
    btnRemove.classList.value = "btn btn-primary"
    btnRemove.addEventListener("click", function() {
        message.remove()
    })

    let controler = document.createElement("div")
    controler.classList.add("msg_controler")
    controler.appendChild(audio)
    controler.appendChild(play)
    controler.appendChild(recognize)
    controler.appendChild(translate)
    controler.appendChild(tts)
    controler.appendChild(audioTranslated)
    controler.appendChild(playTranslated)
    controler.appendChild(btnRemove)

    let title = document.createElement("span")
    if (name != "record2")
        title.innerText = +(new Date())+" "+lang1value+" - "+lang2value
    else
        title.innerText = lang2value+" - "+lang1value+" "+ +(new Date())

    let message = document.createElement("div")
    message.classList.add("message")
    if (name == "record2")
        message.classList.add("right")
    message.appendChild(title)
    message.appendChild(recognizedMsg)
    message.appendChild(translatedMsg)
    message.appendChild(error)
    message.appendChild(controler)

    chat.prepend(message)
    chat.scrollTop = chat.scrollHeight
    
    eventGetRecognize()
}

function renderLanguage(languages) {
    for (const l of Object.keys(languages)) {
        const langName = l.split(" ")[0]
        li = document.createElement("li")
        const a1 = document.createElement("a")
        a1.classList.add("dropdown-item")
        a1.href = "#"
        a1.innerText = l
        a1.addEventListener("click", function() {
            lang1value = languages[l]
            btnLang1.innerText = langName
            for (li of lang1.querySelectorAll("li")) {
                li.querySelector("a").classList.remove("active")
            }
            a1.classList.add("active")
        })
        if (lang1value == languages[l]) {
            a1.classList.add("active")
        }
        li.appendChild(a1)
        lang1.append(li)


        const a2 = document.createElement("a")
        a2.classList.add("dropdown-item")
        a2.href = "#"
        a2.innerText = l
        a2.addEventListener("click", function() {
            lang2value = languages[l]
            btnLang2.innerText = langName
            for (li of lang2.querySelectorAll("li")) {
                li.querySelector("a").classList.remove("active")
            }
            a2.classList.add("active")
        })
        if (lang2value == languages[l]) {
            a2.classList.add("active")
        }
        li = document.createElement("li")
        li.appendChild(a2)
        lang2.append(li)
    }
}

getInfo()
// renderLanguage()