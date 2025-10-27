package handlers

import "TRYREST/internal/storage/postgre"

// в одно ведро собрали все хендлеры
type Handler struct {
	UserHandler    *UserHandler
	EventHandler   *EventHandler
	BookingHandler *BookingHandler
}

// инициализирует все под-хендлеры
func NewHandler(storage *postgre.Storage) *Handler {
	return &Handler{
		UserHandler:    NewUserHandler(storage),
		EventHandler:   NewEventHandler(storage),
		BookingHandler: NewBookingHandler(storage),
	}
}
