package recognize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

/**
yandex documentation recognize https://yandex.cloud/ru/docs/speechkit/stt/api/request-api
*/

type speechFormat string
type Language string

const (
	LPCM    speechFormat = "lpcm"
	OggIpus speechFormat = "oggopus"

	LangRu Language = "ru-RU"

	api = "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
)

type Recognizer struct {
	url      string
	folderID string
	apiKey   string
}

type Option func(*Recognizer)

func New(apiKey string, opts ...Option) *Recognizer {
	r := &Recognizer{
		apiKey: apiKey,
		url:    api,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type RecognizeOption func(*map[string]string)

func SetLanguage(l Language) RecognizeOption {
	return func(m *map[string]string) {
		_m := *m
		_m["lang"] = string(l)
		m = &_m
	}
}

func SetFormat(f speechFormat) RecognizeOption {
	return func(m *map[string]string) {
		_m := *m
		_m["format"] = string(f)
		m = &_m
	}
}

func (r *Recognizer) Recognize(data []byte, opts ...RecognizeOption) (string, error) {
	skURL, _ := url.Parse(r.url)
	values := map[string]string{
		"lang":            string(LangRu),
		"topic":           "general",
		"profanityFilter": "false",
		"format":          string(OggIpus),
		"rawResults":      "true",
		"sampleRateHerts": "48000",
	}
	if r.folderID != "" {
		values["folderId"] = r.folderID
	}

	for _, opt := range opts {
		opt(&values)
	}

	urlValues := skURL.Query()
	for k, v := range values {
		urlValues.Add(k, v)
	}
	skURL.RawQuery = urlValues.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, skURL.String(), bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if r.apiKey != "" {
		req.Header.Add("Authorization", "Api-Key "+r.apiKey)
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var res struct {
		Result       string `json:"result"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
	err = json.Unmarshal(bBody, &res)
	if err != nil {
		return "", err
	}

	if res.ErrorMessage != "" {
		return "", fmt.Errorf("yandex error: %s", res.ErrorMessage)
	}

	return res.Result, nil
}
