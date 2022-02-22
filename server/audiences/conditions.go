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
%s
%s
%s
%s
GROUP BY i.profile_id
%s
`
	var events []string
	var having []string
	var havingCond string
	var found []string
	var foundCond string
	for _, e := range c.events {
		events = append(events, e.generate(organizationID))
		if !e.Exclude {
			having = append(having, e.having())
		}
		found = append(found, e.found())
	}

	filterAnonymous := "AND toUInt8(i.is_anonymous) = 0"
	if c.includeAnonymous {
		filterAnonymous = ""
	}

	if len(having) > 0 {
		havingCond = "HAVING " + strings.Join(having, "\nAND ")
	}

	if len(found) > 0 {
		foundCond = strings.Join(found, "\nAND ")
	}

	return fmt.Sprintf(
		query,
		organizationID.String(),
		strings.Join(events, ""),
		foundCond,
		c.funnelOrder(),
		c.traitConditions(found),
		c.eventConditions(),
		filterAnonymous,
		havingCond,
	)
}

func (c conditions) traitConditions(found []string) string {
	traits := c.traits.generate()
	if len(found) > 0 && len(traits) > 0 {
		return fmt.Sprintf("AND %s", traits)
	}
	return c.traits.generate()
}

func (c conditions) funnelOrder() string {
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
	if len(times) > 0 {
		return " AND " + strings.Join(times, " \nAND ")
	}
	return ""
}

func (c conditions) eventConditions() string {
	var conditions []string
	for _, e := range c.events {
		conds := e.Properties.generate(false, e.InternalID.String())
		if len(conds) > 0 {
			conditions = append(conditions, conds)
		}
	}
	if len(conditions) > 0 {
		return fmt.Sprintf("AND %s", strings.Join(conditions, " \nAND "))
	}
	return ""
}
