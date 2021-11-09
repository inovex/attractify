package computedtraits

import "fmt"

const mostFrequentQuery = `
SELECT %s(properties, %s) as v, count(*) as c
FROM events
WHERE organization_id = ?
AND event_id = '%s'
%s
%s
GROUP BY v
HAVING c >= %d
ORDER BY c DESC
LIMIT 1
`

func (c *ComputedTrait) mostFrequent() (interface{}, error) {
	query := fmt.Sprintf(
		mostFrequentQuery,
		c.properties.extract(),
		c.properties.property(),
		quote(c.ct.EventID.String()),
		c.generateConditions(),
		c.profileIdentities(),
		c.properties.MinFrequency,
	)

	var (
		res   interface{}
		count int
	)
	row := c.app.Analytics.DB.QueryRowx(query, c.ct.OrganizationID)
	if err := row.Scan(&res, &count); err != nil {
		return nil, err
	}

	return res, nil
}
