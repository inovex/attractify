<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Change Your Organization Details</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="organization" />
          </v-toolbar>
          <Attributes :organization="organization" />
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Help from './Help'
import Attributes from './organization/Attributes'

import { mapActions } from 'vuex'

export default {
  components: { Attributes, Help },
  data() {
    return {
      organization: {}
    }
  },
  methods: {
    ...mapActions('organization', { showOrganization: 'show' })
  },
  async created() {
    this.organization = await this.showOrganization()
  }
}
</script>
