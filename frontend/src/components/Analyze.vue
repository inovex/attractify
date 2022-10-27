<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-toolbar dark>
            <v-toolbar-title>Analyze action</v-toolbar-title>
            <v-spacer></v-spacer>
            <help name="analyze" />
          </v-toolbar>
          <v-card-text>
            <v-row>
              <v-col class="col-lg-6">
                <APISelect
                  label="Action"
                  icon="mdi-ticket-percent"
                  :loadCallback="loadActions"
                  v-model="actionId"
                  @change="render"
                />
              </v-col>
              <v-col class="col-lg-2">
                <v-menu
                  ref="menu.start"
                  v-model="menu.start"
                  :close-on-content-click="false"
                  transition="scale-transition"
                  offset-y
                  max-width="290px"
                  min-width="290px"
                >
                  <template v-slot:activator="{ on }">
                    <v-text-field
                      v-model="range.start"
                      label="Date"
                      persistent-hint
                      prepend-icon="mdi-clock-start"
                      v-on="on"
                    ></v-text-field>
                  </template>
                  <v-date-picker
                    v-model="range.start"
                    no-title
                    @input="menu.start = false"
                    @change="render"
                  ></v-date-picker>
                </v-menu>
              </v-col>
              <v-col class="col-lg-2">
                <v-menu
                  ref="menu.end"
                  v-model="menu.end"
                  :close-on-content-click="false"
                  transition="scale-transition"
                  offset-y
                  max-width="290px"
                  min-width="290px"
                >
                  <template v-slot:activator="{ on }">
                    <v-text-field
                      v-model="range.end"
                      label="Date"
                      persistent-hint
                      prepend-icon="mdi-clock-end"
                      v-on="on"
                    ></v-text-field>
                  </template>
                  <v-date-picker
                    v-model="range.end"
                    no-title
                    @input="menu.end = false"
                    @change="render"
                  ></v-date-picker>
                </v-menu>
              </v-col>
              <v-col class="col-lg-1">
                <v-select :items="rangeTypes" v-model="range.interval" @change="render"></v-select>
              </v-col>
              <v-col>
                <v-btn @click="render" icon>
                  <v-icon>mdi-refresh</v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col class="col-lg-4">
        <v-card>
          <v-card-title>Events</v-card-title>
          <v-card-text>
            <LineChart :chart-data="events" :fontColor="legendFontColor"></LineChart>
          </v-card-text>
          <v-card-text>
            <v-simple-table>
              <template>
                <thead>
                  <tr>
                    <th class="text-left">Event</th>
                    <th class="text-right">Count</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(value, key) in totalEvents" :key="key">
                    <td>{{ capitalizeFirstLetter(key) }}</td>
                    <td class="text-right">{{ value }}</td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col class="col-lg-4">
        <v-card min-height="100%">
          <v-card-title>Rates</v-card-title>
          <v-card-text>
            <DoughnutChart :chart-data="ratesChart" :fontColor="legendFontColor"></DoughnutChart>
          </v-card-text>
          <v-card-text>
            <v-simple-table>
              <template>
                <thead>
                  <tr>
                    <th class="text-left">Rate</th>
                    <th class="text-right">%</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Shown / Delivered</td>
                    <td class="text-right">
                      {{ new Intl.NumberFormat().format((rates.shown / rates.delivered) * 100) }}
                    </td>
                  </tr>
                  <tr>
                    <td>Hidden / Shown</td>
                    <td class="text-right">
                      {{ new Intl.NumberFormat().format((rates.hidden / rates.shown) * 100) }}
                    </td>
                  </tr>
                  <tr>
                    <td>Declined / Shown</td>
                    <td class="text-right">
                      {{ new Intl.NumberFormat().format((rates.declined / rates.shown) * 100) }}
                    </td>
                  </tr>
                  <tr>
                    <td>Accepted / Shown</td>
                    <td class="text-right">
                      {{ new Intl.NumberFormat().format((rates.accepted / rates.shown) * 100) }}
                    </td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col class="col-lg-4">
        <v-card min-height="100%">
          <v-card-title>Reach</v-card-title>
          <v-card-text>
            <BarChart :chart-data="reachChart" :fontColor="legendFontColor"></BarChart>
          </v-card-text>
          <v-card-text>
            <v-simple-table>
              <template>
                <thead>
                  <tr>
                    <th class="text-left">Channel</th>
                    <th class="text-right">Events</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(value, key) of reach" :key="key">
                    <td>{{ capitalizeFirstLetter(value.channel) }}</td>
                    <td class="text-right">
                      {{ new Intl.NumberFormat().format(value.total) }}
                    </td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import actionsClient from '../lib/rest/actions.js'
import analyzeClient from '../lib/rest/analyze.js'
import APISelect from './common/APISelect.vue'
import LineChart from './analyze/LineChart.vue'
import DoughnutChart from './analyze/DoughnutChart.vue'
import BarChart from './analyze/BarChart.vue'
import Help from './Help'
import moment from 'moment'

const humanLabels = {
  months: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  weekDays: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
}

