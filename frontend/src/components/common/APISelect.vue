<template>
  <div>
    <v-select
      :dense="dense"
      :items="items"
      :label="label"
      :name="name"
      :loading="isLoading"
      :value="value"
      :placeholder="placeholder"
      :hide-details="hideDetails || false"
      :prepend-icon="icon"
      :outlined="outlined"
      :multiple="multiple"
      @input="input"
      @change="change"
      :return-object="returnObject"
      :rules="rules"
      :clearable="clearable"
    ></v-select>
  </div>
</template>

<script>
export default {
  props: [
    'label',
    'placeholder',
    'icon',
    'name',
    'loadCallback',
    'hideDetails',
    'value',
    'dense',
    'multiple',
    'returnObject',
    'outlined',
    'rules',
    'clearable'
  ],
  data: () => ({
    model: null,
    items: [],
    isLoading: false
  }),
  methods: {
    input(e) {
      this.$emit('input', e)
    },
    change(e) {
      this.$emit('change', e)
    },
    async load() {
      if (this.isLoading) return
      this.isLoading = true

      try {
        this.items = await this.loadCallback()
      } catch (_) {
        _
      } finally {
        this.isLoading = false
      }
    }
  },
  async created() {
    await this.load()
  }
}
</script>
