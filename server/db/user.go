package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

const (
	RoleAdmin     = "admin"
	RoleMarketeer = "marketeer"
)

type User struct {
	ID             uuid.UUID  `db:"id"`
	OrganizationID uuid.UUID  `db:"organization_id"`
	Email          string     `db:"email"`
	Password       []byte     `db:"password"`
	Salt           []byte     `db:"salt"`
	Name           string     `db:"name"`
	Role           string     `db:"role"`
	LoggedOutAt    *time.Time `db:"logged_out_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

type CreateUserParams struct {
	OrganizationID uuid.UUID
	Email          string
	Password       []byte
	Salt           []byte
	Name           string
	Role           string
}

func (d *DB) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	const q = `
INSERT INTO users (
    organization_id,
    email,
    password,
    salt,
    name,
    role
) VALUES (
    $1, lower($2), $3, $4, $5, $6
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Email,
		arg.Password,
		arg.Salt,
		arg.Name,
		arg.Role,
	)
	var u User
	return u, row.StructScan(&u)
}

func (d *DB) CreateCLIUser(arg CreateUserParams) (User, error) {
	const q = `
INSERT INTO users (
    organization_id,
    email,
    password,
    salt,
    name,
    role
) VALUES (
    $1, lower($2), $3, $4, $5, $6
)
RETURNING *
`

	row := d.db.QueryRowx(q,
		arg.OrganizationID,
		arg.Email,
		arg.Password,
		arg.Salt,
		arg.Name,
		arg.Role,
	)
	var u User
	return u, row.StructScan(&u)
}

func (d *DB) DeleteUser(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM users
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	const q = `
SELECT *
FROM users
WHERE id = $1
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, id)
	var u User
	return u, row.StructScan(&u)
}

func (d *DB) GetUserByEmail(ctx context.Context, email string) (User, error) {
	const q = `
SELECT *
FROM users
WHERE email = lower($1)
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, email)
	var u User
	return u, row.StructScan(&u)
}

func (d *DB) GetUserForOrganization(ctx context.Context, orgID, id uuid.UUID) (User, error) {
	const q = `
SELECT *
FROM users
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var u User
	return u, row.StructScan(&u)
}

func (d *DB) GetUsersForOrganization(ctx context.Context, organizationID uuid.UUID) ([]User, error) {
	const q = `
SELECT *
FROM users
WHERE organization_id = $1
`

	rows, err := d.db.QueryxContext(ctx, q, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var u User
		if err := rows.StructScan(&u); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DB) UpdateUserPassword(ctx context.Context, id uuid.UUID, password, salt []byte) error {
	const q = `
UPDATE users
SET password = $1,
    salt = $2,
    updated_at = now()
WHERE id = $3
`

	_, err := d.db.ExecContext(ctx, q, password, salt, id)
	return err
}

type UpdateUserPasswordAndNameParams struct {
	Password []byte
	Salt     []byte
	Name     string
	ID       uuid.UUID
}

func (d *DB) UpdateUserPasswordAndName(ctx context.Context, arg UpdateUserPasswordAndNameParams) error {
	const q = `
UPDATE users
SET password = $1,
    salt = $2,
    name = $3,
    updated_at = now()
WHERE id = $4
`

	_, err := d.db.ExecContext(ctx, q,
		arg.Password,
		arg.Salt,
		arg.Name,
		arg.ID,
	)
	return err
}

type UpdateUserPropertiesParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
	Email          string
	Name           string
}

func (d *DB) UpdateUserProperties(ctx context.Context, arg UpdateUserPropertiesParams) error {
	const q = `
UPDATE users
SET email = lower($1),
    name = $2,
    updated_at = now()
WHERE organization_id = $3
AND id = $4
`

	_, err := d.db.ExecContext(ctx, q,
		arg.Email,
		arg.Name,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

func (d *DB) UpdateUserRole(ctx context.Context, orgID, id uuid.UUID, role string) error {
	const q = `
UPDATE users
SET role = $1,
    updated_at = now()
WHERE organization_id = $2
AND id = $3
`

	_, err := d.db.ExecContext(ctx, q, role, orgID, id)
	return err
}

func (d *DB) UpdateUserLoggedOutAt(ctx context.Context, orgID, id uuid.UUID, loggedOutAt time.Time) error {
	const q = `
UPDATE users
SET logged_out_at = $1,
    updated_at = now()
WHERE organization_id = $2
AND id = $3
`

	_, err := d.db.ExecContext(ctx, q, loggedOutAt, orgID, id)
	return err
}
