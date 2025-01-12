package tipster

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"github.com/playmixer/tipster/internal/adapters/models"
	"github.com/playmixer/tipster/pkg/audio"
	"go.uber.org/zap"
)

var (
	lenRefreshToken        uint  = 50
	lifeTimeCacheRecognize int64 = 60 * 60
)

type recognizerI interface {
	Recognize(data []byte, language string) (string, error)
}

type translatorI interface {
	Translate(sourceLang, targetLang string, text string) (string, error)
}

type cacheI interface {
	Set(ctx context.Context, k string, v any, lifeTime int64)
	Get(ctx context.Context, k string) any
}

type ttsI interface {
	Speech(text string, lang string) ([]byte, error)
}

type notifyI interface {
	SendEmail(to string, subject string, body []byte) error
}

type storeI interface {
	NewOTP(ctx context.Context, email string, otp string) error
	GetOTP(ctx context.Context, email string) (*models.UserOTP, error)
	GetUser(ctx context.Context, userID uint) (*models.User, error)

	NewRefreshToken(ctx context.Context, userID uint, token string, lifeTime int64) error
	UpdRefreshToken(ctx context.Context, userID uint, oldToken, newToken string, lifeTime int64) error
	DelRefreshToken(ctx context.Context, userID uint, refreshToken string) error

	NewRecognize(ctx context.Context, userID uint, fID string, _range float32, language string, text string) (*models.Recognize, error)
	NewTranslate(ctx context.Context, userID uint, fID string, fromText, fromLanguage, toText, toLanguage string) (*models.Translate, error)
	NewSpeech(ctx context.Context, userID uint, fID string, text, language string) (*models.Speech, error)
}

type Tipster struct {
	cfg          Config
	store        storeI
	notification notifyI
	recognizer   recognizerI
	translator   translatorI
	tts          ttsI
	cache        cacheI
	log          *zap.Logger
}

type Option func(*Tipster)

func SetLogger(l *zap.Logger) Option {
	return func(t *Tipster) {
		t.log = l
	}
}

func New(cfg Config, store storeI,
	notification notifyI,
	recognizer recognizerI,
	translator translatorI,
	speech ttsI,
	cache cacheI,
	options ...Option,
) (*Tipster, error) {
	t := &Tipster{
		cfg:          cfg,
		log:          zap.NewNop(),
		store:        store,
		notification: notification,
		recognizer:   recognizer,
		translator:   translator,
		tts:          speech,
		cache:        cache,
	}

	for _, opt := range options {
		opt(t)
	}
	return t, nil
}

func (t *Tipster) SendOTP(ctx context.Context, email string) error {
	var subject string = "OTP"
	var otp string = randomString(t.cfg.OTPLength, numbers)

	err := t.store.NewOTP(ctx, email, otp)
	if err != nil {
		return fmt.Errorf("failed save otp: %w", err)
	}

	tmpl, err := template.ParseFiles("./templates/email/sendOTP.html")
	if err != nil {
		return fmt.Errorf("failed parse otp template: %w", err)
	}
	data := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(data, struct{ Code string }{Code: otp})
	if err != nil {
		return fmt.Errorf("failed execute template: %w", err)
	}

	err = t.notification.SendEmail(email, subject, data.Bytes())
	if err != nil {
		return fmt.Errorf("failed send email: %w", err)
	}
	return nil
}

func (t *Tipster) CheckOTP(ctx context.Context, email string, otp string) (*models.User, error) {
	userOTP, err := t.store.GetOTP(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed get otp: %w", err)
	}
	if userOTP.User.ID == 0 {
		return nil, apperrors.ErrUserNotFoundWithThisEmail
	}

	if userOTP.Code != otp || userOTP.Code == "" {
		return nil, apperrors.ErrPassNotEqual
	}

	return &userOTP.User, nil
}

func (t *Tipster) NewRefreshToken(ctx context.Context, userID uint) (string, error) {
	tkn := randomString(lenRefreshToken, alphabetLower+alphabetUpper+numbers)
	err := t.store.NewRefreshToken(ctx, userID, tkn, 0)
	if err != nil {
		return "", fmt.Errorf("failed create refresh token: %w", err)
	}
	return tkn, nil
}

func (t *Tipster) UpdRefreshToken(ctx context.Context, userID uint, refreshToken string) (string, error) {
	tkn := randomString(lenRefreshToken, alphabetLower+alphabetUpper+numbers)
	err := t.store.UpdRefreshToken(ctx, userID, refreshToken, tkn, 0)
	if err != nil {
		return "", fmt.Errorf("failed upd refresh token: %w", err)
	}

	return tkn, nil
}

func (t *Tipster) DelRefreshToken(ctx context.Context, userID uint, refreshToken string) error {
	err := t.store.DelRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return fmt.Errorf("failed delete refresh token")
	}

	return nil
}

