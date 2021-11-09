<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-form ref="form" v-model="valid">
          <v-card class="elevation-12">
            <v-toolbar dark flat>
              <v-toolbar-title>Activate Account</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
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
                name="passwordConfirmation"
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
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" :disabled="!valid" @click.prevent="activate()" type="submit">Activate</v-btn>
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
      password: '',
      passwordConfirmation: '',
      name: '',
      valid: false,
      rules: {
        required: value => !!value || 'Required.',
        min: v => v.length >= 8 || 'Min 8 characters'
      }
    }
  },
  methods: {
    ...mapActions('user', { activateUser: 'activate' }),
    async activate() {
      if (this.password !== this.passwordConfirmation) {
        this.$notify.error('Your password and the password confirmation do not match.')
        return
      }

      try {
        await this.activateUser({
          token: this.$route.params.token,
          password: this.password,
          name: this.name
        })

        this.$router.push('/user/login?state=activated')
      } catch (e) {
        this.$notify.error('An error has occured, please check your password and try again.')
      }
    }
  }
}
</script>
