<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title v-if="actiontype.id">Edit Action Template</v-toolbar-title>
              <v-toolbar-title v-if="!actiontype.id">Create new Action Template</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="action" />
            </v-toolbar>
            <v-card-text class="raised">
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    :disabled="actiontype.version != 1"
                    label="Type of Action"
                    name="type"
                    prepend-icon="mdi-tune"
                    type="text"
                    @input="changes = true"
                    v-model="actiontype.name"
                    :rules="[rules.required]"
                  />
                </v-col>
                <v-col align-self="center" class="col-lg-6">
                  <v-row align="center" style="height: 100%; top: auto">
                    <v-icon>mdi-timeline-clock-outline</v-icon>
                    <v-card-text style="width: auto; font-size: 16px">Version: {{ actiontype.version }}</v-card-text>
                  </v-row>
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <PropertiesTemplates :properties="actiontype.properties" @change="changes = true" />
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
import UnsavedContent from './UnsavedContent.vue'
import Help from './Help'
import PropertiesTemplates from './action/PropertiesTemplates.vue'
import actionClient from '../lib/rest/actions'

export default {
  components: { Help, UnsavedContent, PropertiesTemplates },
  data() {
    return {
      actiontype: {
        properties: [],
        name: '',
        version: 1
      },
      path: '',
      valid: false,
      changes: false,
      exitUnsaved: false,
      exitUrl: null,
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    cancel() {
      this.$router.push('/actiontypes')
    },
    ...mapActions('actiontypes', ['show', 'create', 'update']),
    async save() {
      try {
        let res = null
        if (this.actiontype.id) {
          res = await this.update(this.actiontype)
        } else {
          res = await this.create(this.actiontype)
        }

        if (res && res.id) {
          this.actiontype.id = res.id
        }
        this.changes = false
        this.$notify.success('Your actiontype has been saved as version ' + this.actiontype.version + '.')
        if (this.exitUnsaved) {
          this.exit()
        }
      } catch (e) {
        this.$notify.error('Could not save actiontype.')
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
        this.$router.push('/actiontypes')
      }
    },
    inUse() {
      let actions = actionClient.list() // TODO: only load actions with actiontypes name
      for (let action in actions) {
        if (action.name == this.actiontype.name && action.version == this.actiontype.version) {
          return true
        }
      }
      return false
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.actiontype = await this.show(id)
        if (this.inUse()) {
          this.actiontype.version++
        }
        delete this.actiontype.trigger
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
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
