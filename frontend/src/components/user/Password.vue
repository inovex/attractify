<template>
  <div>
    <v-card-title>Change password</v-card-title>
    <v-card-text>
      <v-form ref="form" v-model="valid">
        <v-row>
          <v-col class="col-lg-4">
            <v-text-field
              label="Current Password"
              name="password"
              prepend-icon="mdi-lock"
              type="password"
              v-model="password"
              hint="At least 8 characters"
              :rules="[rules.required, rules.min]"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col class="col-lg-4">
            <v-text-field
              label="Your New Password"
              name="newPassword"
              prepend-icon="mdi-lock"
              type="password"
              v-model="newPassword"
              hint="At least 8 characters"
              :rules="[rules.required, rules.min]"
            />
          </v-col>
          <v-col class="col-lg-4">
            <v-text-field
              label="Confirm Your New Password"
              name="newPasswordConfirmation"
              prepend-icon="mdi-lock"
              type="password"
              hint="At least 8 characters"
              v-model="newPasswordConfirmation"
              :rules="[rules.required, rules.min]"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-alert type="info"
              >Your new password must be at least 8 chars long. We recommend a longer password (> 20 chars) with a
              combination of numbers (upper- and lowercase), letters and special chars.</v-alert
            >
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer />
      <v-btn color="primary" :disabled="!valid" @click="change()" rounded>Change Password</v-btn>
    </v-card-actions>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data() {
    return {
      password: '',
      newPassword: '',
      newPasswordConfirmation: '',
      valid: false,
      rules: {
        required: value => !!value || 'Required.',
        min: v => v.length >= 8 || 'Min 8 characters'
      }
    }
  },
  methods: {
    ...mapActions('user', ['updatePassword']),
    async change() {
      if (this.newPassword !== this.newPasswordConfirmation) {
        this.$notify.error('Your password and the password confirmation do not match.')
        return
      }

      try {
        await this.updatePassword({
          oldPassword: this.password,
          newPassword: this.newPassword
        })

        this.$notify.success('Your password has been updated successfully.')
      } catch (e) {
        this.$notify.error('Could not update password.')
      }
    }
  }
}
</script>
