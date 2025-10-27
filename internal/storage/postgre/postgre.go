package postgre

import (
	"TRYREST/internal/models"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq" // инициализация драйвера postgres
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(dsn string, logger *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//проверка подключения
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	// напрямую, но так как появились миграции то ненужно
	//createSQL := `
	//CREATE TABLE IF NOT EXISTS users (
	//    id SERIAL PRIMARY KEY,
	//    name VARCHAR(255) NOT NULL,
	//    email VARCHAR(255) UNIQUE NOT NULL
	//);
	//
	//CREATE TABLE IF NOT EXISTS events (
	//    id SERIAL PRIMARY KEY,
	//    title VARCHAR(255) NOT NULL,
	//    description TEXT
	//);
	//
	//CREATE TABLE IF NOT EXISTS bookings (
	//    id SERIAL PRIMARY KEY,
	//    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
	//    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
	//);
	//`
	//
	//if _, err := db.Exec(createSQL); err != nil {
	//	return nil, fmt.Errorf("%s: %w", op, err)
	//}

	return &Storage{
		db:  db,
		log: logger,
	}, nil
}

func (s *Storage) GetAllUsers() ([]models.User, error) {
	const op = "storage.postgre.GetAllUsers"
	rows, err := s.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		s.log.Error("Failed to query users", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			s.log.Error("Failed to close rows", slog.String("op", op), slog.Any("error", cerr))
		}
	}()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			s.log.Error("Failed to scan user", slog.String("op", op), slog.Any("error", err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		s.log.Error("Error iterating rows", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return users, nil
}

func (s *Storage) GetUserByID(id int64) (models.User, error) {
	const op = "storage.postgre.GetUserByID"
	var user models.User
	err := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("%s: user not found", op)
	}
	if err != nil {
		s.log.Error("Failed to query user by ID", slog.String("op", op), slog.Any("error", err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) AddUser(name, email string) (int64, error) {
	const op = "storage.postgres.AddUser"
	var id int64
	// используем QueryRow + RETURNING id
	err := s.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", name, email).Scan(&id)
	if err != nil {
		s.log.Error("Failed to insert user", slog.String("op", op), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UpdateUser(id int64, name, email string) error {
	const op = "storage.postgre.UpdateUser"
	result, err := s.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", name, email, id)
	if err != nil {
		s.log.Error("Failed to update user", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: user not found", op)
	}
	return nil
}

func (s *Storage) DeleteUser(id int64) error {
	const op = "storage.postgre.DeleteUser"
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		s.log.Error("Failed to delete user", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: user not found", op)
	}
	return nil
}

func (s *Storage) GetAllEvents() ([]models.Event, error) {
	const op = "storage.postgre.GetAllEvents"
	rows, err := s.db.Query("SELECT id, title, description FROM events")
	if err != nil {
		s.log.Error("Failed to query events", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			s.log.Error("Failed to close rows", slog.String("op", op), slog.Any("error", cerr))
		}
	}()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description); err != nil {
			s.log.Error("Failed to scan event", slog.String("op", op), slog.Any("error", err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		s.log.Error("Error iterating rows", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return events, nil
}

func (s *Storage) GetEventByID(id int64) (models.Event, error) {
	const op = "storage.postgre.GetEventByID"
	var event models.Event
	err := s.db.QueryRow("SELECT id, title, description FROM events WHERE id = $1", id).Scan(&event.ID, &event.Title, &event.Description)
	if err == sql.ErrNoRows {
		return models.Event{}, fmt.Errorf("%s: event not found", op)
	}
	if err != nil {
		s.log.Error("Failed to query event by ID", slog.String("op", op), slog.Any("error", err))
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	return event, nil
}

func (s *Storage) AddEvent(title, description string) (int64, error) {
	const op = "storage.postgres.AddEvent"
	var id int64
	err := s.db.QueryRow("INSERT INTO events (title, description) VALUES ($1, $2) RETURNING id", title, description).Scan(&id)
	if err != nil {
		s.log.Error("Failed to insert event", slog.String("op", op), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UpdateEvent(id int64, title, description string) error {
	const op = "storage.postgre.UpdateEvent"
	result, err := s.db.Exec("UPDATE events SET title = $1, description = $2 WHERE id = $3", title, description, id)
	if err != nil {
		s.log.Error("Failed to update event", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: event not found", op)
	}
	return nil
}

func (s *Storage) DeleteEvent(id int64) error {
	const op = "storage.postgre.DeleteEvent"
	result, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		s.log.Error("Failed to delete event", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: event not found", op)
	}
	return nil
}

func (s *Storage) GetAllBookings() ([]models.Booking, error) {
	const op = "storage.postgre.GetAllBookings"
	rows, err := s.db.Query("SELECT id, event_id, user_id FROM bookings")
	if err != nil {
		s.log.Error("Failed to query bookings", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			s.log.Error("Failed to close rows", slog.String("op", op), slog.Any("error", cerr))
		}
	}()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.EventID, &booking.UserID); err != nil {
			s.log.Error("Failed to scan booking", slog.String("op", op), slog.Any("error", err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		s.log.Error("Error iterating rows", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return bookings, nil
}

func (s *Storage) GetBookingByID(id int64) (models.Booking, error) {
	const op = "storage.postgre.GetBookingByID"
	var booking models.Booking
	err := s.db.QueryRow("SELECT id, event_id, user_id FROM bookings WHERE id = $1", id).Scan(&booking.ID, &booking.EventID, &booking.UserID)
	if err == sql.ErrNoRows {
		return models.Booking{}, fmt.Errorf("%s: booking not found", op)
	}
	if err != nil {
		s.log.Error("Failed to query booking by ID", slog.String("op", op), slog.Any("error", err))
		return models.Booking{}, fmt.Errorf("%s: %w", op, err)
	}
	return booking, nil
}

func (s *Storage) AddBooking(eventID, userID int64) (int64, error) {
	const op = "storage.postgres.AddBooking"
	var id int64
	err := s.db.QueryRow("INSERT INTO bookings (event_id, user_id) VALUES ($1, $2) RETURNING id", eventID, userID).Scan(&id)
	if err != nil {
		s.log.Error("Failed to insert booking", slog.String("op", op), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UpdateBooking(id, eventID, userID int64) error {
	const op = "storage.postgre.UpdateBooking"
	result, err := s.db.Exec("UPDATE bookings SET event_id = $1, user_id = $2 WHERE id = $3", eventID, userID, id)
	if err != nil {
		s.log.Error("Failed to update booking", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: booking not found", op)
	}
	return nil
}

func (s *Storage) DeleteBooking(id int64) error {
	const op = "storage.postgre.DeleteBooking"
	result, err := s.db.Exec("DELETE FROM bookings WHERE id = $1", id)
	if err != nil {
		s.log.Error("Failed to delete booking", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("Failed to check rows affected", slog.String("op", op), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: booking not found", op)
	}
	return nil
}

func (s *Storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	// sql.DB.Close() закрывает пул соединений и освобождает ресурсы
	return s.db.Close()
}
