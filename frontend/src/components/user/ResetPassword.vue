<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-form ref="form" v-model="valid">
          <v-card class="elevation-12">
            <v-toolbar dark flat>
              <v-toolbar-title>Reset Password</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
              Please enter your email address. We'll send you a reset code that you can use to reset your password.
            </v-card-text>
            <v-card-text>
              <v-text-field
                label="Your Email Address"
                name="email"
                prepend-icon="mdi-email"
                type="text"
                v-model="email"
                v-if="!resetMode"
                :rules="[rules.required]"
                autofocus
              />

              <v-text-field
                label="Your Password"
                name="password"
                prepend-icon="mdi-lock"
                type="password"
                v-model="password"
                v-if="resetMode"
                :rules="[rules.required, rules.min]"
              />

              <v-text-field
                label="Confirm Your Password"
                name="passwordConfirmationField"
                prepend-icon="mdi-lock"
                type="password"
                hint="At least 8 characters"
                v-model="passwordConfirmation"
                v-if="resetMode"
                :rules="[rules.required, rules.min]"
              />

              <v-alert type="info"
                >Your new password must be at least 8 chars long. We recommend a longer password (> 20 chars) with a
                combination of numbers (upper- and lowercase), letters and special chars.</v-alert
              >
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn
                color="primary"
                @click.prevent="requestCode()"
                :disabled="!valid"
                v-if="!resetMode"
                type="submit"
                rounded
                >Request Code</v-btn
              >
              <v-btn
                color="primary"
                @click.prevent="setNewPassword()"
                :disabled="!valid"
                v-if="resetMode"
                type="submit"
                rounded
                >Reset Password</v-btn
              >
            </v-card-actions>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data() {
    return {
      email: '',
      password: '',
      passwordConfirmation: '',
      valid: false,
      resetMode: false,
      token: '',
      rules: {
        required: value => !!value || 'Required.',
        min: v => v.length >= 8 || 'Min 8 characters',
        exact: v => v.length == 6 || 'Should be 6 characters'
      }
    }
  },
  methods: {
    ...mapActions('user', { reset: 'reset', resetPassword: 'resetPassword' }),
    async requestCode() {
      try {
        await this.reset(this.email)
        this.$notify.success('A reset code has been sent to your email address.')
      } catch (e) {
        this.$notify.error('Could not send reset code.')
      }
    },
    async setNewPassword() {
      if (this.password !== this.passwordConfirmation) {
        this.$notify.error('Your password and the password confirmation do not match.')
        return
      }

      try {
        await this.resetPassword({
          token: this.token,
          password: this.password
        })

        this.$router.push('/user/login?state=resetted')
      } catch (e) {
        this.$notify.error('Could not update password.')
      }
    }
  },
  created() {
    this.token = this.$route.params.token
    if (this.token && this.token.length > 0) {
      this.resetMode = true
    }
  }
}
</script>
