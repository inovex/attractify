<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Action Types</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="actionTypes" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="actionTypes">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="edit(item)">
                <v-icon title="Edit type">mdi-pencil</v-icon>
              </v-btn>
              <v-btn icon @click="archive(item)">
                <v-icon title="Archive type">mdi-archive</v-icon>
              </v-btn>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.version="{ item }">
              <span>{{ item.version }}</span>
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
import client from '../lib/rest/actionTypes.js'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      actionTypes: [],
      dialog: false,
      infoDialog: false,
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Version', value: 'version' },
        { text: 'Created', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ]
    }
  },
  methods: {
    create() {
      this.$router.push({ path: '/action-type' })
    },
    edit(actiontype) {
      this.$router.push({ path: `/action-type/${actiontype.id}` })
    },
    async archive(action) {
      if (confirm('Do you really want to archive this type?')) {
        await client.delete(action.name)
        this.actionTypes = await client.list()
        this.cleanupTypeList()
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    cleanupTypeList() {
      let lastTypeName = ''
      for (let i = this.actionTypes.length - 1; i >= 0; i -= 1) {
        let item = this.actionTypes[i]
        if (lastTypeName == item.name || item.isArchived) {
          this.actionTypes.splice(i, 1)
        }
        lastTypeName = item.name
      }
    }
  },
  async created() {
    this.actionTypes = await client.list()
    this.cleanupTypeList()
  }
}
</script>
