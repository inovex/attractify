<template>
  <v-card outlined>
    <v-card-text>
      <v-row>
        <v-col class="col-lg-3">
          <v-select
            dense
            label="Channel"
            :items="channels"
            v-model="condition.channel"
            :value="condition.channel"
            prepend-icon="mdi-cellphone"
          ></v-select>
        </v-col>
        <v-col class="col-lg-3" v-if="condition.channel">
          <APISelect
            dense
            label="Context Key"
            :loadCallback="listContextKeys"
            :value="condition.key || ''"
            @change="
              e => {
                setProperty(condition, e)
              }
            "
            return-object
          />
        </v-col>
        <v-col class="col-lg-2" v-if="condition.key">
          <v-select
            dense
            :items="getOperators(condition.type)"
            label="Operator"
            :value="condition.operator"
            v-model="condition.operator"
            :rules="[rules.required]"
          ></v-select>
        </v-col>
        <v-col v-if="condition.operator">
          <v-text-field
            v-if="showNumberField(condition)"
            dense
            label="Value"
            name="value"
            type="text"
            v-model.number="condition.value"
            :rules="[rules.required]"
          />
          <div v-else>
            <v-text-field
              v-if="showValue(condition.operator)"
              dense
              label="Value"
              name="value"
              type="text"
              v-model="condition.value"
              :rules="[rules.required]"
            />
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import APISelect from '../common/APISelect.vue'
import contextClient from '../../lib/rest/contexts.js'
import channelClient from '../../lib/rest/channels.js'

export default {
  components: { APISelect },
  props: ['condition'],
  data() {
    return {
      channels: [],
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
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    showNumberField(condition) {
      if (['integer', 'float'].indexOf(condition.type) > -1) {
        return true
      }

      if (
        condition.type === 'dateTime' &&
        ['within_last_days', 'within_next_days', 'before_last_days', 'after_next_days'].indexOf(condition.operator) > -1
      ) {
        return true
      }

      return false
    },
    showValue(operator) {
      return ['exists', 'not_exists', 'true', 'false'].indexOf(operator) === -1
    },
    setProperty(condition, e) {
      condition.key = e.value
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
    },
    listContextKeys() {
      return contextClient.listProperties(this.condition.channel)
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
