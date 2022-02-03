package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Key       []byte    `db:"key"`
	Timezone  string    `db:"timezone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrganizationAuth struct {
	OrganizationID uuid.UUID `db:"id"`
	Key            []byte    `db:"key"`
	Timezone       string    `db:"timezone"`
	Channel        string    `db:"channel"`
}

func (d *DB) CreateOrganization(ctx context.Context, name, timezone string, key []byte) (Organization, error) {
	const q = `
INSERT INTO organizations (
    name,
    key,
    timezone
) VALUES ($1, $2, $3)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q, name, key, timezone)
	var o Organization
	return o, row.StructScan(&o)
}

func (d *DB) CreateCLIOrganization(name string, timezone string, key []byte) (Organization, error) {
	const q = `
INSERT INTO organizations (
    name,
    key,
    timezone
) VALUES ($1, $2, $3)
RETURNING *
`

	row := d.db.QueryRowx(q, name, key, timezone)
	var o Organization
	return o, row.StructScan(&o)
}

func (d *DB) GetKeyForOrganization(ctx context.Context, id uuid.UUID) ([]byte, error) {
	const q = `
SELECT key
FROM organizations
WHERE id = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, id)
	var key []byte
	err := row.Scan(&key)
	return key, err
}

func (d *DB) GetOrganization(ctx context.Context, id uuid.UUID) (Organization, error) {
	const q = `
SELECT *
FROM organizations
WHERE id = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, id)
	var o Organization
	return o, row.StructScan(&o)
}

func (d *DB) GetOrganizationByToken(ctx context.Context, token string) (OrganizationAuth, error) {
	const q = `
SELECT organizations.id, organizations.key, organizations.timezone, auth_tokens.channel
FROM organizations
LEFT JOIN auth_tokens
ON auth_tokens.organization_id = organizations.id
WHERE auth_tokens.token = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, token)
	var o OrganizationAuth
	return o, row.StructScan(&o)
}

func (d *DB) UpdateOrganization(ctx context.Context, id uuid.UUID, name, timezone string) error {
	const q = `
UPDATE organizations
SET name = $1, timezone = $2, updated_at = now()
WHERE id = $3
`

	_, err := d.db.ExecContext(ctx, q, name, timezone, id)
	return err
}
