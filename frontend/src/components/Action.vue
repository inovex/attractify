<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title v-if="action.id">Edit Action</v-toolbar-title>
              <v-toolbar-title v-if="!action.id">Create new Action</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="action" />
            </v-toolbar>
            <v-card-text class="raised">
              <v-card-title>Name & Type</v-card-title>
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="Name of Action"
                    name="name"
                    prepend-icon="mdi-pencil"
                    type="text"
                    @input="changes = true"
                    v-model="action.name"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    dense
                    prepend-icon="mdi-tune"
                    :items="actionTypeSelector"
                    label="Type of Action"
                    @change="selectType"
                    v-model="action.type"
                    :rules="[rules.required]"
                  ></v-select>
                </v-col>
                <v-col class="col-lg-3">
                  <v-select
                    dense
                    prepend-icon="mdi-tune"
                    :items="versionSelector"
                    label="Version of Type"
                    @change="selectVersion"
                    v-model="action.version"
                    :rules="[rules.required]"
                  ></v-select>
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text class="raised">
              <v-row>
                <v-col class="col-lg-3">
                  <v-text-field
                    label="Tag"
                    prepend-icon="mdi-tag"
                    type="text"
                    ref="tag"
                    v-model="tag"
                    @input="changes = true"
                    append-icon="mdi-plus"
                    @click:append="
                      action.tags.push(tag.toLowerCase())
                      tag = ''
                    "
                    @keyup.enter.prevent="
                      action.tags.push(tag.toLowerCase())
                      tag = ''
                    "
                  />
                </v-col>
                <v-col>
                  <v-chip
                    v-for="(view, index) in action.tags"
                    v-bind:key="index"
                    class="ma-2"
                    close
                    @click:close="action.tags.splice(index, 1)"
                    >{{ view }}</v-chip
                  >
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text class="raised">
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    label="State"
                    :items="state"
                    v-model="action.state"
                    prepend-icon="mdi-play-pause"
                    @change="changes = true"
                  ></v-select>
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <v-tabs v-model="tabs" centered>
                <v-tab href="#properties">Properties</v-tab>
                <v-tab href="#targeting">Targeting</v-tab>
                <v-tab href="#capping">Capping</v-tab>
                <v-tab href="#hooks">Hooks</v-tab>
                <v-tab href="#testUsers">Testusers</v-tab>
              </v-tabs>

              <v-tabs-items v-model="tabs">
                <v-tab-item value="properties">
                  <Properties
                    :properties="action.properties"
                    :typeProperties="action.typeProperties"
                    @change="changes = true"
                  />
                </v-tab-item>

                <v-tab-item value="targeting">
                  <Targeting :targeting="action.targeting" @change="changes = true" />
                </v-tab-item>

                <v-tab-item value="capping">
                  <Capping :capping="action.capping" @change="changes = true" />
                </v-tab-item>

                <v-tab-item value="hooks">
                  <Hooks :hooks="action.hooks" @change="changes = true" />
                </v-tab-item>

                <v-tab-item value="testUsers">
                  <TestUsers :test-users="action.testUsers" @change="changes = true" />
                </v-tab-item>
              </v-tabs-items>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>

    <v-col class="sticky text-center">
      <v-spacer />
      <v-btn rounded elevation="2" @click="cancel()">Cancel</v-btn>
      <v-btn
        rounded
        elevation="2"
        color="primary"
        style="color: var(--v-buttontext-base)"
        :disabled="!valid"
        @click="save()"
        >Save</v-btn
      >
    </v-col>

    <v-dialog v-model="exitUnsaved" max-width="700px" closeable>
      <UnsavedContent :valid="valid" @cancel="cancelExit" @save="save" @exit="exit" />
    </v-dialog>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex'
import Properties from './action/Properties.vue'
import Targeting from './action/Targeting.vue'
import Capping from './action/Capping.vue'
import Hooks from './action/Hooks.vue'
import TestUsers from './action/TestUsers.vue'
import UnsavedContent from './UnsavedContent.vue'
import actionTypesClient from '../lib/rest/actionTypes'
import Help from './Help'

export default {
  components: { Properties, Targeting, Capping, Hooks, TestUsers, Help, UnsavedContent },
  data() {
    return {
      tabs: '',
      action: {
        state: 'inactive',
        type: '',
        version: 1,
        tags: [],
        properties: [],
        typeProperties: [],
        targeting: {
          channels: [],
          start: {},
          end: {},
          traitConditions: [],
          contextConditions: []
        },
        capping: [],
        hooks: [],
        testUsers: []
      },
      state: [
        { text: 'Inactive', value: 'inactive' },
        { text: 'Staging', value: 'staging' },
        { text: 'Active', value: 'active' }
      ],
      tag: '',
      path: '',
      valid: false,
      changes: false,
      exitUnsaved: false,
      exitUrl: null,
      rules: {
        required: (value) => !!value || 'Required.'
      },
      actionTypes: [],
      actionTypeSelector: [],
      versionSelector: []
    }
  },
  methods: {
    cancel() {
      this.$router.push('/actions')
    },
    ...mapActions('actions', ['show', 'create', 'update']),
    async save() {
      try {
        let res = null
        if (this.action.id) {
          res = await this.update(this.action)
        } else {
          res = await this.create(this.action)
        }

        if (res && res.id) {
          this.action.id = res.id
        }

        this.$notify.success('Your action has been saved.')
        if (this.exitUnsaved) {
          this.exit()
        }
      } catch (e) {
        this.$notify.error('Could not save action.')
      }
    },
    cancelExit() {
      this.exitUnsaved = false
      this.exitUrl = null
    },
    exit() {
      this.changes = false
      if (this.exitUrl) {
        this.$router.push(this.exitUrl)
      } else {
        this.$router.push('/actions')
      }
    },
    selectType() {
      this.changes = true

      this.versionSelector = []
      this.actionTypes.forEach((actionType) => {
        if (actionType.name == this.action.type) {
          this.versionSelector.push(actionType.version)
        }
      })

      this.action.version = this.versionSelector[this.versionSelector.length - 1]
      this.selectVersion()
    },
    selectVersion() {
      let currentVersion
      this.actionTypes.every((actionType) => {
        if (this.action.version == actionType.version && this.action.type == actionType.name) {
          currentVersion = actionType
          return false
        }
        return true
      })

      this.action.typeProperties = currentVersion.properties
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.action = await this.show(id)
        delete this.action.trigger
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }

    actionTypesClient.list().then((actionTypes) => {
      this.actionTypes = actionTypes
      this.actionTypes.forEach((actionType) => {
        this.actionTypeSelector.push(actionType.name)
      })
    })
  },
  beforeRouteLeave(to, from, next) {
    if (this.changes) {
      this.exitUnsaved = true
      this.exitUrl = to.path
      return false
    }
    next()
  }
}
</script>

<style scoped>
.raised {
  background-color: rgba(0, 0, 0, 0.05);
}
.sticky {
  margin: 0.5rem 0 0 0;
  position: -webkit-sticky;
  position: sticky;
  bottom: 1rem;
  z-index: 1;
}
.sticky button {
  margin: 0 0.5rem;
}
</style>
