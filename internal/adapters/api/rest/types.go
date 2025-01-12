package rest

type sendOTPRequest struct {
	Email string `json:"email"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type textTranslateRequest struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	Text       string `json:"text"`
	FrontendID string `json:"frontendID"`
}

type textSpeechRequest struct {
	Lang       string `json:"lang"`
	Text       string `json:"text"`
	FrontendID string `json:"frontendID"`
}
