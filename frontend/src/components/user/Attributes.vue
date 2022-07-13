<template>
  <div>
    <v-card-title>Personal details</v-card-title>
    <v-card-text>
      <v-form ref="form" v-model="valid">
        <v-text-field
          label="Your Email"
          name="email"
          prepend-icon="mdi-email"
          type="text"
          v-model="user.email"
          :rules="[rules.required]"
        />

        <v-text-field
          label="Your Name"
          name="name"
          prepend-icon="mdi-account"
          type="text"
          v-model="user.name"
          :rules="[rules.required]"
        />
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer />
      <v-btn color="primary" :disabled="!valid" @click="change()" rounded>Change Details</v-btn>
    </v-card-actions>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  props: ['user'],
  data() {
    return {
      valid: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },
  methods: {
    ...mapActions('user', ['update']),
    async change() {
      try {
        if (this.user.email === this.oldEmail) {
          this.update({ name: this.user.name })
        } else {
          this.update(this.user)
        }

        this.$notify.success('Your details have been changed.')
        this.$bus.$emit('user:update')
      } catch (e) {
        this.$notify.error('An error has occured, please check your details.')
      }
    }
  }
}
</script>
