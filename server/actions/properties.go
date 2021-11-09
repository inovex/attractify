package actions

import (
	"encoding/json"

	"attractify.io/platform/db"
	"github.com/tidwall/gjson"
)

type Variables map[string]interface{}

func (a Action) MapProperties(channel string) map[string]interface{} {
	propMap := map[string]interface{}{}
	for _, s := range a.properties {
		if !a.ForChannel(s.Channels, channel) {
			continue
		}

		if s.Type == "text" {
			propMap[s.Name] = s.Value
		} else if s.Type == "custom_trait" {
			propMap[s.Name] = a.extractVal(a.profile.CustomTraits, s)
		} else if s.Type == "computed_trait" {
			propMap[s.Name] = a.extractVal(a.profile.ComputedTraits, s)
		}
	}

	return propMap
}

func (a Action) extractVal(source json.RawMessage, p db.ActionProperty) interface{} {
	v := gjson.GetBytes(source, p.SourceKey)
	switch p.SourceType {
	case "integer":
		return v.Int()
	case "float":
		return v.Float()
	case "boolean":
		return v.Bool()
	default:
		return v.String()
	}
}

func (a Action) ForChannel(channels []string, channel string) bool {
	for _, ch := range channels {
		if channel == ch {
			return true
		}
	}
	return false
}
