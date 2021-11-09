<script>
import { Bar, mixins } from 'vue-chartjs'

export default {
  extends: Bar,
  props: ['options'],
  mixins: [mixins.reactiveProp],
  methods: {
    render() {
      this.renderChart(this.chartData, {
        responsive: true,
        maintainAspectRatio: false
        // scales: {
        //   yAxes: [
        //     {
        //       //              type: 'logarithmic'
        //     }
        //   ]
        // }
      })
    }
  },
  mounted() {
    this.render()
  }
}
</script>
