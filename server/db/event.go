package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type EventType string

type EventOverview struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type Event struct {
	ID             uuid.UUID       `db:"id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Name           string          `db:"name"`
	Description    string          `db:"description"`
	Version        int             `db:"version"`
	Structure      json.RawMessage `db:"structure"`
	Properties     json.RawMessage `db:"properties"`
	JSONSchema     json.RawMessage `db:"json_schema"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

type CreateEventParams struct {
	OrganizationID uuid.UUID
	Name           string
	Description    string
	Structure      json.RawMessage
	JSONSchema     json.RawMessage
	Properties     json.RawMessage
}

func (d *DB) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	const q = `
INSERT INTO events (
	organization_id,
    name,
	description,
	structure,
	json_schema,
	properties
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Description,
		arg.Structure,
		arg.JSONSchema,
		arg.Properties,
	)
	var t Event
	return t, row.StructScan(&t)
}

func (d *DB) DeleteEvent(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM events
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetEvent(ctx context.Context, orgID, id uuid.UUID) (Event, error) {
	const q = `
SELECT *
FROM events
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var e Event
	return e, row.StructScan(&e)
}

func (d *DB) GetEventByName(ctx context.Context, orgID uuid.UUID, name string) (Event, error) {
	const q = `
SELECT *
FROM events
WHERE organization_id = $1
AND name = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, name)
	var e Event
	return e, row.StructScan(&e)
}

func (d *DB) GetEvents(ctx context.Context, orgID uuid.UUID) ([]Event, error) {
	const q = `
SELECT *
FROM events
WHERE organization_id = $1
ORDER BY name ASC
`

	var items []Event
	err := d.db.SelectContext(ctx, &items, q, orgID)
	return items, err
}

type UpdateEventParams struct {
	Name           string
	Description    string
	Structure      json.RawMessage
	JSONSchema     json.RawMessage
	Properties     json.RawMessage
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateEvent(ctx context.Context, arg UpdateEventParams) error {
	const q = `
UPDATE events
SET name = $1,
	description = $2,
	structure = $3,
	json_schema = $4,
	properties = $5,
	updated_at = now()
WHERE organization_id = $6
AND id = $7
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.Name,
		arg.Description,
		arg.Structure,
		arg.JSONSchema,
		arg.Properties,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

func (d *DB) GetEventEvents(ctx context.Context, orgID uuid.UUID) ([]EventOverview, error) {
	const q = `
SELECT id, name
FROM events
WHERE organization_id = $1
`

	var items []EventOverview
	if err := d.db.SelectContext(ctx, &items, q, orgID); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) GetEventProperties(ctx context.Context, orgID, id uuid.UUID) (json.RawMessage, error) {
	const q = `
SELECT properties
FROM events
WHERE organization_id = $1
AND ID = $2
`

	var structure json.RawMessage
	row := d.db.QueryRowxContext(ctx, q, orgID, id)

	if err := row.Scan(&structure); err != nil {
		return nil, err
	}
	return structure, nil
}
