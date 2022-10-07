package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ActionType struct {
	ID             uuid.UUID       `db:"id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Name           string          `db:"name"`
	Version        int             `db:"version"`
	Properties     json.RawMessage `db:"properties"`
	Archived       bool            `db:"archived"`
	CreatedAt      time.Time       `db:"created_at"`
}

type ActionTypeProperty struct {
	Channels   []string `json:"channels"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	SourceKey  string   `json:"sourceKey"`
	SourceType string   `json:"sourceType"`
}

type ActionTypeCount struct {
	Total *int64 `json:"total"`
	User  *int64 `json:"user"`
}

type ActionTypeDateTime struct {
	Date *string `json:"date"`
	Time *string `json:"time"`
}

type CreateActionTypeParams struct {
	OrganizationID uuid.UUID
	Name           string
	Version        int
	Properties     json.RawMessage
}

func (d *DB) CreateActionType(ctx context.Context, arg CreateActionTypeParams) (ActionType, error) {
	const q = `
INSERT INTO actiontypes (
    organization_id,
	name,
	version,
    properties
) VALUES (
    $1, $2, $3, $4
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Version,
		arg.Properties,
	)
	var a ActionType
	return a, row.StructScan(&a)
}

func (d *DB) ArchiveActionType(ctx context.Context, orgID uuid.UUID, name string) error {
	const q = `
UPDATE actiontypes 
SET archived='true'
WHERE organization_id = $1
AND name = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, name)
	return err
}

func (d *DB) UnArchiveActionType(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
UPDATE actiontypes 
SET archived=false
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetActionTypeByUUID(ctx context.Context, orgID, id uuid.UUID) (ActionType, error) {
	const q = `
SELECT *
FROM actiontypes
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var c ActionType
	return c, row.StructScan(&c)
}

func (d *DB) GetActionTypes(ctx context.Context, orgID uuid.UUID) ([]ActionType, error) {
	const q = `
SELECT *
FROM actiontypes
WHERE organization_id = $1
ORDER BY name, version
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID)
}

func (d *DB) GetActionTypesByName(ctx context.Context, orgID uuid.UUID, name string) ([]ActionType, error) {
	const q = `
SELECT *
FROM actiontypes
WHERE organization_id = $1
AND name = $2
ORDER BY version
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID, name)
}

func (d *DB) GetActionTypesByNameAndVersion(ctx context.Context, orgID uuid.UUID, name string, version int) ([]ActionType, error) {
	const q = `
SELECT *
FROM actiontypes
WHERE organization_id = $1
AND name = $2
AND version = $3
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID, name, version)
}
