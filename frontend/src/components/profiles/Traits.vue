<template>
  <div>
    <v-row>
      <v-col>
        <v-list v-if="hasTraits()">
          <v-card-text>
            <!-- TODO: show values and make them editable -->
            <Structure :structure="traits" />
          </v-card-text>
          <v-list-item v-for="(value, key) in traits" :key="key">
            <v-list-item-content>
              <strong>{{ key }}</strong>
            </v-list-item-content>
            <v-list-item-content>{{ value }}</v-list-item-content>
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
    return {}
  },
  methods: {
    hasTraits() {
      if (!this.traits) {
        return false
      }
      return Object.keys(this.traits).length > 0
    }
  }
}
</script>
