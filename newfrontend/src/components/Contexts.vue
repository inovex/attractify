<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Context Definitions</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="contexts" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>

          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="contexts">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="edit(item)"> <v-icon title="Edit event">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete event">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:no-data>No Context Definitions Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/contexts'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      contexts: [],
      headers: [
        { text: 'Channel', value: 'channel' },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    create() {
      this.$router.push({ path: '/context' })
    },
    edit(contextSpec) {
      this.$router.push({ path: `/context/${contextSpec.id}` })
    },
    remove(contextSpec) {
      if (confirm('Do you really want to delete this context?')) {
        client.delete(contextSpec.id)

        let idx = this.contexts.findIndex(es => es.id === contextSpec.id)
        this.contexts.splice(idx, 1)
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
      this.contexts = await client.list()
    } catch (_) {
      _
    }
  }
}
</script>
