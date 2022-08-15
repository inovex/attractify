<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Action template overview</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="actiontemplates" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="actiontemplates">
            <template v-slot:item.actiontemplate="{ item }">
              <v-btn icon @click="edit(item)">
                <v-icon title="Edit action">mdi-pencil</v-icon>
              </v-btn>
              <v-btn icon @click="remove(item)" v-if="!item.active">
                <v-icon title="Delete action">mdi-delete</v-icon>
              </v-btn>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:no-data>No Templates Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/actions.js'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      actiontemplates: [],
      dialog: false,
      infoDialog: false,
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ]
    }
  },
  methods: {
    create() {
      this.$router.push({ path: '/actiontemplate' })
    },
    edit(action) {
      this.$router.push({ path: `/actiontemplate/${action.id}` })
    },
    async remove(action) { // TODO: check if action template is in use and cannot be deleted
      if (confirm('Do you really want to delete this action?')) {
        await client.delete(action.id)
        this.actions = await client.list()
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
    this.actions = await client.list()
  }
}
</script>
