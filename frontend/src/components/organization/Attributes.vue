<template>
  <div>
    <v-card-title>Name</v-card-title>
    <v-card-text>
      <v-form ref="form" v-model="valid">
        <v-text-field
          label="Organization Name"
          name="name"
          prepend-icon="mdi-account"
          type="text"
          v-model="organization.name"
          :rules="[rules.required]"
        />
      </v-form>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-title>Timezone</v-card-title>
    <v-card-text>
      <v-alert type="info">Be careful with changing yout timezone, as this will affect all event trackings.</v-alert>
    </v-card-text>
    <v-card-text>
      <v-form ref="form" v-model="valid">
        <v-select
          label="Your Timezone"
          prepend-icon="mdi-clock"
          v-model="organization.timezone"
          :items="timezones"
          :rules="[rules.required]"
        />
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer />
      <v-btn color="primary" :disabled="!valid" @click="change()" text>Change Details</v-btn>
    </v-card-actions>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  props: ['organization'],
  data() {
    return {
      valid: false,
      timezones: ['Europe/Berlin'],
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    ...mapActions('organization', ['update']),
    async change() {
      try {
        await this.update({
          name: this.organization.name,
          timezone: this.organization.timezone
        })

        this.$notify.success('Changes have been saved successfully.')
      } catch (e) {
        this.$notify.error('Could not save attributes.')
      }
    }
  }
}
</script>
