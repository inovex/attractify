<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
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
                <v-col class="col-6">
                  <v-autocomplete
                    v-model="selectedProfile"
                    :items="foundProfiles"
                    :loading="isLoading"
                    :search-input.sync="userSearch"
                    prepend-icon="mdi-account"
                    hide-no-data
                    item-text="name"
                    label="User ID"
                    required
                    return-object
                  ></v-autocomplete>
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-6">
                  <v-textarea
                    :rules="[rules.json]"
                    v-model="user.customTraits"
                    multi-line
                    label="Custom Properties"
                    required
                  ></v-textarea>
                </v-col>
                <v-col class="col-6">
                  <v-textarea
                    :rules="[rules.json]"
                    v-model="user.computedTraits"
                    label="Computed Properties"
                    required
                  ></v-textarea>
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-6">
                  <v-textarea :rules="[rules.json]" v-model="user.context" label="Context" required></v-textarea>
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-3">
                  <v-btn
                    rounded
                    color="primary"
                    style="color: var(--v-buttontext-base)"
                    :disabled="!valid || !selectedProfile"
                    @click="startSimulation()"
                    >Check Actions</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <v-expansion-panels focusable>
                <v-expansion-panel v-for="(action, index) in computedActions" :key="index">
                  <v-expansion-panel-header disable-icon-rotate>
                    {{ action.name }}
                    <template v-slot:actions>
                      <v-icon :color="listIcons[action.display].color"> {{ listIcons[action.display].icon }} </v-icon>
                    </template>
                  </v-expansion-panel-header>

                  <v-expansion-panel-content style="padding-top: 1rem">
                    <v-btn
                      rounded
                      color="primary"
                      style="color: var(--v-buttontext-base); margin-bottom: 1rem"
                      @click="gotoAction(action.id)"
                      >Goto Action</v-btn
                    >
                    <v-expansion-panels focusable>
                      <v-expansion-panel v-for="(step, index) in action.steps" :key="'test' + index">
                        <v-expansion-panel-header disable-icon-rotate>
                          {{ step.name }}
                          <template v-slot:actions>
                            <v-icon :color="listIcons[!step.blocking].color">
                              {{ listIcons[!step.blocking].icon }}
                            </v-icon>
                          </template>
                        </v-expansion-panel-header>

                        <v-expansion-panel-content> Info: {{ step.info }} </v-expansion-panel-content>
                      </v-expansion-panel>
                    </v-expansion-panels>
                  </v-expansion-panel-content>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>
  
<script>
import actionClient from '../lib/rest/actions'
import profilesClient from '../lib/rest/profiles'
import Help from './Help'

export default {
  components: { Help },
  data() {
    return {
      errorMessages: '',
      selectedProfile: null,
      userSearch: '',
      user: {
        customTraits: '{}',
        computedTraits: '{}',
        context: '{}',
        channel: ''
      },
      listIcons: {
        false: {
          icon: 'mdi-alert-circle',
          color: 'error'
        },
        true: {
          icon: 'mdi-check',
          color: 'teal'
        }
      },
      computedActions: [],
      actionId: '',
      foundProfiles: [],
      isLoading: false,
      valid: false,
      rules: {
        json: (value) => {
          try {
            JSON.parse(value)
          } catch (e) {
            return 'This textfield must contain a valid json object.'
          }
          return true
        }
      }
    }
  },
  methods: {
    async startSimulation() {
      let params = {
        userId: this.user.userId,
        customTraits: JSON.parse(this.user.customTraits),
        computedTraits: JSON.parse(this.user.computedTraits),
        context: JSON.parse(this.user.context),
        channel: this.user.channel
      }
      await actionClient.simulate(params).then((res) => {
        this.computedActions = res
      })
    },
    gotoAction(id) {
      let route = this.$router.resolve({ path: '/action/' + id })
      window.open(route.href)
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
          setTimeout(() => {
            this.isLoading = false
          }, 200)
        })
    },
    foundProfiles() {
      if (this.foundProfiles.length === 1) {
        this.selectedProfile = this.foundProfiles[0]
      }
    },
    selectedProfile(profile) {
      this.user.customTraits = JSON.stringify(profile.value.customTraits)
      this.user.computedTraits = JSON.stringify(profile.value.computedTraits)
      this.user.channel = profile.value.channel
      this.user.userId = profile.value.userId
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.userSearch = id
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
  }
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
