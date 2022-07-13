<template>
  <div>
    <v-card-title>Target Channels</v-card-title>
    <v-card-subtitle>Define where the actions should be displayed.</v-card-subtitle>
    <v-card-text>
      <v-row>
        <v-col class="col-lg-4">
          <v-select
            label="Channels"
            :items="channels"
            v-model="targeting.channels"
            prepend-icon="mdi-cellphone"
            multiple
            clearable
          ></v-select>
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-title>Audiences</v-card-title>
    <v-card-subtitle>Define where the actions should be displayed.</v-card-subtitle>
    <v-card-text>
      <v-row>
        <v-col class="col-lg-4">
          <APISelect
            label="Audiences"
            icon="mdi-account-group"
            :loadCallback="loadAudiences"
            v-model="targeting.audiences"
            multiple
          />
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-title>Time Range</v-card-title>
    <v-card-subtitle>Set a start and end date as well as a start and/or end time.</v-card-subtitle>
    <v-card-text>
      <v-row>
        <v-col class="col-lg-2">
          <v-dialog
            ref="dialogStartDate"
            v-model="modalStartDate"
            :return-value.sync="targeting.start.date"
            persistent
            width="290px"
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="targeting.start.date"
                label="Start Date"
                prepend-icon="mdi-calendar"
                readonly
                clearable
                v-on="on"
              ></v-text-field>
            </template>
            <v-date-picker v-model="targeting.start.date" type="date" scrollable>
              <v-spacer></v-spacer>
              <v-btn rounded @click="modalStartDate = false">Cancel</v-btn>
              <v-btn rounded color="primary" @click="$refs.dialogStartDate.save(targeting.start.date)">OK</v-btn>
            </v-date-picker>
          </v-dialog>
        </v-col>
        <v-col class="col-lg-2">
          <v-dialog
            ref="dialogStartTime"
            v-model="modalStartTime"
            :return-value.sync="targeting.start.time"
            persistent
            width="290px"
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="targeting.start.time"
                label="Start Time"
                prepend-icon="mdi-clock"
                readonly
                clearable
                v-on="on"
              ></v-text-field>
            </template>
            <v-time-picker v-if="modalStartTime" v-model="targeting.start.time" full-width format="24hr">
              <v-spacer></v-spacer>
              <v-btn rounded @click="modalStartTime = false">Cancel</v-btn>
              <v-btn rounded color="primary" @click="$refs.dialogStartTime.save(targeting.start.time)">OK</v-btn>
            </v-time-picker>
          </v-dialog>
        </v-col>
        <v-col class="col-lg-2">
          <v-dialog
            ref="dialogEndDate"
            v-model="modalEndDate"
            :return-value.sync="targeting.end.date"
            persistent
            width="290px"
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="targeting.end.date"
                label="End Date"
                prepend-icon="mdi-calendar"
                readonly
                clearable
                v-on="on"
              ></v-text-field>
            </template>
            <v-date-picker v-model="targeting.end.date" type="date" scrollable>
              <v-spacer></v-spacer>
              <v-btn rounded @click="modalEndDate = false">Cancel</v-btn>
              <v-btn rounded color="primary" @click="$refs.dialogEndDate.save(targeting.end.date)">OK</v-btn>
            </v-date-picker>
          </v-dialog>
        </v-col>
        <v-col class="col-lg-2">
          <v-dialog
            ref="dialogEndTime"
            v-model="modalEndTime"
            :return-value.sync="targeting.end.time"
            persistent
            width="290px"
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="targeting.end.time"
                label="End Time"
                prepend-icon="mdi-clock"
                readonly
                clearable
                v-on="on"
              ></v-text-field>
            </template>
            <v-time-picker v-if="modalEndTime" v-model="targeting.end.time" full-width format="24hr">
              <v-spacer></v-spacer>
              <v-btn rounded @click="modalEndTime = false">Cancel</v-btn>
              <v-btn rounded color="primary" @click="$refs.dialogEndTime.save(targeting.end.time)">OK</v-btn>
            </v-time-picker>
          </v-dialog>
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-title>
      Trait Conditions
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" icon>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item v-for="(item, index) in traitTypes" :key="index" @click="addTraitCondition(item.value)">
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-card-title>
    <v-card-subtitle>
      If you want to filter actions based on their trait properties, you can specify these conditions here.
    </v-card-subtitle>
    <v-card-text>
      <v-row v-for="(condition, index) of targeting.traitConditions" :key="index">
        <v-col>
          <TraitCondition :condition="condition" />
        </v-col>
        <v-col class="col-lg-1">
          <v-btn icon @click="removeTraitCondition(index)">
            <v-icon>mdi-trash-can-outline</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-title>
      Context Conditions
      <v-btn icon @click="addContextCondition">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-subtitle>
      If you want to filter actions based on their context properties, you can specify these conditions here.
    </v-card-subtitle>
    <v-card-text>
      <v-row v-for="(condition, index) of targeting.contextConditions" :key="index">
        <v-col>
          <ContextCondition :condition="condition" />
        </v-col>
        <v-col class="col-lg-1">
          <v-btn icon @click="removeContextCondition(index)">
            <v-icon>mdi-trash-can-outline</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </div>
</template>

<script>
import channelClient from '../../lib/rest/channels.js'
import client from '../../lib/rest/audiences.js'
import APISelect from '../common/APISelect.vue'
import TraitCondition from './TraitCondition.vue'
import ContextCondition from './ContextCondition.vue'

export default {
  components: { APISelect, TraitCondition, ContextCondition },
  props: ['targeting'],
  data() {
    return {
      channels: [],
      modalStartDate: false,
      modalStartTime: false,
      modalEndDate: false,
      modalEndTime: false,
      traitTypes: [
        { title: 'Custom Trait', value: 'custom' },
        { title: 'Computed Trait', value: 'computed' }
      ],
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    async loadAudiences() {
      const res = await client.list()
      return res.map(e => {
        return { text: `${e.name} (${e.profileCount} profiles)`, value: e.id }
      })
    },
    addContextCondition() {
      this.targeting.contextConditions.push({})
    },
    removeContextCondition(index) {
      this.targeting.contextConditions.splice(index, 1)
    },
    addTraitCondition(source) {
      this.targeting.traitConditions = this.targeting.traitConditions || []
      this.targeting.traitConditions.push({ source: source })
    },
    removeTraitCondition(index) {
      this.targeting.traitConditions.splice(index, 1)
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
