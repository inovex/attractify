<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Reactions</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="reactions" />
          </v-toolbar>
          <v-data-table
            :headers="headers"
            :items="reactions"
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
                      label="Action"
                      :loadCallback="loadActions"
                      @change="setAction"
                      icon="mdi-ticket-percent"
                      clearable
                      return-object
                    />
                  </v-col>
                  <v-col>
                    <v-select
                      :items="eventList"
                      label="Events"
                      @change="resetAndLoad"
                      prepend-icon="mdi-bell"
                      v-model="events"
                      multiple
                    ></v-select>
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
                        @change="resetAndLoad"
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
                      @keydown.enter="resetAndLoad"
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
                <v-icon title="Show reaction">mdi-magnify</v-icon>
              </v-btn>
              <v-btn @click="remove(item)" icon>
                <v-icon title="Delete reaction">mdi-delete</v-icon>
              </v-btn>
            </template>
            <template v-slot:item.event="{ item }">
              <span>{{ item.event }}</span>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.duration="{ item }">
              <span>{{ item.duration }}s</span>
            </template>
            <template v-slot:no-data>No Reactions Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
    <v-dialog v-model="dialog" max-width="700px" closeable>
      <v-card>
        <v-card-title>
          <span class="headline">Reaction details</span>
        </v-card-title>
        <v-card-text>
          {{ action.profileId }}
          <h4>Properties</h4>
          <v-card outlined class="pa-2">
            <pre>{{ action.properties }}</pre>
          </v-card>
          <br />
          <h4>Context</h4>
          <v-card outlined class="pa-2">
            <pre style="overflow: auto;">{{ action.context }}</pre>
          </v-card>
          <br />
          <h4>Result</h4>
          <v-card outlined class="pa-2">
            <pre style="overflow: auto;">{{ action.result }}</pre>
          </v-card>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import Help from './Help'
import APISelect from './common/APISelect.vue'
import actionsClient from '../lib/rest/actions.js'
import reactionsClient from '../lib/rest/reactions.js'
import moment from 'moment'

export default {
  components: { APISelect, Help },
  data() {
    return {
      dialog: false,
      action: {},
      reactions: [],
      userId: '',
      events: [],
      actionId: null,
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
      eventList: [
        { text: 'Delivered', value: 'delivered' },
        { text: 'Shown', value: 'shown' },
        { text: 'Hidden', value: 'hidden' },
        { text: 'Declined', value: 'declined' },
        { text: 'Accepted', value: 'accepted' }
      ],
      headers: [
        {
          text: 'Event',
          align: 'left',
          value: 'event'
        },
        { text: 'Channel', value: 'channel' },
        { text: 'Date', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ]
    }
  },
  methods: {
    setAction(e) {
      if (e !== undefined) {
        this.actionId = e.value
      } else {
        this.actionId = null
      }
      this.resetAndLoad()
    },
    async loadActions() {
      return await actionsClient.listNames()
    },
    details(action) {
      this.action = action
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
        events: this.events.join(','),
        start: this.range.start,
        end: this.range.end
      }

      if (this.userId && this.userId.length > 0) {
        params.userId = this.userId
      }

      if (this.actionId) {
        params.actionId = this.actionId
      }

      if (params.itemsPerPage > 0) {
        try {
          const res = await reactionsClient.list(params)
          this.reactions = res.reactions
          this.totalItems = res.count
        } catch (_) {
          _
        }
      }
    },
    async remove(action) {
      try {
        await reactionsClient.delete(action.id)
        this.$notify.success('The reaction has been removed successfully.')
        this.load()
      } catch (_) {
        _
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
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
    this.load()
  }
}
</script>
