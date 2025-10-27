package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"TRYREST/internal/models"
	"TRYREST/internal/storage/postgre"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	storage *postgre.Storage
}

func NewBookingHandler(storage *postgre.Storage) *BookingHandler {
	return &BookingHandler{storage: storage}
}

func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookings, err := h.storage.GetAllBookings()
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) GetBookingById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.storage.GetBookingByID(id)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBooking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.storage.AddBooking(newBooking.EventID, newBooking.UserID)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}
	newBooking.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBooking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var updatedBooking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&updatedBooking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.storage.UpdateBooking(id, updatedBooking.EventID, updatedBooking.UserID); err != nil {
		if err.Error() == "storage.postgre.UpdateBooking: booking not found" {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update booking", http.StatusInternalServerError)
		return
	}
	updatedBooking.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBooking)
}

func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteBooking(id); err != nil {
		if err.Error() == "storage.postgre.DeleteBooking: booking not found" {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
