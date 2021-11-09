<template>
  <v-container>
    <v-row>
      <v-col>
        <v-form ref="form" v-model="valid" on>
          <v-card>
            <v-toolbar dark>
              <v-toolbar-title>New GDPR/DSGVO Request</v-toolbar-title>
              <v-spacer></v-spacer>
              <help name="privacy" />
            </v-toolbar>
            <v-card-text>
              <v-card-title>User</v-card-title>
              <v-row>
                <v-col class="col-lg-6">
                  <v-text-field
                    label="User ID"
                    name="name"
                    prepend-icon="mdi-account"
                    type="text"
                    v-model="userId"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="col-lg-3">
                  <v-select
                    label="Operation"
                    :items="operations"
                    v-model="operation"
                    prepend-icon="mdi-calendar-check"
                  ></v-select>
                </v-col>
              </v-row>
              <v-row v-if="operation === 'export'">
                <v-col class="col-lg-3">
                  <v-text-field
                    label="Email Address of User"
                    name="name"
                    prepend-icon="mdi-email"
                    type="text"
                    v-model="email"
                    :rules="[rules.required]"
                  />
                </v-col>
              </v-row>
            </v-card-text>

            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" :disabled="!valid" @click="run()" text>Create Request</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Locked Profile Identities</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialog" max-width="700px" closeable>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">Create locked identity</span>
                </v-card-title>
                <v-form ref="form" v-model="valid">
                  <v-card-text>
                    <v-row>
                      <v-col></v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-text-field label="User ID" prepend-icon="mdi-account" type="text" v-model="lockedUserId" />
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn text @click="close()">Close</v-btn>
                    <v-btn text color="primary" :disabled="!valid" @click="create()">Create</v-btn>
                  </v-card-actions>
                </v-form>
              </v-card>
            </v-dialog>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="lockedProfileIdentities">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="remove(item.id)"> <v-icon>mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:no-data>No Locked Identities Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import client from '../lib/rest/privacy.js'
import moment from 'moment'

export default {
  components: { Help },
  data() {
    return {
      operation: '',
      email: '',
      userId: '',
      dialog: false,
      lockedUserId: '',
      lockedProfileIdentities: [],
      headers: [
        {
          text: 'User ID',
          align: 'left',
          value: 'userId'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      operations: [
        { text: 'Export Data', value: 'export' },
        { text: 'Delete Data', value: 'deletion' }
      ],
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    async run() {
      if (this.operation === 'export') {
        try {
          await client.export({ userId: this.userId, email: this.email })
          this.$notify.success('Your export request has been finished.')
        } catch (_) {
          this.$notify.error('Your export request could not be processed.')
        }
      } else if (this.operation === 'deletion') {
        try {
          await client.deletion({ userId: this.userId })
          this.$notify.success('Your deletion request has been finished.')
        } catch (_) {
          this.$notify.error('Your deletion request could not be processed.')
        }
      }
    },
    async remove(id) {
      if (confirm('Do you really want to delete this locked identity?')) {
        await client.deleteLockedIdentities(id)
        this.list()
      }
    },
    async create() {
      try {
        await client.createLockedIdentities({ userId: this.lockedUserId })
        this.$notify.success('Identity lock successfully created.')
        this.list()
        this.close()
      } catch (_) {
        this.$notify.error('Could not create identity lock.')
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    close() {
      this.dialog = false
      this.user = {}
      this.$refs.form.reset()
    },
    async list() {
      try {
        this.lockedProfileIdentities = await client.lockedIdentities()
      } catch (_) {
        this.$notify.error('Could not load locked identities.')
      }
    }
  },
  created() {
    this.list()
  }
}
</script>
