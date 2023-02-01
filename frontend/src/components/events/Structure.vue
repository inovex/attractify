<template>
  <div>
    <v-card outlined>
      <v-card-text>
        <div v-if="structure && structure.length === 0">No properties defined yet</div>
        <v-treeview :items="structure" open-all hoverable dense open-on-click>
          <template slot="label" slot-scope="{ item }">
            <span class="structure__property">{{ item.name }}</span>
            <em class="structure__hints ml-2"
              >{{ item.properties.type }}{{ item.properties.isRequired ? ', required' : '' }}</em
            >
          </template>
          <template slot="append" slot-scope="{ item }">
            <v-btn @click="openAddChildDialog(item)" small icon>
              <v-icon title="Add child property">mdi-plus</v-icon>
            </v-btn>
            <v-btn @click="openEditDialog(item)" small icon>
              <v-icon title="Edit property">mdi-pencil</v-icon>
            </v-btn>
            <v-btn @click="removeChild(item.id)" small icon>
              <v-icon title="Delete property">mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-treeview>
      </v-card-text>
    </v-card>
    <br />
    <v-btn @click="openAddRootDialog()" small text>
      <v-icon>mdi-plus</v-icon>
      <span>Add property</span>
    </v-btn>
    <v-dialog v-model="dialog" max-width="700px" closeable>
      <v-card>
        <v-card-title>
          <span class="headline">Property Details</span>
          <v-spacer></v-spacer>
          <help name="properties" />
        </v-card-title>
        <v-form ref="form" v-model="valid">
          <v-card-text>
            <v-row>
              <v-col>
                <v-text-field
                  label="Key"
                  prepend-icon="mdi-key"
                  type="text"
                  v-model="property.key"
                  :rules="[rules.required]"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-select
                  label="Type"
                  :items="types"
                  v-model="property.type"
                  prepend-icon="mdi-ticket-percent"
                  :rules="[rules.required]"
                ></v-select>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-text-field label="Regex Pattern" prepend-icon="mdi-regex" type="text" v-model="property.pattern" />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-switch label="Required" prepend-icon="mdi-debug-step-over" v-model="property.isRequired" />
              </v-col>
            </v-row>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn rounded @click="cancel()">Cancel</v-btn>
            <v-btn rounded color="primary" style="color: var(--v-buttontext-base)" :disabled="!valid" @click="save()"
              >Save</v-btn
            >
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import Help from '../Help'
export default {
  components: { Help },
  props: ['structure'],
  data() {
    return {
      dialog: false,
      property: { isRequired: false },
      saveCallback: null,
      types: [
        { text: 'Object', value: 'object' },
        { text: 'String', value: 'string' },
        { text: 'Integer', value: 'integer' },
        { text: 'Float', value: 'float' },
        { text: 'Boolean', value: 'boolean' },
        { text: 'DateTime', value: 'dateTime' },
        { text: 'Array', value: 'array' }
      ],
      valid: false,
      rules: {
        required: (value) => !!value || 'Required.'
      }
    }
  },
  methods: {
    openDialog(property = null) {
      if (property) {
        this.property = property
      }
      this.dialog = true
    },
    cancel() {
      this.dialog = false
      this.property = { isRequired: false }
      this.$refs.form.reset()
    },
    save() {
      this.$emit('savecallback')
      this.saveCallback(this.property)
      this.cancel()
    },
    openAddRootDialog() {
      this.saveCallback = (property) => {
        this.addNode(this.structure, property)
      }
      this.openDialog()
    },
    openAddChildDialog(item) {
      this.saveCallback = (property) => {
        item.properties.type = 'object'
        if (!item.children) {
          this.$set(item, 'children', [])
        }
        this.addNode(item.children, property)
      }
      this.openDialog()
    },
    openEditDialog(item) {
      this.saveCallback = (property) => {
        item.id = property.key
        item.name = property.key
        item.properties = property
      }
      this.openDialog(item.properties)
    },
    addNode(item, property) {
      item.push({
        id: property.key,
        name: property.key,
        properties: {
          type: property.type,
          key: property.key,
          pattern: property.pattern,
          isRequired: property.isRequired
        }
      })
    },
    removeChild(id, items = null) {
      if (!items) {
        items = this.structure
      }

      this.$emit('savecallback')
      let i = 0
      for (let n of items) {
        if (n.id === id) {
          items = items.splice(i, 1)
          return
        }

        if (n.children && n.children.length > 0) {
          this.removeChild(id, n.children)
        }

        i++
      }
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
