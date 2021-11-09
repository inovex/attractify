package actions

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type targetingCondition struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

func (tc targetingCondition) eval(properties json.RawMessage) bool {
	v := gjson.GetBytes(properties, tc.Key)
	if tc.Type != "not_exists" && !v.Exists() {
		return false
	}

	switch tc.Type {
	case "string":
		return tc.stringCondition(v)
	case "integer":
		return tc.numberCondition(v)
	case "float":
		return tc.numberCondition(v)
	case "boolean":
		return tc.booleanCondition(v)
	case "array":
		return tc.arrayCondition(v)
	case "dateTime":
		return tc.dateTimeCondition(v)
	}
	return true
}

func (tc targetingCondition) stringCondition(v gjson.Result) bool {
	switch tc.Operator {
	case "equals":
		return v.String() == tc.Value.(string)
	case "not_equals":
		return v.String() != tc.Value.(string)
	case "contains":
		return strings.Contains(v.String(), tc.Value.(string))
	case "does_not_contain":
		return !strings.Contains(v.String(), tc.Value.(string))
	case "starts_with":
		return strings.HasPrefix(v.String(), tc.Value.(string))
	case "ends_with":
		return strings.HasSuffix(v.String(), tc.Value.(string))
	case "exists":
		return v.Value() != nil
	case "not_exists":
		return v.Value() == nil
	}
	return false
}

func (tc targetingCondition) numberCondition(v gjson.Result) bool {
	switch tc.Operator {
	case "equals":
		return v.Float() == tc.Value.(float64)
	case "not_equals":
		return v.Float() != tc.Value.(float64)
	case "less_than":
		return v.Float() < tc.Value.(float64)
	case "greater_than":
		return v.Float() > tc.Value.(float64)
	case "less_than_or_equal":
		return v.Float() <= tc.Value.(float64)
	case "greater_than_or_equal":
		return v.Float() >= tc.Value.(float64)
	case "exists":
		return v.Value() != nil
	case "not_exists":
		return v.Value() == nil
	}
	return false
}

func (tc targetingCondition) booleanCondition(v gjson.Result) bool {
	switch tc.Operator {
	case "true":
		return v.Bool()
	case "false":
		return !v.Bool()
	case "exists":
		return v.Value() != nil
	case "not_exists":
		return v.Value() == nil
	}
	return false
}
func (tc targetingCondition) arrayCondition(v gjson.Result) bool {
	switch tc.Operator {
	case "contains":
		for _, x := range v.Array() {
			if x.String() == tc.Value.(string) {
				return true
			}
		}
		return false
	case "does_not_contain":
		for _, x := range v.Array() {
			if x.String() == tc.Value.(string) {
				return false
			}
		}
		return true
	case "exists":
		return v.Value() != nil
	case "not_exists":
		return v.Value() == nil
	}
	return false
}

func (tc targetingCondition) dateTimeCondition(v gjson.Result) bool {
	cdt, err := time.Parse(time.RFC3339, v.String())
	if err != nil {
		return false
	}

	switch tc.Operator {
	case "before_date":
		vdt, err := time.Parse(time.RFC3339, tc.Value.(string))
		if err != nil {
			return false
		}
		return cdt.Before(vdt)
	case "after_date":
		vdt, err := time.Parse(time.RFC3339, tc.Value.(string))
		if err != nil {
			return false
		}
		return cdt.After(vdt)
	case "within_last_days":
		return cdt.After(time.Now().AddDate(0, 0, -int(tc.Value.(float64))))
	case "within_next_days":
		return cdt.Before(time.Now().AddDate(0, 0, int(tc.Value.(float64))))
	case "before_last_days":
		return cdt.Before(time.Now().AddDate(0, 0, -int(tc.Value.(float64))))
	case "after_next_days":
		return cdt.After(time.Now().AddDate(0, 0, int(tc.Value.(float64))))
	case "exists":
		return v.Value() != nil
	case "not_exists":
		return v.Value() == nil
	}
	return false
}
