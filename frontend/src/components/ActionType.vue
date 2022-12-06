<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title v-if="actionType.id">Edit Action Type</v-toolbar-title>
              <v-toolbar-title v-if="!actionType.id">Create new Action Type</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="actionType" />
            </v-toolbar>
            <v-card-text class="raised">
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    :disabled="exists"
                    label="Type of Action"
                    name="type"
                    prepend-icon="mdi-tune"
                    type="text"
                    @input="changes = true"
                    v-model="actionType.name"
                    :rules="[rules.required]"
                  />
                </v-col>
                <v-col align-self="center" class="col-lg-6">
                  <v-row align="center" style="height: 100%; top: auto">
                    <v-icon>mdi-timeline-clock-outline</v-icon>
                    <v-card-text style="width: auto; font-size: 16px">Version: {{ actionType.version }}</v-card-text>
                  </v-row>
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <Properties :typeProperties="null" :properties="actionType.properties" @change="changes = true" />
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
        :disabled="!valid || !changes"
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
import Properties from './action/Properties.vue'
import actionTypeClient from '../lib/rest/actionTypes'

export default {
  components: { Help, UnsavedContent, Properties },
  data() {
    return {
      actionType: {
        properties: [],
        name: '',
        version: 1,
        isInUse: false
      },
      path: '',
      exists: false,
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
      this.$router.push('/action-types')
    },
    ...mapActions('actionTypes', ['show', 'create', 'update']),
    async save() {
      try {
        let newRoute = false
        if (this.actionType.id && !this.actionType.isInUse) {
          await this.update(this.actionType)
        } else {
          let exists = await actionTypeClient.nameExists(this.actionType.name)
          if (
            !this.actionType.id &&
            exists &&
            !confirm(
              'This type already exists. Do you want to create a new version of this type and take older versions out of the archive?'
            )
          ) {
            return
          }
          await this.create(this.actionType).then((res) => {
            if (res && res.id) {
              if (this.actionType.id) {
                this.actionType.version = this.actionType.version + 1
              } else {
                newRoute = true
              }
              this.actionType.id = res.id
            }
          })
        }

        this.changes = false
        this.$notify.success('Your action type has been saved.')

        if (this.exitUnsaved) {
          this.exit()
        }
        if (newRoute) {
          this.$router.push(`/action-type/${this.actionType.id}`)
        }
      } catch (e) {
        this.$notify.error('Could not save action type.')
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
        this.$router.push('/action-types')
      }
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      this.exists = true
      try {
        this.actionType = await this.show(id)
        delete this.actionType.trigger
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
