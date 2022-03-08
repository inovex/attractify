<template>
  <div>
    <v-card-title>
      Testusers
      <v-btn icon @click="addTestuser">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-subtitle
      >Here you can add testusers that do not have to be in any audience or can skip targeting at all.</v-card-subtitle
    >
    <v-card-text>
      <v-card outlined v-for="(testUser, index) in testUsers" v-bind:key="index" class="mb-4">
        <v-card-text>
          <v-row>
            <v-col class="col-lg-2">
              <v-select
                :items="channels"
                label="Channel"
                :value="testUser.channels"
                v-model="testUsers[index].channels"
                multiple
              ></v-select>
            </v-col>
            <v-col class="col-lg-3">
              <v-text-field label="User ID" type="text" v-model="testUsers[index].userId" :rules="[rules.required]" />
            </v-col>
            <v-col class="col-lg-4">
              <v-text-field label="Description" type="text" v-model="testUsers[index].description" />
            </v-col>
            <v-col class="col-lg-2">
              <v-switch toggle label="Skip targeting" v-model="testUsers[index].skipTargeting" />
            </v-col>
            <v-col class="col-lg-1">
              <v-btn icon @click="removeTestuser(index)">
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
  props: ['testUsers'],
  data() {
    return {
      channels: [],
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    addTestuser() {
      this.testUsers.push({})
    },
    removeTestuser(index) {
      this.testUsers.splice(index, 1)
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
