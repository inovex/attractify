<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>API Access</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="api" />
          </v-toolbar>
          <v-card-title>Auth tokens for client integrations</v-card-title>
          <v-card-text>
            <v-row>
              <v-col class="col-lg-3">
                <v-select :items="channels" label="Channel" prepend-icon="mdi-devices" v-model="channel"></v-select>
              </v-col>
              <v-col>
                <v-btn @click="createToken" :disabled="channel === ''" text color="primary">Create Token</v-btn>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-data-table disable-pagination hide-default-footer :headers="headers" :items="authTokens">
                  <template v-slot:item.action="{ item }">
                    <v-btn icon @click="deleteToken(item.id)"> <v-icon>mdi-delete</v-icon> </v-btn>&nbsp;
                  </template>
                  <template v-slot:item.token="{ item }">
                    <pre>{{ item.token }}</pre>
                  </template>
                  <template v-slot:item.createdAt="{ item }">
                    <span>{{ formatDate(item.createdAt) }}</span>
                  </template>
                  <template v-slot:no-data>No auth tokens available</template>
                </v-data-table>
              </v-col>
            </v-row>
          </v-card-text>
          <v-divider></v-divider>
          <v-card-title>Webhook signature key</v-card-title>
          <v-card-text>
            <v-row>
              <v-col>
                Show the key used for creating signatures (displayed as hex, please use as binary).
              </v-col>
            </v-row>
            <v-row>
              <v-col class="col-lg-3">
                <v-text-field
                  label="Password to unlock"
                  prepend-icon="mdi-lock"
                  type="password"
                  v-model="password"
                  :rules="[rules.required]"
                />
              </v-col>
              <v-col>
                <v-btn @click="showKey()" text color="primary">Show Signature Key</v-btn>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-card v-if="key">
                  <v-card-text>
                    <pre wrap="wrap">{{ key }}</pre>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
          <v-divider></v-divider>
          <v-card-title>Platform API access</v-card-title>
          <v-card-text>
            <v-row>
              <v-col>
                Here you can generate a token that can be used for the SDK or client.
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-btn @click="getToken()" color="primary" text>Generate Token</v-btn>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-card v-if="jwt">
                  <v-card-text>
                    <pre wrap="wrap">{{ jwt }}</pre>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>

          <!-- <v-divider></v-divider>
          <v-card-text>
            <v-form ref="form" v-model="valid" on>
              <v-row>
                <v-col>Here you can simulate Webhook Requests.</v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-select
                    label="Event"
                    :items="events"
                    v-model="webhook.event"
                    prepend-icon="mdi-bell"
                    :rules="[rules.required]"
                    clearable
                  ></v-select>
                </v-col>
                <v-col>
                  <v-select
                    label="Select campaign to test webhook"
                    :items="campaigns"
                    v-model="webhook.campaignId"
                    prepend-icon="mdi-ticket-percent"
                    :rules="[rules.required]"
                    clearable
                  ></v-select>
                </v-col>
                <v-col>
                  <v-select
                    label="Channel"
                    :items="channels"
                    v-model="webhook.channel"
                    prepend-icon="mdi-cellphone"
                    :rules="[rules.required]"
                    clearable
                  ></v-select>
                </v-col>
                <v-col>
                  <v-text-field
                    label="User ID"
                    v-model="userId"
                    prepend-icon="mdi-numeric"
                    :rules="[rules.required]"
                    clearable
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-textarea
                    label="User properties as JSON"
                    type="text"
                    v-model="properties"
                    prepend-icon="mdi-code-json"
                    :rules="[rules.json]"
                    clearable
                  />
                </v-col>
              </v-row>
              <v-row v-if="webhookResult">
                <v-col>
                  HTTP Status code:
                  <strong>
                    <pre>{{webhookResult.statusCode}}</pre>
                  </strong>
                </v-col>
                <v-col>
                  Signature:
                  <pre>{{webhookResult.signature}}</pre>
                </v-col>
              </v-row>
              <v-row v-if="webhookResult">
                <v-col>
                  HTTP Response Body:
                  <pre>{{webhookResult.body}}</pre>
                </v-col>
              </v-row>
              <v-row v-if="webhookError">
                <v-col>
                  <v-alert type="error">Could not trigger webhook: {{webhookResult.error}}</v-alert>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-btn @click="triggerWebhook()" :disabled="!valid">Test webhook</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card-text>-->
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex'
import moment from 'moment'
import channelClient from '../lib/rest/channels.js'
import tokenClient from '../lib/rest/authTokens.js'
import Help from './Help'

export default {
  components: { Help },
  data() {
    return {
      jwt: '',
      key: '',
      password: '',
      type: '',
      userId: '',
      properties: '',
      webhook: {},
      valid: false,
      webhookResult: null,
      webhookError: false,
      authTokens: [],
      headers: [
        {
          text: 'Channel',
          align: 'left',
          value: 'channel'
        },
        {
          text: 'Token',
          value: 'token'
        },
        { text: 'Created', value: 'createdAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      events: [
        { text: 'close', value: 'close' },
        { text: 'dismiss', value: 'dismiss' },
        { text: 'accept', value: 'accept' }
      ],
      channels: [],
      channel: '',
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
    ...mapActions('organization', { generateToken: 'token', getKey: 'key' }),
    async getToken() {
      const res = await this.generateToken()
      this.jwt = res.token
    },
    async showKey() {
      try {
        const res = await this.getKey(this.password)
        this.key = res.key
      } catch (error) {
        this.$notify.error('Wrong password.')
      }
    },
    async triggerWebhook() {
      this.webhookError = false
      this.webhook.properties = JSON.parse(this.properties)
      this.webhook.userId = this.userId
      try {
        const res = await this.testWebhook(this.webhook)
        this.webhookResult = res
      } catch (error) {
        this.webhookError = true
      }
    },
    formatDate(date) {
      return moment(date).format('YYYY-MM-DD, HH:mm')
    },
    async createToken() {
      try {
        await tokenClient.create({ channel: this.channel })
        this.$notify.success('Token successfully created.')
      } catch (_) {
        _
      }

      this.loadTokens()
    },
    async deleteToken(id) {
      if (confirm('Do you really want to delete this token?')) {
        await tokenClient.delete(id)
        this.loadTokens()
      }
    },
    async loadTokens() {
      this.authTokens = await tokenClient.list()
    }
  },
  async created() {
    this.channels = await channelClient.select()
    this.loadTokens()
  }
}
</script>
