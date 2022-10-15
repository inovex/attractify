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
                <v-col class="col-5">
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
                <v-col class="col-lg-3">
                  <v-btn
                    rounded
                    color="primary"
                    style="color: var(--v-buttontext-base)"
                    :disabled="selectedProfile == null"
                    @click="startSimulation()"
                    >Check Actions</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <div class="grids">
                <v-card outlined v-for="index in 10" :key="index">
                  <v-toolbar dark>
                    <v-toolbar-title>Action {{ index }}</v-toolbar-title>
                  </v-toolbar>
                  <br />
                  <v-row>
                    <v-column class="col-5">
                      <v-card-text>
                        <p style="font-size: 16px">Status:</p>
                      </v-card-text>
                    </v-column>
                    <v-column class="col-5">
                      <v-chip class="ma-2" color="green">Displayed</v-chip>
                    </v-column>
                    <v-column></v-column>
                  </v-row>
                  <v-card-text>
                    <p class="text--primary">Errors: None</p>
                  </v-card-text>
                  <v-card-actions>
                    <v-btn>Details</v-btn>
                  </v-card-actions>
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
  row-gap: 2rem;
}

div.grids > div {
  width: 30%;
}

.error-box {
  border: 1px solid red !important;
}
</style>