export default {
  components: { APISelect, LineChart, DoughnutChart, BarChart, Help },
  props: {
    darkmode: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      actionId: null,
      rangeTypes: [
        { text: 'Year', value: 'year' },
        { text: 'Month', value: 'month' },
        { text: 'Day', value: 'day' },
        { text: 'Week Day', value: 'week_day' },
        { text: 'Day', value: 'day' },
        { text: 'Hour', value: 'hour' }
      ],
      menu: {
        start: false,
        end: false
      },
      range: {
        start: moment().startOf('month').format('YYYY-MM-DD'),
        end: moment().endOf('month').format('YYYY-MM-DD'),
        interval: 'day'
      },
      events: {},
      rates: {},
      reach: [],
      ratesChart: {},
      reachChart: {},
      totalEvents: {
        delivered: 0,
        shown: 0,
        hidden: 0,
        declined: 0,
        accepted: 0
      },
      averageInteractionTime: {
        show: 0.0,
        close: 0.0,
        dismiss: 0.0,
        accept: 0.0
      },
      legendFontColor: 'rgba(255,255,255,0.7)'
    }
  },
  watch: {
    range() {
      this.render()
    },
    darkmode() {
      this.updateChartColor()
    }
  },
  methods: {
    async loadActions() {
      const res = await actionsClient.list()
      return res.map((e) => {
        return { text: e.name, value: e.id }
      })
    },
    async render() {
      if (!this.actionId) {
        return
      }

      let start = moment(this.range.start).startOf('day').utc().format('YYYY-MM-DD HH:mm:ss')
      let end = moment(this.range.end).endOf('day').utc().format('YYYY-MM-DD HH:mm:ss')
      let range = {
        start: start,
        end: end,
        interval: this.range.interval
      }
      try {
        const events = await analyzeClient.events(this.actionId, range)
        this.rates = await analyzeClient.rates(this.actionId, range)
        this.reach = await analyzeClient.reach(this.actionId, range)

        this.prepareEvents(events)
        this.prepareRates(this.rates)
        this.prepareReach(this.reach)
      } catch (e) {
        this.$notify.info('No data available for the selected filters.')
      }
    },
    updateChartColor() {
      this.legendFontColor = this.darkmode ? 'rgba(255,255,255,0.7)' : 'rgba(0,0,0,0.6)'
    },
    prepareEvents(events) {
      let datasets = {
        delivered: {
          label: 'Delivered',
          borderColor: 'rgba(235, 235, 235, 0.3)',
          fill: false,
          data: []
        },
        shown: {
          label: 'Shown',
          borderColor: 'rgba(52, 171, 235, 0.3)',
          fill: false,
          data: []
        },
        hidden: {
          label: 'Hidden',
          borderColor: 'rgba(232, 235, 52, 0.3)',
          fill: false,
          data: []
        },
        declined: {
          label: 'Declined',
          borderColor: 'rgba(255, 82, 82, 0.3)',
          fill: false,
          data: []
        },
        accepted: {
          label: 'Accepted',
          borderColor: 'rgba(76, 175, 80, 0.3)',
          fill: false,
          data: []
        }
      }

      this.totalEvents = {
        delivered: 0,
        shown: 0,
        hidden: 0,
        declined: 0,
        accepted: 0
      }
      for (const e of events) {
        this.totalEvents[e.event] += e.total
      }

      this.events = this.prepareData(events, datasets, 'total')
    },
    prepareRates(rates) {
      this.ratesChart = {
        labels: ['Delivered', 'Shown', 'Hidden', 'Declined', 'Accepted'],
        datasets: [
          {
            backgroundColor: [
              'rgba(235, 235, 235, 0.3)',
              'rgba(52, 171, 235, 0.3)',
              'rgba(232, 235, 52, 0.3)',
              'rgba(255, 82, 82, 0.3)',
              'rgba(76, 175, 80, 0.3)'
            ],
            borderColor: 'rgba(255, 255, 255, 0.3)',
            data: [rates.delivered, rates.shown, rates.hidden, rates.declined, rates.accepted]
          }
        ]
      }
    },
    prepareReach(reach) {
      let labels = []
      let data = []
      let backgroundColors = []

      for (const e of reach) {
        labels.push(e.channel)
        data.push(e.total)
        backgroundColors.push('hsla(' + 360 * Math.random() + ', 50%, 50%, 0.3)')
      }

      this.reachChart = {
        labels: labels,
        datasets: [
          {
            label: 'Channels',
            backgroundColor: backgroundColors,
            data: data
          }
        ]
      }
    },
    prepareData(data, datasets, key) {
      const intervals = new Set()
      let interval = this.range.interval
      if (interval === 'week_day') {
        interval = 'weekDay'
      }

      for (const item of data) {
        intervals.add(item[interval])
      }

      for (let i of intervals) {
        for (let n of ['delivered', 'shown', 'hidden', 'declined', 'accepted']) {
          let found = false
          for (const item of data) {
            if (item.event === n && item[interval] === i) {
              datasets[item.event].data.push(item[key])
              found = true
            }
          }
          if (!found) {
            datasets[n].data.push(0)
          }
        }
      }

      let labels = Array.from(intervals)
      if (interval === 'hour') {
        labels = labels.map((val) => `${val}h`)
      }

      if (interval === 'month') {
        labels = labels.map((val) => humanLabels.months[val - 1])
      }

      if (interval === 'weekDay') {
        labels = labels.map((val) => humanLabels.weekDays[val - 1])
      }

      return {
        labels: labels,
        datasets: Object.values(datasets)
      }
    },
    capitalizeFirstLetter(string) {
      return string.charAt(0).toUpperCase() + string.slice(1)
    },
    hashCode(input) {
      var hash = 0,
        i,
        chr,
        len
      if (input.length === 0) return hash
      for (i = 0, len = input.length; i < len; i++) {
        chr = input.charCodeAt(i)
        hash = (hash << 5) - hash + chr
        hash |= 0
      }
      return hash
    }
  },
  async created() {
    this.actionId = this.$route.params.id
    if (this.actionId) {
      this.render()
      this.updateChartColor()
    }
  }
}
</script>
