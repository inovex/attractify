<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Action Simulation</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="actionsimulation" />
              <!-- TODO: write help page -->
            </v-toolbar>
            <v-card-text class="raised">
              <v-card-title>Action & User</v-card-title>
              <v-row>
                <v-col class="col-lg-6">
                  <APISelect
                    label="Action"
                    icon="mdi-ticket-percent"
                    :loadCallback="loadActions"
                    v-model="actionId"
                    :valid="true"
                  />
                </v-col>
                <v-col class="col-lg-6">
                  <!-- TODO: Find a better way to handle showing the profile ID when searching for a user_id or identity_id-->
                  <v-autocomplete
                    v-model="selectedProfile"
                    :items="foundProfiles"
                    :loading="isLoading"
                    :search-input.sync="userSearch"
                    :filter="
                      () => {
                        return true
                      }
                    "
                    hide-no-data
                    item-text="name"
                    label="Profile ID"
                    return-object
                  ></v-autocomplete>
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-6">
                  <v-btn
                    rounded
                    color="primary"
                    style="color: var(--v-buttontext-base)"
                    :disabled="actionId == '' || selectedProfile == null"
                    @click="startSimulation()"
                    >Start Simulation</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>
  
<script>
import APISelect from './common/APISelect.vue'
import actionsClient from '../lib/rest/actions'
import profilesClient from '../lib/rest/profiles'
import Help from './Help'

export default {
  components: { APISelect, Help },
  data() {
    return {
      selectedProfile: null,
      userSearch: '',
      user: {},
      actionId: '',
      foundProfiles: [],
      isLoading: false
    }
  },
  methods: {
    // TODO: simulation,
    async loadActions() {
      const res = await actionsClient.list()
      return res.map((e) => {
        return { text: e.name, value: e.id }
      })
    },
    startSimulation() {
      console.log('TODO: start simulation')
    }
  },
  watch: {
    async userSearch(u) {
      if (this.isLoading || u == null || u == '' || this.selectedProfile != null) return

      this.isLoading = true
      this.foundProfiles = []

      profilesClient
        .search(u)
        .then((res) => {
          res.map((e) => {
            this.foundProfiles.push({ name: e.id, value: e })
          })
        })
        .finally(() => {
          this.isLoading = false
        })
    }
  },
  async created() {}
}
</script>
