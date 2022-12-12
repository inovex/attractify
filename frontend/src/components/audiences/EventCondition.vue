<template>
  <div>
    <v-app-bar
      flat
      @click="selectCompareEvent(event)"
      :class="{ event__selector: true, 'lighten-5': event.parentId, 'lighten-3': !event.parentId }"
    >
      <v-toolbar-title class="pl-0">
        <span v-if="!event.exclude && !event.parentId">performed</span>
        <span v-if="event.exclude && !event.parentId">did not perform</span>
        <span v-if="!event.exclude && event.parentId">and then performed</span>
        <span v-if="event.exclude && event.parentId">and then did not perform</span>
      </v-toolbar-title>
      <div class="mt-6 ml-4">
        <APISelect
          dense
          :loadCallback="loadEvents"
          :value="event.id || ''"
          @change="setEvent"
          outlined
          placeholder="event"
          return-object
        />
      </div>
    </v-app-bar>
    <v-card-text>
      <v-row>
        <v-col cols="auto" dense>
          <div class="mt-3">in channel</div>
        </v-col>
        <v-col dense>
          <v-select
            dense
            :items="channels"
            multiple
            outlined
            :hide-details="true"
            placeholder="channel"
            :value="event.channels"
            @change="changes"
            v-model="event.channels"
          ></v-select>
        </v-col>
        <v-col dense>
          <v-select
            dense
            outlined
            :hide-details="true"
            :items="eventOperators"
            :value="event.operator"
            @change="changes"
            v-model="event.operator"
            :rules="[rules.required]"
          ></v-select>
        </v-col>
        <v-col cols="1" dense>
          <v-text-field
            dense
            outlined
            name="count"
            type="number"
            @input="changes"
            :hide-details="true"
            v-model.number="event.count"
            :rules="[rules.numberRequired]"
          />
        </v-col>
        <v-col cols="auto" dense>
          <div class="mt-3">time(s)</div>
        </v-col>
        <v-col cols="2" dense>
          <v-select
            dense
            outlined
            :hide-details="true"
            :items="timeWindowTypes"
            :value="event.timeWindowOperator"
            @change="changes"
            v-model="event.timeWindowOperator"
            :rules="[rules.required]"
          ></v-select>
        </v-col>
        <v-col dense>
          <v-text-field
            dense
            outlined
            v-if="event.timeWindowOperator === 'within'"
            label="days"
            name="start"
            type="text"
            @input="changes"
            :hide-details="true"
            v-model.number="event.timeWindowStart"
            :rules="[rules.required]"
          />
          <v-text-field
            dense
            outlined
            v-if="event.timeWindowOperator === 'in_between'"
            label="start days"
            name="start"
            type="text"
            @input="changes"
            :hide-details="true"
            v-model.number="event.timeWindowStart"
            :rules="[rules.required]"
          />
        </v-col>
        <v-col dense>
          <v-text-field
            dense
            outlined
            v-if="event.timeWindowOperator === 'in_between'"
            label="end days"
            name="end"
            type="text"
            @input="changes"
            v-model.number="event.timeWindowEnd"
            :rules="[rules.required]"
          />
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-text>
      <PropertyCondition
        v-for="(property, key) of event.properties"
        :key="key"
        :property="property"
        :event="event"
        :delete-callback="
          () => {
            deleteProperty(key)
            changes()
          }
        "
      />
      <v-row>
        <v-col cols="auto" v-if="event.id">
          <v-menu offset-y>
            <template v-slot:activator="{ on }">
              <v-btn v-on="on" rounded>
                <v-icon>mdi-plus</v-icon>
                compare property
              </v-btn>
            </template>
            <v-list>
              <v-list-item v-for="(item, index) in propertyTargets" :key="index" @click="addProperty(item.value)">
                <v-list-item-title>{{ item.title }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>
      </v-row>
    </v-card-text>
  </div>
</template>

<script>
import PropertyCondition from './PropertyCondition'
import APISelect from '../common/APISelect.vue'
import client from '../../lib/rest/events.js'
import channelClient from '../../lib/rest/channels.js'

export default {
  components: { APISelect, PropertyCondition },
  props: ['event'],
  data() {
    return {
      channels: [],
      propertyTargets: [
        { title: 'with static value', value: 'static' },
        { title: 'with custom trait', value: 'custom_trait' },
        { title: 'with computed trait', value: 'computed_trait' },
        { title: 'with another funnel property', value: 'funnel_property' }
      ],
      timeWindowTypes: [
        { text: 'any time', value: 'any_time' },
        { text: 'within', value: 'within' },
        { text: 'in between', value: 'in_between' }
      ],
      eventOperators: [
        { text: 'exactly', value: 'exactly' },
        { text: 'less than', value: 'less_than' },
        { text: 'more than', value: 'more_than' },
        { text: 'less or exactly', value: 'less_or_exactly' },
        { text: 'more or exactly', value: 'more_or_exactly' }
      ],
      useTimeWindow: false,
      valid: false,
      rules: {
        required: (value) => !!value || 'Required.',
        numberRequired: (value) => value > -1
      }
    }
  },
  methods: {
    setEvent(e) {
      this.event.properties = []
      this.event.id = e.value || null
      this.$forceUpdate()
      this.changes()
    },
    addProperty(target) {
      this.event.properties.push({ target: target })
    },
    selectCompareEvent(event) {
      this.$bus.$emit('audience:event:selectCompareEvent', { id: event.id, internalId: event.internalId })
    },
    deleteProperty(index) {
      this.event.properties.splice(index, 1)
    },
    loadEvents() {
      return client.listEventNames()
    },
    changes() {
      this.$emit('change')
    }
  },
  async created() {
    this.useTimeWindow = !!this.event.timeWindowOperator
    this.channels = await channelClient.select()
  }
}
</script>

<style>
.event__selector:hover {
  cursor: pointer;
}
</style>
