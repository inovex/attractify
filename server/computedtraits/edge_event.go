package computedtraits

import (
	"fmt"
	"time"
)

const edgeEventQuery = `
SELECT %s(properties, %s), created_at
FROM events
WHERE organization_id = ?
AND event_id = '%s'
%s
%s
ORDER BY created_at %s
LIMIT 1
`

func (c *ComputedTrait) edgeEvent(first bool) (interface{}, error) {
	order := "DESC"
	if first {
		order = "ASC"
	}
	query := fmt.Sprintf(
		edgeEventQuery,
		c.properties.extract(),
		c.properties.property(),
		quote(c.ct.EventID.String()),
		c.generateConditions(),
		c.profileIdentities(),
		order,
	)

	var (
		res       interface{}
		createdAt time.Time
	)
	row := c.app.Analytics.DB.QueryRowx(query, c.ct.OrganizationID)
	if err := row.Scan(&res, &createdAt); err != nil {
		return nil, err
	}
	if c.properties.UseTimestamp {
		return createdAt.UTC().Format(time.RFC3339), nil
	}

	return res, nil
}
