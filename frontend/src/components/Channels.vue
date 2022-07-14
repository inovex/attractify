<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Channels</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="channels" />
            <v-dialog v-model="dialog" max-width="700px" closeable>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">Edit channel</span>
                </v-card-title>
                <v-form ref="form" v-model="valid">
                  <v-card-text>
                    <v-row>
                      <v-col>
                        <v-text-field label="Name" prepend-icon="mdi-text" type="text" v-model="channel.name" />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-text-field
                          :disabled="!!this.channel.id"
                          label="Key"
                          prepend-icon="mdi-code"
                          type="text"
                          v-model="channel.key"
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn rounded @click="close()">Close</v-btn>
                    <v-btn rounded color="primary" style="color: var(--v-buttontext-base)" :disabled="!valid" @click="save()">Save</v-btn>
                  </v-card-actions>
                </v-form>
              </v-card>
            </v-dialog>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="channels">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="edit(item)"> <v-icon title="Edit channel">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete channel">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:no-data>No Channels Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import client from '../lib/rest/channels.js'
import Help from './Help'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      channels: [],
      dialog: false,
      channel: {},
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        {
          text: 'Key',
          value: 'key'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  watch: {
    dialog(val) {
      val || this.close()
    }
  },
  methods: {
    async save() {
      try {
        let res = null
        if (this.channel.id) {
          res = await client.update(this.channel)
        } else {
          res = await client.create(this.channel)
          this.load()
        }

        if (res && res.id) {
          this.channel.id = res.id
        }

        this.$notify.success('Channel successfully saved.')
        this.close()
      } catch (e) {
        this.$notify.error('Could not save channel.')
      }
    },
    edit(channel) {
      this.channel = channel
      this.dialog = true
    },
    async remove(channel) {
      if (confirm('Do you really want to delete this channel?')) {
        await client.delete(channel.id)
        this.load()
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    timeAgo(date) {
      return moment(date).fromNow()
    },
    close() {
      this.dialog = false
      this.channel = {}
      this.$refs.form.reset()
    },
    async load() {
      this.channels = await client.list()
    }
  },
  async created() {
    this.load()
  }
}
</script>
