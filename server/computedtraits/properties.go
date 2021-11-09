package computedtraits

import (
	"fmt"
	"strings"
)

type properties struct {
	Property        string `json:"property"`
	Type            string `json:"type"`
	AggregationType string `json:"aggregationType"`
	MinFrequency    int    `json:"minFrequency"`
	UseTimestamp    bool   `json:"useTimestamp"`
}

func (p properties) property() string {
	key := quote(p.Property)
	return fmt.Sprintf("'%s'", strings.ReplaceAll(key, ".", "', '"))
}

func (p properties) extract() string {
	extract := ""
	switch p.Type {
	case "string":
		extract = "JSONExtractString"
	case "integer":
		extract = "JSONExtractInt"
	case "float":
		extract = "JSONExtractFloat"
	case "boolean":
		extract = "JSONExtractBool"
	case "dateTime":
		extract = "JSONExtractString"
	}
	return extract
}
