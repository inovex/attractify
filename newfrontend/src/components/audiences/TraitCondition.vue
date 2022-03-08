<template>
  <v-card-text>
    <v-row>
      <v-col cols="4">
        <APISelect
          v-if="trait.source === 'custom'"
          dense
          outlined
          placeholder="trait"
          :loadCallback="listCustomTraitProperties"
          :value="trait.key || ''"
          :hide-details="true"
          @change="
            (e) => {
              setProperty(e)
            }
          "
          return-object
        />

        <APISelect
          v-if="trait.source === 'computed'"
          dense
          outlined
          placeholder="trait"
          :loadCallback="listComputedTraits"
          :hide-details="true"
          :value="trait.key || ''"
          @change="
            (e) => {
              setProperty(e)
            }
          "
          return-object
        />
      </v-col>
      <v-col cols="2" dense v-if="trait.key">
        <v-select
          dense
          :items="getOperators(trait.dataType)"
          placeholder="operator"
          outlined
          :hide-details="true"
          :value="trait.operator"
          v-model="trait.operator"
          :rules="[rules.required]"
          @change="
            (t) => {
              setOperator(t)
            }
          "
        ></v-select>
      </v-col>
      <v-col dense cols="auto" v-if="trait.operator && showValue(trait.operator)">
        <v-menu offset-y>
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" text>
              <v-icon>mdi-plus</v-icon>
              compare to
            </v-btn>
          </template>
          <v-list>
            <v-list-item v-for="(item, index) in traitTargets" :key="index" @click="setTarget(item.value)">
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
      <v-col v-if="trait.target === 'static' && showValue(trait.operator)">
        <v-text-field
          outlined
          dense
          placeholder="Value"
          name="value"
          :hide-details="true"
          type="text"
          v-model="trait.value"
          :rules="[rules.required]"
        />
      </v-col>
      <v-col dense cols="auto" v-if="trait.target === 'funnel_property' && trait.operator && showValue(trait.operator)">
        <v-btn @click="waitingForCompare = !waitingForCompare">
          <span v-if="!waitingForCompare">Select event</span>
          <span v-if="waitingForCompare">Cancel select</span>
        </v-btn>
      </v-col>
      <v-col dense v-if="trait.target === 'funnel_property' && trait.operator && showValue(trait.operator)">
        <v-select
          dense
          :items="compareEventProperties"
          outlined
          placeholder="property"
          :hide-details="true"
          :rules="[rules.required]"
          v-model="trait.eventProperty"
        ></v-select>
      </v-col>
    </v-row>
  </v-card-text>
</template>

<script>
import APISelect from '../common/APISelect.vue'

import client from '../../lib/rest/events.js'
import customTraitsClient from '../../lib/rest/customTraits.js'
import computedTraitsClient from '../../lib/rest/computedTraits.js'

import operators from '../../lib/operators'

export default {
  components: { APISelect },
  props: ['trait'],
  data() {
    return {
      waitingForCompare: false,
      compareEventId: null,
      compareEventProperties: [],
      traitTargets: [
        { title: 'with static value', value: 'static' },
        { title: 'with another funnel property', value: 'funnel_property' }
      ],
      valid: false,
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    setProperty(e) {
      this.trait.key = e.value
      this.trait.value = null
      this.trait.eventId = null
      this.trait.eventProperty = null
      this.trait.operator = null
      this.trait.target = null

      if (this.trait.source === 'computed') {
        this.trait.dataType = this.computedTraitType(e.type, e.propertyType)
      } else {
        this.trait.dataType = e.type
      }
      this.$forceUpdate()
    },
    setTarget(target) {
      this.trait.target = target
      this.trait.value = null
      this.trait.eventId = null
      this.trait.eventProperty = null
      this.$forceUpdate()
    },
    setOperator() {
      this.$forceUpdate()
    },
    getOperators(type) {
      switch (type) {
        case 'string':
          return operators.stringOperators
        case 'integer':
          return operators.numberOperators
        case 'float':
          return operators.numberOperators
        case 'boolean':
          return operators.booleanOperators
        case 'dateTime':
          return operators.dateTimeOperators
        case 'array':
          return operators.arrayOperators
      }
    },
    computedTraitType(type, propertyType) {
      switch (type) {
        case 'count_events':
          return 'integer'
        case 'unique_list_count':
          return 'integer'
        case 'aggregation':
          return 'float'
        case 'unique_list':
          return 'array'
        case 'first_event':
          return propertyType
        case 'last_event':
          return propertyType
        case 'most_frequent':
          return propertyType
      }

      return propertyType
    },
    showValue(operator) {
      return ['exists', 'not_exists', 'true', 'false'].indexOf(operator) === -1
    },
    async loadCompareEventProperties() {
      if (!this.trait.eventId) {
        return
      }
      let properties = await client.listProperties(this.trait.eventId)
      this.compareEventProperties = properties.filter((p) => {
        if (p.type === this.trait.dataType) {
          return true
        }
        if (this.trait.dataType === 'float' && p.type === 'integer') {
          return true
        }

        return false
      })
    },
    receiveCompareEvent(event) {
      if (this.waitingForCompare) {
        this.waitingForCompare = false
        this.trait.eventId = event.id
        this.trait.internalEventId = event.internalId
        this.trait.eventProperty = null

        this.loadCompareEventProperties()
      }
    },
    listCustomTraitProperties() {
      return customTraitsClient.listProperties()
    },
    listComputedTraits() {
      return computedTraitsClient.listTraits()
    }
  },
  created() {
    this.$bus.on('audience:event:selectCompareEvent', this.receiveCompareEvent)
    this.loadCompareEventProperties()
  }
}
</script>
