<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>User Profiles</v-toolbar-title>

            <v-spacer></v-spacer>
            <help name="profiles" />
            <v-btn icon @click="load()">
              <v-icon>mdi-refresh</v-icon>
            </v-btn>
          </v-toolbar>
          <v-card-text>
            <v-row>
              <v-col class="col-lg-6">
                <v-text-field
                  v-model="search.userId"
                  label="User ID"
                  persistent-hint
                  prepend-icon="mdi-account-search"
                  @keydown.enter="searchProfile"
                  clearable
                ></v-text-field>
              </v-col>
              <v-col class="col-lg-6" style="display: flex; align-items: center">
                <v-btn icon @click="searchProfile()"> <v-icon title="Search">mdi-magnify</v-icon> </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="profiles">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="show(item)"> <v-icon title="Show profile">mdi-magnify</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="refreshComputedTraits(item)">
                <v-icon title="Refresh computed traits">mdi-account-convert</v-icon> </v-btn
              >&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete profile">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:no-data>No User Profiles Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/profiles.js'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      profiles: [],
      headers: [
        {
          text: 'ID',
          align: 'left',
          value: 'id'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      search: {}
    }
  },
  methods: {
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    timeAgo(date) {
      return moment(date).fromNow()
    },
    show(profile) {
      this.$router.push({ path: `/profile/${profile.id}` })
    },
    async refreshComputedTraits(profile) {
      try {
        await client.refreshComputedTraits(profile.id)
        this.$notify.success('Computed traits have been refreshed.')
      } catch (_) {
        this.$notify.error('Could not refresh computed traits.')
      }
    },
    async remove(profile) {
      if (confirm('Do you really want to delete this profile and all its events and identities?')) {
        try {
          await client.delete(profile.id)
          this.profiles = await client.list()
          this.$notify.success('Profile has been deleted.')
        } catch (_) {
          this.$notify.error('Could not delete profile.')
        }
      }
    },
    async load() {
      try {
        this.profiles = await client.list()
      } catch (_) {
        _
      }
    },
    async searchProfile() {
      try {
        this.profiles = await client.search(this.search.userId)
      } catch (_) {
        this.profiles = {}
      }
    }
  },
  async created() {
    this.load()
  }
}
</script>
