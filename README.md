for i in `gcloud compute instances list | awk '$0=$6'`; do echo "'$i:2112'",; done

docker run -p 9090:9090 -d --volume prometheus-data:/prometheus-data -v "$(pwd)"/prometheus.yaml:/etc/prometheus/prometheus.yml prom/prometheus

docker run -d --network host --name grafana grafana/grafana:6.5.0
