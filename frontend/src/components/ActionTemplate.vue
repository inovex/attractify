<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title v-if="actiontemplate.id">Edit Action Template</v-toolbar-title>
              <v-toolbar-title v-if="!actiontemplate.id">Create new Action Template</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="action" />
            </v-toolbar>
            <v-card-text class="raised">
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="Type of Action"
                    name="type"
                    prepend-icon="mdi-tune"
                    type="text"
                    @input="changes = true"
                    v-model="actiontemplate.type"
                    :rules="[rules.required]"
                  />
                </v-col>
                <v-col align-self="center" class="col-lg-6">
                  <v-row align="center" style="height: 100%; top: auto">
                    <v-icon>mdi-timeline-clock-outline</v-icon>
                    <v-card-text style="width: auto; font-size: 16px"
                      >Version: {{ actiontemplate.version }}</v-card-text
                    >
                  </v-row>
                </v-col>
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <PropertiesTemplates :properties="actiontemplate.properties" @change="changes = true" />
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
      actiontemplate: {
        properties: [],
        type: '',
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
      this.$router.push('/actiontemplates')
    },
    ...mapActions('actiontemplates', ['show', 'create', 'update']),
    async save() {
      try {
        if (this.inUse()) {
          this.actiontemplate.version++
        }
        let res = null
        if (this.actiontemplate.id) {
          res = await this.update(this.actiontemplate)
        } else {
          res = await this.create(this.actiontemplate)
        }

        if (res && res.id) {
          this.actiontemplate.id = res.id
        }

        this.$notify.success('Your actiontype has been saved as version ' + this.actiontemplate.version + '.')
        if (this.exitUnsaved) {
          this.exit()
        }
      } catch (e) {
        this.$notify.error('Could not save actiontemplate.')
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
        this.$router.push('/actiontemplates')
      }
    },
    inUse() {
      let actions = actionClient.list()
      for (let action in actions) {
        if (action.type == this.type) {
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
        this.actiontemplate = await this.show(id)
        delete this.actiontemplate.trigger
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
