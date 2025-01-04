package tipster

import (
	"context"
	"fmt"

	"go.uber.org/zap"
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

type Tipster struct {
	recognizer recognizerI
	translator translatorI
	tts        ttsI
	cache      cacheI
	log        *zap.Logger
}

type Option func(*Tipster)

func SetLogger(l *zap.Logger) Option {
	return func(t *Tipster) {
		t.log = l
	}
}

func New(recognizer recognizerI, translator translatorI, speech ttsI, cache cacheI, options ...Option) (*Tipster, error) {
	t := &Tipster{
		log:        zap.NewNop(),
		recognizer: recognizer,
		translator: translator,
		tts:        speech,
		cache:      cache,
	}

	for _, opt := range options {
		opt(t)
	}
	return t, nil
}

func (t *Tipster) Recognize(ctx context.Context, data []byte, language string) (string, error) {
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

	t.log.Debug("recognize caching", zap.String("key", key))
	t.cache.Set(ctx, key, text, 0)

	return text, nil
}

func (t *Tipster) Translate(ctx context.Context, sourceLanguage, targetLanguage string, text string) (string, error) {
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

	t.log.Debug("translate caching", zap.String("key", key))
	t.cache.Set(ctx, key, res, 0)

	return res, nil
}

func (t *Tipster) Speech(ctx context.Context, text, lang string) ([]byte, error) {
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

	t.log.Debug("speech caching", zap.String("key", key))
	t.cache.Set(ctx, key, string(res), 0)

	return res, nil
}
