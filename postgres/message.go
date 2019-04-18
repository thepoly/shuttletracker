package postgres

import (
	"database/sql"

	"github.com/thepoly/shuttletracker"
)

// MessageService implements a mock of shuttletracker.MessageService.
type MessageService struct {
	db *sql.DB
}

func (ms *MessageService) initializeSchema(db *sql.DB) error {
	ms.db = db
	schema := `
CREATE TABLE IF NOT EXISTS messages (
	id bool PRIMARY KEY DEFAULT true CHECK (id = true),
	message text,
	enabled bool NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now()
);`
	_, err := ms.db.Exec(schema)
	return err
}

// Message returns the Message.
func (ms *MessageService) Message() (*shuttletracker.Message, error) {
	query := "SELECT message, enabled, created, updated FROM messages;"
	row := ms.db.QueryRow(query)
	message := &shuttletracker.Message{}
	err := row.Scan(&message.Message, &message.Enabled, &message.Created, &message.Updated)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrMessageNotFound
	} else if err != nil {
		return nil, err
	}
	return message, nil
}

// SetMessage updates the Message.
func (ms *MessageService) SetMessage(message *shuttletracker.Message) error {
	statement := "INSERT INTO messages (message, enabled, updated) VALUES ($1, $2, now())" +
		" ON CONFLICT (id) DO UPDATE SET message = excluded.message, enabled = excluded.enabled, updated = excluded.updated" +
		" RETURNING created, updated;"
	row := ms.db.QueryRow(statement, message.Message, message.Enabled)
	return row.Scan(&message.Created, &message.Updated)
}
