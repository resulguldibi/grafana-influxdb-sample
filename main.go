package main

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	//"time"

)

type Node struct {
	Id   int
	Name string
	Ip   string
}

type Service struct {
	Id   int
	Name string
	Uri  string
}

func main() {

	//produce sample metric data

	client := http.Client{}



	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		for {

			time.Sleep(time.Duration(time.Millisecond * 3000))

			randomData := getRandomMetricData()
			resp, err := client.Post("http://influxdb:8086/write?db=serviceMetricsDb","text/plain; charset=utf-8", strings.NewReader(randomData.(string)))

			if err != nil{
				fmt.Println("error -> ", err)
			}

			if resp != nil{

				if  resp.StatusCode != http.StatusNoContent{
					responseBody, err2 := ioutil.ReadAll(resp.Body)

					if err2 != nil{
						fmt.Println("err2 -> ", err2)
					}

					fmt.Println(string(responseBody))
				}

				resp.Body.Close()
			}
		}
	}()


	<-quit

	client.CloseIdleConnections()

}

func getRandomMetricData() interface{} {
	nodeList := make([]interface{}, 0)
	nodeList = append(nodeList, &Node{Id: 1, Name: "node1", Ip: "127.0.0.1"})
	nodeList = append(nodeList, &Node{Id: 2, Name: "node2", Ip: "127.0.0.2"})
	nodeList = append(nodeList, &Node{Id: 3, Name: "node3", Ip: "127.0.0.3"})

	serviceList := make([]interface{}, 0)
	serviceList = append(serviceList, &Service{Id: 1, Name: "service-1", Uri: "/test-service-1"})
	serviceList = append(serviceList, &Service{Id: 2, Name: "service-2", Uri: "/test-service-2"})
	serviceList = append(serviceList, &Service{Id: 3, Name: "service-3", Uri: "/test-service-3"})

	node := getRandomData(nodeList).(*Node)
	service := getRandomData(serviceList).(*Service)
	requestId := uuid.New().String()

	var status bool = false
	if generateRandomInt(0, 10)%2 == 0 {
		status = true
	}

	backendLatency := generateRandomInt(100, 200)
	gatewayLatency := generateRandomInt(100, 200)

	//return fmt.Sprintf(`request,id=%s,nodeId=%d,nodeName=%s,nodeIp=%s,serviceId=%d,serviceName=%s,serviceUri=%s,isPolicySuccessful=%t,isPolicyViolation=%t,isRoutingFailure=%t totalFrontendLatency=%d,totalBackendLatency=%d %d000000`, requestId, node.Id, node.Name, node.Ip, service.Id, service.Name, service.Uri, status, status, status, gatewayLatency, backendLatency, time.Now().Unix())

	return fmt.Sprintf(`request,id=%s,nodeId=%d,nodeName=%s,nodeIp=%s,serviceId=%d,serviceName=%s,serviceUri=%s,isPolicySuccessful=%t,isPolicyViolation=%t,isRoutingFailure=%t totalFrontendLatency=%d,totalBackendLatency=%d`, requestId, node.Id, node.Name, node.Ip, service.Id, service.Name, service.Uri, status, status, status, gatewayLatency, backendLatency)

}

func getRandomData(list []interface{}) interface{} {
	index := generateRandomInt(0, len(list))
	return list[index]
}

func generateRandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}
