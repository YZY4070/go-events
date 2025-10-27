package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"TRYREST/internal/models"
	"TRYREST/internal/storage/postgre"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	storage *postgre.Storage
}

func NewEventHandler(storage *postgre.Storage) *EventHandler {
	return &EventHandler{storage: storage}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events, err := h.storage.GetAllEvents()
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := h.storage.GetEventByID(id)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.storage.AddEvent(newEvent.Title, newEvent.Description)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	newEvent.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var updatedEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.storage.UpdateEvent(id, updatedEvent.Title, updatedEvent.Description); err != nil {
		if err.Error() == "storage.postgre.UpdateEvent: event not found" {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}
	updatedEvent.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedEvent)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteEvent(id); err != nil {
		if err.Error() == "storage.postgre.DeleteEvent: event not found" {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
