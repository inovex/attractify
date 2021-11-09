package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type CustomTraitsProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type CustomTraits struct {
	OrganizationID uuid.UUID       `db:"organization_id"`
	Structure      json.RawMessage `db:"structure"`
	JSONSchema     json.RawMessage `db:"json_schema"`
	Properties     json.RawMessage `db:"properties"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

type UpsertCustomTraitsParams struct {
	OrganizationID uuid.UUID
	Structure      json.RawMessage
	JSONSchema     json.RawMessage
	Properties     json.RawMessage
}

func (d *DB) UpsertCustomTraits(ctx context.Context, arg UpsertCustomTraitsParams) (CustomTraits, error) {
	const q = `
INSERT INTO custom_traits (
	organization_id,
	structure,
	json_schema,
	properties
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (organization_id)
DO UPDATE
SET structure = $2,
	json_schema = $3,
	properties = $4,
	updated_at = now()
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Structure,
		arg.JSONSchema,
		arg.Properties,
	)
	var t CustomTraits
	return t, row.StructScan(&t)
}

func (d *DB) DeleteCustomTraits(ctx context.Context, orgID uuid.UUID) error {
	const q = `
DELETE FROM custom_traits
WHERE organization_id = $1
`

	_, err := d.db.ExecContext(ctx, q, orgID)
	return err
}

func (d *DB) GetCustomTraits(ctx context.Context, orgID uuid.UUID) (CustomTraits, error) {
	const q = `
SELECT *
FROM custom_traits
WHERE organization_id = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID)
	var e CustomTraits
	return e, row.StructScan(&e)
}

func (d *DB) GetCustomTraitsProperties(ctx context.Context, orgID uuid.UUID) (json.RawMessage, error) {
	const q = `
SELECT properties
FROM custom_traits
WHERE organization_id = $1
`

	var structure json.RawMessage
	row := d.db.QueryRowxContext(ctx, q, orgID)
	if err := row.Scan(&structure); err != nil {
		return nil, err
	}
	return structure, nil
}

func (d *DB) GetCustomTraitsSchema(ctx context.Context, orgID uuid.UUID) (json.RawMessage, error) {
	const q = `
SELECT json_schema
FROM custom_traits
WHERE organization_id = $1
`

	var schema json.RawMessage
	row := d.db.QueryRowxContext(ctx, q, orgID)
	if err := row.Scan(&schema); err != nil {
		return nil, err
	}
	return schema, nil
}
