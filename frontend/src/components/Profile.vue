<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>User Profile Details</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="profile" />
            <v-btn icon @click="load()">
              <v-icon>mdi-reload</v-icon>
            </v-btn>
          </v-toolbar>
          <v-card-text>
            <v-tabs v-model="tabs" centered>
              <v-tab href="#identities">Identities</v-tab>
              <!-- TODO: make identities deleteable -->
              <v-tab href="#customTraits">Custom Traits</v-tab>
              <!-- TODO: make them editable -->
              <v-tab href="#computedTraits">Computed Traits</v-tab>
              <!-- TODO: make them editable -->
              <v-tab href="#events">Latest Events</v-tab>
              <!-- TODO: make events deletable (all and single events) -->
              <!-- TODO: add reactions and make them deletable aswell -->
            </v-tabs>
            <!-- TODO: make whole profile deletable (events, reactions, identities) -->

            <v-tabs-items v-model="tabs">
              <v-tab-item value="identities">
                <Identities :profile="profile" />
              </v-tab-item>

              <v-tab-item value="customTraits">
                <Traits :traits="profile.customTraits" />
              </v-tab-item>

              <v-tab-item value="computedTraits">
                <Traits :traits="profile.computedTraits" />
              </v-tab-item>

              <v-tab-item value="events">
                <Events :profile="profile" />
              </v-tab-item>
            </v-tabs-items>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import profiles from '../lib/rest/profiles'
import moment from 'moment'

import Traits from './profiles/Traits.vue'
import Identities from './profiles/Identities.vue'
import Events from './profiles/Events.vue'

export default {
  components: { Traits, Identities, Events, Help },
  data() {
    return {
      tabs: '',
      profile: {},
      identities: []
    }
  },
  methods: {
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    timeAgo(date) {
      return moment(date).fromNow()
    },
    async load() {
      const id = this.$route.params.id
      if (id) {
        try {
          this.profile = await profiles.show(id)
        } catch (error) {
          this.$router.push({ path: '/404' })
        }
      }
    }
  },
  async created() {
    this.load()
  }
}
</script>
