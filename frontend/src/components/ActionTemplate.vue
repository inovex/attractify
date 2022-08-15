<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title v-if="id">Edit Action Template</v-toolbar-title>
              <v-toolbar-title v-if="!id">Create new Action Template</v-toolbar-title>
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
              </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <Properties :properties="actiontemplate.properties" @change="changes=true" />
          </v-card>
        </v-form>
      </v-col>
    </v-row>

    <v-col class="sticky text-center">
      <v-spacer />
      <v-btn rounded elevation="2" @click="cancel()">Cancel</v-btn>
      <v-btn rounded elevation="2" color="primary" style="color: var(--v-buttontext-base)" :disabled="!valid" @click="save()">Save</v-btn>
    </v-col>

    <v-dialog v-model="exitUnsaved" max-width="700px" closeable>
      <UnsavedContent :valid="valid" @cancel="cancelExit" @save="save" @exit="exit"/>
    </v-dialog>
  </v-container>
</template>

<script>
import Properties from './action/Properties.vue'
import UnsavedContent from './UnsavedContent.vue'
import Help from './Help'

export default {
  components: { Help, UnsavedContent, Properties },
  data() {
    return {
      actiontemplate: {
        properties: [],
      },
      path: '',
      valid: false,
      changes: false,
      exitUnsaved: false,
      exitUrl: null,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    cancel(){
      this.$router.push('/actiontemplates')
    },
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
        if(this.exitUnsaved){
          this.exit()
        }
      } catch (e) {
        this.$notify.error('Could not save action.')
      }
    },
    cancelExit(){
      this.exitUnsaved = false
      this.exitUrl = null
    },
    exit(){
      this.changes = false
      if(this.exitUrl){
        this.$router.push(this.exitUrl)
      }else{
        this.$router.push('/actions')
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
  },
  beforeRouteLeave(to, from, next) {
    if(this.changes){
      this.exitUnsaved = true
      this.exitUrl = to.path
      return false
    }
    next()
  },
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
.sticky button{
  margin: 0 0.5rem;
}
</style>
