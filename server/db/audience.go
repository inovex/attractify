package db

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	TimeWindowTypeAnytime   = "anytime"
	TimeWindowTypeWithin    = "within"
	TimeWindowTypeInBetween = "in_between"
)

type Audience struct {
	ID               uuid.UUID       `db:"id"`
	OrganizationID   uuid.UUID       `db:"organization_id"`
	Name             string          `db:"name"`
	Description      string          `db:"description"`
	IncludeAnonymous bool            `db:"include_anonymous"`
	Events           json.RawMessage `db:"events"`
	Traits           json.RawMessage `db:"traits"`
	ProfileCount     int             `db:"profile_count"`
	CurrentSetID     uuid.NullUUID   `db:"current_set_id"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
	RefreshedAt      time.Time       `db:"refreshed_at"`
}

type CreateAudienceParams struct {
	OrganizationID   uuid.UUID
	Name             string
	Description      string
	IncludeAnonymous bool
	Events           json.RawMessage
	Traits           json.RawMessage
}

func (d *DB) CreateAudience(ctx context.Context, arg CreateAudienceParams) (Audience, error) {
	const q = `
INSERT INTO audiences (
	organization_id,
  name,
	description,
	include_anonymous,
	events,
	traits
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Description,
		arg.IncludeAnonymous,
		arg.Events,
		arg.Traits,
	)
	var a Audience
	return a, row.StructScan(&a)
}

type DeleteAudienceParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) DeleteAudience(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM audiences
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetAudience(ctx context.Context, orgID, id uuid.UUID) (Audience, error) {
	const q = `
SELECT *
FROM audiences
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var a Audience
	return a, row.StructScan(&a)
}

func (d *DB) GetAudiences(ctx context.Context, organizationID uuid.UUID) ([]Audience, error) {
	const q = `
SELECT *
FROM audiences
WHERE organization_id = $1
`

	var items []Audience
	if err := d.db.SelectContext(ctx, &items, q, organizationID); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) GetAllAudiences(ctx context.Context) ([]Audience, error) {
	const q = `
SELECT *
FROM audiences
`

	var items []Audience
	if err := d.db.SelectContext(ctx, &items, q); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) ValidateAudience(ctx context.Context, orgID uuid.UUID, ids []string) error {
	const q = `
SELECT count(*)
FROM audiences
WHERE organization_id = ?
AND id IN (?)
`

	query, args, err := sqlx.In(q, orgID, ids)
	if err != nil {
		return err
	}

	query = d.db.Rebind(query)
	row := d.db.QueryRowxContext(ctx, query, args...)

	var count int
	row.Scan(&count)

	if count != len(ids) {
		return errors.New("not all audience IDs found for organization")
	}
	return nil
}

type UpdateAudienceParams struct {
	Name             string
	Description      string
	IncludeAnonymous bool
	Events           json.RawMessage
	Traits           json.RawMessage
	OrganizationID   uuid.UUID
	ID               uuid.UUID
}

func (d *DB) UpdateAudience(ctx context.Context, arg UpdateAudienceParams) error {
	const q = `
UPDATE audiences
SET name = $1,
	description = $2,
	include_anonymous = $3,
	events = $4,
	traits = $5,
	updated_at = now()
WHERE organization_id = $6
AND id = $7
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.Name,
		arg.Description,
		arg.IncludeAnonymous,
		arg.Events,
		arg.Traits,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

type UpdateAudienceProfilesParams struct {
	CurrentSetID   uuid.UUID
	ProfileCount   int
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateAudienceProfiles(ctx context.Context, arg UpdateAudienceProfilesParams) error {
	const q = `
UPDATE audiences
SET current_set_id = $1,
	profile_count = $2,
	updated_at = now(),
	refreshed_at = now()
WHERE organization_id = $3
AND id = $4
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.CurrentSetID,
		arg.ProfileCount,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}
