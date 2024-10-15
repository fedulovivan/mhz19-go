### Prometheus
  - https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
  - https://prometheus.io/docs/guides/go-application/ 
  - https://blog.yakunin.dev/promethues-%D0%BE%D0%BF%D0%BE%D0%B2%D0%B5%D1%89%D0%B5%D0%BD%D0%B8%D0%B5-%D1%87%D0%B5%D1%80%D0%B5%D0%B7-telegram/
  - https://habr.com/ru/articles/709204/
  - https://medium.com/@mertcakmak2/custom-prometheus-metrics-in-golang-and-send-alert-to-slack-with-grafana-99a27dffe430
  - https://habr.com/ru/articles/492742/
  - https://www.dmosk.ru/miniinstruktions.php?mini=prometheus-stack-docker
  - https://grafana.com/blog/2020/06/23/how-to-visualize-prometheus-histograms-in-grafana/
  - https://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach/24326540#24326540
  - https://www.youtube.com/watch?v=WQmpeOvCCUY
  - https://medium.com/mercari-engineering/have-you-been-using-histogram-metrics-correctly-730c9547a7a9
  - https://medium.com/@kylelzk/prometheus-practical-lab-how-to-create-metrics-with-go-client-api-b4119c4f8755
  - What is Vectors https://satyanash.net/software/2021/01/04/understanding-prometheus-range-vectors.html
  - Deleting Metric https://medium.com/@burakceviz97/prometheus-metric-deletion-guide-8866bc5434ff

### Prometheus API
- http://localhost:9092/api/v1/labels
- http://localhost:9092/api/v1/label/rule_name/values
- curl -XPOST 'http://localhost:9092/api/v1/admin/tsdb/delete_series?match[]=mhz19_errors'
- curl -XPOST -g 'http://localhost:9092/api/v1/admin/tsdb/delete_series?match[]={instance="host.docker.internal:7070"'


### Grafana API
- curl http://localhost:3002/api/dashboards/uid/ee0kwndjlf6kge
