<template>
  <v-card outlined>
    <v-card-text>
      <h4>Event condition</h4>
      <v-row>
        <v-col>
          <v-select
            dense
            label="Property"
            :items="properties"
            v-model="condition.property"
            return-object
            @change="
              e => {
                setProperty(condition, e)
              }
            "
          ></v-select>
        </v-col>
        <v-col class="col-lg-2">
          <v-select
            dense
            :items="getOperators(condition.type)"
            label="Operator"
            :value="condition.operator"
            v-model="condition.operator"
            :rules="[rules.required]"
          ></v-select>
        </v-col>
        <v-col>
          <v-text-field
            dense
            v-if="showValue(condition.operator)"
            label="Value"
            name="value"
            type="text"
            v-model="condition.value"
            :rules="[rules.required]"
          />
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: ['properties', 'condition'],
  data() {
    return {
      stringOperators: [
        { text: 'Equals', value: 'equals' },
        { text: 'Not Equals', value: 'not_equals' },
        { text: 'Contains', value: 'contains' },
        { text: 'Does Not Contain ', value: 'does_not_contain' },
        { text: 'Starts With ', value: 'starts_with' },
        { text: 'Ends With', value: 'ends_with' },
        { text: 'Exists', value: 'exists' },
        { text: 'Not Exists', value: 'not_exists' }
      ],
      numberOperators: [
        { text: 'Equals', value: 'equals' },
        { text: 'Not Equals', value: 'not_equals' },
        { text: 'Less Than', value: 'less_than' },
        { text: 'Greater Than', value: 'greater_than' },
        { text: 'Less Than Or Equal', value: 'less_than_or_equal' },
        { text: 'Greater Than Or Equal', value: 'greater_than_or_equal' },
        { text: 'Exists', value: 'exists' },
        { text: 'Not Exists', value: 'not_exists' }
      ],
      booleanOperators: [
        { text: 'True ', value: 'true' },
        { text: 'False', value: 'false' },
        { text: 'Exists', value: 'exists' },
        { text: 'Not Exists', value: 'not_exists' }
      ],
      dateTimeOperators: [
        { text: 'Exists', value: 'exists' },
        { text: 'Not Exists', value: 'not_exists' },
        { text: 'Before Date ', value: 'before_date' },
        { text: 'After Date', value: 'after_date' },
        { text: 'Within Last Days ', value: 'within_last_days' },
        { text: 'Within Next Days ', value: 'within_next_days' },
        { text: 'Before Last Days ', value: 'before_last_days' },
        { text: 'After Next Days', value: 'after_next_days' }
      ],
      arrayOperators: [
        { text: 'Contains', value: 'contains' },
        { text: 'Does Not Contain ', value: 'does_not_contain' },
        { text: 'Exists', value: 'exists' },
        { text: 'Not Exists', value: 'not_exists' }
      ],
      valid: false,
      rules: {
        required: value => !!value || 'Required.',
        numberRequired: value => value > -1
      }
    }
  },
  methods: {
    showValue(operator) {
      return ['exists', 'not_exists', 'true', 'false'].indexOf(operator) === -1
    },
    setProperty(condition, e) {
      condition.property = e.value
      condition.type = e.type
      this.$forceUpdate()
    },
    getOperators(type) {
      switch (type) {
        case 'string':
          return this.stringOperators
        case 'integer':
          return this.numberOperators
        case 'float':
          return this.numberOperators
        case 'boolean':
          return this.booleanOperators
        case 'dateTime':
          return this.dateTimeOperators
        case 'array':
          return this.arrayOperators
      }
    }
  }
}
</script>
