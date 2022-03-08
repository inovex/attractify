<script>
import { LineChart /*, mixins*/ } from "vue-chart-3";

export default {
  extends: LineChart,
  props: ["options"],
  /*mixins: [mixins.reactiveProp],*/
  methods: {
    render() {
      this.renderChart(this.chartData, {
        responsive: true,
        maintainAspectRatio: false,
      });
    },
  },
  mounted() {
    this.render();
  },
};
</script>
