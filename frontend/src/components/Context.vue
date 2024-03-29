<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Context Definition</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="context" />
            </v-toolbar>
            <v-card-text>
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    label="Channel"
                    :items="channels"
                    v-model="context.channel"
                    prepend-icon="mdi-cellphone"
                    @change="changes=true"
                  ></v-select>
                </v-col>
              </v-row>
            </v-card-text>
            <div>
              <v-card-title>Properties</v-card-title>
              <v-card-subtitle>
                Here you can define all allowed context properties.
              </v-card-subtitle>
              <v-card-text>
                <Structure :structure="context.structure" @savecallback="onUpdate"/>
              </v-card-text>
            </div>
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
import Help from './Help'
import Structure from './events/Structure.vue'
import UnsavedContent from './UnsavedContent.vue'
import client from '../lib/rest/contexts'
import channelClient from '../lib/rest/channels.js'

export default {
  components: { Structure, Help, UnsavedContent },
  data() {
    return {
      channels: [],
      context: {
        structure: []
      },
      types: [
        { text: 'String', value: 'string' },
        { text: 'Integer', value: 'integer' },
        { text: 'Float', value: 'float' },
        { text: 'Bool', value: 'bool' },
        { text: 'DateTime', value: 'datetime' },
        { text: 'Array', value: 'array' },
        { text: 'Object', value: 'object' }
      ],
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
      this.$router.push('/contexts')
    },
    addCondition() {
      this.context.conditions.push({})
    },
    deleteCondition(index) {
      this.context.conditions.splice(index, 1)
    },
    onUpdate(){
      if(this.valid){
        this.save()
      }
    },
    async save() {
      try {
        let res = null
        if (this.context.id) {
          res = await client.update(this.context)
        } else {
          res = await client.create(this.context)
        }

        if (res && res.id) {
          this.context.id = res.id
        }

        this.$notify.success('The structure has been saved successfully.')
        this.changes = false
        if(this.exitUnsaved){
          this.exit()
        }
      } catch (e) {
        this.$notify.error('Could not save structure.')
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
        this.$router.push('/contexts')
      }
    },
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.context = await client.show(id)
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
    this.channels = await channelClient.select()
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
