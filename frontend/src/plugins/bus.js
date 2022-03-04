/*import Vue from 'vue'

const vueBus = new Vue()

const BusPlugin = {
  install(instance) {
    instance.prototype.$bus = vueBus
    instance.bus = vueBus
  }
}

export default BusPlugin
*/
import mitt from 'mitt'

const emitter = mitt()

const BusPlugin = {
  install(app) {
    app.config.globalProperties.$bus = emitter
  }
}

export { BusPlugin }
