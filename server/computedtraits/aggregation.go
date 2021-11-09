package computedtraits

import (
	"errors"
	"fmt"
	"math"
)

const aggregationQuery = `
SELECT %s(JSONExtractFloat(properties, %s)) as aggregation
FROM events
WHERE organization_id = ?
AND event_id = '%s'
%s
%s
`

func (c *ComputedTrait) aggregation() (float64, error) {
	query := fmt.Sprintf(
		aggregationQuery,
		quote(c.properties.AggregationType),
		c.properties.property(),
		quote(c.ct.EventID.String()),
		c.generateConditions(),
		c.profileIdentities(),
	)

	var aggregate float64
	row := c.app.Analytics.DB.QueryRowx(query, c.ct.OrganizationID)
	if err := row.Scan(&aggregate); err != nil {
		return 0, err
	}

	if math.IsNaN(aggregate) {
		return 0, errors.New("could not calculate result")
	}

	return math.Round(aggregate*1000) / 1000, nil
}
