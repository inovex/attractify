<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>Custom Traits Definition</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="customTraits" />
            </v-toolbar>
            <div>
              <v-card-title>Properties</v-card-title>
              <v-card-subtitle>Here you can define all allowed custom traits properties.</v-card-subtitle>
              <v-card-text>
                <Structure :structure="customTraits.structure" />
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
import client from '../lib/rest/customTraits'

export default {
  components: { Structure, Help },
  data() {
    return {
      customTraits: {
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
      this.customTraits.conditions.push({})
    },
    deleteCondition(index) {
      this.customTraits.conditions.splice(index, 1)
    },
    async save() {
      try {
        await client.upsert(this.customTraits)
        this.$notify.success('The structure has been saved successfully.')
      } catch (e) {
        this.$notify.error('Could not save structure.')
      }
    }
  },
  async created() {
    try {
      this.customTraits = await client.show()
    } catch (error) {
      this.customTraits = {
        structure: []
      }
    }
  }
}
</script>
