package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

const (
	TraitTypeString = "string"
	TraitTypeNumber = "number"
	TraitTypeBool   = "bool"
	TraitTypeObject = "object"
	TraitTypeArray  = "array"
)

type TraitType string

type Trait interface{}

type Traits map[string]Trait

type Profile struct {
	ID             uuid.UUID       `db:"id" json:"-"`
	OrganizationID uuid.UUID       `db:"organization_id" json:"-"`
	CustomTraits   json.RawMessage `db:"custom_traits" json:"customTraits"`
	ComputedTraits json.RawMessage `db:"computed_traits" json:"computedTraits"`
	CreatedAt      time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updatedAt"`
}

type CreateProfileParams struct {
	OrganizationID uuid.UUID
	CustomTraits   json.RawMessage
	CreatedAt      time.Time
}

func (d *DB) CreateProfile(ctx context.Context, arg CreateProfileParams) (Profile, error) {
	const q = `
INSERT INTO profiles (
    organization_id,
	custom_traits,
	created_at
) VALUES (
    $1, $2, $3
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.CustomTraits,
		arg.CreatedAt.UTC(),
	)
	var p Profile
	return p, row.StructScan(&p)
}

func (d *DB) GetProfilesForOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]Profile, error) {
	const q = `
SELECT *
FROM profiles
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3
`

	var items []Profile
	return items, d.db.SelectContext(ctx, &items, q, orgID, limit, offset)
}

func (d *DB) GetProfile(ctx context.Context, orgID, id uuid.UUID) (Profile, error) {
	const q = `
SELECT *
FROM profiles
WHERE organization_id = $1
AND id = $2
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var p Profile
	return p, row.StructScan(&p)
}

func (d *DB) SearchByUserID(ctx context.Context, orgID uuid.UUID, id string) ([]ProfileIdentityWithTraits, error) {
	q := `
SELECT i.*, p.custom_traits, p.computed_traits
FROM profiles p 
INNER JOIN profile_identities i 
ON p.id = i.profile_id
WHERE p.organization_id = $1
AND i.user_id = $2
`

	var items []ProfileIdentityWithTraits
	return items, d.db.SelectContext(ctx, &items, q, orgID, id)
}

type UpdateProfileParams struct {
	CustomTraits   json.RawMessage
	ComputedTraits json.RawMessage
	OrganizationID uuid.UUID
	ID             uuid.UUID
	UpdatedAt      time.Time
}

func (d *DB) UpdateProfile(ctx context.Context, arg UpdateProfileParams) error {
	const q = `
UPDATE profiles
SET custom_traits = $1,
	computed_traits = $2,
    updated_at = $3
WHERE organization_id = $4
AND id = $5
`

	_, err := d.db.ExecContext(
		ctx,
		q,
		arg.CustomTraits,
		arg.ComputedTraits,
		arg.UpdatedAt,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

type UpdateProfileComputedTraitsParams struct {
	ComputedTraits json.RawMessage
	UpdatedAt      time.Time
	ID             uuid.UUID
}

func (d *DB) UpdateProfileComputedTraits(ctx context.Context, arg UpdateProfileComputedTraitsParams) error {
	const q = `
UPDATE profiles
SET computed_traits = $1,
    updated_at = $2
WHERE id = $3
`

	_, err := d.db.ExecContext(ctx, q,
		arg.ComputedTraits,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}

func (d *DB) DeleteProfile(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM profiles
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

type ProfileBucket struct {
	Bucket time.Time `db:"bucket"`
	Count  int       `db:"count"`
}

func (d *DB) GetNewProfilesLast24h(ctx context.Context, organizationID uuid.UUID) ([]ProfileBucket, error) {
	const q = `
SELECT now() - (
		(8 - width_bucket(CAST(extract('hour' FROM created_at) AS DECIMAL), 0.0, 24.0, 8)
	) * interval '3 hours') AS bucket,
	count(*) AS count
FROM profiles
WHERE created_at >= now() - interval '1 day'
AND organization_id = $1
GROUP BY bucket
ORDER BY bucket ASC
`

	var items []ProfileBucket
	return items, d.db.SelectContext(ctx, &items, q, organizationID)
}
