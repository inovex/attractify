package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Channel struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Name           string    `db:"name"`
	Key            string    `db:"key"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type CreateChannelParams struct {
	OrganizationID uuid.UUID
	Name           string
	Key            string
}

func (d *DB) CreateChannel(ctx context.Context, arg CreateChannelParams) (Channel, error) {
	const q = `
INSERT INTO channels (
    organization_id,
	name,
	key
) VALUES (
    $1, $2, $3
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Key,
	)
	var t Channel
	return t, row.StructScan(&t)
}

type DeleteChannelParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) DeleteChannel(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM channels
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetChannels(ctx context.Context, orgID uuid.UUID) ([]Channel, error) {
	const q = `
SELECT *
FROM channels
WHERE organization_id = $1
`

	var items []Channel
	err := d.db.SelectContext(ctx, &items, q, orgID)
	return items, err
}

type GetChannelParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) GetChannel(ctx context.Context, orgID, id uuid.UUID) (Channel, error) {
	const q = `
SELECT *
FROM channels
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var t Channel
	return t, row.StructScan(&t)
}

type UpdateChannelParams struct {
	Name           string
	Key            string
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateChannel(ctx context.Context, arg UpdateChannelParams) error {
	const q = `
UPDATE channels
SET name = $1,
	  key = $2,
    updated_at = now()
WHERE organization_id = $3
AND id = $4
`

	_, err := d.db.ExecContext(ctx, q,
		arg.Name,
		arg.Key,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}
