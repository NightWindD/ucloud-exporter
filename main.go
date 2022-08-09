package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test1/collector"
	"test1/exporter"
	"time"
)

func main() {
	publicKey := flag.String("pubkey", os.Getenv("UCLOUD_PUB_KEY"), "ucloud api public key")
	privateKey := flag.String("privatekey", os.Getenv("UCLOUD_PRIVATE_KEY"), "ucloud api private key")
	projectid := flag.String("projectId", os.Getenv("ucloud_project_id"), "confirm the project")
	host := flag.String("host", "0.0.0.0", "服务监听地址")
	port := flag.Int("port", 9200, "服务监听端口")
	service := flag.String("service", "acs_cdn", "输出Metrics的服务，默认为全部")
	rangeTime := flag.Int64("rangeTime", 3000, "rangeTime")
	delayTime := flag.Int64("delayTime", 60, "delayTime")
	tickerTime := flag.Int("tickerTime", 10, "tickerTime")
	flag.Parse()
	serviceArr := strings.Split(*service, ",")
	cfg := ucloud.NewConfig()
	cfg.BaseUrl = "https://api.ucloud.cn"

	cred := auth.NewCredential()
	cred.PublicKey = *publicKey
	cred.PrivateKey = *privateKey
	//var projectId string
	projectId := *projectid

	uCdnClient := ucdn.NewClient(&cfg, &cred)
	var infoCount *int
	var infoList [20]*string
	infoCount = &collector.RetrieveInfoList(projectId, uCdnClient).TotalCount
	for i := 0; i < *infoCount; i++ {
		infoList[i] = &collector.RetrieveInfoList(projectId, uCdnClient).DomainInfoList[i].Domain
	}

	ticker := time.NewTicker(time.Duration(*tickerTime) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				infoCount = &collector.RetrieveInfoList(projectId, uCdnClient).TotalCount
				for i := 0; i < *infoCount; i++ {
					infoList[i] = &collector.RetrieveInfoList(projectId, uCdnClient).DomainInfoList[i].Domain
				}
				log.Println(*infoCount, time.Now())
			}
		}
	}()
	*infoCount = 3
	for _, ae := range serviceArr {
		switch ae {
		case "acs_cdn":
			cdn := exporter.CdnCloudExporter(infoCount, infoList, projectId, *rangeTime, *delayTime, uCdnClient)
			prometheus.MustRegister(cdn)
		default:
			log.Println("暂不支持该服务，请根据提示选择服务。")
		}
	}
	listenAddress := net.JoinHostPort(*host, strconv.Itoa(*port))
	log.Println(listenAddress)
	log.Println("Running on", listenAddress)
	http.Handle("/metrics", promhttp.Handler()) //注册
	log.Fatal(http.ListenAndServe(listenAddress, nil))

}
