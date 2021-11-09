import Vue from 'vue'
import VueRouter from 'vue-router'

import Dashboard from '../components/Dashboard.vue'
import Action from '../components/Action.vue'
import Actions from '../components/Actions.vue'
import Reactions from '../components/Reactions.vue'
import Analyze from '../components/Analyze.vue'

import Organization from '../components/Organization.vue'
import Register from '../components/organization/Register.vue'

import User from '../components/User.vue'
import Login from '../components/user/Login.vue'
import Activate from '../components/user/Activate.vue'
import ResetPassword from '../components/user/ResetPassword.vue'

import NotFound from '../components/NotFound.vue'
import API from '../components/API.vue'
import Users from '../components/Users.vue'
import store from '../store/index.js'

import Profiles from '../components/Profiles.vue'
import Profile from '../components/Profile.vue'

import Contexts from '../components/Contexts.vue'
import Context from '../components/Context.vue'
import Event from '../components/Event.vue'
import Events from '../components/Events.vue'
import EventLog from '../components/EventLog.vue'

import Channels from '../components/Channels.vue'

import Audiences from '../components/Audiences.vue'
import Audience from '../components/Audience.vue'

import CustomTraits from '../components/CustomTraits.vue'
import ComputedTraits from '../components/ComputedTraits.vue'
import ComputedTrait from '../components/ComputedTrait.vue'

import Privacy from '../components/Privacy.vue'

Vue.use(VueRouter)

function requireAuth(to, from, next) {
  const user = store.getters['user/get']

  if (user) {
    next()
  } else {
    next({
      path: '/user/login',
      query: { redirect: to.fullPath }
    })
  }
}

export default new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Dashboard, beforeEnter: requireAuth },
    { path: '/actions', component: Actions, beforeEnter: requireAuth },
    { path: '/reactions', component: Reactions, beforeEnter: requireAuth },
    { path: '/action/:id?', component: Action, beforeEnter: requireAuth },
    { path: '/analyze/:id?', component: Analyze, beforeEnter: requireAuth },
    {
      path: '/organization',
      component: Organization,
      beforeEnter: requireAuth
    },
    { path: '/api', component: API, beforeEnter: requireAuth },
    { path: '/users', component: Users, beforeEnter: requireAuth },
    { path: '/event-log', component: EventLog, beforeEnter: requireAuth },

    { path: '/contexts', component: Contexts, beforeEnter: requireAuth },
    { path: '/context/:id?', component: Context, beforeEnter: requireAuth },
    {
      path: '/events',
      component: Events,
      beforeEnter: requireAuth
    },
    {
      path: '/event/:id?',
      component: Event,
      beforeEnter: requireAuth
    },

    {
      path: '/channels',
      component: Channels,
      beforeEnter: requireAuth
    },

    { path: '/organization/register', component: Register },

    { path: '/user', component: User },
    { path: '/user/login', component: Login },
    { path: '/user/reset-password', component: ResetPassword },
    { path: '/user/reset-password/:token?', component: ResetPassword },
    { path: '/user/activate/:token', component: Activate },

    { path: '/profiles', component: Profiles },
    { path: '/profile/:id?', component: Profile },
    { path: '*', component: NotFound },

    { path: '/audiences', component: Audiences },
    { path: '/audience/:id?', component: Audience },

    { path: '/custom-traits', component: CustomTraits },
    { path: '/computed-traits', component: ComputedTraits },
    { path: '/computed-trait/:id?', component: ComputedTrait },

    { path: '/privacy', component: Privacy }
  ]
})
