package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type InvalidEventType string

type InvalidEvent struct {
	ID             uuid.UUID       `db:"id"`
	EventID        uuid.UUID       `db:"event_id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Channel        string          `db:"channel"`
	Name           string          `db:"name"`
	Properties     json.RawMessage `db:"properties"`
	Context        json.RawMessage `db:"context"`
	Type           string          `db:"type"`
	CreatedAt      time.Time       `db:"created_at"`
}

type CreateInvalidEventParams struct {
	EventID        uuid.UUID
	OrganizationID uuid.UUID
	Channel        string
	Name           string
	Properties     json.RawMessage
	Context        json.RawMessage
	Type           string
	CreatedAt      time.Time
}

func (d *DB) CreateInvalidEvent(ctx context.Context, arg CreateInvalidEventParams) (InvalidEvent, error) {
	const q = `
INSERT INTO invalid_events (
	event_id,
	organization_id,
	channel,
	name,
	properties,
	context,
	type,
	created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.EventID,
		arg.OrganizationID,
		arg.Channel,
		arg.Name,
		arg.Properties,
		arg.Context,
		arg.Type,
		arg.CreatedAt,
	)
	var t InvalidEvent
	return t, row.StructScan(&t)
}

func (d *DB) UpdateInvalidEvent(ctx context.Context, newName string, orgID, eventId uuid.UUID) error {
	const q = `
UPDATE invalid_events
SET name = $1
WHERE organization_id = $2
AND event_id = $3
`
	_, err := d.db.ExecContext(ctx, q, newName, orgID, eventId)
	return err
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
