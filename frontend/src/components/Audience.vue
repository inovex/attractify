<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Audience</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="audience" />
            </v-toolbar>
            <v-card-text>
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="Name"
                    name="name"
                    prepend-icon="mdi-pencil"
                    type="text"
                    v-model="audience.name"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-10">
                  <v-text-field
                    label="Description"
                    name="type"
                    prepend-icon="mdi-text"
                    type="text"
                    v-model="audience.description"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-10">
                  <v-checkbox
                    label="Include anonymous profiles"
                    prepend-icon="mdi-incognito"
                    v-model="audience.includeAnonymous"
                  />
                </v-col>
              </v-row>
            </v-card-text>

            <v-card-title>
              Select all profiles, who...
              <v-menu offset-y>
                <template v-slot:activator="{ on }">
                  <v-btn v-on="on" icon>
                    <v-icon>mdi-plus</v-icon>
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item v-for="(item, index) in conditionTypes" :key="index" @click="addCondition(item.value)">
                    <v-list-item-title>{{ item.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-card-title>
            <v-card-text>
              <v-row>
                <v-col>
                  <v-row v-for="(event, key) of audience.events" :key="`event-${key}`">
                    <v-col>
                      <v-card outlined :class="event.parentId ? 'ml-8' : ''" elevation="1" tile>
                        <EventCondition :event="event" />
                        <v-divider></v-divider>
                        <v-card-actions>
                          <v-btn rounded @click="addFunnelEvent(event.internalId, false)">
                            <v-icon>mdi-filter-plus-outline</v-icon>
                            <span class="ml-2">and then performed...</span>
                          </v-btn>
                          <v-btn rounded @click="addFunnelEvent(event.internalId, true)">
                            <v-icon>mdi-filter-minus-outline</v-icon>
                            <span class="ml-2">and then didn't perform...</span>
                          </v-btn>
                          <v-spacer></v-spacer>
                          <v-btn icon @click="deleteCondition('events', key)">
                            <v-icon>mdi-trash-can-outline</v-icon>
                          </v-btn>
                        </v-card-actions>
                      </v-card>
                    </v-col>
                  </v-row>
                  <v-row v-for="(trait, key) of audience.traits" :key="`trait-${key}`">
                    <v-col>
                      <v-card outlined elevation="1" tile>
                        <v-app-bar flat color="blue-grey" dense class="lighten-4">
                          <v-toolbar-title class="pl-0 text-s">
                            <span v-if="trait.source === 'custom'">have custom trait</span>
                            <span v-if="trait.source === 'computed'">have computed trait</span>
                          </v-toolbar-title>
                          <v-spacer></v-spacer>
                          <v-btn icon @click="deleteCondition('traits', key)">
                            <v-icon>mdi-trash-can-outline</v-icon>
                          </v-btn>
                        </v-app-bar>
                        <TraitCondition :trait="trait" />
                      </v-card>
                    </v-col>
                  </v-row>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
    <v-overlay :value="showLoading">
      <v-progress-circular indeterminate size="64"></v-progress-circular>
      <span class="ml-4">Loading audience data, this might take a while...</span>
    </v-overlay>

    <v-dialog v-model="dialog" max-width="700px" closeable>
      <v-card>
        <v-card-title>
          <span class="headline">Audience Preview</span>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col> Preview run in {{ previewResult.executionTime }}ms. </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-simple-table dense :fixed-header="true">
                <template v-slot:default>
                  <thead>
                    <tr>
                      <th class="text-left">ID</th>
                      <th class="text-left">Details</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item, key) in previewResult.profiles" :key="key">
                      <td>{{ item.id }}</td>
                      <td>
                        <v-btn icon :to="{ path: `/profile/${item.id}` }">
                          <v-icon>mdi-magnify</v-icon>
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn rounded @click="dialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-col class="sticky text-center">
      <v-spacer />
      <v-btn rounded elevation="2" @click="cancel()">Cancel</v-btn>
      <v-btn rounded elevation="2" color="secondary" :disabled="!valid" @click="refresh()">Refresh</v-btn>
      <v-btn rounded elevation="2" color="secondary" :disabled="!valid" @click="preview()">Preview</v-btn>
      <v-btn rounded elevation="2" color="primary" :disabled="!valid" @click="save()">Save</v-btn>
    </v-col>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/audiences'
import TraitCondition from './audiences/TraitCondition.vue'
import EventCondition from './audiences/EventCondition.vue'
import { v4 as uuid } from 'uuid'

export default {
  components: { TraitCondition, EventCondition, Help },
  data() {
    return {
      audience: {
        events: [],
        traits: []
      },
      conditionTypes: [
        { title: 'performed event', value: 'event' },
        { title: 'have custom trait', value: 'custom_trait' },
        { title: 'have computed trait', value: 'computed_trait' }
      ],
      valid: false,
      showLoading: false,
      dialog: false,
      previewResult: {},
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    cancel(){
      this.$router.push('/audiences')
    },
    show(item) {
      this.$router.push({ path: `/profile/${item.id}` })
    },
    addCondition(type) {
      let condition = {}
      if (type === 'event') {
        condition.internalId = uuid()
        condition.timeWindowOperator = 'any_time'
        condition.operator = 'exactly'
        condition.properties = []
        this.audience.events.push(condition)
      } else {
        if (type === 'custom_trait') {
          condition.source = 'custom'
        } else {
          condition.source = 'computed'
        }
        this.audience.traits.push(condition)
      }
    },
    addFunnelEvent(parentId = null, exclude = false) {
      let condition = {
        internalId: uuid(),
        parentId: parentId,
        exclude: exclude,
        timeWindowOperator: 'any_time',
        operator: 'exactly',
        properties: []
      }

      let idx = this.audience.events.map((e) => e.internalId).indexOf(parentId)
      if (idx > -1) {
        this.audience.events.splice(idx + 1, 0, condition)
      } else {
        this.audience.events.push(condition)
      }
    },
    deleteCondition(type, index) {
      this.audience[type].splice(index, 1)
    },
    async save() {
      try {
        let res = null
        if (this.audience.id) {
          res = await client.update(this.audience)
        } else {
          res = await client.create(this.audience)
        }

        if (res && res.id) {
          this.audience.id = res.id
        }

        this.$notify.success('The audience has been saved successfully.')
      } catch (e) {
        this.$notify.error('Could not save audience.')
      }
    },
    async preview() {
      this.showLoading = true
      try {
        this.previewResult = await client.preview(this.audience)
        this.dialog = true
      } catch (_) {
        this.$notify.error('Could not preview audience.')
      }
      this.showLoading = false
    },
    async refresh() {
      this.showLoading = true
      try {
        let res = await client.refresh(this.audience.id)
        this.$notify.success(`The audience has been successfully refreshed with ${res.count} profiles.`)
      } catch (_) {
        this.$notify.error('Could not refresh audience.')
      } finally {
        this.showLoading = false
      }
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.audience = await client.show(id)
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
  }
}
</script>

<style scoped>
.sticky {
  margin: 0.5rem 0 0 0;
  position: -webkit-sticky;
  position: sticky;
  bottom: 1rem;
  z-index: 1;
}
.sticky button{
  margin: 0 0.5rem;
}
</style>
