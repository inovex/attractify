import { createApp } from 'vue'

const vue = createApp()

const methods = {
  success: message => {
    vue.$bus.emit('notify', { type: 'success', message: message })
  },
  error: message => {
    vue.$bus.emit('notify', { type: 'error', message: message })
  },
  info: message => {
    vue.$bus.emit('notify', { type: 'cyan', message: message })
  }
}

export default {
  install(Vue) {
    Vue.provide.$notify = methods
    Vue.notify = methods
  }
}
