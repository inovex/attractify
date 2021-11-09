<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-form ref="form" v-model="valid">
          <v-card class="elevation-12">
            <v-toolbar dark flat>
              <v-toolbar-title>Register new Organization</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
              <v-text-field
                label="Your Email Address"
                name="email"
                prepend-icon="mdi-email"
                type="text"
                v-model="email"
                :rules="[rules.required]"
                autofocus
              />

              <v-text-field
                label="Your Password"
                name="password"
                prepend-icon="mdi-lock"
                type="password"
                v-model="password"
                hint="At least 8 characters"
                :rules="[rules.required, rules.min]"
              />

              <v-text-field
                label="Confirm Your Password"
                name="passwordConfirmationField"
                prepend-icon="mdi-lock"
                type="password"
                hint="At least 8 characters"
                v-model="passwordConfirmation"
                :rules="[rules.required, rules.min]"
              />

              <v-text-field
                label="Your Name"
                name="name"
                prepend-icon="mdi-account"
                type="text"
                v-model="name"
                :rules="[rules.required]"
              />

              <v-select
                label="Your Timezone"
                prepend-icon="mdi-clock"
                v-model="timezone"
                :items="timezones"
                :rules="[rules.required]"
              />

              <v-checkbox v-model="terms" :rules="[rules.checked]">
                <span slot="label">
                  Accept
                  <router-link to="/terms">Terms of Use</router-link>
                </span>
              </v-checkbox>
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" :disabled="!valid" @click.prevent="register()" type="submit">Register</v-btn>
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
      name: '',
      timezone: '',
      timezones: ['Europe/Berlin'],
      valid: false,
      terms: false,
      rules: {
        checked: value => value || 'Required.',
        required: value => !!value || 'Required.',
        min: v => v.length >= 8 || 'Min 8 characters'
      }
    }
  },
  methods: {
    ...mapActions('organization', ['signUp']),
    async register() {
      if (this.password !== this.passwordConfirmation) {
        this.$notify.error('Your password and the password confirmation do not match.')
        return
      }

      try {
        await this.signUp({
          email: this.email,
          password: this.password,
          name: this.name,
          timezone: this.timezone
        })

        this.$router.push('/user/login?state=activated')
      } catch (e) {
        this.$notify.error('An error has occured, please check your email and password and try again.')
      }
    }
  }
}
</script>
