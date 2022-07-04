package analytics

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Properties map[string]interface{}

type Event struct {
	ID             uuid.UUID `db:"id" json:"-"`
	OrganizationID uuid.UUID `db:"organization_id" json:"-"`
	IdentityID     uuid.UUID `db:"identity_id" json:"-"`
	EventID        uuid.UUID `db:"event_id" json:"eventId"`
	Channel        string    `db:"channel" json:"channel"`
	Context        string    `db:"context" json:"context"`
	Properties     string    `db:"properties" json:"properties"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

const createEvent = `
INSERT INTO events (
	id,
	organization_id,
	identity_id,
	event_id,
	channel,
	context,
	properties,
	created_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?
)`

type CreateEventParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	IdentityID     uuid.UUID
	EventID        uuid.UUID
	Channel        string
	Context        string
	Properties     string
	CreatedAt      time.Time
}

func (a Analytics) CreateEvent(arg CreateEventParams) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	tx.Exec(`SET insert_distributed_sync = 1;`)
	_, err = tx.Exec(createEvent,
		arg.ID,
		arg.OrganizationID,
		arg.IdentityID,
		arg.EventID,
		arg.Channel,
		arg.Context,
		arg.Properties,
		arg.CreatedAt.UTC(),
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (a Analytics) CreateEvents(arg []CreateEventParams) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createEvent)
	if err != nil {
		return err
	}
	for _, e := range arg {
		_, err = stmt.Exec(
			e.ID,
			e.OrganizationID,
			e.IdentityID,
			e.EventID,
			e.Channel,
			e.Context,
			e.Properties,
			e.CreatedAt.UTC(),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

const getEvents = `
SELECT *
FROM events
WHERE organization_id = ?
AND created_at BETWEEN ? AND ?
ORDER BY created_at DESC
LIMIT ?, ?
`

const getEventsForIDs = `
SELECT *
FROM events
WHERE organization_id = ?
AND event_id IN (?)
AND created_at BETWEEN ? AND ?
ORDER BY created_at DESC
LIMIT ?, ?
`

const getEventsForIdentity = `
SELECT *
FROM events
WHERE organization_id = ?
AND identity_id IN (?)
AND created_at BETWEEN ? AND ?
ORDER BY created_at DESC
LIMIT ?, ?
`

const getEventsForIDsAndIdentity = `
SELECT *
FROM events
WHERE organization_id = ?
AND event_id IN (?)
AND identity_id IN (?)
AND created_at BETWEEN ? AND ?
ORDER BY created_at DESC
LIMIT ?, ?
`

type GetEventsParams struct {
	OrganizationID uuid.UUID
	EventIDs       []uuid.UUID
	IdentityIDs    []uuid.UUID
	Start          time.Time
	End            time.Time
	Limit          int
	Offset         int
}

func (a Analytics) GetEvents(arg GetEventsParams) ([]Event, error) {
	var items []Event

	if len(arg.EventIDs) == 0 && len(arg.IdentityIDs) == 0 {
		return items, a.DB.Select(
			&items,
			getEvents,
			arg.OrganizationID,
			arg.Start.Format("2006-01-02 15:04:05"),
			arg.End.Format("2006-01-02 15:04:05"),
			arg.Offset,
			arg.Limit,
		)
	} else if len(arg.EventIDs) > 0 && len(arg.IdentityIDs) == 0 {
		query, args, err := sqlx.In(getEventsForIDs,
			arg.OrganizationID,
			arg.EventIDs,
			arg.Start.Format("2006-01-02 15:04:05"),
			arg.End.Format("2006-01-02 15:04:05"),
			arg.Offset,
			arg.Limit,
		)
		if err != nil {
			return nil, err
		}

		query = a.DB.Rebind(query)
		return items, a.DB.Select(&items, query, args...)
	} else if len(arg.EventIDs) == 0 && len(arg.IdentityIDs) > 0 {
		query, args, err := sqlx.In(getEventsForIdentity,
			arg.OrganizationID,
			arg.IdentityIDs,
			arg.Start.Format("2006-01-02 15:04:05"),
			arg.End.Format("2006-01-02 15:04:05"),
			arg.Offset,
			arg.Limit,
		)
		if err != nil {
			return nil, err
		}

		query = a.DB.Rebind(query)
		return items, a.DB.Select(&items, query, args...)
	}

	query, args, err := sqlx.In(getEventsForIDsAndIdentity,
		arg.OrganizationID,
		arg.EventIDs,
		arg.IdentityIDs,
		arg.Start.Format("2006-01-02 15:04:05"),
		arg.End.Format("2006-01-02 15:04:05"),
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}

	query = a.DB.Rebind(query)
	return items, a.DB.Select(&items, query, args...)
}

const getEventsForIdentities = `
SELECT *
FROM events
WHERE organization_id = ?
AND identity_id IN (?)
ORDER BY created_at DESC
`

type GetEventsForIdentitiesParams struct {
	OrganizationID uuid.UUID
	IdentityIDs    []uuid.UUID
}

func (a Analytics) GetEventsForIdentities(arg GetEventsForIdentitiesParams) ([]Event, error) {
	var items []Event
	query, args, err := sqlx.In(getEventsForIdentities,
		arg.OrganizationID,
		arg.IdentityIDs,
	)
	if err != nil {
		return nil, err
	}

	query = a.DB.Rebind(query)
	return items, a.DB.Select(&items, query, args...)
}

const getEventCount = `
SELECT count(*)
FROM events
WHERE organization_id = ?
AND created_at BETWEEN ? AND ?
`

func (a Analytics) GetEventCount(arg GetEventsParams) (int, error) {
	row := a.DB.QueryRowx(getEventCount, arg.OrganizationID, arg.Start, arg.End)
	var count int
	return count, row.Scan(&count)
}

const getLastDaysEventCount = `
SELECT count(*)
FROM events
WHERE organization_id = ?
AND created_at >= subtractDays(now(), 1)
`

func (a Analytics) GetLastDaysEventCount(organizationID uuid.UUID) (int, error) {
	row := a.DB.QueryRowx(getLastDaysEventCount, organizationID)
	var count int
	return count, row.Scan(&count)
}

const deleteEvent = `
ALTER TABLE events
DELETE
WHERE organization_id = ?
AND id = ?
`

type DeleteEventParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (a Analytics) DeleteEvent(arg DeleteEventParams) error {
	_, err := a.DB.Exec(deleteEvent, arg.OrganizationID, arg.ID)
	return err
}

const deleteEventByIdentityID = `
ALTER TABLE events
DELETE
WHERE organization_id = ?
AND identity_id = ?
`

type DeleteEventByIdentityIDParams struct {
	OrganizationID uuid.UUID
	IdentityID     uuid.UUID
}

func (a Analytics) DeleteEventByIdentityID(arg DeleteEventByIdentityIDParams) error {
	_, err := a.DB.Exec(deleteEventByIdentityID, arg.OrganizationID, arg.IdentityID)
	return err
}

const deleteEventsByIdentityIDs = `
ALTER TABLE events
DELETE
WHERE organization_id = ?
AND identity_id IN (?)
`

type DeleteEventsByIdentityIDsParams struct {
	OrganizationID uuid.UUID
	IdentityIDs    []uuid.UUID
}

func (a Analytics) DeleteEventsByIdentityIDs(arg DeleteEventsByIdentityIDsParams) error {
	query, args, err := sqlx.In(deleteEventsByIdentityIDs,
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
