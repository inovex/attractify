<script>
import { Line, mixins } from 'vue-chartjs'

export default {
  extends: Line,
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
