package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

const (
	StateInactive = "inactive"
	StateStaging  = "staging"
	StateActive   = "active"

	GroupAll  = "all"
	GroupUser = "user"

	ActionHookTypeWebhook    = "webhook"
	ActionHookTypeTrackEvent = "track_event"

	TraitConditionTypeCustom   = "custom"
	TraitConditionTypeComputed = "computed"
)

type ActionState string

type Action struct {
	ID             uuid.UUID       `db:"id"`
	OrganizationID uuid.UUID       `db:"organization_id"`
	Name           string          `db:"name"`
	Type           string          `db:"type"`
	Tags           json.RawMessage `db:"tags"`
	Properties     json.RawMessage `db:"properties"`
	Targeting      json.RawMessage `db:"targeting"`
	Capping        json.RawMessage `db:"capping"`
	Hooks          json.RawMessage `db:"hooks"`
	State          ActionState     `db:"state"`
	TestUsers      json.RawMessage `db:"test_users"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

type ActionTags []string
type ActionHookType string

type ActionProperty struct {
	Channels   []string `json:"channels"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	SourceKey  string   `json:"sourceKey"`
	SourceType string   `json:"sourceType"`
}

type ActionCount struct {
	Total *int64 `json:"total"`
	User  *int64 `json:"user"`
}

type ActionDateTime struct {
	Date *string `json:"date"`
	Time *string `json:"time"`
}

type ActionCapping struct {
	Channels []string `json:"channels"`
	Event    string   `json:"event"`
	Group    string   `json:"group"`
	Count    int      `json:"count"`
	Within   int      `json:"within"`
}

type ContextCondition struct {
	Channel  string      `json:"channel"`
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type TraitCondition struct {
	Source   string      `json:"source"`
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type ActionTargeting struct {
	Audiences         []string           `json:"audiences"`
	Channels          []string           `json:"channels"`
	TraitConditions   []TraitCondition   `json:"traitConditions"`
	ContextConditions []ContextCondition `json:"contextConditions"`
	Start             ActionDateTime     `json:"start"`
	End               ActionDateTime     `json:"end"`
}

type ActionProperties struct {
	Channels []string `json:"channels"`
	URLs     []string `json:"urls"`
}

type ActionHook struct {
	Channels   []string        `json:"channels"`
	Event      string          `json:"event"`
	Type       ActionHookType  `json:"type"`
	Properties json.RawMessage `json:"properties"`
}

type ActionTestUser struct {
	Channels      []string `json:"channels"`
	UserID        string   `json:"userId"`
	Description   string   `json:"description"`
	SkipTargeting bool     `json:"skipTargeting"`
}

type CreateActionParams struct {
	OrganizationID uuid.UUID
	Name           string
	Type           string
	Tags           json.RawMessage
	State          ActionState
	Properties     json.RawMessage
	Targeting      json.RawMessage
	Capping        json.RawMessage
	Hooks          json.RawMessage
	TestUsers      json.RawMessage
}

func (d *DB) CreateAction(ctx context.Context, arg CreateActionParams) (Action, error) {
	const q = `
INSERT INTO actions (
    organization_id,
    name,
	type,
	tags,
	state,
    properties,
    targeting,
    capping,
	hooks,
	test_users
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *
`

	row := d.db.QueryRowxContext(ctx, q,
		arg.OrganizationID,
		arg.Name,
		arg.Type,
		arg.Tags,
		arg.State,
		arg.Properties,
		arg.Targeting,
		arg.Capping,
		arg.Hooks,
		arg.TestUsers,
	)
	var a Action
	return a, row.StructScan(&a)
}

func (d *DB) DeleteAction(ctx context.Context, orgID, id uuid.UUID) error {
	const q = `
DELETE FROM actions
WHERE organization_id = $1
AND id = $2
`

	_, err := d.db.ExecContext(ctx, q, orgID, id)
	return err
}

func (d *DB) GetAction(ctx context.Context, orgID, id uuid.UUID) (Action, error) {
	const q = `
SELECT *
FROM actions
WHERE organization_id = $1
AND id = $2
LIMIT 1
`

	row := d.db.QueryRowxContext(ctx, q, orgID, id)
	var c Action
	return c, row.StructScan(&c)
}

func (d *DB) GetActions(ctx context.Context, orgID uuid.UUID) ([]Action, error) {
	const q = `
SELECT *
FROM actions
WHERE organization_id = $1
`

	var items []Action
	return items, d.db.SelectContext(ctx, &items, q, orgID)
}

func (d *DB) GetActiveActionsCount(ctx context.Context, orgID uuid.UUID) (int, error) {
	const q = `
SELECT count(*)
FROM actions
WHERE organization_id = $1
AND state = 'active'
LIMIT 1
`

	row := d.db.QueryRowContext(ctx, q, orgID)
	var count int
	return count, row.Scan(&count)
}

type UpdateActionParams struct {
	Name           string
	Type           string
	Tags           json.RawMessage
	State          ActionState
	Properties     json.RawMessage
	Targeting      json.RawMessage
	Capping        json.RawMessage
	Hooks          json.RawMessage
	TestUsers      json.RawMessage
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateAction(ctx context.Context, arg UpdateActionParams) error {
	const q = `
UPDATE actions
SET name = $1,
	type = $2,
	tags = $3,
	state = $4,
    properties = $5,
    targeting = $6,
    capping = $7,
	hooks = $8,
	test_users = $9,
    updated_at = now()
WHERE organization_id = $10
AND id = $11
`

	_, err := d.db.ExecContext(ctx, q,
		arg.Name,
		arg.Type,
		arg.Tags,
		arg.State,
		arg.Properties,
		arg.Targeting,
		arg.Capping,
		arg.Hooks,
		arg.TestUsers,
		arg.OrganizationID,
		arg.ID,
	)
	return err
}

type UpdateActionStateParams struct {
	State          ActionState
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (d *DB) UpdateActionState(ctx context.Context, orgID, id uuid.UUID, state ActionState) error {
	const q = `
UPDATE actions
SET state = $1, updated_at = now()
WHERE organization_id = $2
AND id = $3
`

	_, err := d.db.ExecContext(ctx, q, state, orgID, id)
	return err
}
