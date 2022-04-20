<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Invalid Events</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="invalid-events" />
          </v-toolbar>

          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="eventSpecs">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="remove(item)"> <v-icon title="Delete event">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>

            <template v-slot:item.type="{ item }">
              <!--<v-icon v-if="item.type === 'properties'">mdi-event-note</v-icon>
              <v-icon v-if="item.type === 'context'">mdi-event-note</v-icon>-->
              <span v-if="item.type === 'properties'">Property</span>
              <span v-if="item.type === 'context'">Context</span>
            </template>
            <template v-slot:item.properties="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.context="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.error="{ item }">
              <span>{{ item.description }}</span>
            </template>

            <template v-slot:item.name="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:no-data>No Event Definitions Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/invalidEvents'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      eventSpecs: [],
      headers: [
        { text: 'Type', value: 'type' },
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Error', value: 'error', sortable: false },
        { text: 'Properties', value: 'properties', sortable: false },
        { text: 'Context', value: 'context', sortable: false },
        { text: 'Created', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      valid: false,
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    remove(eventSpec) {
      if (confirm('Do you really want to delete this event?')) {
        client.delete(eventSpec.id)

        let idx = this.eventSpecs.findIndex((es) => es.id === eventSpec.id)
        this.eventSpecs.splice(idx, 1)
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    timeAgo(date) {
      return moment(date).fromNow()
    }
  },
  async created() {
    try {
      this.eventSpecs = await client.list()
    } catch (_) {
      _
    }
  }
}
</script>
