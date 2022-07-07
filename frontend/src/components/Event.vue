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
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" text :disabled="!valid" @click="save()">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
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
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
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
  }
}
</script>
