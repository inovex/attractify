<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Users</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="users" />
            <v-dialog v-model="dialog" max-width="700px" closeable>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">Edit user {{ user.name }}</span>
                </v-card-title>
                <v-form ref="form" v-model="valid">
                  <v-card-text>
                    <v-row>
                      <v-col></v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-text-field
                          label="Email"
                          prepend-icon="mdi-email"
                          type="text"
                          v-model="user.email"
                          :disabled="userLoaded"
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-select
                          label="Role"
                          :items="roles"
                          v-model="user.role"
                          prepend-icon="mdi-ticket-percent"
                          :rules="[rules.required]"
                        ></v-select>
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn text @click="close()">Close</v-btn>
                    <v-btn text color="primary" :disabled="!valid" @click="save()">Save</v-btn>
                  </v-card-actions>
                </v-form>
              </v-card>
            </v-dialog>
          </v-toolbar>
          <v-data-table disable-pagination hide-default-footer :headers="headers" :items="allusers">
            <template v-slot:item.action="{ item }">
              <v-btn icon @click="resendInvitiation(item)">
                <v-icon title="Resend invitation email">mdi-email-sync</v-icon> </v-btn
              >&nbsp; <v-btn icon @click="edit(item)"> <v-icon title="Edit user">mdi-pencil</v-icon> </v-btn>&nbsp;
              <v-btn icon @click="remove(item)"> <v-icon title="Delete user">mdi-delete</v-icon> </v-btn>&nbsp;
            </template>
            <template v-slot:item.createdAt="{ item }">
              <span>{{ formatDate(item.createdAt) }}</span>
            </template>
            <template v-slot:item.updatedAt="{ item }">
              <span>{{ timeAgo(item.updatedAt) }}</span>
            </template>
            <template v-slot:no-data>No Users Available</template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import { mapActions, mapGetters } from 'vuex'
import moment from 'moment'
import Vue from 'vue'

export default {
  components: { Help },
  data() {
    return {
      dialog: false,
      user: {},
      headers: [
        {
          text: 'Name',
          align: 'left',
          value: 'name'
        },
        { text: 'Email', value: 'email' },
        { text: 'Role', value: 'role' },
        { text: 'Created', value: 'createdAt' },
        { text: 'Updated', value: 'updatedAt' },
        { text: 'Actions', value: 'action', align: 'right', sortable: false }
      ],
      valid: false,
      roles: [
        { text: 'Admin', value: 'admin' },
        { text: 'Marketeer', value: 'marketeer' }
      ],
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
    ...mapActions('users', ['list', 'show', 'create', 'resend', 'update', 'delete']),
    async save() {
      try {
        let res = null
        if (this.user.id) {
          res = await this.update(this.user)
        } else {
          res = await this.create(this.user)
        }

        if (res && res.id) {
          this.user.id = res.id
        }

        this.$notify.success('User successfully saved.')
        this.close()
      } catch (e) {
        this.$notify.error('Could not save user.')
      }

      this.list()
    },
    edit(user) {
      this.user = Vue.util.extend({}, user)
      this.dialog = true
    },
    async resendInvitiation(user) {
      try {
        await this.resend(user.id)
        this.$notify.success('Invitation successfully resent.')
      } catch (_) {
        this.$notify.error('Could not resent invitation.')
      }
    },
    remove(user) {
      if (confirm('Do you really want to delete this user?')) {
        this.delete(user.id)
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
      this.user = {}
      this.$refs.form.reset()
    }
  },
  computed: {
    ...mapGetters('users', { allusers: 'all' }),
    userLoaded() {
      return this.user.id && this.user.id.length > 0
    }
  },
  async created() {
    this.list()
  }
}
</script>
