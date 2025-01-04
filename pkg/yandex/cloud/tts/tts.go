package tts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type speechFormat string
type Language string

const (
	LPCM    speechFormat = "lpcm"
	OggIpus speechFormat = "oggopus"
	MP3     speechFormat = "mp3"

	LangRu Language = "ru-RU"
	LangEu Language = "en-EU"
	LangKz Language = "kk-KK"
	LangDe Language = "de-DE"
	LangHe Language = "he-IL" // иврит
	LangUz Language = "uz-UZ"

	api = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"
)

type TTS struct {
	api      string
	apiKey   string
	folderID string
}

var (
	voices = map[string]string{
		"ru-RU": "filipp",
		"en-EU": "john",
		"kk-KK": "madi",
		"de-DE": "lea",
		"he-IL": "naomi",
	}
	languages = map[string]Language{
		"ru": LangRu,
		"en": LangEu,
		"kk": LangKz,
		"de": LangDe,
		// "he": LangHe,
		"uz": LangUz,
	}
)

func New(apiKey string) *TTS {
	t := &TTS{
		api:    api,
		apiKey: apiKey,
	}

	return t
}

type SpeechOption func(*map[string]string)

func SetLanguage(l Language) SpeechOption {
	return func(m *map[string]string) {
		_m := *m
		_m["lang"] = string(l)
		m = &_m
	}
}

func SetFormat(f speechFormat) SpeechOption {
	return func(m *map[string]string) {
		_m := *m
		_m["format"] = string(f)
		m = &_m
	}
}

func ParseLang(l string) *Language {
	if strings.Contains(l, "-") {
		_l := strings.Split(l, "-")
		l = _l[0]
	}
	if lang, ok := languages[l]; ok {
		return &lang
	}

	return nil
}

func (r *TTS) Speech(text string, opts ...SpeechOption) ([]byte, error) {
	skURL, _ := url.Parse(r.api)
	values := map[string]string{
		"lang":            string(LangRu),
		"format":          string(OggIpus),
		"sampleRateHertz": "48000",
		"text":            text,
	}
	if r.folderID != "" {
		values["folderID"] = r.folderID
	}

	for _, opt := range opts {
		opt(&values)
	}

	if voice, ok := voices[values["lang"]]; ok {
		values["voice"] = voice
	}

	urlValues := skURL.Query()
	for k, v := range values {
		urlValues.Add(k, v)
	}
	skURL.RawQuery = urlValues.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, skURL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	if r.apiKey != "" {
		req.Header.Add("Authorization", "Api-Key "+r.apiKey)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var res struct {
			ErrorCode    string `json:"error_code"`
			ErrorMessage string `json:"error_message"`
		}
		err = json.Unmarshal(bBody, &res)
		if err != nil {
			return nil, err
		}

		if res.ErrorMessage != "" {
			return nil, fmt.Errorf("yandex error: %s", res.ErrorMessage)
		}
	}

	return bBody, nil
}
