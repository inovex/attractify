package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type AudienceProfile struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	AudienceID     uuid.UUID `db:"audience_id"`
	ProfileID      uuid.UUID `db:"profile_id"`
	SetID          uuid.UUID `db:"set_id"`
	CreatedAt      time.Time `db:"created_at"`
}

type CreateAudienceProfilesParams struct {
	OrganizationID uuid.UUID
	AudienceID     uuid.UUID
	ProfileID      uuid.UUID
	SetID          uuid.UUID
}

func (d *DB) CreateAudienceProfiles(ctx context.Context, args []CreateAudienceProfilesParams) error {
	const q = `
INSERT INTO audience_profiles (
    organization_id,
    audience_id,
    profile_id,
    set_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *
`

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	for _, arg := range args {
		if _, err := stmt.ExecContext(ctx,
			arg.OrganizationID,
			arg.AudienceID,
			arg.ProfileID,
			arg.SetID,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) DeleteAudienceProfilesBySetID(ctx context.Context, orgID, setID uuid.UUID) error {
	const q = `
DELETE FROM audience_profiles
WHERE organization_id = $1
AND set_id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, setID)
	return err
}

func (d *DB) GetAudienceProfile(ctx context.Context, orgID, profileID uuid.UUID, audienceIDs []string) (*AudienceProfile, error) {
	const q = `
SELECT *
FROM audience_profiles
WHERE organization_id = ?
AND audience_id IN(?)
AND profile_id = ?
LIMIT 1
`

	query, args, err := sqlx.In(q, orgID, audienceIDs, profileID)
	if err != nil {
		return nil, err
	}

	query = d.db.Rebind(query)
	row := d.db.QueryRowx(query, args...)

	var t AudienceProfile
	return &t, row.StructScan(&t)
}

func (d *DB) DeleteProfileFromAudiences(ctx context.Context, orgID, profileID uuid.UUID) error {
	const q = `
DELETE FROM audience_profiles
WHERE organization_id = $1
AND profile_id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, profileID)
	return err
}
