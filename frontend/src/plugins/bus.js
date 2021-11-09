import Vue from 'vue'

const vueBus = new Vue()

const BusPlugin = {
  install(instance) {
    instance.prototype.$bus = vueBus
    instance.bus = vueBus
  }
}

export default BusPlugin
