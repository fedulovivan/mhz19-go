### Prometheus
  - Google: Monitoring Distributed Systems https://sre.google/sre-book/monitoring-distributed-systems/
  - Golang Application monitoring using Prometheus https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
  - Instrumenting A Go Application For Prometheus https://prometheus.io/docs/guides/go-application/ 
  - Promethues оповещение через Telegram https://blog.yakunin.dev/promethues-%D0%BE%D0%BF%D0%BE%D0%B2%D0%B5%D1%89%D0%B5%D0%BD%D0%B8%D0%B5-%D1%87%D0%B5%D1%80%D0%B5%D0%B7-telegram/
  - Основы мониторинга (обзор Prometheus и Grafana) https://habr.com/ru/articles/709204/
  - Creating Custom Prometheus Metrics in Golang and Sending Alerts to Slack https://medium.com/@mertcakmak2/custom-prometheus-metrics-in-golang-and-send-alert-to-slack-with-grafana-99a27dffe430
  - Визуализируем данные Node JS приложения с помощью Prometheus + Grafana https://habr.com/ru/articles/492742/
  - Prometheus + Grafana + Alertmanager в Docker https://www.dmosk.ru/miniinstruktions.php?mini=prometheus-stack-docker
  - How to visualize Prometheus histograms in Grafana https://grafana.com/blog/2020/06/23/how-to-visualize-prometheus-histograms-in-grafana/
  - Connect to host from docker container https://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach/24326540#24326540
  - Prometheus и PromQL основы сбора метрик https://www.youtube.com/watch?v=WQmpeOvCCUY
  - Have You Been Using Histogram Metrics Correctly https://medium.com/mercari-engineering/have-you-been-using-histogram-metrics-correctly-730c9547a7a9
  - How to create metrics with Go client API https://medium.com/@kylelzk/prometheus-practical-lab-how-to-create-metrics-with-go-client-api-b4119c4f8755
  - What is Vectors https://satyanash.net/software/2021/01/04/understanding-prometheus-range-vectors.html
  - Deleting Series https://medium.com/@burakceviz97/prometheus-metric-deletion-guide-8866bc5434ff
  - Writing An Exporter Or Custom Collector https://prometheus.io/docs/instrumenting/writing_exporters

### TODO
  - try: https://prometheus.io/docs/guides/node-exporter/ for macmini host
    - https://krsnachalise.medium.com/installing-node-exporter-in-linux-machines-d85e81d8808d

### Prometheus API
- http://localhost:9092/api/v1/labels
- http://localhost:9092/api/v1/label/rule_name/values
- curl -XPOST 'http://localhost:9092/api/v1/admin/tsdb/delete_series?match[]=mhz19_errors'
- curl -XPOST -g 'http://localhost:9092/api/v1/admin/tsdb/delete_series?match[]={instance="host.docker.internal:7070"}'

### Grafana API
- curl http://localhost:3002/api/dashboards/uid/ee0kwndjlf6kge
