<template>
    <div>
      <v-card>
        <v-card-title>
          <span class="headline">Unsaved Changes</span>
          <v-spacer></v-spacer>
        </v-card-title>
        <v-card-text>
            You have unsaved changes. Do you really want to exit?
        </v-card-text>
        <v-form ref="form">
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn rounded @click="cancel()">Cancel</v-btn>
            <v-btn rounded color="primary" style="color: var(--v-buttontext-base)" @click="exit()">Exit</v-btn>
            <v-btn rounded color="primary" style="color: var(--v-buttontext-base)" :disabled="!valid" @click="save()">Save & Exit</v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </div>
</template>

<script>
export default {
  props: {
    valid: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    cancel(){
        this.$emit('cancel')
    },
    save() {
      this.$emit('save')
    },
    exit() {
      this.$emit('exit')
    }
  },
  async created() {}
}
</script>

<style scoped>
.structure__property {
  font-family: Consolas, Menlo, Monaco, Lucida Console, Liberation Mono, DejaVu Sans Mono, Bitstream Vera Sans Mono,
    Courier New, monospace;
  font-weight: bold;
}

.structure__hints {
  opacity: 0.4;
}
</style>
