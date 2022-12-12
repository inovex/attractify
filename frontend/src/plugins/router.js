import Vue from 'vue'
import VueRouter from 'vue-router'

import Dashboard from '../components/Dashboard.vue'
import Action from '../components/Action.vue'
import Actions from '../components/Actions.vue'
import ActionSimulation from '../components/ActionSimulation.vue'
import ActionType from '../components/ActionType.vue'
import ActionTypes from '../components/ActionTypes.vue'
import Reactions from '../components/Reactions.vue'
import Analyze from '../components/Analyze.vue'

import Organization from '../components/Organization.vue'

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
import InvalidEvents from '../components/InvalidEvents.vue'
import EventLog from '../components/EventLog.vue'

import Channels from '../components/Channels.vue'

import Audiences from '../components/Audiences.vue'
import Audience from '../components/Audience.vue'

import CustomTraits from '../components/CustomTraits.vue'
import ComputedTraits from '../components/ComputedTraits.vue'
import ComputedTrait from '../components/ComputedTrait.vue'

import Privacy from '../components/Privacy.vue'

Vue.use(VueRouter)

const stayLoggedInTime = 86400000;

function requireAuth(to, from, next) {
  const user = store.getters['user/get']

  if (user) {
    if (user.timestamp > ((Date.now() - stayLoggedInTime))) {
      user.timestamp = Date.now()
      next()
      return
    }
    Vue.prototype.$bus.$emit('user:logout')
  }
  next({
    path: '/user/login',
    query: { redirect: to.fullPath }
  })

}

export default new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Dashboard, beforeEnter: requireAuth },
    { path: '/actions', component: Actions, beforeEnter: requireAuth },
    { path: '/action-simulation/:id?', component: ActionSimulation, beforeEnter: requireAuth },
    { path: '/action-types', component: ActionTypes, beforeEnter: requireAuth },
    { path: '/reactions', component: Reactions, beforeEnter: requireAuth },
    { path: '/action/:id?', component: Action, beforeEnter: requireAuth },
    { path: '/action-type/:id?', component: ActionType, beforeEnter: requireAuth },
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
      path: '/invalid-events',
      component: InvalidEvents,
      beforeEnter: requireAuth
    },
    {
      path: '/channels',
      component: Channels,
      beforeEnter: requireAuth
    },

    { path: '/profiles', component: Profiles, beforeEnter: requireAuth },
    { path: '/profile/:id?', component: Profile, beforeEnter: requireAuth },
    { path: '*', component: NotFound, beforeEnter: requireAuth },

    { path: '/audiences', component: Audiences, beforeEnter: requireAuth },
    { path: '/audience/:id?', component: Audience, beforeEnter: requireAuth },

    { path: '/custom-traits', component: CustomTraits, beforeEnter: requireAuth },
    { path: '/computed-traits', component: ComputedTraits, beforeEnter: requireAuth },
    { path: '/computed-trait/:id?', component: ComputedTrait, beforeEnter: requireAuth },

    { path: '/privacy', component: Privacy, beforeEnter: requireAuth },

    { path: '/user', component: User, beforeEnter: requireAuth },
    { path: '/user/login', component: Login },
    { path: '/user/reset-password', component: ResetPassword },
    { path: '/user/reset-password/:token?', component: ResetPassword },
    { path: '/user/activate/:token', component: Activate }

  ]
})
