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
	IsArchived     bool            `db:"is_archived"`
	IsInUse        bool            `db:"is_in_use"`
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

type UpdateActionTypeParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Properties     json.RawMessage
}

func (d *DB) CreateActionType(ctx context.Context, arg CreateActionTypeParams) (ActionType, error) {
	const q = `
INSERT INTO action_types(
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

func (d *DB) UpdateActionType(ctx context.Context, args UpdateActionTypeParams) error {
	const q = `
UPDATE action_types
SET properties = $1
WHERE organization_id = $2
AND id = $3
`

	_, err := d.db.ExecContext(ctx, q, args.Properties, args.OrganizationID, args.ID)
	return err
}

func (d *DB) ArchiveActionType(ctx context.Context, orgID uuid.UUID, name string) error {
	const q = `
UPDATE action_types
SET is_archived = true
WHERE organization_id = $1
AND name = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, name)
	return err
}

func (d *DB) UnArchiveActionType(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
UPDATE action_types
SET is_archived = false
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetActionTypeByUUID(ctx context.Context, orgID, id uuid.UUID) (ActionType, error) {
	const q = `
SELECT DISTINCT t.*, 
CASE WHEN a.id is not null THEN 'True' ELSE 'False' END as is_in_use
FROM action_types t LEFT JOIN actions a ON a.type_id = t.id
WHERE t.organization_id = $1
AND t.id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var c ActionType
	return c, row.StructScan(&c)
}

func (d *DB) GetActionTypes(ctx context.Context, orgID uuid.UUID) ([]ActionType, error) {
	const q = `
SELECT DISTINCT t.*, 
CASE WHEN a.id is not null THEN 'True' ELSE 'False' END as is_in_use
FROM action_types t LEFT JOIN actions a ON a.type_id = t.id
WHERE t.organization_id = $1
ORDER BY t.name, t.version
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID)
}

func (d *DB) GetActionTypesByName(ctx context.Context, orgID uuid.UUID, name string) ([]ActionType, error) {
	const q = `
SELECT DISTINCT t.*, 
CASE WHEN a.id is not null THEN 'True' ELSE 'False' END as is_in_use
FROM action_types t LEFT JOIN actions a ON a.type_id = t.id
WHERE t.organization_id = $1
AND t.name = $2
ORDER BY t.version
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID, name)
}

func (d *DB) GetActionTypesByNameAndVersion(ctx context.Context, orgID uuid.UUID, name string, version int) ([]ActionType, error) {
	const q = `
SELECT DISTINCT t.*, 
CASE WHEN a.id is not null THEN 'True' ELSE 'False' END as is_in_use
FROM action_types t LEFT JOIN actions a ON a.type_id = t.id
WHERE t.organization_id = $1
AND t.name = $2
AND t.version = $3
`

	var items []ActionType
	return items, d.db.SelectContext(ctx, &items, q, orgID, name, version)
}

func (d *DB) GetNewActionTypeVersion(ctx context.Context, orgID uuid.UUID, name string) (int, error) {
	const q = `
SELECT COUNT(*)
FROM action_types t 
WHERE t.organization_id = $1
AND t.name = $2
`

	var version []int
	return (version[0] + 1), d.db.SelectContext(ctx, &version, q, orgID, name)
}
