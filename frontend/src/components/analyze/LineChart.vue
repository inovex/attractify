<script>
import { Line, mixins } from 'vue-chartjs'

export default {
  extends: Line,
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
