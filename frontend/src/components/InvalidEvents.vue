<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Invalid Events</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="invalidevents" />
            <v-btn @click="load()" icon>
              <v-icon>mdi-refresh</v-icon>
            </v-btn>
          </v-toolbar>

          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="eventSpecs">
            <template v-slot:item.action="{ item }">
              <v-btn @click="details(item)" icon>
                <v-icon title="Show event">mdi-magnify</v-icon>
              </v-btn>
              <v-btn icon @click="remove(item)"> <v-icon title="Delete event">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>

            <template v-slot:item.type="{ item }">
              <span v-if="item.type === 'properties'">Property</span>
              <span v-if="item.type === 'context'">Context</span>
            </template>

            <template v-slot:item.name="{ item }">
              <span>{{ item.name }}</span>
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:no-data>No Invalid Events Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="dialog" max-width="700px" closeable>
      <v-card>
        <v-card-title>
          <span class="headline">Event details</span>
        </v-card-title>
        <v-card-text>
          <h4>{{ detailView.displayName }}</h4>
          <v-card outlined class="pa-2">
            <pre style="overflow: auto; display: grid">
              <h4 v-if="detailView.displaySchema.valid">Valid</h4>
              <p v-if="detailView.displaySchema.valid" style="color: green">{{ detailView.displaySchema.valid }}</p>
              <h4 v-if="detailView.displaySchema.invalid">Wrong value</h4>
              <p v-if="detailView.displaySchema.invalid" style="color: red">{{ detailView.displaySchema.invalid }}</p>
              <h4 v-if="detailView.displaySchema.notSet">Missing</h4>
              <p v-if="detailView.displaySchema.notSet" style="color: grey">{{ detailView.displaySchema.notSet }}</p>
            </pre>
          </v-card>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import Help from './Help'
import invalidEventClient from '../lib/rest/invalidEvents'
import eventClient from '../lib/rest/events'
import contextClient from '../lib/rest/contexts'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      dialog: false,
      event: {},
      detailView: {
        displayName: '',
        displaySchema: {
          valid: {},
          invalid: {},
          notSet: {}
        }
      },
      detailErrors: {},
      eventSpecs: [],
      headers: [
        { text: 'Source', value: 'type' },
        {
          text: 'Event name',
          align: 'left',
          value: 'name'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      valid: false,
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    remove(eventSpec) {
      invalidEventClient.delete(eventSpec.id)

      let idx = this.eventSpecs.findIndex((es) => es.id === eventSpec.id)
      this.eventSpecs.splice(idx, 1)
    },
    details(event) {
      this.event = event

      const json = event.type === 'properties' ? event.properties : event.context

      if (event.type === 'properties') {
        try {
          const properties = eventClient.list()
          properties.then((properties) => {
            const schema = properties.find((e) => e.organizationID === event.organizationID && e.name === event.name)
            const schemaJSON = this.getJSONFromArray(schema.structure)
            var result = this.getValidateJSON(json, schemaJSON, {})

            result.invalid = this.addInvalidJSON(json, result.notSet, result.valid, schemaJSON)

            this.detailView = {
              displayName: 'Properties',
              displaySchema: result
            }
          })
        } catch (e) {
          this.$notify.error('Event for this error does not exist anymore.')
        }
      } else {
        try {
          const contexts = contextClient.list()
          contexts.then((contexts) => {
            const schema = contexts.find(
              (c) => c.organizationID === event.organizationID && c.channel === event.channel
            ).structure
            const schemaJSON = this.getJSONFromArray(schema)
            var result = this.getValidateJSON(json, schemaJSON, {})
            result.invalid = this.addInvalidJSON(json, result.notSet, result.valid, schemaJSON)

            this.detailView = {
              displayName: 'Context',
              displaySchema: result
            }
          })
        } catch (e) {
          this.$notify.error('Context could not be found.')
        }
      }
      this.dialog = true
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    timeAgo(date) {
      return moment(date).fromNow()
    },
    async load() {
      try {
        this.eventSpecs = await invalidEventClient.list()
      } catch (_) {
        _
      }
    },
    getJSONFromArray(array) {
      var json = {}
      array.forEach((element) => {
        if (element.properties.type === 'object') {
          json[element.name] = this.getJSONFromArray(element.children)
        } else {
          json[element.name] = element.properties.type
        }
      })
      return json
    },
    getValidateJSON(json, schema, result) {
      for (let elem in schema) {
        if (!json[elem]) {
          if (!result['notSet']) result['notSet'] = {}
          result['notSet'][elem] = schema[elem]
          continue
        }

        if (typeof schema[elem] === 'string') {
          const jsonType = typeof json[elem]
          if (jsonType.replace('number', 'integer') !== schema[elem]) {
            continue
          }
          if (!result['valid']) result['valid'] = {}
          result['valid'][elem] = json[elem]
          continue
        }

        var recursive = this.getValidateJSON(json[elem], schema[elem], {})
        for (let rec in recursive) {
          if (!result[rec]) result[rec] = {}
          result[rec][elem] = recursive[rec]
        }
      }
      return result
    },
    addInvalidJSON(json, notSet, valid, schema) {
      var result
      for (let elem in json) {
        if (typeof json[elem] === 'object') {
          var recursive = this.addInvalidJSON(
            json[elem],
            notSet ? notSet[elem] : 'empty',
            valid ? valid[elem] : 'empty',
            schema ? schema[elem] : 'empty'
          )

          for (let rec in recursive) {
            if (!result) result = {}
            if (!result[elem]) result[elem] = {}
            result[elem][rec] = recursive[rec]
          }
          continue
        }
        if ((!notSet || !notSet[elem]) && (!valid || !valid[elem])) {
          if (!result) result = {}
          if (schema !== 'empty' && schema[elem]) {
            result[elem] = (typeof json[elem]).replace('number', 'integer') + ' (correct type: ' + schema[elem] + ')'
            continue
          }
          result[elem] = json[elem]
        }
      }
      return result
    }
  },
  async created() {
    this.load()
  }
}
</script>
