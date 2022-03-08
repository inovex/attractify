<script>
import { Doughnut, mixins } from "vue-chart-3";

export default {
  extends: Doughnut,
  props: ["options"],
  mixins: [mixins.reactiveProp],
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
