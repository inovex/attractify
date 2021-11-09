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
                    v-model="action.name"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="Type of Action"
                    name="type"
                    prepend-icon="mdi-tune"
                    type="text"
                    v-model="action.type"
                    :rules="[rules.required]"
                  />
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
                  <Properties :properties="action.properties" />
                </v-tab-item>

                <v-tab-item value="targeting">
                  <Targeting :targeting="action.targeting" />
                </v-tab-item>

                <v-tab-item value="capping">
                  <Capping :capping="action.capping" />
                </v-tab-item>

                <v-tab-item value="hooks">
                  <Hooks :hooks="action.hooks" />
                </v-tab-item>

                <v-tab-item value="testUsers">
                  <TestUsers :test-users="action.testUsers" />
                </v-tab-item>
              </v-tabs-items>
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" :disabled="!valid" @click="save()" text>Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex'
import Properties from './action/Properties.vue'
import Targeting from './action/Targeting.vue'
import Capping from './action/Capping.vue'
import Hooks from './action/Hooks.vue'
import TestUsers from './action/TestUsers.vue'
import Help from './Help'

export default {
  components: { Properties, Targeting, Capping, Hooks, TestUsers, Help },
  data() {
    return {
      tabs: '',
      action: {
        state: 'inactive',
        tags: [],
        properties: [],
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
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
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
      } catch (e) {
        this.$notify.error('Could not save action.')
      }
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
  }
}
</script>

<style scoped>
.raised {
  background-color: rgba(0, 0, 0, 0.05);
}
</style>
