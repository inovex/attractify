<template>
  <v-row>
    <v-col cols="auto" dense>
      <div class="mt-3">and where</div>
    </v-col>
    <v-col dense cols="3">
      <APISelect
        dense
        outlined
        placeholder="property"
        :loadCallback="loadProperties(event.id)"
        :value="property.key || ''"
        return-object
        :hide-details="true"
        @change="
          e => {
            setProperty(e)
          }
        "
      />
    </v-col>
    <v-col dense cols="2" v-if="property.key">
      <v-select
        dense
        outlined
        :hide-details="true"
        :items="getOperators(property.dataType, property.target)"
        :value="property.operator"
        v-model="property.operator"
        :rules="[rules.required]"
        @change="
          t => {
            setOperator(t)
          }
        "
      ></v-select>
    </v-col>
    <v-col dense cols="auto" v-if="property.operator && showValue(property.operator, property.target)">
      <div v-if="property.target === 'static'" class="mt-3">static value</div>
      <div v-if="property.target === 'custom_trait'" class="mt-3">custom trait</div>
      <div v-if="property.target === 'computed_trait'" class="mt-3">computed trait</div>
      <div v-if="property.target === 'funnel_property'" class="mt-3">funnel event property</div>
    </v-col>
    <v-col dense v-if="property.target === 'static' && property.operator">
      <v-text-field
        dense
        outlined
        v-if="showValue(property.operator, property.target)"
        name="value"
        type="text"
        :hide-details="true"
        v-model="property.value"
        :rules="[rules.required]"
      />
    </v-col>
    <v-col dense v-if="property.target === 'custom_trait' && property.operator">
      <APISelect
        dense
        outlined
        placeholder="custom trait key"
        :hide-details="true"
        v-if="showValue(property.operator, property.target)"
        :loadCallback="listCustomTraitProperties(property.dataType)"
        :value="property.traitKey || ''"
        @change="
          t => {
            setCompareProperty(t)
          }
        "
        :rules="[rules.required]"
        return-object
      />
    </v-col>
    <v-col dense v-if="property.target === 'computed_trait' && property.operator">
      <APISelect
        dense
        outlined
        v-if="showValue(property.operator, property.target)"
        placeholder="computed trait key"
        :loadCallback="listComputedTraits(property.dataType)"
        :value="property.traitKey || ''"
        :hide-details="true"
        :rules="[rules.required]"
        @change="
          t => {
            setCompareProperty(t)
          }
        "
        return-object
      />
    </v-col>
    <v-col dense cols="auto" v-if="property.target === 'funnel_property' && property.operator">
      <v-btn @click="waitingForCompare = !waitingForCompare">
        <span v-if="!waitingForCompare">Select event</span>
        <span v-if="waitingForCompare">Cancel select</span>
      </v-btn>
    </v-col>
    <v-col dense v-if="property.target === 'funnel_property' && property.operator">
      <v-select
        dense
        :items="compareEventProperties"
        outlined
        placeholder="property"
        :hide-details="true"
        :rules="[rules.required]"
        v-model="property.eventProperty"
      ></v-select>
    </v-col>
    <v-col cols="auto" dense>
      <v-btn icon @click="deleteCallback()">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script>
import APISelect from '../common/APISelect.vue'

import client from '../../lib/rest/events.js'
import customTraitsClient from '../../lib/rest/customTraits.js'
import computedTraitsClient from '../../lib/rest/computedTraits.js'

import operators from '../../lib/operators'

export default {
  components: { APISelect },
  props: ['property', 'event', 'deleteCallback'],
  data() {
    return {
      waitingForCompare: false,
      compareEventId: null,
      compareEventProperties: [],
      valid: false,
      rules: {
        required: value => !!value || 'Required.',
        numberRequired: value => value > -1
      }
    }
  },
  methods: {
    setProperty(e) {
      this.property.key = e.value
      this.property.dataType = e.type
      this.property.traitKey = null
      this.property.value = null
      this.property.eventId = null
      this.property.eventProperty = null
      this.property.operator = null
      this.$forceUpdate()
    },
    setOperator() {
      this.$forceUpdate()
    },
    setCompareProperty(e) {
      this.property.traitKey = e.value
    },
    setCompareEventProperty(e) {
      this.property.eventProperty = e.value
    },
    showValue(operator, target) {
      if (target !== 'static' && ['true', 'false'].indexOf(operator) > -1) {
        return true
      }
      return ['exists', 'not_exists', 'true', 'false'].indexOf(operator) === -1
    },
    receiveCompareEvent(event) {
      if (this.waitingForCompare) {
        this.waitingForCompare = false
        this.property.eventId = event.id
        this.property.internalEventId = event.internalId
        this.property.eventProperty = null

        this.loadCompareEventProperties()
      }
    },
    async loadCompareEventProperties() {
      if (!this.property.eventId) {
        return
      }

      let properties = await client.listProperties(this.property.eventId)
      this.compareEventProperties = properties.filter(p => p.type === this.property.dataType)
    },
    getOperators(type, target) {
      if (target === 'static') {
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
      } else {
        switch (type) {
          case 'string':
            return operators.stringComparisonOperators
          case 'integer':
            return operators.numberComparisonOperators
          case 'float':
            return operators.numberComparisonOperators
          case 'boolean':
            return operators.booleanComparisonOperators
          case 'dateTime':
            return operators.dateTimeComparisonOperators
          case 'array':
            return operators.arrayComparisonOperators
        }
      }
    },
    mapType(type) {
      if (type === 'integer' || type === 'float') {
        return 'float'
      }
      return type
    },
    loadProperties(id) {
      return () => {
        return client.listProperties(id)
      }
    },
    listCustomTraitProperties(type) {
      return async () => {
        let traits = await customTraitsClient.listProperties()
        return traits.filter(t => t.type === type)
      }
    },
    listComputedTraits(type) {
      return async () => {
        let traits = await computedTraitsClient.listTraits()
        return traits.filter(t => {
          if (t.propertyType === 'float' && type === 'integer') {
            return true
          }
          if (t.propertyType === 'integer' && type === 'float') {
            return true
          }

          return t.propertyType === type
        })
      }
    }
  },
  created() {
    this.$bus.$on('audience:event:selectCompareEvent', this.receiveCompareEvent)
    this.loadCompareEventProperties()
  }
}
</script>
