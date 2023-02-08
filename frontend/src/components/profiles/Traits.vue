<template>
  <div>
    <v-row>
      <v-col>
        <v-list v-if="hasTraits()">
          <v-list-item v-for="(value, key) in traits" :key="key">
            <v-list-item-content>
              <v-text-field
                :rules="[() => !!value || 'This field is required']"
                :label="key"
                :disabled="currentEdit != key"
                dense
                v-model="traits[key]"
              ></v-text-field>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn v-if="currentEdit != key" icon @click="edit(key)">
                <v-icon title="edit trait">mdi-pencil</v-icon>
              </v-btn>
              <div v-else>
                <v-btn rounded @click="cancel()"> cancel </v-btn>
                <v-btn icon @click="save(key)">
                  <v-icon title="save trait">mdi-content-save</v-icon>
                </v-btn>
              </div>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        <div v-else>There are no traits available</div>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: ['traits'],
  data() {
    return {
      currentEdit: null
    }
  },
  methods: {
    hasTraits() {
      if (!this.traits) {
        return false
      }
      return Object.keys(this.traits).length > 0
    },
    cancel() {
      this.currentEdit = null
    },
    edit(key) {
      if (key != null) {
        this.$emit('updatedTraits')
      }
      this.currentEdit = key
      console.log(this.traits[key])
    },
    save(key) {
      this.currentEdit = null
      console.log(key)
      this.$emit('updatedTraits')
    }
  }
}
</script>
