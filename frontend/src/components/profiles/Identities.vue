<template>
  <div>
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">ID</th>
            <th class="text-left">Type</th>
            <th class="text-left">Is Anonymous</th>
            <th class="text-left">Channel</th>
            <th class="text-left">Created At</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in identities" :key="item.id">
            <td>{{ item.userId }}</td>
            <td>{{ item.type }}</td>
            <td>{{ item.isAnonymous }}</td>
            <td>{{ item.channel }}</td>
            <td>{{ item.createdAt }}</td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>
import profiles from '../../lib/rest/profiles'

export default {
  props: ['profile'],
  data() {
    return {
      identities: []
    }
  },
  methods: {},
  async created() {
    const id = this.$route.params.id
    if (id) {
      try {
        this.identities = await profiles.listIdentities(id)
      } catch (_) {
        _
      }
    }
  }
}
</script>
