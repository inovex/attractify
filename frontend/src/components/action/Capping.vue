<template>
  <div>
    <v-card-title>
      Capping Conditions
      <v-btn icon @click="addCondition">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-subtitle>Here you can define different capping conditions.</v-card-subtitle>
    <v-card-text>
      <v-card outlined v-for="(slot, index) in capping" v-bind:key="index" class="mb-4">
        <v-card-text>
          <v-row>
            <v-col class="col-lg-3">
              <v-select :items="channels" label="Channel" v-model="capping[index].channels" multiple></v-select>
            </v-col>
            <v-col class="col-lg-2">
              <v-select :items="events" label="Event" v-model="capping[index].event"></v-select>
            </v-col>
            <v-col class="col-lg-2">
              <v-select
                :items="group"
                label="Target group"
                v-model="capping[index].group"
                :rules="[rules.required]"
              ></v-select>
            </v-col>
            <v-col class="col-lg-2">
              <v-text-field
                label="Per User Show Count"
                prepend-icon="mdi-number"
                type="number"
                clearable
                v-model.number="capping[index].count"
                :rules="[rules.required]"
              />
            </v-col>
            <v-col class="col-lg-2">
              <v-text-field
                label="Days within"
                prepend-icon="mdi-calendar"
                type="number"
                clearable
                v-model.number="capping[index].within"
              />
            </v-col>
            <v-col class="col-lg-1">
              <v-btn icon @click="removeCondition(index)">
                <v-icon>mdi-trash-can-outline</v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-card-text>
  </div>
</template>

<script>
import channelClient from '../../lib/rest/channels.js'

export default {
  props: ['capping'],
  data() {
    return {
      channels: [],
      events: [
        { text: 'Shown', value: 'shown' },
        { text: 'Hidden', value: 'hidden' },
        { text: 'Declined', value: 'declined' },
        { text: 'Accepted', value: 'accepted' }
      ],
      group: [
        { text: 'All users', value: 'all' },
        { text: 'Current User', value: 'user' }
      ],
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    addCondition() {
      this.capping.push({})
    },
    removeCondition(index) {
      this.capping.splice(index, 1)
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
