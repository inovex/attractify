package analytics

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	ReactionEventDelivered = "delivered"
	ReactionEventShown     = "shown"
	ReactionEventHidden    = "hidden"
	ReactionEventDeclined  = "declined"
	ReactionEventAccepted  = "accepted"
)

type ReactionAggregation struct {
	Channel  string  `db:"channel"`
	Event    string  `db:"event"`
	Total    int     `db:"total"`
	Duration float64 `db:"duration"`
	Year     int     `db:"year"`
	Month    int     `db:"month"`
	Day      int     `db:"day"`
	WeekDay  int     `db:"weekDay"`
	Hour     int     `db:"hour"`
}

type ReactionBucket struct {
	Bucket time.Time `db:"bucket"`
	Count  int       `db:"count"`
}

type Reaction struct {
	ID             uuid.UUID `db:"id" json:"-"`
	OrganizationID uuid.UUID `db:"organization_id" json:"-"`
	ActionID       uuid.UUID `db:"action_id" json:"actionId"`
	IdentityID     uuid.UUID `db:"identity_id" json:"-"`
	Channel        string    `db:"channel" json:"channel"`
	Event          string    `db:"event" json:"event"`
	Context        string    `db:"context" json:"context"`
	Properties     string    `db:"properties" json:"properties"`
	Result         string    `db:"result" json:"result"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	FullCount      int       `db:"full_count"`
}

const getReactions = `
SELECT *, count(*) over() AS full_count
FROM reactions
WHERE organization_id = ?
%s
%s
%s
AND created_at BETWEEN ? AND ?
ORDER BY created_at DESC
LIMIT ?, ?
`

type GetReactionsParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Events         []string
	UserID         string
	Start          time.Time
	End            time.Time
	Limit          int
	Offset         int
}

func (a Analytics) GetReactions(arg GetReactionsParams) ([]Reaction, error) {
	var (
		items           []Reaction
		queryAction     = ""
		queryEvents     = ""
		queryIdentities = ""
		args            = []interface{}{
			arg.OrganizationID,
		}
	)

	if arg.ActionID != uuid.Nil {
		queryAction = "AND action_id = ?"
		args = append(args, arg.ActionID)
	}

	if len(arg.Events) > 0 {
		ev := []string{}
		for _, e := range arg.Events {
			ev = append(ev, escapeParameter(e))
		}
		queryEvents = fmt.Sprintf("AND event IN ('%s')", strings.Join(ev, "','"))
	}

	if len(arg.UserID) > 0 {
		queryIdentities = fmt.Sprintf("AND identity_id GLOBAL IN (SELECT id FROM identities WHERE organization_id = '%s' AND user_id = '%s')", arg.OrganizationID, arg.UserID)
	}

	args = append(args,
		arg.Start.Format("2006-01-02 15:04:05"),
		arg.End.Format("2006-01-02 15:04:05"),
		arg.Offset,
		arg.Limit)

	query := fmt.Sprintf(getReactions, queryAction, queryEvents, queryIdentities)
	return items, a.DB.Select(&items, query, args...)
}

const createReaction = `
INSERT INTO reactions (
	organization_id,
	action_id,
	identity_id,
	channel,
	event,
	context,
	properties,
	result,
	created_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?
)`

type CreateReactionParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	IdentityID     uuid.UUID
	Channel        string
	Event          string
	Context        string
	Properties     string
	Result         string
	CreatedAt      time.Time
}

func (a Analytics) CreateReaction(arg CreateReactionParams) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(createReaction,
		arg.OrganizationID,
		arg.ActionID,
		arg.IdentityID,
		arg.Channel,
		arg.Event,
		arg.Context,
		arg.Properties,
		arg.Result,
		arg.CreatedAt.UTC(),
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (a Analytics) CreateReactions(arg []CreateReactionParams) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createReaction)
	if err != nil {
		return err
	}
	for _, a := range arg {
		_, err = stmt.Exec(
			a.OrganizationID,
			a.ActionID,
			a.IdentityID,
			a.Channel,
			a.Event,
			a.Context,
			a.Properties,
			a.Result,
			a.CreatedAt.UTC(),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

const getReactionCount = `
SELECT count(*)
FROM reactions
WHERE organization_id = ?
AND action_id = ?
%s
%s
AND event = ?
%s
`

type GetReactionCountParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Channels       []string
	Event          string
	IsUser         bool
	ProfileID      uuid.UUID
	Within         int
}

func escapeParameter(v string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(v)
}

func (a Analytics) GetReactionCount(arg GetReactionCountParams) (int, error) {
	var count int

	var (
		identity = ""
		channels = ""
		within   = ""
	)
	if arg.IsUser {
		identity = fmt.Sprintf(
			"AND identity_id GLOBAL IN (SELECT id FROM identities WHERE organization_id = '%s' AND profile_id = '%s')",
			arg.OrganizationID.String(),
			arg.ProfileID.String(),
		)
	}
	if len(arg.Channels) > 0 {
		var chs []string
		for _, c := range arg.Channels {
			chs = append(chs, escapeParameter(c))
		}
		channels = fmt.Sprintf("AND channel IN ('%s')", strings.Join(chs, "','"))
	}
	if arg.Within > 0 {
		within = fmt.Sprintf("AND created_at BETWEEN subtractDays(now(), %d) AND now()", arg.Within)
	}
	query := fmt.Sprintf(getReactionCount, identity, channels, within)
	row := a.DB.QueryRowx(query,
		arg.OrganizationID,
		arg.ActionID,
		arg.Event,
	)

	return count, row.Scan(&count)
}

const deleteReaction = `
ALTER TABLE reactions
DELETE
WHERE organization_id = ?
AND id = ?
`

type DeleteReactionParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (a Analytics) DeleteReaction(arg DeleteReactionParams) error {
	_, err := a.DB.Exec(deleteReaction, arg.OrganizationID, arg.ID)
	return err
}

const getReactionEventCount = `
SELECT event, count(*) as total
FROM reactions
WHERE organization_id = ?
AND action_id = ?
AND created_at BETWEEN ? AND ?
GROUP BY event
`

type GetReactionEventCountParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Timezone       string
	Start          string
	End            string
}

func (a Analytics) GetReactionEventCount(arg GetReactionEventCountParams) ([]ReactionAggregation, error) {
	var res []ReactionAggregation
	if err := a.DB.Select(&res, getReactionEventCount,
		arg.OrganizationID,
		arg.ActionID,
		arg.Start,
		arg.End,
	); err != nil {
		return nil, err
	}

	return res, nil
}

const getReactionCountForInterval = `
SELECT event, count(*) as total, %s AS %s
FROM reactions
WHERE organization_id = ?
AND action_id = ?
AND created_at BETWEEN ? AND ?
GROUP BY event, %s
ORDER BY %s ASC
`

type GetReactionEventCountForIntervalParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Timezone       string
	Start          string
	End            string
	Interval       string
}

func (a Analytics) GetEventOverview(arg GetReactionEventCountForIntervalParams) ([]ReactionAggregation, error) {
	aggregate := ""
	interval := arg.Interval
	switch arg.Interval {
	case "year":
		aggregate = "toYear(created_at, ?)"
	case "month":
		aggregate = "toMonth(created_at, ?)"
	case "day":
		aggregate = "toDayOfMonth(created_at, ?)"
	case "week_day":
		aggregate = "toDayOfWeek(created_at, ?)"
		interval = "weekDay"
	case "hour":
		aggregate = "toHour(created_at, ?)"
	default:
		aggregate = "toDayOfMonth(created_at, ?)"
		interval = "day"
	}

	query := fmt.Sprintf(getReactionCountForInterval, aggregate, interval, interval, interval)

	var res []ReactionAggregation
	if err := a.DB.Select(&res, query,
		arg.Timezone,
		arg.OrganizationID,
		arg.ActionID,
		arg.Start,
		arg.End,
	); err != nil {
		return nil, err
	}

	return res, nil
}

const getReactionChannelCount = `
SELECT channel, count(*) as total
FROM reactions
WHERE organization_id = ?
AND action_id = ?
AND created_at BETWEEN ? AND ?
GROUP BY channel
`

type GetReactionChannelCountParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Timezone       string
	Start          string
	End            string
}

func (a Analytics) GetReactionChannelCount(arg GetReactionChannelCountParams) ([]ReactionAggregation, error) {
	var res []ReactionAggregation
	if err := a.DB.Select(&res, getReactionChannelCount,
		arg.OrganizationID,
		arg.ActionID,
		arg.Start,
		arg.End,
	); err != nil {
		return nil, err
	}

	return res, nil
}

const getReactionChannelDeliveries = `
SELECT channel, count(*) as total
FROM reactions
WHERE organization_id = ?
AND action_id = ?
AND created_at BETWEEN ? AND ?
AND event = 'delivered'
GROUP BY channel
`

type GetReactionChannelDeliveriesParams struct {
	OrganizationID uuid.UUID
	ActionID       uuid.UUID
	Timezone       string
	Start          string
	End            string
}

func (a Analytics) GetReactionChannelDeliveries(arg GetReactionChannelDeliveriesParams) ([]ReactionAggregation, error) {
	var res []ReactionAggregation
	if err := a.DB.Select(&res, getReactionChannelDeliveries,
		arg.OrganizationID,
		arg.ActionID,
		arg.Start,
		arg.End,
	); err != nil {
		return nil, err
	}

	return res, nil
}

const deleteReactionByIdentityIDs = `
ALTER TABLE reactions
DELETE
WHERE organization_id = ?
AND identity_id IN (?)
`

type DeleteReactionByIdentityIDsParams struct {
	OrganizationID uuid.UUID
	IdentityIDs    []uuid.UUID
}

func (a Analytics) DeleteReactionByIdentityIDs(arg DeleteReactionByIdentityIDsParams) error {
	query, args, err := sqlx.In(deleteReactionByIdentityIDs,
		arg.OrganizationID,
		arg.IdentityIDs,
	)
	if err != nil {
		return err
	}
	query = a.DB.Rebind(query)
	_, err = a.DB.Exec(query, args...)
	return err
}

const getReactionsForIdentities = `
SELECT *
FROM reactions
WHERE organization_id = ?
AND identity_id IN (?)
ORDER BY created_at DESC
`

type GetReactionsForIdentitiesParams struct {
	OrganizationID uuid.UUID
	IdentityIDs    []uuid.UUID
}

func (a Analytics) GetReactionsForIdentities(arg GetReactionsForIdentitiesParams) ([]Reaction, error) {
	var items []Reaction
	query, args, err := sqlx.In(getReactionsForIdentities,
		arg.OrganizationID,
		arg.IdentityIDs,
	)
	if err != nil {
		return nil, err
	}

	query = a.DB.Rebind(query)
	return items, a.DB.Select(&items, query, args...)
}

const getReactionsForInterval = `
SELECT bucket, count
FROM
(
	SELECT arrayJoin(timeSlots(subtractDays(now(), 1), toUInt32(86400), 10800)) AS bucket
) AS buckets
LEFT JOIN (
	SELECT toStartOfInterval(created_at, INTERVAL 3 hour) AS bucket,
	       count(*) AS count
	FROM reactions
	WHERE organization_id = ?
    AND created_at >= subtractDays(now(), 1)
    GROUP BY bucket
) reactions
ON buckets.bucket = reactions.bucket
`

func (a Analytics) GetReactionsForInterval(organizationID uuid.UUID) ([]ReactionBucket, error) {
	var res []ReactionBucket
	if err := a.DB.Select(&res, getReactionsForInterval, organizationID); err != nil {
		return nil, err
	}

	return res, nil
}
