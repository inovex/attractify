<template>
  <div>
    <v-card-title>
      Hooks
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" icon>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item v-for="(item, index) in types" :key="index" @click="addReaction(item.value)">
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-card-title>
    <v-card-subtitle>
      If you want to define what happens after a user has interacted with an action, you can define an hook.
    </v-card-subtitle>
    <v-card outlined v-for="(value, index) in hooks" v-bind:key="index" class="mb-4">
      <v-card-text v-if="value.type === 'execute_webhook'">
        <h4>Webhook</h4>
        <v-row>
          <v-col class="col-lg-3">
            <v-select
              :items="channels"
              multiple
              label="Channels"
              prepend-icon="mdi-bell"
              @change="changes"
              v-model="hooks[index].channels"
              clearable
            ></v-select>
          </v-col>
          <v-col class="col-lg-3">
            <v-select
              :items="events"
              label="Subscribed event"
              prepend-icon="mdi-bell"
              @change="changes"
              v-model="hooks[index].event"
              :rules="[rules.required]"
            ></v-select>
          </v-col>
          <v-col>
            <v-text-field
              label="Webhook URL"
              type="text"
              v-model="hooks[index].properties.url"
              @input="changes"
              prepend-icon="mdi-link"
              :rules="[rules.required]"
              clearable
            />
          </v-col>
          <v-col class="col-lg-2">
            <v-switch toggle label="Return result" v-model="hooks[index].properties.returnResult" />
          </v-col>
          <v-col class="col-lg-1">
            <v-btn icon @click="removeReaction(index)">
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-textarea
              label="Additional properties as JSON"
              type="text"
              v-model="hooks[index].properties.params"
              prepend-icon="mdi-code-json"
              @input="changes"
              :rules="[rules.json]"
              clearable
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-text v-if="value.type === 'track_event'">
        <h4>Track Event</h4>
        <v-row>
          <v-col class="col-lg-3">
            <v-select
              :items="channels"
              multiple
              label="Channels"
              prepend-icon="mdi-bell"
              @change="changes"
              v-model="hooks[index].channels"
              clearable
            ></v-select>
          </v-col>
          <v-col class="col-lg-3">
            <v-select
              :items="events"
              label="Subscribed event"
              prepend-icon="mdi-bell"
              @change="changes"
              v-model="hooks[index].event"
              :rules="[rules.required]"
            ></v-select>
          </v-col>
          <v-col>
            <APISelect
              label="Event Name"
              :loadCallback="loadEvents"
              v-model="hooks[index].properties.eventId"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col class="col-lg-1">
            <v-btn icon @click="removeReaction(index)">
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import client from '../../lib/rest/events.js'
import channelClient from '../../lib/rest/channels.js'
import APISelect from '../common/APISelect.vue'

export default {
  components: { APISelect },
  props: ['hooks'],
  data() {
    return {
      events: [
        { text: 'Shown', value: 'shown' },
        { text: 'Hidden', value: 'hidden' },
        { text: 'Declined', value: 'declined' },
        { text: 'Accepted', value: 'accepted' }
      ],
      types: [
        { title: 'Webhook', value: 'execute_webhook' },
        { title: 'Track Event', value: 'track_event' }
      ],
      channels: [],
      rules: {
        required: value => !!value || 'Required.',
        json: value => {
          try {
            if (!value) {
              return true
            }

            JSON.parse(value)
            return true
          } catch (e) {
            return 'Invalid JSON.'
          }
        }
      }
    }
  },
  methods: {
    addReaction(type) {
      this.hooks.push({ type: type, properties: {} })
    },
    removeReaction(index) {
      this.hooks.splice(index, 1)
    },
    loadEvents() {
      return client.listEventNames()
    },
    changes(){
      this.$emit('change')
    }
  },
  async created() {
    this.channels = await channelClient.select()
  }
}
</script>
