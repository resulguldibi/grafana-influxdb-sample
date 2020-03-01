
--influx db ayağa kaldırmak için
docker run -d -p 8086:8086 -p 2003:2003 -v /var/lib/influxdb:/var/lib/influxdb -e INFLUXDB_GRAPHITE_ENABLED=true -e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=admin123 influxdb

--influxdb container'ı içerisinde https://s3.amazonaws.com/noaa.water-database/NOAA_data.txt adresinden indirilen ve /var/lib/influxdb path'ine kopyalanan dosya container içerisindeki volume de görülür ve import edilebilir.
root@f7be069aac2e:/var/lib/influxdb# influx -import -path=NOAA_data.txt -precision=s -database=NOAA_water_database

--import edilen datayı kontrol etmek için influx komutu çalıştırılır ve sonrasında "show databases" komutu çalıştırılarak ilgili veri tabanları görülebilir.
root@f7be069aac2e:/var/lib/influxdb# influx
Connected to http://localhost:8086 version 1.7.10
InfluxDB shell version: 1.7.10
> show databases;
name: databases
name
----
_internal
NOAA_water_database


--use NOAA_water_database komutu ile NOAA_water_database isimli veri tabanı seçilir.
--show measurements komutu ile veri tabanı içerisindeki tabloya karşılık olarak düşünülebilecek olan veriler görülebilir.

--grafana ayağa kaldırmak için

docker run -d -p 3000:3000 grafana/grafana


--influxdb de örnek veritabanı oluşturulması.

curl -i -XPOST http://localhost:8086/query --data-urlencode "q=CREATE DATABASE serviceMetricsDb"


--örnek metric datası

curl -i -XPOST 'http://localhost:8086/write?db=serviceMetricsDb' --data-binary 'request,id=d7c9613a-8d27-4391-90f3-966f92614b52,nodeId=3,nodeName=node3,nodeIp=127.0.0.3,serviceId=1,serviceName=service-1,serviceUri=/test-service-1,isPolicySuccessful=false,isPolicyViolation=false,isRoutingFailure=false totalFrontendLatency=11,totalBackendLatency=159'


request,id=${requestId},nodeId=${nodeId.result},nodeName=${nodeName.result},nodeIp=${nodeIp.result},serviceId=${serviceId.result},serviceName=${serviceName.result},serviceUri=${serviceUri.result},isPolicySuccessful=${isPolicySuccessful.result},isPolicyViolation=${isPolicyViolation.result},isRoutingFailure=${isRoutingFailure.result} totalFrontendLatency=${totalFrontendLatency.result},totalBackendLatency=${totalBackendLatency.result} ${time.result}000000


--grafana ya da influx db containerları restart olduklarında data kaybı olmaması ıcın ilgili verileri (influxdb database volume, grafana dashboard ve datasources  configleri) hostta bir dizine volumelemek gerekiyor.




--grafana container'ı ayaga kalktıktan sonra rest api isteği ile datasource set edilir.

POST /api/datasources HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Authorization: Basic YWRtaW46YWRtaW4xMjM=

{
    "orgId": 1,
    "name": "serviceMetrics",
    "type": "influxdb",
    "typeLogoUrl": "public/app/plugins/datasource/influxdb/img/influxdb_logo.svg",
    "access": "direct",
    "url": "http://localhost:8086",
    "password": "admin123",
    "user": "admin",
    "database": "serviceMetricsDb",
    "basicAuth": false,
    "isDefault": true,
    "jsonData": {
        "httpMode": "GET"
    },
    "readOnly": false
}


--grafana container'ı ayaga kalktıktan sonra rest api isteği ile dashboard set edilir.

POST /api/dashboards/db/ HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Authorization: Basic YWRtaW46YWRtaW4xMjM=

