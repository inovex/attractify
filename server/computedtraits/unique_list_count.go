package computedtraits

import "fmt"

const uniqueListCountQuery = `
SELECT arrayUniq(groupArray(%s(properties, %s)))
FROM events
WHERE organization_id = ?
AND event_id = '%s'
%s
%s
`

func (c *ComputedTrait) uniqueListCount() (int, error) {
	query := fmt.Sprintf(
		uniqueListCountQuery,
		c.properties.extract(),
		c.properties.property(),
		quote(c.ct.EventID.String()),
		c.generateConditions(),
		c.profileIdentities(),
	)

	var count int
	row := c.app.Analytics.DB.QueryRowx(query, c.ct.OrganizationID)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
