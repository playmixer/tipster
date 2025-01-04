package rest

type textTranslateRequest struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	Text       string `json:"text"`
}

type textSpeechRequest struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}
