<template>
  <div>
    <v-form ref="form" v-model="valid" on>
      <v-card-title>Simulation</v-card-title>
      <v-card-subtitle>Check whether your action behaves how you want it to.</v-card-subtitle>
      <v-card-text>
        <v-row>
          <v-col class="col-lg-4">
            <v-select
              label="User"
              :loading="isLoading"
              :items="testUserList"
              v-model="selectedProfile"
              prepend-icon="mdi-account"
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-text>
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
              :disabled="!valid || !user.userId"
              @click="startSimulation()"
              >Check Action</v-btn
            >
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-text>
        <div style="text-align: center" v-if="computedAction">
          <v-icon :size="40" :color="listIcons[computedAction.display ? 'correct' : 'error'].color">
            {{ listIcons[computedAction.display ? 'correct' : 'error'].icon }}
          </v-icon>
        </div>
        <v-expansion-panels v-if="computedAction" focusable>
          <v-expansion-panel v-for="(step, index) in computedAction.steps" :key="'test' + index">
            <v-expansion-panel-header disable-icon-rotate>
              {{ step.name }}
              <template v-slot:actions>
                <v-icon :color="listIcons[step.state].color">
                  {{ listIcons[step.state].icon }}
                </v-icon>
              </template>
            </v-expansion-panel-header>

            <v-expansion-panel-content>
              <v-card-text>Info: {{ step.info }} </v-card-text>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card-text>
    </v-form>
  </div>
</template>
  
<script>
import actionClient from '../../lib/rest/actions'
import profilesClient from '../../lib/rest/profiles'

export default {
  components: {},
  props: ['action'],
  data() {
    return {
      listIcons: {
        error: {
          icon: 'mdi-alert-circle',
          color: 'error'
        },
        correct: {
          icon: 'mdi-check',
          color: 'teal'
        },
        server_error: {
          icon: 'mdi-alert-circle',
          color: 'orange'
        }
      },
      testUserList: [],
      selectedProfile: null,
      user: {
        customTraits: '{}',
        computedTraits: '{}',
        context: '{}',
        channel: '',
        userId: ''
      },
      computedAction: null,
      valid: false,
      isLoading: false,
      rules: {
        required: (value) => !!value || 'Required.',
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
        channel: this.user.channel,
        actionId: this.action.id
      }
      await actionClient.simulate(params).then((res) => {
        this.computedAction = res
      })
    }
  },
  watch: {
    selectedProfile(profile) {
      this.user.customTraits = JSON.stringify(profile.customTraits)
      this.user.computedTraits = JSON.stringify(profile.computedTraits)
      this.user.channel = profile.channel
      this.user.userId = profile.userId
    }
  },
  async created() {
    this.action.testUsers.forEach((testUser) => {
      this.isLoading = true
      profilesClient
        .searchWithTraits(testUser.userId)
        .then((res) => {
          res.map((e) => {
            this.testUserList.push({ text: e.userId + ' - ' + testUser.description, value: e })
          })
        })
        .finally(() => {
          this.isLoading = false
        })
    })
  }
}
</script>
  