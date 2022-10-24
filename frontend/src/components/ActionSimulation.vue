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
                    required
                    hide-no-data
                    item-text="name"
                    label="User ID"
                    return-object
                  ></v-autocomplete>
                </v-col>
              </v-row>
              <!-- are Channel, Type, Tags important? -->
              <v-row>
                <v-col class="col-6">
                  <v-textarea
                    :rules="[rules.json]"
                    v-model="user.customProperties"
                    multi-line
                    label="Custom Properties"
                    required
                  ></v-textarea>
                </v-col>
                <v-col class="col-6">
                  <v-textarea
                    :rules="[rules.json]"
                    v-model="user.computedProperties"
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
                    :disabled="!valid"
                    @click="startSimulation()"
                    >Check Actions</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <v-expansion-panels focusable>
                <v-expansion-panel v-for="(action, index) in computedActions" v-bind:key="index">
                  <v-expansion-panel-header disable-icon-rotate>
                    {{ action.name }}
                    <template v-slot:actions>
                      <v-icon :color="listIcons[action.state].color"> {{ listIcons[action.state].icon }} </v-icon>
                    </template>
                  </v-expansion-panel-header>

                  <v-expansion-panel-content>
                    <div v-for="(step, index) in action.steps" :key="step + index">
                      {{ step.name }} => {{ step.userValue }} compared to {{ step.dataValue }} <br />
                      Error: {{ step.error }}
                    </div>
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
import actionsClient from '../lib/rest/actions'
import profilesClient from '../lib/rest/profiles'
import Help from './Help'

export default {
  components: { Help },
  data() {
    return {
      errorMessages: '',
      selectedProfile: null,
      userSearch: '',
      user: {},
      listIcons: {
        error: {
          icon: 'mdi-alert-circle',
          color: 'error'
        },
        correct: {
          icon: 'mdi-check',
          color: 'teal'
        }
      },
      computedActions: [
        {
          name: 'Action 1',
          state: 'correct',
          steps: [
            {
              name: 'Check 1',
              userValue: '10',
              dataValue: '20',
              error: ''
            },
            {
              name: 'Check 2',
              userValue: '10',
              dataValue: '20',
              error: ''
            }
          ]
        },
        {
          name: 'Action 2',
          state: 'error',
          steps: [
            {
              name: 'Check 1',
              userValue: '10',
              dataValue: '20',
              error: 'Value of user not high enough'
            },
            {
              name: 'Check 2',
              userValue: '10',
              dataValue: '20',
              error: 'Value of user not high enough'
            }
          ]
        }
      ],
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
    // TODO: simulation, when profile is selected load traits for that profile into json field
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
