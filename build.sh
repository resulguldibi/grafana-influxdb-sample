docker build -t resulguldibi/grafana-influxdb-sample .
docker-compose up -d
sleep 5
curl -H "Content-Type: application/json" -H "Authorization: Basic YWRtaW46YWRtaW4=" -d "@gateway-service-metrics-datasource.json" -X POST  http://localhost:3000/api/datasources
curl -H "Content-Type: application/json" -H "Authorization: Basic YWRtaW46YWRtaW4=" -d "@gateway-service-metrics-dashboard.json" -X POST  http://localhost:3000/api/dashboards/db/


#status_code=$(curl --write-out %{http_code} --silent --output /dev/null -H "Content-Type: application/json" -H "Authorization: Basic YWRtaW46YWRtaW4=" -d "@gateway-service-metrics-datasource.json" -X POST  http://localhost:3000/api/datasources)

#if [[ "$status_code" -ne 200 ]] ; then
#  echo "Site status changed to $status_code"
#else
#  exit 0
#fi
