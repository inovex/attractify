<template>
  <div>
    <v-card-title>
      Properties
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" icon>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item v-for="(item, index) in types" :key="index" @click="addProperty(item.value)">
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-card-title>
    <v-card-subtitle>
      Action Properties are used to assign values to names and use them later in
      the actions.
    </v-card-subtitle>
    <v-card outlined v-for="(prop, index) in properties" v-bind:key="index" class="mb-4">
      <v-card-text>
        <h4 v-if="prop.type === 'text'">Text</h4>
        <h4 v-if="prop.type === 'custom_trait'">Custom Trait</h4>
        <h4 v-if="prop.type === 'computed_trait'">Computed Trait</h4>
        <v-row>
          <v-col class="col-lg-3">
            <v-select
              dense
              :items="channels"
              label="Channels"
              :value="prop.channels"
              v-model="properties[index].channels"
              multiple
            ></v-select>
          </v-col>
          <v-col class="col-lg-2">
            <v-text-field
              dense
              label="Name"
              type="text"
              v-model="properties[index].name"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col>
            <v-text-field
              dense
              v-if="prop.type === 'text'"
              label="Value"
              :type="prop.type"
              v-model="properties[index].value"
              :rules="[rules.required]"
            />

            <APISelect
              v-if="prop.type === 'custom_trait'"
              dense
              label="Trait Key"
              :loadCallback="listCustomTraitKeys"
              :value="prop.sourceKey || ''"
              @change="(e) => {setProperty(prop, e)}"
              return-object
            />

            <APISelect
              v-if="prop.type === 'computed_trait'"
              dense
              label="Trait Key"
              :loadCallback="listComputedTraitKeys"
              :value="prop.sourceKey || ''"
              @change="(e) => {setProperty(prop, e)}"
              return-object
            />
          </v-col>
          <v-col class="col-lg-1">
            <v-btn icon @click="removeProperty(index)">
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import channelClient from '../../lib/rest/channels.js'
import APISelect from '../common/APISelect.vue'
import customTraitsClient from '../../lib/rest/customTraits.js'
import computedTraitsClient from '../../lib/rest/computedTraits.js'

export default {
  components: { APISelect },
  props: ['properties'],
  data() {
    return {
      types: [
        { title: 'Text', value: 'text' },
        { title: 'Custom Trait', value: 'custom_trait' },
        { title: 'Computed Trait', value: 'computed_trait' }
      ],
      channels: [],
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    computedTraitType(type, sourceType) {
      switch (type) {
        case 'count_events':
          return 'integer'
        case 'unique_list_count':
          return 'integer'
        case 'aggregation':
          return 'float'
        case 'unique_list':
          return 'array'
        case 'first_event':
          return sourceType
        case 'last_event':
          return sourceType
        case 'most_frequent':
          return sourceType
      }

      return sourceType
    },
    setProperty(property, val) {
      property.sourceKey = val.value
      if (property.type === 'custom_trait') {
        property.sourceType = val.type
      } else if (property.type === 'computed_trait') {
        property.sourceType = this.computedTraitType(val.type, val.sourceType)
      }
    },
    addProperty(type) {
      this.properties.push({ type: type })
    },
    removeProperty(index) {
      this.properties.splice(index, 1)
    },
    listCustomTraitKeys() {
      return customTraitsClient.listProperties()
    },
    listComputedTraitKeys() {
      return computedTraitsClient.listTraits()
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
