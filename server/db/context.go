package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ContextProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type Context struct {
	ID             uuid.UUID       `db:"id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Channel        string          `db:"channel"`
	Structure      json.RawMessage `db:"structure"`
	JSONSchema     json.RawMessage `db:"json_schema"`
	Properties     json.RawMessage `db:"properties"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

type CreateContextParams struct {
	OrganizationID uuid.UUID
	Channel        string
	Structure      json.RawMessage
	JSONSchema     json.RawMessage
	Properties     json.RawMessage
}

func (d *DB) CreateContext(ctx context.Context, arg CreateContextParams) (Context, error) {
	const q = `
INSERT INTO contexts (
	organization_id,
	channel,
	structure,
	json_schema,
	properties
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Channel,
		arg.Structure,
		arg.JSONSchema,
		arg.Properties,
	)
	var t Context
	return t, row.StructScan(&t)
}

type UpdateContextParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Structure      json.RawMessage
	JSONSchema     json.RawMessage
	Properties     json.RawMessage
}

func (d *DB) UpdateContext(ctx context.Context, arg UpdateContextParams) error {
	const q = `
UPDATE contexts
SET structure = $1,
	json_schema = $2,
	properties = $3,
	updated_at = now()
WHERE organization_id = $4
AND id = $5
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.Structure,
		arg.JSONSchema,
		arg.Properties,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

func (d *DB) DeleteContext(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM contexts
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetContexts(ctx context.Context, orgID uuid.UUID) ([]Context, error) {
	const q = `
SELECT *
FROM contexts
WHERE organization_id = $1
`

	var items []Context
	err := d.db.SelectContext(ctx, &items, q, orgID)
	return items, err
}

func (d *DB) GetContextByID(ctx context.Context, orgID, id uuid.UUID) (Context, error) {
	const q = `
SELECT *
FROM contexts
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var e Context
	return e, row.StructScan(&e)
}

func (d *DB) GetContextProperties(ctx context.Context, orgID uuid.UUID, channel string) (json.RawMessage, error) {
	const q = `
SELECT properties
FROM contexts
WHERE organization_id = $1
AND channel = $2
`

	var structure json.RawMessage
	row := d.db.QueryRowxContext(ctx, q, orgID, channel)
	if err := row.Scan(&structure); err != nil {
		return nil, err
	}
	return structure, nil
}

func (d *DB) GetContextSchema(ctx context.Context, orgID uuid.UUID, channel string) (json.RawMessage, error) {
	const q = `
SELECT json_schema
FROM contexts
WHERE organization_id = $1
AND channel = $2
`

	var schema json.RawMessage
	row := d.db.QueryRowxContext(ctx, q, orgID, channel)
	if err := row.Scan(&schema); err != nil {
		return nil, err
	}
	return schema, nil
}
