<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Audiences</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="audiences" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="audiences">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="refresh(item)"> <v-icon title="Refresh audience">mdi-account-convert</v-icon> </v-btn
              >&nbsp; <v-btn icon @click="edit(item)"> <v-icon title="Edit audience">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete audience">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.refreshedAt="{ item }">
              <span>{{ timeAgo(item.refreshedAt) }}</span>
            </template>
            <template v-slot:no-data>No Audiences are available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/audiences'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      audiences: [],
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Description', value: 'description' },
        { text: 'Matched profiles', value: 'profileCount' },
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
    async load() {
      try {
        this.audiences = await client.list()
      } catch (_) {
        _
      }
    },
    create() {
      this.$router.push({ path: '/audience' })
    },
    edit(audience) {
      this.$router.push({ path: `/audience/${audience.id}` })
    },
    async refresh(audience) {
      try {
        let res = await client.refresh(audience.id)
        this.$notify.success(`The audience has been successfully refreshed with ${res.count} profiles.`)
        this.load()
      } catch (_) {
        this.$notify.error('Could not refresh audience.')
      }
    },
    remove(audience) {
      if (confirm('Do you really want to delete this audience?')) {
        client.delete(audience.id)

        let idx = this.audiences.findIndex(es => es.id === audience.id)
        this.audiences.splice(idx, 1)
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
    this.load()
  }
}
</script>
