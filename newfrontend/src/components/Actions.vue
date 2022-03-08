<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Action overview</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="actions" />
            <v-btn icon @click="create()">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="actions">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="analyze(item)">
                <v-icon title="Analyze action">mdi-chart-bar</v-icon>
              </v-btn>
              <v-btn icon @click="toggleState(item)">
                <v-icon title="Change action state">mdi-play-pause</v-icon>
              </v-btn>
              <v-btn icon @click="edit(item)">
                <v-icon title="Edit action">mdi-pencil</v-icon>
              </v-btn>
              <v-btn icon @click="remove(item)" v-if="!item.active">
                <v-icon title="Delete action">mdi-delete</v-icon>
              </v-btn>
              <v-btn icon @click="duplicate(item)">
                <v-icon title="Duplicate action">mdi-content-copy</v-icon>
              </v-btn>
            </template>
            <template v-slot:item.state="{ item }">
              <span v-if="item.state === 'inactive'"> <v-icon size="medium">mdi-pause</v-icon>Inactive </span>
              <span v-if="item.state === 'staging'"> <v-icon size="medium">mdi-test-tube</v-icon>Staging </span>
              <span v-if="item.state === 'active'"> <v-icon size="medium">mdi-run</v-icon>Active </span>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:item.tags="{ item }">
              <v-chip v-for="(view, index) in item.tags" v-bind:key="index" class="mr-1" small>{{ view }}</v-chip>
            </template>
            <template v-slot:item.targeting.channels="{ item }">
              <span>{{ item.targeting.channels.join(', ') }}</span>
            </template>
            <template v-slot:no-data>No Actions Available</template>
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
      actions: [],
      dialog: false,
      infoDialog: false,
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'State', value: 'state' },
        { text: 'Channels', value: 'targeting.channels' },
        { text: 'Tags', value: 'tags' },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ]
    }
  },
  methods: {
    create() {
      this.$router.push({ path: '/action' })
    },
    analyze(action) {
      this.$router.push({ path: `/analyze/${action.id}` })
    },
    edit(action) {
      this.$router.push({ path: `/action/${action.id}` })
    },
    async remove(action) {
      if (confirm('Do you really want to delete this action?')) {
        await client.delete(action.id)
        this.actions = await client.list()
      }
    },
    async duplicate(action) {
      if (confirm('Do you really want to duplicate this action?')) {
        await client.duplicate(action.id)
        this.actions = await client.list()
      }
    },
    toggleState(action) {
      if (action.state === 'inactive') {
        action.state = 'staging'
      } else if (action.state === 'staging') {
        action.state = 'active'
      } else if (action.state === 'active') {
        action.state = 'inactive'
      }
      client.updateState({
        id: action.id,
        state: action.state
      })
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
