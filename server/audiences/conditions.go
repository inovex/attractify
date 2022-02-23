package audiences

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

type conditions struct {
	events           events
	traits           traits
	includeAnonymous bool
}

func (c conditions) BuildQuery(organizationID uuid.UUID) string {
	const query = `
SELECT i.profile_id AS id
FROM (
	SELECT * FROM full_identities
	WHERE organization_id = '%s'
) i
%s
WHERE %s

GROUP BY i.profile_id
%s
`
	var whereCond []string
	var events []string
	var having []string
	var havingCond string
	for _, e := range c.events {
		events = append(events, e.generate(organizationID))
		if !e.Exclude {
			having = append(having, e.having())
		}
		whereCond = append(whereCond, e.found())
	}

	if !c.includeAnonymous {
		whereCond = append(whereCond, "AND toUInt8(i.is_anonymous) = 0")
	}

	whereCond = append(whereCond, c.funnelOrder()...)
	whereCond = append(whereCond, c.traitConditions()...)
	whereCond = append(whereCond, c.eventConditions()...)

	if len(having) > 0 {
		havingCond = "HAVING " + strings.Join(having, "\nAND ")
	}

	return fmt.Sprintf(
		query,
		organizationID.String(),
		strings.Join(events, ""),
		strings.Join(whereCond, "\nAND "),
		havingCond,
	)
}

func (c conditions) traitConditions() []string {
	return c.traits.generate()
}

func (c conditions) funnelOrder() []string {
	var times []string
	var lastEventID string
	for _, e := range c.events {
		if e.ParentID != uuid.Nil {
			times = append(
				times,
				fmt.Sprintf(
					"IF(%s.found = 1 AND %s.found = 1, %s.created_at < %s.created_at, 1)",
					sqlID(lastEventID),
					sqlID(e.InternalID.String()),
					sqlID(lastEventID),
					sqlID(e.InternalID.String()),
				),
			)
		}
		lastEventID = e.InternalID.String()
	}
	return times
}

func (c conditions) eventConditions() []string {
	var conditions []string
	for _, e := range c.events {
		conds := e.Properties.generate(false, e.InternalID.String())
		if len(conds) > 0 {
			conditions = append(conditions, conds)
		}
	}
	return conditions
}
