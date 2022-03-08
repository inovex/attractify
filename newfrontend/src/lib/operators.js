export default {
  stringOperators: [
    { text: 'equals', value: 'equals' },
    { text: 'not equals', value: 'not_equals' },
    { text: 'contains', value: 'contains' },
    { text: 'does not contain ', value: 'does_not_contain' },
    { text: 'starts with ', value: 'starts_with' },
    { text: 'ends with', value: 'ends_with' },
    { text: 'exists', value: 'exists' },
    { text: 'does not exist', value: 'not_exists' }
  ],
  numberOperators: [
    { text: 'equals', value: 'equals' },
    { text: 'not equals', value: 'not_equals' },
    { text: 'is less than', value: 'less_than' },
    { text: 'is greater than', value: 'greater_than' },
    { text: 'is less than or equal', value: 'less_than_or_equal' },
    { text: 'is greater than or equal', value: 'greater_than_or_equal' },
    { text: 'exists', value: 'exists' },
    { text: 'does not exist', value: 'not_exists' }
  ],
  booleanOperators: [
    { text: 'is true ', value: 'true' },
    { text: 'is false', value: 'false' },
    { text: 'exists', value: 'exists' },
    { text: 'does not exist', value: 'not_exists' }
  ],
  dateTimeOperators: [
    { text: 'exists', value: 'exists' },
    { text: 'does not exist', value: 'not_exists' },
    { text: 'before date ', value: 'before_date' },
    { text: 'after date', value: 'after_date' },
    { text: 'within last days ', value: 'within_last_days' },
    { text: 'within next days ', value: 'within_next_days' },
    { text: 'before last days ', value: 'before_last_days' },
    { text: 'after next days', value: 'after_next_days' }
  ],
  arrayOperators: [
    { text: 'contains', value: 'contains' },
    { text: 'does not contain ', value: 'does_not_contain' },
    { text: 'exists', value: 'exists' },
    { text: 'does not exist', value: 'not_exists' }
  ],
  stringComparisonOperators: [
    { text: 'equals', value: 'equals' },
    { text: 'not equals', value: 'not_equals' },
    { text: 'contains', value: 'contains' },
    { text: 'does not contain ', value: 'does_not_contain' },
    { text: 'starts with ', value: 'starts_with' },
    { text: 'ends with', value: 'ends_with' }
  ],
  numberComparisonOperators: [
    { text: 'equals', value: 'equals' },
    { text: 'not equals', value: 'not_equals' },
    { text: 'is less than', value: 'less_than' },
    { text: 'is greater than', value: 'greater_than' },
    { text: 'is less than or equal', value: 'less_than_or_equal' },
    { text: 'is greater than or equal', value: 'greater_than_or_equal' }
  ],
  booleanComparisonOperators: [
    { text: 'equals', value: 'equals' },
    { text: 'not equals', value: 'not_equals' }
  ],
  dateTimeComparisonOperators: [
    { text: 'before date ', value: 'before_date' },
    { text: 'after date', value: 'after_date' },
    { text: 'equals', value: 'equals' }
  ],
  arrayComparisonOperators: [
    { text: 'contains', value: 'contains' },
    { text: 'does not contain ', value: 'does_not_contain' }
  ]
}