func (t *Tipster) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := t.store.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %w", err)
	}

	return user, nil
}

func (t *Tipster) Recognize(ctx context.Context, userID uint, fID string, data []byte, language string) (string, error) {
	meta, err := audio.ReadWavMetadata(data)
	if err != nil {
		return "", fmt.Errorf("failed read record metadata: %w", err)
	}
	if meta.Duration < time.Second/2 {
		return "", apperrors.ErrShortRecord
	}

	if meta.Duration > time.Second*10 {
		return "", apperrors.ErrLongRecord
	}

	key := fmt.Sprintf("%s_%s", hashBtoS(data), language)
	_text := t.cache.Get(ctx, key)
	if _text != nil {
		if text, ok := _text.(string); ok {
			t.log.Debug("recognize get from cache", zap.String("key", key))
			return text, nil
		}
	}

	text, err := t.recognizer.Recognize(data, language)
	if err != nil {
		return "", fmt.Errorf("failed recognize: %w", err)
	}

	_, err = t.store.NewRecognize(ctx, userID, fID, float32(meta.Duration)/float32(time.Second), language, text)
	if err != nil {
		t.log.Error("failed save recognize meta", zap.Error(err))
	}

	t.log.Debug("recognize caching", zap.String("key", key))
	t.cache.Set(ctx, key, text, lifeTimeCacheRecognize)

	return text, nil
}

func (t *Tipster) Translate(ctx context.Context, userID uint, fID string, sourceLanguage, targetLanguage string, text string) (string, error) {
	key := fmt.Sprintf("%s_%s_%s", sourceLanguage, targetLanguage, text)
	_res := t.cache.Get(ctx, key)
	if _res != nil {
		if res, ok := _res.(string); ok {
			t.log.Debug("translate get from cache", zap.String("key", key))
			return res, nil
		}
	}

	res, err := t.translator.Translate(sourceLanguage, targetLanguage, text)
	if err != nil {
		return "", fmt.Errorf("failed translate text `%s`: %w", text, err)
	}

	_, err = t.store.NewTranslate(ctx, userID, fID, text, sourceLanguage, res, targetLanguage)
	if err != nil {
		t.log.Error("failed save translate meta", zap.Error(err))
	}

	t.log.Debug("translate caching", zap.String("key", key))
	t.cache.Set(ctx, key, res, 0)

	return res, nil
}

func (t *Tipster) Speech(ctx context.Context, userID uint, fID string, text, lang string) ([]byte, error) {
	key := fmt.Sprintf("%s_%s", lang, text)
	_res := t.cache.Get(ctx, key)
	if _res != nil {
		if res, ok := _res.(string); ok {
			t.log.Debug("speech get from cache", zap.String("key", key))
			return []byte(res), nil
		}
	}

	res, err := t.tts.Speech(text, lang)
	if err != nil {
		return nil, fmt.Errorf("failed generate speech: %w", err)
	}
	_, err = t.store.NewSpeech(ctx, userID, fID, text, lang)
	if err != nil {
		t.log.Error("failed save translate meta", zap.Error(err))
	}

	t.log.Debug("speech caching", zap.String("key", key))
	t.cache.Set(ctx, key, string(res), 0)

	return res, nil
}

func (t *Tipster) GetLanguages(ctx context.Context) map[string]string {
	return map[string]string{
		"RU (русский)":                   "ru-RU",
		"EN (английский)":                "en-US",
		"DE (немецкий)":                  "de-DE",
		"ES (испанский)":                 "es-ES",
		"FI (финский)":                   "fi-FI",
		"FR (французский)":               "fr-FR",
		"HE (иврит)":                     "he-HE",
		"IT (итальянский)":               "it-IT",
		"KZ (казахский)":                 "kk-KZ",
		"NL (голландский)":               "nl-NL",
		"PL (польский)":                  "pl-PL",
		"PT (португальский)":             "pt-PT",
		"BR (бразильский португальский)": "pt-BR",
		"SE (шведский)":                  "sv-SE",
		"TR (турецкий)":                  "tr-TR",
		"UZ (узбекский)":                 "uz-UZ",
	}
}

func (t *Tipster) GetRecognizeLanguages(ctx context.Context) []string {
	return []string{
		"ru-RU",
		"en-US",
		"de-DE",
		"es-ES",
		"fi-FI",
		"fr-FR",
		"he-HE",
		"it-IT",
		"kk-KZ",
		"nl-NL",
		"pl-PL",
		"pt-PT",
		"pt-BR",
		"sv-SE",
		"tr-TR",
		"uz-UZ",
	}
}

func (t *Tipster) GetSpeechLanguages(ctx context.Context) []string {
	return []string{
		"ru-RU",
		"en-US",
		"de-DE",
		"kk-KZ",
		"uz-UZ",
	}
}
