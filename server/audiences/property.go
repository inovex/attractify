package audiences

import (
	"fmt"
	"strings"
)

type property struct {
	Source          string `json:"source"`
	Target          string `json:"target"`
	Key             string `json:"key"`
	Operator        string `json:"operator"`
	Value           string `json:"value"`
	DataType        string `json:"dataType"`
	EventID         string `json:"eventId"`
	InternalEventID string `json:"internalEventId"`
	EventProperty   string `json:"eventProperty"`
	TraitKey        string `json:"traitKey"`
}

func (p property) generateForEvents(eventID string) string {
	sourceKey := fmt.Sprintf("%s.properties", sqlID(eventID))
	switch p.Target {
	case "static":
		return p.condition("properties", p.keyPath(p.Key), p.value())
	case "custom_trait":
		return p.targetCondition(sourceKey, p.keyPath(p.Key), "i.custom_traits", p.keyPath(p.TraitKey))
	case "computed_trait":
		return p.targetCondition(sourceKey, p.keyPath(p.Key), "i.computed_traits", p.keyPath(p.TraitKey))
	case "funnel_property":
		k := fmt.Sprintf("%s.properties", sqlID(p.InternalEventID))
		return p.targetCondition(sourceKey, p.keyPath(p.Key), k, p.keyPath(p.EventProperty))
	}
	return ""
}

func (p property) generateForTraits(traitType string) string {
	switch p.Target {
	case "static":
		return p.condition(traitType, p.keyPath(p.Key), p.value())
	case "funnel_property":
		k := fmt.Sprintf("%s.properties", sqlID(p.InternalEventID))
		return p.targetCondition(traitType, p.keyPath(p.Key), k, p.keyPath(p.EventProperty))
	}
	return ""
}

func (p property) keyPath(key string) string {
	key = escape(key)
	return fmt.Sprintf("'%s'", strings.ReplaceAll(key, ".", "', '"))
}

func (p property) value() string {
	return escape(p.Value)
}

func (p property) condition(source, key, value string) string {
	switch p.DataType {
	case "string":
		return p.stringCondition(source, key, value)
	case "integer":
		return p.integerCondition(source, key, value)
	case "float":
		return p.floatCondition(source, key, value)
	case "boolean":
		return p.booleanCondition(source, key)
	case "dateTime":
		return p.dateTimeCondition(source, key, value)
	case "array":
		return p.arrayCondition(source, key, value)
	}
	return ""
}

func (p property) stringCondition(source, key, value string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) = '%s'", source, key, value)
	case "not_equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) != '%s'", source, key, value)
	case "contains":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%%%%%s%%%%'", source, key, value)
	case "does_not_contain":
		return fmt.Sprintf("JSONExtractString(%s, %s) NOT LIKE '%%%%%s%%%%'", source, key, value)
	case "starts_with":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%s%%'", source, key, value)
	case "ends_with":
		return fmt.Sprintf("JSONExtractString(%s, %s) LIKE '%%%s'", source, key, value)
	default:
		return ""
	}
}

func (p property) integerCondition(source, key, value string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "equals":
		return fmt.Sprintf("JSONExtractInt(%s, %s) = %s", source, key, value)
	case "not_equals":
		return fmt.Sprintf("JSONExtractInt(%s, %s) != %s", source, key, value)
	case "less_than":
		return fmt.Sprintf("JSONExtractInt(%s, %s) < %s", source, key, value)
	case "greater_than":
		return fmt.Sprintf("JSONExtractInt(%s, %s) > %s", source, key, value)
	case "less_than_or_equal":
		return fmt.Sprintf("JSONExtractInt(%s, %s) <= %s", source, key, value)
	case "greater_than_or_equal":
		return fmt.Sprintf("JSONExtractInt(%s, %s) >= %s", source, key, value)
	default:
		return ""
	}
}

func (p property) floatCondition(source, key, value string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "equals":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) = %s", source, key, value)
	case "not_equals":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) != %s", source, key, value)
	case "less_than":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) < %s", source, key, value)
	case "greater_than":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) > %s", source, key, value)
	case "less_than_or_equal":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) <= %s", source, key, value)
	case "greater_than_or_equal":
		return fmt.Sprintf("JSONExtractFloat(%s, %s) >= %s", source, key, value)
	default:
		return ""
	}
}

func (p property) booleanCondition(source, key string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "true":
		return fmt.Sprintf("JSONExtractBool(%s, %s) = 1", source, key)
	case "false":
		return fmt.Sprintf("JSONExtractBool(%s, %s) = 0", source, key)
	default:
		return ""
	}
}

