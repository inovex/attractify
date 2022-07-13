<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-img class="logo" :src="require('../../assets/logo_black.svg')" c />
        <v-form ref="form" v-model="valid">
          <v-card class="elevation-12">
            <v-toolbar dark flat>
              <v-toolbar-title>Please log in</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
              <v-alert type="success" v-if="state === 'activated'"
                >Your account has been activated. You can log in now.</v-alert
              >

              <v-alert type="success" v-if="state === 'resetted'"
                >Your password has been resetted. You can log in now.</v-alert
              >

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
                :rules="[rules.required]"
                v-model="password"
              />
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" @click.prevent="login()" :disabled="!valid" type="submit" rounded>Login</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
        <br />
        <div class="user__links">
          <router-link :to="{ path: '/user/reset-password' }">Forgot password?</router-link>
        </div>
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
      valid: false,
      success: false,
      rules: {
        required: (value) => !!value || 'Required.',
        min: (v) => v.length >= 8 || 'Min 8 characters'
      }
    }
  },
  methods: {
    ...mapActions('user', ['signIn']),
    async login() {
      this.success = false
      this.state = ''

      try {
        await this.signIn({
          email: this.email,
          password: this.password
        })

        if (this.$route.query.redirect && this.$route.query.redirect.indexOf('/') === 0) {
          this.$router.push(this.$route.query.redirect)
        } else {
          this.$router.push('/')
        }

        this.$bus.$emit('user:update')
      } catch (e) {
        this.$notify.error('Wrong username or password.')
      }
    }
  },
  created() {
    this.state = this.$route.query.state
  }
}
</script>

<style scoped>
.logo {
  width: 60%;
  margin: 0 auto;
  margin-top: 20%;
  margin-bottom: 20%;
}

.user__links {
  text-align: center;
  font-size: 0.9em;
}

.user__links a {
  text-decoration: none;
}
</style>
