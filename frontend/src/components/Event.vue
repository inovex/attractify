<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Event Definition</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="event" />
            </v-toolbar>
            <v-card-text>
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="Name"
                    name="name"
                    prepend-icon="mdi-pencil"
                    type="text"
                    v-model="event.name"
                    :rules="[rules.required]"
                    :v-on:keyUp="changes=true"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-10">
                  <v-text-field
                    label="Description"
                    name="type"
                    prepend-icon="mdi-text"
                    type="text"
                    v-model="event.description"
                    :v-on:keyUp="changes=true"
                  />
                </v-col>
              </v-row>
            </v-card-text>
            <div>
              <v-card-title>Properties</v-card-title>
              <v-card-subtitle>
                Here you can define all allowed event properties.
              </v-card-subtitle>
              <v-card-text>
                <Structure :structure="event.structure" v-on:savecallback="onUpdate"/>
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
  </v-container>
</template>

<script>
import Help from './Help'
import Structure from './events/Structure.vue'
import client from '../lib/rest/events'

export default {
  components: { Structure, Help },
  data() {
    return {
      event: {
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
      changes: false,
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    cancel(){
      this.$router.push('/events')
    },
    addCondition() {
      this.event.conditions.push({})
    },
    deleteCondition(index) {
      this.event.conditions.splice(index, 1)
    },
    onUpdate(){
      if(this.valid){
        this.save()
      }
    },
    async save() {
      try {
        let res = null
        if (this.event.id) {
          res = await client.update(this.event)
        } else {
          res = await client.create(this.event)
        }

        if (res && res.id) {
          this.event.id = res.id
        }

        this.$notify.success('The structure has been saved successfully.')
        this.changes = false
      } catch (e) {
        this.$notify.error('Could not save structure.')
      }
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.event = await client.show(id)
      } catch (error) {
        this.$router.push({ path: '/404' })
      }
    }
  },
  beforeRouteLeave(to, from, next) {
    if(this.changes){
      var leave = confirm('You have unsaved changes. Are you sure you want to leave this page?');
      if(!leave){
        return false
      }
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
