package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/playmixer/tipster/internal/adapters/api/rest"
	"github.com/playmixer/tipster/internal/adapters/cache"
	"github.com/playmixer/tipster/internal/adapters/logger"
	"github.com/playmixer/tipster/internal/adapters/notification"
	"github.com/playmixer/tipster/internal/adapters/recognizer"
	"github.com/playmixer/tipster/internal/adapters/storage"
	"github.com/playmixer/tipster/internal/adapters/translator/yandex"
	"github.com/playmixer/tipster/internal/adapters/tts"
	"github.com/playmixer/tipster/internal/core/config"
	"github.com/playmixer/tipster/internal/core/tipster"
	"go.uber.org/zap"
)

var (
	shutdownDelay = time.Second * 2
)

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("failed initialize config: %w", err)
	}

	lgr, err := logger.New(cfg.LogLVL)
	if err != nil {
		return fmt.Errorf("failed initialize logger: %w", err)
	}

	rcgz, err := recognizer.New(cfg.RecognizerName, cfg.Recognizer, lgr)
	if err != nil {
		return fmt.Errorf("failed initialize recognize: %w", err)
	}

	trns := yandex.New(cfg.Translator.Yandex)

	cch := cache.New(ctx, cfg.CacheName, cfg.Cache, lgr)

	speech, err := tts.New(cfg.TTSName, cfg.TTS)
	if err != nil {
		return fmt.Errorf("failed initialize tts: %w", err)
	}

	notify, err := notification.New(cfg.Notify)
	if err != nil {
		return fmt.Errorf("failed initialize notify: %w", err)
	}

	store, err := storage.New(cfg.Store, lgr)
	if err != nil {
		return fmt.Errorf("failed initialize storage: %w", err)
	}

	service, err := tipster.New(cfg.Tipster, store, notify, rcgz, trns, speech, cch, tipster.SetLogger(lgr))
	if err != nil {
		return fmt.Errorf("failed initialize tipster: %w", err)
	}

	server, err := rest.New(service, rest.SetAddress(cfg.Address), rest.Logger(lgr))
	if err != nil {
		return fmt.Errorf("failed initalize server: %w", err)
	}
	lgr.Info("Starting")
	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.Error("failed run server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	lgr.Info("Stopping...")
	ctxShutdown, stop := context.WithTimeout(context.Background(), shutdownDelay)
	defer stop()

	// выключаем http сервер
	if err := server.Shutdown(ctxShutdown); err != nil {
		lgr.Error("Server Shutdown with error", zap.Error(err))
	}

	<-ctxShutdown.Done()
	lgr.Info("Service stopped")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
