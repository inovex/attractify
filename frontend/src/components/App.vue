<template>
  <v-app id="app">
    <v-snackbar v-model="showNotification" :top="true" :color="notification.type" :timeout="4000" :vertical="true">
      {{ notification.message }}
      <v-btn text @click="showNotification = false">Close</v-btn>
    </v-snackbar>
    <v-navigation-drawer v-if="isLoggedIn" v-model="drawer" :clipped="$vuetify.breakpoint.lgAndUp" app dark>
      <v-toolbar-title @click="$router.push({ path: '/' })">
        <v-img class="header__logo" :src="require('../assets/logo.svg')" c />
      </v-toolbar-title>

      <v-list>
        <template v-for="(item, i) in items">
          <template v-if="item.link">
            <v-list-item if="isVisibleForRole(item.role)" :key="item.text" :to="item.link" link>
              <v-list-item-action>
                <v-icon>{{ item.icon }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>{{ item.text }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
          <v-divider v-else-if="item.divider" :key="i" />
          <v-subheader v-if="item.header" :key="i">{{ item.text }}</v-subheader>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container fluid>
        <v-row justify="end">
          <v-col cols="auto" class="mr-3">
            <v-tooltip v-if="!$vuetify.theme.dark" bottom>
              <template v-slot:activator="{ on }">
                <v-btn v-on="on" width="48" height="48" fab @click="darkMode">
                  <v-icon color="primary">mdi-white-balance-sunny</v-icon>
                </v-btn>
              </template>
              <span>Dark Mode On</span>
            </v-tooltip>

            <v-tooltip v-else bottom>
              <template v-slot:activator="{ on }">
                <v-btn v-on="on" width="48" height="48" fab @click="darkMode">
                  <v-icon color="primary">mdi-moon-waning-crescent</v-icon>
                </v-btn>
              </template>
              <span>Dark Mode Off</span>
            </v-tooltip>
          </v-col>
          <v-col cols="auto" class="mr-3">
            <v-menu bottom v-if="isLoggedIn">
              <template v-slot:activator="{ on }">
                <v-avatar
                  class="avatar__icon"
                  color="primary"
                  style="color: var(--v-buttontext-base)"
                  size="48"
                  v-on="on"
                  >{{ initials }}</v-avatar
                >
              </template>
              <v-list>
                <v-list-item :to="{ path: '/user' }">
                  <v-list-item-title> <v-icon left>mdi-account</v-icon>My Profile </v-list-item-title>
                </v-list-item>
                <v-list-item @click="logout()">
                  <v-list-item-title> <v-icon left>mdi-logout</v-icon>Log Out </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <template v-if="$route.matched.length">
          <router-view v-if="renderComponent" :darkmode="this.darkmode"></router-view>
        </template>
      </v-container>
    </v-main>
    <v-footer class="footer">
      <v-col class="text-center" cols="12">
        {{ new Date().getFullYear() }} —
        <strong>Attractify</strong>
        — &nbsp;
        <a class="footer__link" href="https://www.inovex.de/de/impressum/" target="_blank">Legal</a>
        -
        <a class="footer__link" href="https://www.inovex.de/de/datenschutz/" target="_blank">Privacy</a>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'

export default {
  data() {
    return {
      renderComponent: true,
      darkmode: false,
      user: {},
      showNotification: false,
      notification: {},
      isLoggedIn: false,
      drawer: null,
      items: [
        { icon: 'mdi-monitor-dashboard', text: 'Dashboard', link: '/' },
        { divider: true },
        { header: true, text: 'Events' },
        {
          icon: 'mdi-devices',
          text: 'Channels',
          link: '/channels'
        },
        {
          icon: 'mdi-clipboard-list-outline',
          text: 'Events',
          link: '/events'
        },
        {
          icon: 'mdi-clipboard-alert-outline',
          text: 'Invalid Events',
          link: '/invalid-events'
        },
        {
          icon: 'mdi-badge-account-outline',
          text: 'Contexts',
          link: '/contexts'
        },

        { icon: 'mdi-table-large', text: 'Event Log', link: '/event-log' },
        { divider: true },
        { header: true, text: 'Personas' },
        {
          icon: 'mdi-account-search',
          text: 'User Profiles',
          link: '/profiles'
        },
        {
          icon: 'mdi-account-key',
          text: 'Custom Traits',
          link: '/custom-traits'
        },
        {
          icon: 'mdi-account-star',
          text: 'Computed Traits',
          link: '/computed-traits'
        },
        {
          icon: 'mdi-account-group',
          text: 'Audiences',
          link: '/audiences'
        },

        { divider: true },
        { header: true, text: 'Actions' },
        { icon: 'mdi-ticket-percent', text: 'Actions', link: '/actions' },
        { icon: 'mdi-clipboard-text-outline', text: 'Action Types', link: '/action-types' },
        {
          icon: 'mdi-gesture-tap',
          text: 'Reactions',
          link: '/reactions'
        },
        { divider: true },
        { header: true, text: 'Insights' },
        { icon: 'mdi-chart-bar', text: 'Analyze', link: '/analyze' },

        { divider: true },
        { header: true, text: 'Settings' },
        { icon: 'mdi-api', text: 'API Access', link: '/api' },
        {
          icon: 'mdi-account-multiple',
          text: 'User Management',
          link: '/users',
          role: 'admin'
        },
        {
          icon: 'mdi-cog',
          text: 'Your Organization',
          link: '/organization',
          role: 'admin'
        },
        {
          icon: 'mdi-shield-account',
          text: 'GDPR/DSGVO',
          link: '/privacy'
        }
      ]
    }
  },
  methods: {
    ...mapActions('user', ['signOut']),
    refresh() {
      if (this.get !== null) {
        this.isLoggedIn = true
        this.user = this.get
      }
    },
    notify(params) {
      this.showNotification = true
      this.notification = params
    },
    async logout() {
      try {
        await this.signOut()
        this.isLoggedIn = false
        this.$router.push({ path: '/user/login' })
      } catch (error) {
        this.$router.push({ path: '/user/login' })
      }

      window.localStorage.clear()
      window.sessionStorage.clear()
    },
    isVisibleForRole(role) {
      if (role === 'admin') {
        return this.user.role === 'admin'
      }
      return true
    },
    darkMode() {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark
      this.darkmode = this.$vuetify.theme.dark
      window.localStorage.setItem('darkmode', this.$vuetify.theme.dark)
      if (window.location.pathname == '/') {
        this.forceRerender()
      }
    },
    forceRerender() {
      this.renderComponent = false

      this.$nextTick(() => {
        this.renderComponent = true
      })
    }
  },
  computed: {
    ...mapGetters('user', ['get']),
    initials() {
      let name = (this.user.name || this.user.email).split(' ')
      return `${(name[0] || ['U'])[0]}${(name[1] || [''])[0]}`.toUpperCase()
    }
  },
  created() {
    this.$bus.$on('user:update', this.refresh)
    this.$bus.$on('notify', this.notify)
    this.$bus.$on('user:logout', this.logout)
    if (window.localStorage.getItem('darkmode') != null) {
      this.$vuetify.theme.dark = JSON.parse(window.localStorage.getItem('darkmode'))
    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      this.$vuetify.theme.dark = true
    }
    this.darkmode = this.$vuetify.theme.dark
    this.refresh()
  }
}
</script>

<style scoped>
.header__logo {
  cursor: pointer;
  width: 75%;
  margin: 12px 0 12px 20px;
}

.avatar__icon {
  cursor: pointer;
}

.avatar__icon:hover {
  opacity: 0.8;
}

.footer {
  font-size: 0.8em;
}

.footer__link {
  text-decoration: none;
}
</style>
