<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Change Your account details</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="user" />
          </v-toolbar>
          <Attributes :user="user" />
          <v-divider></v-divider>
          <Password />
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import Attributes from './user/Attributes'
import Password from './user/Password'

import { mapActions } from 'vuex'

export default {
  components: { Attributes, Password, Help },
  data() {
    return {
      user: {}
    }
  },
  methods: {
    ...mapActions('user', { showUser: 'show' })
  },
  async created() {
    this.user = await this.showUser()
  }
}
</script>
