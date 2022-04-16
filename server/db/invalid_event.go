package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type InvalidEventType string

type InvalidEventOverview struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type InvalidEventProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type InvalidEvent struct {
	ID             uuid.UUID       `db:"id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Name           string          `db:"name"`
	Properties     json.RawMessage `db:"properties"`
	Context        json.RawMessage `db:"context"`
	Error          string          `db:"error"`
	CreatedAt      time.Time       `db:"created_at"`
}

type CreateInvalidEventParams struct {
	OrganizationID uuid.UUID
	Name           string
	Properties     json.RawMessage
	Context        json.RawMessage
	Error          string
	CreatedAt      time.Time
}

func (d *DB) CreateInvalidEvent(ctx context.Context, arg CreateInvalidEventParams) (InvalidEvent, error) {
	const q = `
INSERT INTO invalid_events (
	organization_id,
    name,
	properties,
	context,
	error,
	created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Properties,
		arg.Context,
		arg.Error,
		arg.CreatedAt,
	)
	var t InvalidEvent
	return t, row.StructScan(&t)
}

func (d *DB) DeleteInvalidEvent(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM invalid_events
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetInvalidEvents(ctx context.Context, orgID uuid.UUID) ([]InvalidEvent, error) {
	const q = `
SELECT *
FROM invalid_events
WHERE organization_id = $1
ORDER BY created_at DESC
`

	var items []InvalidEvent
	err := d.db.SelectContext(ctx, &items, q, orgID)
	return items, err
}
