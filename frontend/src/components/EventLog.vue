<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Event Log</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="eventLog" />
          </v-toolbar>
          <v-data-table
            :headers="headers"
            :items="events"
            :options.sync="options"
            :disable-sort="true"
            :server-items-length="totalItems"
            :footer-props="{
              'items-per-page-options': [5, 10, 25, 50]
            }"
          >
            <template v-slot:top>
              <v-container>
                <v-row>
                  <v-col>
                    <APISelect
                      label="Event"
                      :loadCallback="loadEventList"
                      @change="setEvent"
                      icon="mdi-bell"
                      multiple
                      clearable
                      return-object
                    />
                  </v-col>
                  <v-col>
                    <v-menu
                      ref="menu.start"
                      v-model="menu.start"
                      :close-on-content-click="false"
                      transition="scale-transition"
                      offset-y
                      max-width="290px"
                      min-width="290px"
                    >
                      <template v-slot:activator="{ on }">
                        <v-text-field
                          v-model="range.start"
                          label="Date"
                          persistent-hint
                          prepend-icon="mdi-clock-start"
                          v-on="on"
                          clearable
                        ></v-text-field>
                      </template>
                      <v-date-picker
                        v-model="range.start"
                        no-title
                        @input="menu.start = false"
                        @change="resetAndLoad()"
                      ></v-date-picker>
                    </v-menu>
                  </v-col>
                  <v-col>
                    <v-menu
                      ref="menu.end"
                      v-model="menu.end"
                      :close-on-content-click="false"
                      transition="scale-transition"
                      offset-y
                      max-width="290px"
                      min-width="290px"
                    >
                      <template v-slot:activator="{ on }">
                        <v-text-field
                          v-model="range.end"
                          label="Date"
                          persistent-hint
                          prepend-icon="mdi-clock-end"
                          v-on="on"
                          clearable
                        ></v-text-field>
                      </template>
                      <v-date-picker
                        v-model="range.end"
                        no-title
                        @input="menu.end = false"
                        @change="resetAndLoad"
                      ></v-date-picker>
                    </v-menu>
                  </v-col>
                  <v-col>
                    <v-text-field
                      label="User ID"
                      prepend-icon="mdi-account"
                      type="text"
                      v-model="userId"
                      @keydown.enter="resetAndLoad()"
                      clearable
                    />
                  </v-col>
                  <v-col class="col-lg-1">
                    <v-btn @click="resetAndLoad()" icon>
                      <v-icon>mdi-refresh</v-icon>
                    </v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </template>
            <template v-slot:item.action="{ item }">
              <v-btn @click="details(item)" icon>
                <v-icon title="Show event">mdi-magnify</v-icon>
              </v-btn>
              <v-btn @click="remove(item)" icon>
                <v-icon title="Delete event">mdi-delete</v-icon>
              </v-btn>
            </template>
            <template v-slot:item.name="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.duration="{ item }">
              <span>{{ item.duration }}s</span>
            </template>
            <template v-slot:no-data>No Events Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
    <v-dialog v-model="dialog" max-width="700px" closeable>
      <v-card>
        <v-card-title>
          <span class="headline">Event details</span>
        </v-card-title>
        <v-card-text>
          {{ event.profileId }}
          <h4>Properties</h4>
          <v-card outlined class="pa-2">
            <pre style="overflow: auto">{{ event.properties }}</pre>
          </v-card>
          <br />
          <h4>Context</h4>
          <v-card outlined class="pa-2">
            <pre style="overflow: auto">{{ event.context }}</pre>
          </v-card>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import Help from './Help'
import APISelect from './common/APISelect.vue'
import eventsClient from '../lib/rest/events.js'
import eventLogClient from '../lib/rest/eventLog.js'
import moment from 'moment'

export default {
  components: { APISelect, Help },
  data() {
    return {
      dialog: false,
      eventMap: {},
      events: [],
      event: {},
      userId: '',
      selectedEvents: [],
      menu: {
        start: false,
        end: false
      },
      range: {
        start: '',
        end: '',
        interval: 'day'
      },
      totalItems: 0,
      options: {},
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Channel', value: 'channel' },
        { text: 'Date', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ]
    }
  },
  methods: {
    setEvent(e) {
      this.selectedEvents = e.map((ev) => ev.value)
      this.resetAndLoad()
    },
    async loadEvents() {
      let events = await eventsClient.listEventNames()
      for (let e of events) {
        this.eventMap[e.value] = e.text
      }
    },
    async loadEventList() {
      return await eventsClient.listEventNames()
    },
    details(event) {
      this.event = event
      this.dialog = true
    },
    resetAndLoad() {
      this.options.page = 1
      return this.load()
    },
    async load() {
      const params = {
        page: this.options.page,
        itemsPerPage: this.options.itemsPerPage,
        events: this.selectedEvents.join(','),
        start: this.range.start,
        end: this.range.end
      }

      if (this.userId && this.userId.length > 0) {
        params.userId = this.userId
      }

      if (this.event && this.event.length > 0) {
        params.event = this.event
      }

      if (params.itemsPerPage > 0) {
        try {
          const res = await eventLogClient.list(params)
          console.log(res)
          this.events = res.events
          for (let e of this.events) {
            e.name = this.resolveEventName(e.eventId)
          }
          this.totalItems = res.count
        } catch (_) {
          _
        }
      }
    },
    async remove(event) {
      try {
        await eventLogClient.delete(event.id)
        this.$notify.success('The event has been removed successfully.')
        this.load()
      } catch (_) {
        _
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    resolveEventName(eventId) {
      return this.eventMap[eventId]
    }
  },
  watch: {
    all() {
      this.totalItems = this.total
    },
    options() {
      this.load()
    }
  },
  async created() {
    await this.loadEvents()
    this.load()
  }
}
</script>