{
    "dashboard": {
        "__inputs": [
            {
                "name": "DS_SERVICEMETRICS",
                "label": "serviceMetrics",
                "description": "",
                "type": "datasource",
                "pluginId": "influxdb",
                "pluginName": "InfluxDB"
            }
        ],
        "__requires": [
            {
                "type": "grafana",
                "id": "grafana",
                "name": "Grafana",
                "version": "5.0.4"
            },
            {
                "type": "panel",
                "id": "graph",
                "name": "Graph",
                "version": "5.0.0"
            },
            {
                "type": "datasource",
                "id": "influxdb",
                "name": "InfluxDB",
                "version": "5.0.0"
            },
            {
                "type": "panel",
                "id": "table",
                "name": "Table",
                "version": "5.0.0"
            }
        ],
        "annotations": {
            "list": [
                {
                    "builtIn": 1,
                    "datasource": "-- Grafana --",
                    "enable": true,
                    "hide": true,
                    "iconColor": "rgba(0, 211, 255, 1)",
                    "name": "Annotations & Alerts",
                    "type": "dashboard"
                }
            ]
        },
        "editable": true,
        "gnetId": null,
        "graphTooltip": 0,
        "id": null,
        "iteration": 1524608661775,
        "links": [],
        "panels": [
            {
                "columns": [
                    {
                        "text": "Total",
                        "value": "total"
                    }
                ],
                "datasource": "serviceMetrics",
                "fontSize": "100%",
                "gridPos": {
                    "h": 5,
                    "w": 7,
                    "x": 0,
                    "y": 0
                },
                "id": 17,
                "links": [],
                "pageSize": null,
                "scroll": false,
                "showHeader": true,
                "sort": {
                    "col": 0,
                    "desc": false
                },
                "styles": [
                    {
                        "alias": "Time",
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 0,
                        "link": false,
                        "pattern": "Time",
                        "preserveFormat": false,
                        "type": "number"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 0,
                        "pattern": "Total",
                        "thresholds": [
                            ""
                        ],
                        "type": "number",
                        "unit": "none"
                    },
                    {
                        "alias": "Status",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 2,
                        "pattern": "Metric",
                        "thresholds": [],
                        "type": "string",
                        "unit": "short"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "decimals": 2,
                        "pattern": "/.*/",
                        "thresholds": [],
                        "type": "number",
                        "unit": "short"
                    }
                ],
                "targets": [
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Routing Failure"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isRoutingFailure",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "B",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Policy Violation"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicyViolation",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "C",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Success"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicySuccessful",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "D",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Total"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "timeFrom": "24h",
                "timeShift": null,
                "title": "Service Status",
                "transform": "timeseries_aggregations",
                "type": "table"
            },
            {
                "aliasColors": {},
                "bars": true,
                "dashLength": 10,
                "dashes": false,
                "datasource": "serviceMetrics",
                "decimals": 0,
                "fill": 1,
                "gridPos": {
                    "h": 6,
                    "w": 17,
                    "x": 7,
                    "y": 0
                },
                "id": 9,
                "legend": {
                    "alignAsTable": true,
                    "avg": false,
                    "current": false,
                    "hideEmpty": false,
                    "hideZero": false,
                    "max": true,
                    "min": true,
                    "rightSide": true,
                    "show": true,
                    "sideWidth": 200,
                    "total": false,
                    "values": true
                },
                "lines": false,
                "linewidth": 1,
                "links": [],
                "nullPointMode": "null",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "spaceLength": 10,
                "stack": false,
                "steppedLine": false,
                "targets": [
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "null"
                                ],
                                "type": "fill"
                            }
                        ],
                        "hide": false,
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "query": "SELECT \"value\" FROM \"totalFrontendLatency\" WHERE $timeFilter GROUP BY time(5s), \"isPolicySuccessful\", \"isPolicyViolation\", \"isRoutingFailure\" fill(null)",
                        "rawQuery": false,
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "mean"
                                },
                                {
                                    "params": [
                                        "Front End"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "null"
                                ],
                                "type": "fill"
                            }
                        ],
                        "hide": false,
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "query": "SELECT \"value\" FROM \"totalFrontendLatency\" WHERE $timeFilter GROUP BY time(5s), \"isPolicySuccessful\", \"isPolicyViolation\", \"isRoutingFailure\" fill(null)",
                        "rawQuery": false,
                        "refId": "B",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalBackendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "mean"
                                },
                                {
                                    "params": [
                                        "Back End"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "isPolicySuccessful",
                                "operator": "=",
                                "value": "true"
                            }
                        ]
                    }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Service Latency",
                "tooltip": {
                    "shared": true,
                    "sort": 0,
                    "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                    "buckets": null,
                    "mode": "time",
                    "name": null,
                    "show": true,
                    "values": []
                },
                "yaxes": [
                    {
                        "decimals": 0,
                        "format": "none",
                        "label": "Response Time (ms)",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": true
                    },
                    {
                        "decimals": 0,
                        "format": "none",
                        "label": "Response Time (ms)",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": false
                    }
                ]
            },
            {
                "columns": [
                    {
                        "text": "Total",
                        "value": "total"
                    }
                ],
                "datasource": "serviceMetrics",
                "fontSize": "100%",
                "gridPos": {
                    "h": 4,
                    "w": 7,
                    "x": 0,
                    "y": 5
                },
                "hideTimeOverride": false,
                "id": 20,
                "links": [],
                "pageSize": null,
                "scroll": true,
                "showHeader": true,
                "sort": {
                    "col": 0,
                    "desc": true
                },
                "styles": [
                    {
                        "alias": "Time",
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "pattern": "Time",
                        "type": "date"
                    },
                    {
                        "alias": "Service",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 2,
                        "pattern": "Metric",
                        "thresholds": [],
                        "type": "string",
                        "unit": "short"
                    },
                    {
                        "alias": "",
                        "colorMode": "value",
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 0,
                        "pattern": "Total",
                        "thresholds": [],
                        "type": "number",
                        "unit": "none"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "decimals": 2,
                        "pattern": "/.*/",
                        "thresholds": [],
                        "type": "number",
                        "unit": "short"
                    }
                ],
                "targets": [
                    {
                        "alias": "$tag_serviceName",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "serviceName"
                                ],
                                "type": "tag"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicyViolation",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "timeFrom": "24h",
                "title": "Policy Violation",
                "transform": "timeseries_aggregations",
                "type": "table"
            },
            {
                "aliasColors": {},
                "bars": false,
                "dashLength": 10,
                "dashes": false,
                "datasource": "serviceMetrics",
                "decimals": 2,
                "fill": 1,
                "gridPos": {
                    "h": 7,
                    "w": 17,
                    "x": 7,
                    "y": 6
                },
                "id": 10,
                "legend": {
                    "alignAsTable": true,
                    "avg": false,
                    "current": false,
                    "hideEmpty": false,
                    "hideZero": false,
                    "max": true,
                    "min": true,
                    "rightSide": true,
                    "show": true,
                    "sideWidth": 200,
                    "total": false,
                    "values": true
                },
                "lines": true,
                "linewidth": 1,
                "links": [],
                "nullPointMode": "null",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "spaceLength": 10,
                "stack": false,
                "steppedLine": false,
                "targets": [
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        " / 5"
                                    ],
                                    "type": "math"
                                },
                                {
                                    "params": [
                                        "Message Rate"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Message Rate",
                "tooltip": {
                    "shared": true,
                    "sort": 0,
                    "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                    "buckets": null,
                    "mode": "time",
                    "name": null,
                    "show": true,
                    "values": []
                },
                "yaxes": [
                    {
                        "decimals": 2,
                        "format": "none",
                        "label": "Message Rate (requests/sec)",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": true
                    },
                    {
                        "decimals": 2,
                        "format": "none",
                        "label": "Message Rate (requests/sec)",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": false
                    }
                ]
            },
            {
                "columns": [
                    {
                        "text": "Total",
                        "value": "total"
                    }
                ],
                "datasource": "serviceMetrics",
                "fontSize": "100%",
                "gridPos": {
                    "h": 4,
                    "w": 7,
                    "x": 0,
                    "y": 9
                },
                "id": 21,
                "links": [],
                "pageSize": null,
                "scroll": true,
                "showHeader": true,
                "sort": {
                    "col": 0,
                    "desc": true
                },
                "styles": [
                    {
                        "alias": "Time",
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "pattern": "Time",
                        "type": "date"
                    },
                    {
                        "alias": "Service",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 2,
                        "pattern": "Metric",
                        "thresholds": [],
                        "type": "string",
                        "unit": "short"
                    },
                    {
                        "alias": "",
                        "colorMode": "value",
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 0,
                        "pattern": "Total",
                        "thresholds": [],
                        "type": "number",
                        "unit": "none"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "decimals": 2,
                        "pattern": "/.*/",
                        "thresholds": [],
                        "type": "number",
                        "unit": "short"
                    }
                ],
                "targets": [
                    {
                        "alias": "$tag_serviceName",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "serviceName"
                                ],
                                "type": "tag"
                            },
                            {
                                "params": [
                                    "isPolicyViolation"
                                ],
                                "type": "tag"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isRoutingFailure",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "timeFrom": "24h",
                "title": "Routing Failure",
                "transform": "timeseries_aggregations",
                "type": "table"
            },
            {
                "columns": [
                    {
                        "text": "Total",
                        "value": "total"
                    }
                ],
                "datasource": "serviceMetrics",
                "fontSize": "100%",
                "gridPos": {
                    "h": 6,
                    "w": 7,
                    "x": 0,
                    "y": 13
                },
                "id": 19,
                "links": [],
                "pageSize": null,
                "scroll": true,
                "showHeader": true,
                "sort": {
                    "col": 0,
                    "desc": true
                },
                "styles": [
                    {
                        "alias": "Time",
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "pattern": "Time",
                        "type": "date"
                    },
                    {
                        "alias": "Service",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 2,
                        "pattern": "Metric",
                        "thresholds": [],
                        "type": "string",
                        "unit": "short"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "dateFormat": "YYYY-MM-DD HH:mm:ss",
                        "decimals": 0,
                        "pattern": "Total",
                        "thresholds": [],
                        "type": "number",
                        "unit": "none"
                    },
                    {
                        "alias": "",
                        "colorMode": null,
                        "colors": [
                            "rgba(245, 54, 54, 0.9)",
                            "rgba(237, 129, 40, 0.89)",
                            "rgba(50, 172, 45, 0.97)"
                        ],
                        "decimals": 2,
                        "pattern": "/.*/",
                        "thresholds": [],
                        "type": "number",
                        "unit": "short"
                    }
                ],
                "targets": [
                    {
                        "alias": "$tag_serviceName",
                        "groupBy": [
                            {
                                "params": [
                                    "$__interval"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "serviceName"
                                ],
                                "type": "tag"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalFrontendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicySuccessful",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "timeFrom": "24h",
                "title": "Success",
                "transform": "timeseries_aggregations",
                "type": "table"
            },
            {
                "aliasColors": {},
                "bars": true,
                "dashLength": 10,
                "dashes": false,
                "datasource": "serviceMetrics",
                "decimals": 0,
                "fill": 1,
                "gridPos": {
                    "h": 6,
                    "w": 17,
                    "x": 7,
                    "y": 13
                },
                "id": 4,
                "legend": {
                    "alignAsTable": true,
                    "avg": false,
                    "current": false,
                    "hideEmpty": false,
                    "hideZero": false,
                    "max": false,
                    "min": false,
                    "rightSide": true,
                    "show": true,
                    "sideWidth": 200,
                    "total": true,
                    "values": true
                },
                "lines": false,
                "linewidth": 1,
                "links": [],
                "nullPointMode": "null",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "spaceLength": 10,
                "stack": true,
                "steppedLine": false,
                "targets": [
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "null"
                                ],
                                "type": "fill"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "A",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalBackendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Routing Failure"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isRoutingFailure",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "null"
                                ],
                                "type": "fill"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "B",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalBackendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Policy Violation"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicyViolation",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    },
                    {
                        "alias": "$col",
                        "groupBy": [
                            {
                                "params": [
                                    "5s"
                                ],
                                "type": "time"
                            },
                            {
                                "params": [
                                    "null"
                                ],
                                "type": "fill"
                            }
                        ],
                        "measurement": "request",
                        "orderByTime": "ASC",
                        "policy": "default",
                        "refId": "C",
                        "resultFormat": "time_series",
                        "select": [
                            [
                                {
                                    "params": [
                                        "totalBackendLatency"
                                    ],
                                    "type": "field"
                                },
                                {
                                    "params": [],
                                    "type": "count"
                                },
                                {
                                    "params": [
                                        "Success"
                                    ],
                                    "type": "alias"
                                }
                            ]
                        ],
                        "tags": [
                            {
                                "key": "isPolicySuccessful",
                                "operator": "=",
                                "value": "true"
                            },
                            {
                                "condition": "AND",
                                "key": "nodeName",
                                "operator": "=~",
                                "value": "/^$nodeName$/"
                            },
                            {
                                "condition": "AND",
                                "key": "serviceName",
                                "operator": "=~",
                                "value": "/^$serviceName$/"
                            }
                        ]
                    }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Service Status",
                "tooltip": {
                    "shared": true,
                    "sort": 0,
                    "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                    "buckets": null,
                    "mode": "time",
                    "name": null,
                    "show": true,
                    "values": []
                },
                "yaxes": [
                    {
                        "decimals": 0,
                        "format": "none",
                        "label": "Number of Messages",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": true
                    },
                    {
                        "decimals": 0,
                        "format": "none",
                        "label": "Number of Messages",
                        "logBase": 1,
                        "max": null,
                        "min": "0",
                        "show": false
                    }
                ]
            }
        ],
        "refresh": "5s",
        "schemaVersion": 16,
        "style": "dark",
        "tags": [],
        "templating": {
            "list": [
                {
                    "allValue": null,
                    "current": {},
                    "datasource": "serviceMetrics",
                    "hide": 0,
                    "includeAll": true,
                    "label": "Gateway Node",
                    "multi": false,
                    "name": "nodeName",
                    "options": [],
                    "query": "SHOW TAG VALUES  WITH KEY = \"nodeName\" ",
                    "refresh": 1,
                    "regex": "",
                    "sort": 1,
                    "tagValuesQuery": "",
                    "tags": [],
                    "tagsQuery": "",
                    "type": "query",
                    "useTags": false
                },
                {
                    "allValue": null,
                    "current": {},
                    "datasource": "serviceMetrics",
                    "hide": 0,
                    "includeAll": true,
                    "label": "Published Service",
                    "multi": false,
                    "name": "serviceName",
                    "options": [],
                    "query": "SHOW TAG VALUES  WITH KEY = \"serviceName\" ",
                    "refresh": 1,
                    "regex": "",
                    "sort": 1,
                    "tagValuesQuery": "",
                    "tags": [],
                    "tagsQuery": "",
                    "type": "query",
                    "useTags": false
                }
            ]
        },
        "time": {
            "from": "now-15m",
            "to": "now"
        },
        "timepicker": {
            "refresh_intervals": [
                "5s",
                "10s",
                "30s",
                "1m",
                "5m",
                "15m",
                "30m",
                "1h",
                "2h",
                "1d"
            ],
            "time_options": [
                "5m",
                "15m",
                "1h",
                "6h",
                "12h",
                "24h",
                "2d",
                "7d",
                "30d"
            ]
        },
        "timezone": "",
        "title": "Gateway Service Metrics",
        "uid": "-0QQRAWiz",
        "version": 3
    },
    "overwrite": false,
    "message": ""
}
