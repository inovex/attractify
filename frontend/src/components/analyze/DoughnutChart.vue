<script>
import { Doughnut, mixins } from 'vue-chartjs'

export default {
  extends: Doughnut,
  props: ['options', 'fontColor'],
  mixins: [mixins.reactiveProp],
  methods: {
    render() {
      this.renderChart(this.chartData, {
        responsive: true,
        maintainAspectRatio: false,
        legend: {
          labels: {
            fontColor: this.fontColor
          }
        }
      })
    }
  },
  mounted() {
    this.render()
  },
  watch: {
    fontColor() {
      this.render()
    }
  }
}
</script>
