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
              <v-card-title>User</v-card-title>
              <v-row>
                <v-col class="col-lg-6">
                  <v-autocomplete
                    v-model="selectedProfile"
                    :items="foundProfiles"
                    :loading="isLoading"
                    :search-input.sync="userSearch"
                    prepend-icon="mdi-account"
                    hide-no-data
                    item-text="name"
                    label="User ID"
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
                    >Check Actions</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <div class="grids">
                <v-card class="pa-2" outlined elevated="5">
                  Action 1 <br />
                  Status: Displayed <br />
                  Errors: None <br />
                  <v-btn>Details</v-btn>
                </v-card>
                <v-card class="pa-2" outlined elevated="5">
                  Action 2 <br />
                  Status: Displayed <br />
                  Errors: None <br />
                  <v-btn>Details</v-btn>
                </v-card>
                <v-card class="pa-2" outlined elevated="5">
                  Action 3 <br />
                  Status: Displayed <br />
                  Errors: None <br />
                  <v-btn>Details</v-btn>
                </v-card>
                <v-card class="pa-2" outlined elevated="5">
                  Action 4 <br />
                  Status: Displayed <br />
                  Errors: None <br />
                  <v-btn>Details</v-btn>
                </v-card>
                <v-card class="pa-2" outlined elevated="5">
                  Action 5 <br />
                  Status: Displayed <br />
                  Errors: None <br />
                  <v-btn>Details</v-btn>
                </v-card>
                <v-card class="pa-2" style="border-color: red" outlined elevated="5">
                  Action 6 <br />
                  Status: Not Displayed <br />
                  Errors: Capping <br />
                  <v-btn>Details</v-btn>
                </v-card>
              </div>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>
  
<script>
import actionsClient from '../lib/rest/actions'
import profilesClient from '../lib/rest/profiles'
import Help from './Help'

export default {
  components: { Help },
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
            this.foundProfiles.push({ name: e.userId, value: e })
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

<style>
div.grids {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
}
</style>
