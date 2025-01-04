const constraints = { audio: true, video: false }
let stream = null
navigator.mediaDevices.getUserMedia(constraints)
 .then((_stream) => { stream = _stream })
 // если возникла ошибка, значит, либо пользователь отказал в доступе,
 // либо запрашиваемое медиа-устройство не обнаружено
 .catch((err) => { console.error(`Not allowed or not found: ${err}`) })


let chunks = []
let mediaRecorder = null
let audioBlob = null

const error_view = document.querySelector("#error")
console.error = function(err) {
  error_view.innerText = err
}

function startRecord(e, name, event) {
    return async function() {
        // проверяем поддержку
        if (!navigator.mediaDevices && !navigator.mediaDevices.getUserMedia) {
            return console.warn('Not supported')
        }

        let isStart = false
    
        // если запись не запущена
        if (!mediaRecorder && (event == "start" || event == "click")) {
            isStart = true
            console.log("start record")
            try {
                // получаем поток аудио-данных
                const stream = await navigator.mediaDevices.getUserMedia({
                    audio: true
                })
                // создаем экземпляр `MediaRecorder`, передавая ему поток в качестве аргумента
                mediaRecorder = new MediaRecorder(stream)
                // запускаем запись
                mediaRecorder.start()
                // по окончанию записи и наличии данных добавляем эти данные в соответствующий массив
                mediaRecorder.ondataavailable = (e) => {
                    chunks.push(e.data)
                }
                // обработчик окончания записи (см. ниже)
                mediaRecorder.onstop = mediaRecorderStop(name)
                
                e.classList.add("recording")
                e.innerHTML = imgRecord
                e.classList.add("text-danger")
                timerID = setTimeout(timeoutStopRecord(e, name), longTimeVoiceMessage)
            } catch (e) {
                console.error(e)
                e.classList.remove("recording")
                e.innerHTML = imgMicrophon
                e.classList.remove("text-danger")
            }
        }
        if (mediaRecorder && !isStart && (event == "end" || event == "click")){
            // если запись запущена, останавливаем ее
            mediaRecorder.stop()
            console.log("stop record")
            e.classList.remove("recording")
            e.innerHTML = imgMicrophon
            e.classList.remove("text-danger")
            mediaRecorder = null
            clearTimeout(timerID)
            timerID = null
        }
        console.log(e)
    }
}

function mediaRecorderStop(name) {
    return async function() {
        audioBlob = new Blob(chunks, { type: 'audio/webm;codecs=opus' })

        const arrayBuffer = await audioBlob.arrayBuffer();
        const audioContext = new AudioContext();
        const audioBuffer = await audioContext.decodeAudioData(arrayBuffer);
        
        const wavBlob = await convertAudioBufferToWavBlob(audioBuffer);
        console.log(wavBlob)
    
        const src = URL.createObjectURL(wavBlob)
        addAudio(src, name)
        mediaRecorder = null
        chunks = []
    }
}

function convertAudioBufferToWavBlob(audioBuffer) {
    return new Promise(function (resolve) {
        var worker = new Worker('static/js/wave-worker.js');

        worker.onmessage = function (e) {
            var blob = new Blob([e.data.buffer], { type: 'audio/wav' });
            resolve(blob);
        };

        let pcmArrays = [];
        for (let i = 0; i < audioBuffer.numberOfChannels; i++) {
            pcmArrays.push(audioBuffer.getChannelData(i));
        }

        worker.postMessage({
            pcmArrays,
            config: { sampleRate: audioBuffer.sampleRate },
        });
    });
}