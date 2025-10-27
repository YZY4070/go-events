package booker

import (
	"context"
	"net/http"
	"os"
	_ "time"

	"TRYREST/internal/config"
	"TRYREST/internal/handlers"
	"TRYREST/internal/lib/logger/sl"
	"TRYREST/internal/storage/postgre"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New собирает и возвращает подготовленный (но ещё не запущенный) HTTP-сервер,
// инициализированный логгер, cleanup-функцию и ошибку (если что-то пошло не так).
func New(cfg *config.Config) (*http.Server, *slog.Logger, func(ctx context.Context) error, error) {
	log := setupLogger(cfg.Env)
	log.Info("booker initialization start", slog.String("env", cfg.Env))

	storage, err := postgre.New(cfg.StoragePath, log)
	if err != nil {
		log.Error("error creating storage", sl.Err(err))
		return nil, nil, nil, err
	}

	h := handlers.NewHandler(storage)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/users", func(r chi.Router) {
		r.Get("/", h.UserHandler.GetAllUsers)
		r.Post("/", h.UserHandler.CreateUser)
		r.Get("/{id}", h.UserHandler.GetUserByID)
		r.Put("/{id}", h.UserHandler.UpdateUser)
		r.Delete("/{id}", h.UserHandler.DeleteUser)
	})

	router.Route("/events", func(r chi.Router) {
		r.Get("/", h.EventHandler.GetAllEvents)
		r.Post("/", h.EventHandler.CreateEvent)
		r.Get("/{id}", h.EventHandler.GetEventByID)
		r.Put("/{id}", h.EventHandler.UpdateEvent)
		r.Delete("/{id}", h.EventHandler.DeleteEvent)
	})

	router.Route("/bookings", func(r chi.Router) {
		r.Get("/", h.BookingHandler.GetAllBookings)
		r.Post("/", h.BookingHandler.CreateBooking)
		r.Get("/{id}", h.BookingHandler.GetBookingById)
		r.Put("/{id}", h.BookingHandler.UpdateBooking)
		r.Delete("/{id}", h.BookingHandler.DeleteBooking)
	})

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	//это функция очистки ресурсов которая использует общий интерфейс(пока до конца не разобрался)
	cleanup := func(ctx context.Context) error {
		type closer interface {
			Close() error
		}
		if c, ok := any(storage).(closer); ok {
			if err := c.Close(); err != nil {
				log.Error("storage close failed", sl.Err(err))
				return err
			}
		}
		return nil
	}

	return srv, log, cleanup, nil
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case "local":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
}
