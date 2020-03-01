FROM golang:latest

WORKDIR /go/src/github.com/resulguldibi/grafana-influxdb-sample
ADD . /go/src/github.com/resulguldibi/grafana-influxdb-sample

RUN GOPATH=/go

RUN go build -o main .

EXPOSE 8000

ENTRYPOINT ["./main"]
