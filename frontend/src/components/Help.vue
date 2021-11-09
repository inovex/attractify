<template>
  <v-dialog v-model="dialog" max-width="700">
    <template v-slot:activator="{ on, attrs }">
      <v-btn icon v-bind="attrs" v-on="on">
        <v-icon>mdi-help</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">
        {{ content.title }}
      </v-card-title>
      <v-card-text>
        <vue-markdown :breaks="false" class="help__body">{{ content.body }}</vue-markdown>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="dialog = false">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import help from '../help/help'
import VueMarkdown from 'vue-markdown'

export default {
  components: { VueMarkdown },
  props: ['name'],
  data() {
    return {
      dialog: false,
      content: {}
    }
  },
  methods: {
    getTitle() {}
  },
  created() {
    if (help[this.name]) {
      this.content = help[this.name]
    }
  }
}
</script>

<style>
.help__body h1 {
  font-size: 1.2em;
  margin-top: 1em;
  margin-bottom: 0.4em;
}

.help__body ul li {
  margin-left: -10px;
  line-height: 1.5em;
}

.help__body p {
  margin: 1em 0;
  line-height: 1.5em;
}
</style>
