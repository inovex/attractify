package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ProfileIdentityWithTraits struct {
	ID             uuid.UUID       `db:"id" json:"-"`
	OrganizationID uuid.UUID       `db:"organization_id" json:"-"`
	ProfileID      uuid.UUID       `db:"profile_id" json:"-"`
	Channel        string          `db:"channel" json:"channel"`
	Type           string          `db:"type" json:"type"`
	UserID         string          `db:"user_id" json:"userId"`
	CustomTraits   json.RawMessage `db:"custom_traits" json:"customTraits"`
	ComputedTraits json.RawMessage `db:"computed_traits" json:"computedTraits"`
	IsAnonymous    bool            `db:"is_anonymous" json:"isAnonymous"`
	CreatedAt      time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time       `db:"updated_at" json:"deletedAt"`
}

type ProfileIdentity struct {
	ID             uuid.UUID `db:"id" json:"-"`
	OrganizationID uuid.UUID `db:"organization_id" json:"-"`
	ProfileID      uuid.UUID `db:"profile_id" json:"-"`
	Channel        string    `db:"channel" json:"channel"`
	Type           string    `db:"type" json:"type"`
	UserID         string    `db:"user_id" json:"userId"`
	IsAnonymous    bool      `db:"is_anonymous" json:"isAnonymous"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"deletedAt"`
}

type CreateProfileIdentityParams struct {
	OrganizationID uuid.UUID
	ProfileID      uuid.UUID
	Channel        string
	Type           string
	UserID         string
	IsAnonymous    bool
	CreatedAt      time.Time
}

func (d *DB) CreateProfileIdentity(ctx context.Context, arg CreateProfileIdentityParams) (ProfileIdentity, error) {
	const q = `
INSERT INTO profile_identities (
  organization_id,
	profile_id,
	channel,
	type,
	user_id,
	is_anonymous,
	created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.ProfileID,
		arg.Channel,
		arg.Type,
		arg.UserID,
		arg.IsAnonymous,
		arg.CreatedAt.UTC(),
	)
	var p ProfileIdentity
	return p, row.StructScan(&p)
}

func (d *DB) GetProfileIdentitiesForProfile(ctx context.Context, orgID, profileID uuid.UUID) ([]ProfileIdentity, error) {
	const q = `
SELECT *
FROM profile_identities
WHERE organization_id = $1
AND profile_id = $2
`

	var items []ProfileIdentity
	if err := d.db.SelectContext(ctx, &items, q, orgID, profileID); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) GetProfileIdentityForUserID(ctx context.Context, orgID uuid.UUID, userID string) (ProfileIdentity, error) {
	const q = `
SELECT *
FROM profile_identities
WHERE organization_id = $1
AND user_id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, userID)
	var p ProfileIdentity
	return p, row.StructScan(&p)
}

func (d *DB) GetProfileIdentitiesForUserID(ctx context.Context, orgID uuid.UUID, userID string) ([]ProfileIdentity, error) {
	const q = `
SELECT *
FROM profile_identities
WHERE organization_id = $1
AND user_id = $2
`

	var items []ProfileIdentity
	if err := d.db.SelectContext(ctx, &items, q, orgID, userID); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) DeleteProfileIdentitiesByProfileID(ctx context.Context, orgID, profileID uuid.UUID) error {
	const q = `
DELETE FROM profile_identities
WHERE organization_id = $1
AND profile_id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, profileID)
	return err
}

func (d *DB) UpdateProfileIdentitiesWithProfileID(ctx context.Context, orgID, oldID, newID uuid.UUID) error {
	const q = `
UPDATE profile_identities
SET profile_id = $1
WHERE organization_id = $2
AND profile_id = $3
`

	_, err := d.db.ExecContext(ctx, q, newID, orgID, oldID)
	return err
}

type UpdateProfileIdentityParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
	UpdatedAt      time.Time
	Type           string
	IsAnonymous    bool
}

func (d *DB) UpdateProfileIdentity(ctx context.Context, arg UpdateProfileIdentityParams) error {
	const q = `
UPDATE profile_identities
SET type = $1,
	is_anonymous = $2,
	updated_at = $3
WHERE organization_id = $4
AND id = $5
`

	_, err := d.db.ExecContext(ctx, q,
		arg.Type,
		arg.IsAnonymous,
		arg.UpdatedAt,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

func (d *DB) CountProfileIdentities(ctx context.Context, orgID, profileID uuid.UUID) (int, error) {
	const q = `
SELECT count(*)
FROM profile_identities
WHERE organization_id = $1
AND profile_id = $2
`

	row := d.db.QueryRowContext(ctx, q, orgID, profileID)
	var count int
	return count, row.Scan(&count)
}
