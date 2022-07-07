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
import client from '../lib/rest/contexts'
import channelClient from '../lib/rest/channels.js'

export default {
  components: { Structure, Help },
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
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
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
      } catch (e) {
        this.$notify.error('Could not save structure.')
      }
    }
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
  }
}
</script>
