package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type AuthToken struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Channel        string    `db:"channel"`
	Token          string    `db:"token"`
	CreatedAt      time.Time `db:"created_at"`
}

type CreateAuthTokenParams struct {
	OrganizationID uuid.UUID
	Channel        string
	Token          string
}

func (d *DB) CreateAuthToken(ctx context.Context, arg CreateAuthTokenParams) (AuthToken, error) {
	const q = `
INSERT INTO auth_tokens (
    organization_id,
	channel,
	token
) VALUES (
    $1, $2, $3
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Channel,
		arg.Token,
	)
	var t AuthToken
	return t, row.StructScan(&t)
}

func (d *DB) DeleteAuthToken(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM auth_tokens
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetAuthTokens(ctx context.Context, orgID uuid.UUID) ([]AuthToken, error) {
	const q = `
SELECT *
FROM auth_tokens
WHERE organization_id = $1
`

	var items []AuthToken
	err := d.db.SelectContext(ctx, &items, q, orgID)
	return items, err
}

func (d *DB) GetAuthToken(ctx context.Context, token string) (AuthToken, error) {
	const q = `
SELECT *
FROM auth_tokens
WHERE token = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, token)
	var t AuthToken
	return t, row.StructScan(&t)
}
