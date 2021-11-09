<template>
  <div>
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">
              <strong>Name</strong>
            </th>
            <th class="text-left">Channel</th>
            <th class="text-left">Context</th>
            <th class="text-left">Properties</th>
            <th class="text-left">Created At</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in events" :key="item.id">
            <td>
              <span>{{ resolveEventName(item.eventId) }}</span>
            </td>
            <td>{{ item.channel }}</td>
            <td>
              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-icon v-bind="attrs" v-on="on">mdi-magnify</v-icon>
                </template>
                <span
                  ><pre>{{ item.context }}</pre></span
                >
              </v-tooltip>
            </td>
            <td>
              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-icon v-bind="attrs" v-on="on">mdi-magnify</v-icon>
                </template>
                <span
                  ><pre>{{ item.properties }}</pre></span
                >
              </v-tooltip>
            </td>
            <td>{{ item.createdAt }}</td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>
import profiles from '../../lib/rest/profiles'
import eventsClient from '../../lib/rest/events.js'

export default {
  props: ['profile'],
  data() {
    return {
      events: [],
      eventMap: {}
    }
  },
  methods: {
    async loadEvents() {
      let events = await eventsClient.listEventNames()
      for (let e of events) {
        this.eventMap[e.value] = e.text
      }
    },
    resolveEventName(eventId) {
      return this.eventMap[eventId]
    }
  },
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        await this.loadEvents()
        this.events = await profiles.listEvents(id)
      } catch (_) {
        _
      }
    }
  }
}
</script>
