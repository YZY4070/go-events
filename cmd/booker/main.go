package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TRYREST/internal/booker"
	"TRYREST/internal/config"
	"TRYREST/internal/lib/logger/sl"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	srv, log, cleanup, err := booker.New(cfg)
	if err != nil {
		if log != nil {
			log.Error("failed to initialize app", sl.Err(err))
		} else {
			slog.Error("failed to initialize app", slog.Any("err", err))
		}
		os.Exit(1)
	}

	// Создаём контекст, который автоматически отменится при получении SIGINT или SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запускаем сервер в горутине.
	go func() {
		log.Info("starting server", slog.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped unexpectedly", sl.Err(err))
			// при критической ошибке просим прекратить ожидание - отменяем контекст
			stop()
		}
	}()

	// Ждём сигнала отмены (Ctrl+C или kill) в ctx.
	<-ctx.Done()

	// Когда получили сигнал — начинаем shutdown.
	// Дадим 10 секунд на завершение активных запросов и cleanup.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutdown signal received, shutting down...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", sl.Err(err))
	} else {
		log.Info("http server stopped")
	}

	if cleanup != nil {
		if err := cleanup(shutdownCtx); err != nil {
			log.Error("cleanup failed", sl.Err(err))
		} else {
			log.Info("cleanup done")
		}
	}

	log.Info("shutdown complete")
}
