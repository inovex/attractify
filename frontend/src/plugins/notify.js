import Vue from 'vue'

const vue = new Vue()

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
    Vue.prototype.$notify = methods
    Vue.notify = methods
  }
}
