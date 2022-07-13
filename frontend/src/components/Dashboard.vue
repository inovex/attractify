<template>
  <v-container>
    <v-row>
      <v-col class="col-lg-6">
        <v-card>
          <v-card-text>
            <div class="title font-weight-light mb-2">Reactions last 24h</div>
          </v-card-text>

          <v-sheet class="chart" color="white">
            <LineChart :chart-data="reactions" :options="chartOptions"></LineChart>
          </v-sheet>

          <v-card-text>
            <v-icon class="mr-2" small>mdi-information</v-icon>
            <span class="caption grey--text font-weight-light">
              This includes
              <em>show, hide, decline</em> and <em>accept</em> actions.
            </span>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col class="col-lg-6">
        <v-card>
          <v-card-text>
            <div class="title font-weight-light mb-2">New users last 24h</div>
          </v-card-text>

          <v-sheet class="chart" color="white">
            <LineChart :chart-data="profiles" :options="chartOptions"></LineChart>
          </v-sheet>

          <v-card-text>
            <v-icon class="mr-2" small>mdi-information</v-icon>
            <span class="caption grey--text font-weight-light">Users that are completely new.</span>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="mt-4">
      <v-col>
        <v-card>
          <v-card-text>
            <div class="title font-weight-light mb-2">Tracked Events last 24h</div>
          </v-card-text>

          <v-sheet class="pl-4 pt-4 pb-4">
            <div class="white--text">
              <span class="text-h3 count">{{ eventCount }}</span>
              &nbsp;
              <span class="count">Event(s)</span>
            </div>
          </v-sheet>

          <v-card-text>
            <v-icon class="mr-2" small>mdi-information</v-icon>
            <span class="caption grey--text font-weight-light"
              >Total number of new events received in the last 24h.</span
            >
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card>
          <v-card-text>
            <div class="title font-weight-light mb-2">Currently Active Actions</div>
          </v-card-text>

          <v-sheet class="pl-4 pt-4 pb-4">
            <div class="white--text">
              <span class="text-h3 count">{{ actionCount }}</span>
              &nbsp;
              <span class="count">Action(s)</span>
            </div>
          </v-sheet>

          <v-card-text>
            <v-icon class="mr-2" small>mdi-information</v-icon>
            <span class="caption grey--text font-weight-light">Actions that are currently active.</span>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import LineChart from './analyze/LineChart.vue'
import dashboard from '../lib/rest/dashboard'
import moment from 'moment'

var style = getComputedStyle(document.body);
var chartColor = style.getPropertyValue('--v-primary-base');

export default {
  components: { LineChart },
  data: () => ({
    reactions: {},
    profiles: {},
    actionCount: 0,
    eventCount: 0,
    chartOptions: {
      legend: { display: false },
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        xAxes: [
          {
            gridLines: {
              display: false
            },
            ticks: {
              fontColor: chartColor
            }
          }
        ],
        yAxes: [
          {
            gridLines: {
              display: true
            },
            ticks: {
              display: true,
              beginAtZero: true
            }
          }
        ]
      }
    }
  }),
  async created() {
    try {
      let res = await dashboard.load()
      if (res !== null) {
        this.reactions = {
          labels: res.reactions.map(r => {
            return moment(r.bucket).format('ha')
          }),
          datasets: [
            {
              borderColor: chartColor,
              fill: false,
              data: res.reactions.map(r => {
                return r.count
              })
            }
          ]
        }
        this.profiles = {
          labels: res.profiles.map(p => {
            return moment(p.bucket).format('ha')
          }),
          datasets: [
            {
              borderColor: chartColor,
              fill: false,
              data: res.profiles.map(r => {
                return r.count
              })
            }
          ]
        }

        this.profileValues = res.profiles.map(p => {
          return p.count
        })
        this.eventCount = res.events
        this.actionCount = res.actions
      }
    } catch (e) {
      this.$notify.error('Could not load dashboard.')
    }
  }
}
</script>

<style>
.v-sheet--offset {
  top: -24px;
  position: relative;
}

.count {
  color: var(--v-primary-base);
}
.chart {
  max-height: 400px;
}
</style>
