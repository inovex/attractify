<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Computed Trait</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="computedTrait" />
            </v-toolbar>
            <v-card-text>
              <v-row>
                <v-col class="col-lg-4">
                  <v-text-field
                    label="Name"
                    name="name"
                    prepend-icon="mdi-pencil"
                    type="text"
                    v-model="computedTrait.name"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-4">
                  <v-text-field
                    label="Key"
                    name="key"
                    prepend-icon="mdi-text"
                    type="text"
                    v-model="computedTrait.key"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    :disabled="!!computedTrait.id"
                    prepend-icon="mdi-cog"
                    :items="types"
                    label="Type"
                    @change="computedTrait.properties = {}"
                    v-model="computedTrait.type"
                    :rules="[rules.required]"
                  ></v-select>
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    :items="events"
                    label="Event"
                    @change="setEventId"
                    prepend-icon="mdi-bell"
                    v-model="computedTrait.eventId"
                    :rules="[rules.required]"
                    return-object
                  ></v-select>
                </v-col>
              </v-row>
            </v-card-text>

            <div v-if="!!computedTrait.type && !!computedTrait.eventId">
              <v-card-text v-if="computedTrait.type === 'event_counter'">
                <v-row>
                  <v-col>No config needed</v-col>
                </v-row>
              </v-card-text>

              <v-card-text v-if="computedTrait.type === 'aggregation'">
                <v-row>
                  <v-col class="col-lg-3">
                    <v-select
                      label="Property Name"
                      prepend-icon="mdi-text"
                      :items="eventProperties.filter(p => p.type === 'float' || p.type === 'integer')"
                      v-model="computedTrait.properties.property"
                      @change="setProperty"
                      return-object
                    ></v-select>
                  </v-col>
                  <v-col class="col-lg-2">
                    <v-select
                      prepend-icon="mdi-cog"
                      :items="aggregationTypes"
                      label="Aggregation Type"
                      v-model="computedTrait.properties.aggregationType"
                      :rules="[rules.required]"
                    ></v-select>
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-text v-if="computedTrait.type === 'most_frequent'">
                <v-row>
                  <v-col class="col-lg-3">
                    <v-select
                      label="Property Name"
                      prepend-icon="mdi-text"
                      :items="eventProperties.filter(p => ['string', 'integer', 'float', 'array'].indexOf(p.type) > -1)"
                      v-model="computedTrait.properties.property"
                      @change="setProperty"
                      return-object
                    ></v-select>
                  </v-col>
                  <v-col class="col-lg-2">
                    <v-text-field
                      label="Minimum frequency"
                      prepend-icon="mdi-number"
                      type="number"
                      v-model.number="computedTrait.properties.minFrequency"
                      :rules="[rules.required]"
                    />
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-text v-if="computedTrait.type === 'first_event' || computedTrait.type === 'last_event'">
                <v-row>
                  <v-col class="col-lg-4">
                    <v-switch
                      label="Use timestamp instead of property value"
                      v-model="computedTrait.properties.useTimestamp"
                      @change="computedTrait.properties.type = 'dateTime'"
                    ></v-switch>
                  </v-col>
                  <v-col class="col-lg-3" v-if="!computedTrait.properties.useTimestamp">
                    <v-select
                      label="Property Name"
                      prepend-icon="mdi-text"
                      :items="eventProperties"
                      v-model="computedTrait.properties.property"
                      @change="setProperty"
                      return-object
                    ></v-select>
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-text v-if="computedTrait.type === 'unique_list' || computedTrait.type === 'unique_list_count'">
                <v-row>
                  <v-col class="col-lg-3">
                    <v-select
                      label="Property Name"
                      prepend-icon="mdi-text"
                      :items="eventProperties.filter(p => ['string', 'integer', 'float'].indexOf(p.type) > -1)"
                      v-model="computedTrait.properties.property"
                      @change="setProperty"
                      return-object
                    ></v-select>
                  </v-col>
                </v-row>
              </v-card-text>
            </div>

            <div v-if="!!computedTrait.eventId">
              <v-card-title>
                Add event conditions...
                <v-btn icon @click="addCondition()">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-row>
                  <v-col>
                    <v-row v-for="(condition, key) of computedTrait.conditions" :key="key">
                      <v-col>
                        <Condition :properties="eventProperties" :condition="condition" />
                      </v-col>
                      <v-col class="col-lg-1">
                        <v-btn icon @click="deleteCondition(key)">
                          <v-icon>mdi-trash-can-outline</v-icon>
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-col>
                </v-row>
              </v-card-text>
            </div>

            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" :disabled="!valid" @click="save()" text>Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
    <v-overlay :value="showLoading">
      <v-progress-circular indeterminate size="64"></v-progress-circular>
      <span class="ml-4">Loading computed trait data, this might take a while...</span>
    </v-overlay>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/computedTraits.js'
import eventsClient from '../lib/rest/events.js'
import Condition from './computedTraits/Condition.vue'

export default {
  components: { Condition, Help },
  data() {
    return {
      computedTrait: {
        conditions: [],
        properties: {}
      },
      eventId: '',
      events: [],
      eventProperties: [],
      types: [
        { text: 'Count Event (int)', value: 'count_events' },
        { text: 'Aggregation (float)', value: 'aggregation' },
        {
          text: 'Most frequent (timestamp or property value)',
          value: 'most_frequent'
        },
        {
          text: 'First Event (timestamp or property value)',
          value: 'first_event'
        },
        {
          text: 'Last Event (timestamp or property value)',
          value: 'last_event'
        },
        { text: 'Unique List (array)', value: 'unique_list' },
        { text: 'Unique list count (int)', value: 'unique_list_count' }
      ],
      aggregationTypes: [
        { text: 'Sum', value: 'sum' },
        { text: 'Average', value: 'avg' },
        { text: 'Minimum', value: 'min' },
        { text: 'Maximum', value: 'max' }
      ],
      valid: false,
      showLoading: false,
      dialog: false,
      previewResult: {},
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  watch: {
    async computedTrait(ct) {
      if (ct.eventId) {
        let evt = this.events.filter(e => e.value === ct.eventId)
        if (evt && evt.length > 0) {
          this.setEventId(evt[0])
        }
      }
    }
  },
  methods: {
    setProperty(e) {
      this.computedTrait.properties.property = e.value
      this.computedTrait.properties.type = e.type
    },
    async setEventId(e) {
      this.computedTrait.eventId = e.value
      this.eventProperties = await eventsClient.listProperties(e.value)
    },
    addCondition() {
      this.computedTrait.conditions.push({})
    },
    deleteCondition(index) {
      this.computedTrait.conditions.splice(index, 1)
    },
    async save() {
      try {
        let res = null
        if (this.computedTrait.id) {
          res = await client.update(this.computedTrait)
        } else {
          res = await client.create(this.computedTrait)
        }

        if (res && res.id) {
          this.computedTrait.id = res.id
        }

        this.$notify.success('The computed trait has been saved successfully.')
      } catch (e) {
        this.$notify.error('Could not save computed trait.')
      }
    },
    async refresh() {
      this.showLoading = true
      try {
        await client.refresh(this.computedTrait.id)
        this.$notify.success(`The computed trait has been successfully refreshed.`)
      } catch (error) {
        this.$notify.error('Could not refresh computed trait.')
      } finally {
        this.showLoading = false
      }
    }
  },
  async created() {
    const id = this.$route.params.id
    this.events = await eventsClient.listEventNames()
    if (id) {
      try {
        this.computedTrait = await client.show(id)
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
  }
}
</script>
