<script>
import { Bar, mixins } from 'vue-chartjs'

export default {
  extends: Bar,
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
  },
  watch: {
    fontColor() {
      this.render()
    }
  }
}
</script>
