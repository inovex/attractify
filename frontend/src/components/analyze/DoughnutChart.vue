<script>
import { Doughnut, mixins } from 'vue-chartjs'

export default {
  extends: Doughnut,
  props: ['options'],
  mixins: [mixins.reactiveProp],
  methods: {
    render() {
      this.renderChart(this.chartData, {
        responsive: true,
        maintainAspectRatio: false
      })
    }
  },
  mounted() {
    this.render()
  }
}
</script>
