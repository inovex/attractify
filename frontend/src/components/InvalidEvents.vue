<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Invalid Events</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="events" />
          </v-toolbar>

          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="eventSpecs">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="edit(item)"> <v-icon title="Edit event">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete event">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.name="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.description="{ item }">
              <span>{{ item.description }}</span>
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
import client from '../lib/rest/events'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      eventSpecs: [],
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Description', value: 'description' },
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
    create() {
      this.$router.push({ path: '/invalid-events' })
    },
    edit(eventSpec) {
      this.$router.push({ path: `/invalid-events/${eventSpec.id}` })
    },
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
