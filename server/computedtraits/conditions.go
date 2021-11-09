package computedtraits

import (
	"fmt"
	"strings"
)

type condition struct {
	Property string `json:"property"`
	Operator string `json:"operator"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

func (c *ComputedTrait) generateConditions() string {
	var conds []string
	for _, c := range c.conditions {
		conds = append(conds, c.condition("properties"))
	}
	if len(conds) == 0 {
		return ""
	}
	return fmt.Sprintf("AND %s", strings.Join(conds, " AND "))
}

func (c condition) property() string {
	key := quote(c.Property)
	return fmt.Sprintf("'%s'", strings.ReplaceAll(key, ".", "', '"))
}

func (c condition) value() string {
	return quote(c.Value)
}

func quote(v string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(v)
}

func (c condition) condition(key string) string {
	switch c.Type {
	case "string":
		return c.stringOperator(key)
	case "integer":
		return c.integerOperator(key)
	case "float":
		return c.floatOperator(key)
	case "boolean":
		return c.booleanOperator(key)
	case "dateTime":
		return c.dateTimeOperator(key)
	case "array":
		return c.arrayOperator(key)
	}
	return ""
}

func (c condition) stringOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) = '%s'", key, c.property(), c.value())
	case "not_equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) != '%s'", key, c.property(), c.value())
	case "contains":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%%%%%s%%%%'", key, c.property(), c.value())
	case "does_not_contain":
		return fmt.Sprintf("JSONExtractString(%s, %s) NOT LIKE '%%%%%s%%%%'", key, c.property(), c.value())
	case "starts_with":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%s%%'", key, c.property(), c.value())
	case "ends_with":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%%%s'", key, c.property(), c.value())
	default:
		return ""
	}
}

func (c condition) integerOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "equals":
		return fmt.Sprintf("JSONExtractInt(%s, %s) = %s", key, c.property(), c.value())
	case "not_equals":
		return fmt.Sprintf("JSONExtractInt(%s, %s) != %s", key, c.property(), c.value())
	case "less_than":
		return fmt.Sprintf("JSONExtractInt(%s, %s) < %s", key, c.property(), c.value())
	case "greater_than":
		return fmt.Sprintf("JSONExtractInt(%s, %s) > %s", key, c.property(), c.value())
	case "less_than_or_equal":
		return fmt.Sprintf("JSONExtractInt(%s, %s) <= %s", key, c.property(), c.value())
	case "greater_than_or_equal":
		return fmt.Sprintf("JSONExtractInt(%s, %s) >= %s", key, c.property(), c.value())
	default:
		return ""
	}
}

func (c condition) floatOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "equals":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) = %s", key, c.property(), c.value())
	case "not_equals":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) != %s", key, c.property(), c.value())
	case "less_than":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) < %s", key, c.property(), c.value())
	case "greater_than":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) > %s", key, c.property(), c.value())
	case "less_than_or_equal":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) <= %s", key, c.property(), c.value())
	case "greater_than_or_equal":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) >= %s", key, c.property(), c.value())
	default:
		return ""
	}
}

func (c condition) booleanOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "true":
		return fmt.Sprintf("JSONExtractBool(%s, %s) = 1", key, c.property())
	case "false":
		return fmt.Sprintf("JSONExtractBool(%s, %s) = 0", key, c.property())
	default:
		return ""
	}
}

func (c condition) dateTimeOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "before_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < toDateTime(%s)", key, c.property(), c.value())
	case "after_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > toDateTime(%s)", key, c.property(), c.value())
	case "within_last_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > subtractDays(now(), %s)", key, c.property(), c.value())
	case "within_next_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < addDays(now(), %s)", key, c.property(), c.value())
	case "before_last_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < subtractDays(now(), %s)", key, c.property(), c.value())
	case "after_next_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > addDays(now(), %s)", key, c.property(), c.value())
	default:
		return ""
	}
}

func (c condition) arrayOperator(key string) string {
	switch c.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", key, c.property())
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", key, c.property())
	case "contains":
		return fmt.Sprintf("has(JSONExtractArrayRaw(%s, %s), '\"%s\"') = 1", key, c.property(), c.value())
	case "does_not_contain":
		return fmt.Sprintf("has(JSONExtractArrayRaw(%s, %s), '\"%s\"') = 0", key, c.property(), c.value())
	default:
		return ""
	}
}
