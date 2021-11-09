package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

const (
	ComputedTraitTypeCountEvents     = "count_events"
	ComputedTraitTypeAggregation     = "aggregation"
	ComputedTraitTypeMostFrequent    = "most_frequent"
	ComputedTraitTypeFirst           = "first_event"
	ComputedTraitTypeLast            = "last_event"
	ComputedTraitTypeUniqueList      = "unique_list"
	ComputedTraitTypeUniqueListCount = "unique_list_count"
)

type ComputedTraitType string

type ComputedTrait struct {
	ID             uuid.UUID         `db:"id"`
	OrganizationID uuid.UUID         `db:"organization_id"`
	Name           string            `db:"name"`
	Key            string            `db:"key"`
	Type           ComputedTraitType `db:"type"`
	EventID        uuid.UUID         `db:"event_id"`
	Conditions     json.RawMessage   `db:"conditions"`
	Properties     json.RawMessage   `db:"properties"`
	CreatedAt      time.Time         `db:"created_at"`
	UpdatedAt      time.Time         `db:"updated_at"`
	RefreshedAt    time.Time         `db:"refreshed_at"`
}

type CreateComputedTraitParams struct {
	OrganizationID uuid.UUID
	Name           string
	Key            string
	Type           string
	EventID        uuid.UUID
	Conditions     json.RawMessage
	Properties     json.RawMessage
}

func (d *DB) CreateComputedTrait(ctx context.Context, arg CreateComputedTraitParams) (ComputedTrait, error) {
	const q = `
INSERT INTO computed_traits (
	organization_id,
	name,
	key,
	type,
	event_id,
	conditions,
	properties
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Key,
		arg.Type,
		arg.EventID,
		arg.Conditions,
		arg.Properties,
	)
	var a ComputedTrait
	return a, row.StructScan(&a)
}

func (d *DB) DeleteComputedTrait(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM computed_traits
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetComputedTrait(ctx context.Context, orgID, id uuid.UUID) (ComputedTrait, error) {
	const q = `
SELECT *
FROM computed_traits
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var a ComputedTrait
	return a, row.StructScan(&a)
}

func (d *DB) GetComputedTraits(ctx context.Context, organizationID uuid.UUID) ([]ComputedTrait, error) {
	const q = `
SELECT *
FROM computed_traits
WHERE organization_id = $1
`

	var items []ComputedTrait
	if err := d.db.SelectContext(ctx, &items, q, organizationID); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) GetComputedTraitsForEvent(ctx context.Context, orgID, id uuid.UUID) ([]ComputedTrait, error) {
	const q = `
SELECT *
FROM computed_traits
WHERE organization_id = $1
AND event_id = $2
`

	var items []ComputedTrait
	if err := d.db.SelectContext(ctx, &items, q, orgID, id); err != nil {
		return nil, err
	}
	return items, nil
}

type UpdateComputedTraitParams struct {
	Name           string
	Key            string
	EventID        uuid.UUID
	Conditions     json.RawMessage
	Properties     json.RawMessage
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateComputedTrait(ctx context.Context, arg UpdateComputedTraitParams) error {
	const q = `
UPDATE computed_traits
SET name = $1,
	key = $2,
	event_id = $3,
	conditions = $4,
	properties = $5,
	updated_at = now()
WHERE organization_id = $6
AND id = $7
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.Name,
		arg.Key,
		arg.EventID,
		arg.Conditions,
		arg.Properties,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}
