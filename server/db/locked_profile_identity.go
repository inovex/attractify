package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type LockedProfileIdentity struct {
	ID             uuid.UUID `db:"id" json:"-"`
	OrganizationID uuid.UUID `db:"organization_id" json:"-"`
	UserID         string    `db:"user_id" json:"userId"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"deletedAt"`
}

func (d *DB) CreateLockedProfileIdentity(ctx context.Context, orgID uuid.UUID, userID string) (LockedProfileIdentity, error) {
	const q = `
INSERT INTO locked_profile_identities (
    organization_id,
	user_id
) VALUES (
    $1, $2
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q, orgID, userID)
	var p LockedProfileIdentity
	return p, row.StructScan(&p)
}

func (d *DB) GetLockedProfileIdentities(ctx context.Context, organizationID uuid.UUID) ([]LockedProfileIdentity, error) {
	const q = `
SELECT *
FROM locked_profile_identities
WHERE organization_id = $1
`

	var items []LockedProfileIdentity
	if err := d.db.SelectContext(
		ctx,
		&items,
		q,
		organizationID,
	); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) GetLockedProfileIdentityForUserID(ctx context.Context, orgID uuid.UUID, userID string) (LockedProfileIdentity, error) {
	const q = `
SELECT *
FROM locked_profile_identities
WHERE organization_id = $1
AND user_id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, userID)
	var p LockedProfileIdentity
	return p, row.StructScan(&p)
}

func (d *DB) DeleteLockedProfileIdentitiesByID(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM locked_profile_identities
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}
