package audiences

import (
	"strings"
)

type properties []property

func (p properties) generate(onlyStatic bool, eventID string) string {
	var props []string
	for _, prop := range p {
		var c string
		if onlyStatic && prop.Target == "static" {
			c = prop.generateForEvents(eventID)
		} else if !onlyStatic && prop.Target != "static" {
			c = prop.generateForEvents(eventID)
		} else {
			continue
		}
		props = append(props, c)
	}
	return strings.Join(props, "\nAND ")
}
