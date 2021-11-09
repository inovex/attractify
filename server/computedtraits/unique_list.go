package computedtraits

import "fmt"

const uniqueListQuery = `
SELECT arrayDistinct(groupArray(%s(properties, %s)))
FROM events
WHERE organization_id = ?
AND event_id = '%s'
%s
%s
`

func (c *ComputedTrait) uniqueList() (interface{}, error) {
	query := fmt.Sprintf(
		uniqueListQuery,
		c.properties.extract(),
		c.properties.property(),
		quote(c.ct.EventID.String()),
		c.generateConditions(),
		c.profileIdentities(),
	)

	var res interface{}
	row := c.app.Analytics.DB.QueryRowx(query, c.ct.OrganizationID)
	if err := row.Scan(&res); err != nil {
		return nil, err
	}

	return res, nil
}
