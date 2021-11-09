<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Computed Traits</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="computedTraits" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="computedTraits">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="edit(item)"> <v-icon title="Edit computed trait">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete computed trait">mdi-delete</v-icon> </v-btn
              >&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.refreshedAt="{ item }">
              <span>{{ timeAgo(item.refreshedAt) }}</span>
            </template>
            <template v-slot:no-data>No computed traits are available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/computedTraits'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      computedTraits: [],
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        {
          text: 'Key',
          value: 'key'
        },
        {
          text: 'Type',
          value: 'type'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Refreshed', value: 'refreshedAt' },
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
      this.$router.push({ path: '/computed-trait' })
    },
    edit(computedTrait) {
      this.$router.push({ path: `/computed-trait/${computedTrait.id}` })
    },
    remove(computedTrait) {
      if (confirm('Do you really want to delete this computed trait?')) {
        client.delete(computedTrait.id)

        let idx = this.computedTraits.findIndex(es => es.id === computedTrait.id)
        this.computedTraits.splice(idx, 1)
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
      this.computedTraits = await client.list()
    } catch (_) {
      _
    }
  }
}
</script>