func (p property) dateTimeCondition(source, key, value string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "before_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < toDateTime(%s)", source, key, value)
	case "after_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > toDateTime(%s)", source, key, value)
	case "within_last_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > subtractDays(now(), %s)", source, key, value)
	case "within_next_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < addDays(now(), %s)", source, key, value)
	case "before_last_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < subtractDays(now(), %s)", source, key, value)
	case "after_next_days":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > addDays(now(), %s)", source, key, value)
	default:
		return ""
	}
}

func (p property) arrayCondition(source, key, value string) string {
	switch p.Operator {
	case "exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 1", source, key)
	case "not_exists":
		return fmt.Sprintf("JSONHas(%s, %s) = 0", source, key)
	case "contains":
		return fmt.Sprintf("has(JSONExtractArrayRaw(%s, %s), '\"%s\"') = 1", source, key, value)
	case "does_not_contain":
		return fmt.Sprintf("has(JSONExtractArrayRaw(%s, %s), '\"%s\"') = 0", source, key, value)
	default:
		return ""
	}
}

func (p property) targetCondition(source, sourceKey, target, targetKey string) string {
	switch p.DataType {
	case "string":
		return p.stringComparison(source, sourceKey, target, targetKey)
	case "integer":
		return p.integerCondition(source, sourceKey, fmt.Sprintf("JSONExtractInt(%s, %s)", target, targetKey))
	case "float":
		return p.floatCondition(source, sourceKey, fmt.Sprintf("JSONExtractFloat(%s, %s)", target, targetKey))
	case "boolean":
		return p.booleanComparison(source, sourceKey, target, targetKey)
	case "dateTime":
		return p.dateTimeComparisonCondition(source, sourceKey, target, targetKey)
	case "array":
		return p.arrayComparison(source, sourceKey, target, targetKey)
	}
	return ""
}

func (p property) stringComparison(source, sourceKey, target, targetKey string) string {
	switch p.Operator {
	case "equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) = JSONExtractString(%s, %s)", source, sourceKey, target, targetKey)
	case "not_equals":
		return fmt.Sprintf("JSONExtractString(%s, %s) != JSONExtractString(%s, %s)", source, sourceKey, target, targetKey)
	case "contains":
		return fmt.Sprintf("positionCaseInsensitiveUTF8(JSONExtractString(%s, %s), JSONExtractString(%s, %s)) > 0", source, sourceKey, target, targetKey)
	case "does_not_contain":
		return fmt.Sprintf("positionCaseInsensitiveUTF8(JSONExtractString(%s, %s), JSONExtractString(%s, %s)) = 0", source, sourceKey, target, targetKey)
	case "starts_with":
		return fmt.Sprintf("startsWith(JSONExtractString(%s, %s), JSONExtractString(%s, %s)) = 1", source, sourceKey, target, targetKey)
	case "ends_with":
		return fmt.Sprintf("endsWith(JSONExtractString(%s, %s), JSONExtractString(%s, %s)) = 1", source, sourceKey, target, targetKey)
	default:
		return ""
	}
}

func (p property) booleanComparison(source, sourceKey, target, targetKey string) string {
	switch p.Operator {
	case "equals":
		return fmt.Sprintf("JSONExtractBool(%s, %s) = JSONExtractBool(%s, %s)", source, sourceKey, target, targetKey)
	case "not_equals":
		return fmt.Sprintf("JSONExtractBool(%s, %s) @= JSONExtractBool(%s, %s)", source, sourceKey, target, targetKey)
	default:
		return ""
	}
}

func (p property) dateTimeComparisonCondition(source, sourceKey, target, targetKey string) string {
	switch p.Operator {
	case "before_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) < parseDateTimeBestEffort(JSONExtract(%s, %s, 'String'))", source, sourceKey, target, targetKey)
	case "after_date":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) > parseDateTimeBestEffort(JSONExtract(%s, %s, 'String'))", source, sourceKey, target, targetKey)
	case "equals":
		return fmt.Sprintf("parseDateTimeBestEffort(JSONExtract(%s, %s, 'String')) = parseDateTimeBestEffort(JSONExtract(%s, %s, 'String'))", source, sourceKey, target, targetKey)
	default:
		return ""
	}
}

func (p property) arrayComparison(source, sourceKey, target, targetKey string) string {
	switch p.Operator {
	case "contains":
		return fmt.Sprintf("hasAll(JSONExtractArrayRaw(%s, %s), JSONExtractArrayRaw(%s, %s)) = 1", source, sourceKey, target, targetKey)
	case "does_not_contain":
		return fmt.Sprintf("hasAll(JSONExtractArrayRaw(%s, %s), JSONExtractArrayRaw(%s, %s)) = 0", source, sourceKey, target, targetKey)
	default:
		return ""
	}
}
