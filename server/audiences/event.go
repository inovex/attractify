package audiences

import (
	"fmt"
	"strings"

	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type event struct {
	ID                 uuid.UUID  `json:"id"`
	InternalID         uuid.UUID  `json:"internalId"`
	ParentID           uuid.UUID  `json:"parentId"`
	Exclude            bool       `json:"exclude"`
	Channels           []string   `json:"channels"`
	Operator           string     `json:"operator"`
	Count              int        `json:"count"`
	TimeWindowOperator string     `json:"timeWindowOperator"`
	TimeWindowStart    int        `json:"timeWindowStart"`
	TimeWindowEnd      int        `json:"timeWindowEnd"`
	Properties         properties `json:"properties"`
}

func (e event) generate(organizationID uuid.UUID) string {
	q := `
LEFT JOIN 
(
  SELECT identity_id, created_at, properties, 1 AS found
  FROM events
  WHERE event_id = '%s'
	AND organization_id = '%s' 
	%s
	%s
	%s
) %s
ON %s.identity_id = i.id
`

	var properties string
	if len(e.Properties) > 0 {
		p := e.Properties.generate(true, "")
		if len(p) > 0 {
			properties = fmt.Sprintf("AND %s", p)
		}
	}

	return fmt.Sprintf(
		q,
		escape(e.ID.String()),
		organizationID.String(),
		e.channels(),
		e.timeWindow(),
		properties,
		sqlID(e.InternalID.String()),
		sqlID(e.InternalID.String()),
	)
}

func (e event) channels() string {
	if len(e.Channels) == 0 {
		return ""
	}
	var chs []string
	for _, c := range e.Channels {
		chs = append(chs, quote(escape(c)))
	}
	return fmt.Sprintf("AND channel IN (%s)", strings.Join(chs, ","))
}

func (e event) timeWindow() string {
	if e.TimeWindowOperator == db.TimeWindowTypeWithin {
		return fmt.Sprintf(
			"AND created_at BETWEEN subtractDays(now(), %d) AND now()", e.TimeWindowEnd,
		)
	} else if e.TimeWindowOperator == db.TimeWindowTypeInBetween {
		return fmt.Sprintf(
			"AND created_at BETWEEN subtractDays(now(), %d) AND subtractDays(now(), %d)",
			e.TimeWindowEnd, e.TimeWindowStart,
		)
	}
	return ""
}

func (e event) operator() string {
	switch e.Operator {
	case "less_than":
		return "<"
	case "less_or_exactly":
		return "<="
	case "more_or_exactly":
		return ">="
	case "more_than":
		return ">"
	default:
		return "="
	}
}

func (e event) having() string {
	return fmt.Sprintf(
		"COUNT(%s.identity_id) %s %d",
		sqlID(e.InternalID.String()), e.operator(), e.Count,
	)
}

func (e event) found() string {
	i := 1
	if e.Exclude {
		i = 0
	}
	return fmt.Sprintf("%s.found = %d", sqlID(e.InternalID.String()), i)
}
