package exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"log"
	"strconv"
	"test1/collector"
)

const cdnNameSpace = "uCloud"

type CdnExporter struct {
	client               *ucdn.UCDNClient
	infoCount            *int
	infoList             [20]*string
	rangeTime            int64
	delayTime            int64
	projectId            string
	cdnRequestHitRate    *prometheus.Desc
	cdnFlowHitRate       *prometheus.Desc
	cdnBandWidth         *prometheus.Desc
	cdnOriginHttpCode4xx *prometheus.Desc
	cdnOriginHttpCode5xx *prometheus.Desc
	cdn95bandwidth       *prometheus.Desc
}

func CdnCloudExporter(infoCount *int, infoList [20]*string, projectId string, rangeTime int64, delayTime int64, c *ucdn.UCDNClient) *CdnExporter {
	return &CdnExporter{
		client:    c,
		infoList:  infoList,
		infoCount: infoCount,
		rangeTime: rangeTime,
		delayTime: delayTime,
		projectId: projectId,
		cdnRequestHitRate: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "request_hit_rate"),
			"总请求命中率(%)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnBandWidth: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "band_width"),
			"域名带宽(Mbps)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnOriginHttpCode4xx: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "http_code_4XX"),
			"http4XX请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnOriginHttpCode5xx: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "http_code_5XX"),
			"http5XX请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdn95bandwidth: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "95_band_width"),
			"95带宽数据(Mbps)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnFlowHitRate: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "flow_hit_rate"),
			"总流量命中率(%)",
			[]string{
				"instanceId",
			},
			nil,
		),
	}
}

func (e *CdnExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.cdnRequestHitRate
	ch <- e.cdnFlowHitRate
	ch <- e.cdnBandWidth
	ch <- e.cdnOriginHttpCode4xx
	ch <- e.cdn95bandwidth
	ch <- e.cdnOriginHttpCode5xx
}

func (e *CdnExporter) Collect(ch chan<- prometheus.Metric) {

	log.Println(*e.infoCount)
	*e.infoCount = 1
	var requestHitRateData float64
	var requestHitRateCount float64
	var requestHitRateAverage float64
	for _, point := range collector.RetrieveHitRate(e.projectId, e.rangeTime, e.delayTime, e.client).HitRateList {
		requestHitRateData += point.RequestHitRate
		requestHitRateCount++
		requestHitRateAverage = requestHitRateData / requestHitRateCount
		requestHitRateAverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", requestHitRateAverage), 64)
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdnRequestHitRate,
			prometheus.GaugeValue,
			requestHitRateAverage,
			*e.infoList[i],
		)
	}

	var bandWidthData float64
	var bandWidthCount float64
	var bandWidthAverage float64

	for _, point := range collector.RetrieveBandWidth(e.projectId, e.rangeTime, e.delayTime, e.client).BandwidthList {
		bandWidthData += point.CdnBandwidth
		bandWidthCount++
		bandWidthAverage = bandWidthData / bandWidthCount
		bandWidthAverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", bandWidthAverage), 64)
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdnBandWidth,
			prometheus.GaugeValue,
			bandWidthAverage,
			*e.infoList[i],
		)
	}

	var http4xxData int
	var http4xxCount int
	var http4xxAverage int

	for _, point := range collector.RetrieveOriginHttpCode4xx(e.projectId, e.rangeTime, e.delayTime, e.client).HttpCodeDetail {
		http4xxData += point.Http4XX.Total
		http4xxCount++
		http4xxAverage = http4xxData / http4xxCount
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdnOriginHttpCode4xx,
			prometheus.GaugeValue,
			float64(http4xxAverage),
			*e.infoList[i],
		)
	}

	var http5xxData int
	var http5xxCount int
	var http5xxAverage int

	for _, point := range collector.RetrieveOriginHttpCode5xx(e.projectId, e.rangeTime, e.delayTime, e.client).HttpCodeDetail {
		http5xxData += point.Http5XX.Total
		http5xxCount++
		http5xxAverage = http5xxData / http5xxCount
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdnOriginHttpCode5xx,
			prometheus.GaugeValue,
			float64(http5xxAverage),
			*e.infoList[i],
		)
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdn95bandwidth,
			prometheus.GaugeValue,
			collector.Retrieve95BandWidth(e.projectId, e.rangeTime, e.delayTime, e.client).CdnBandwidth,
			*e.infoList[i],
		)
	}

	var flowHitRateData float64
	var flowHitRateCount float64
	var flowHitRateAverage float64
	for _, point := range collector.RetrieveHitRate(e.projectId, e.rangeTime, e.delayTime, e.client).HitRateList {
		flowHitRateData += point.FlowHitRate
		flowHitRateCount++
		flowHitRateAverage = flowHitRateData / flowHitRateCount
		flowHitRateAverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", flowHitRateAverage), 64)
	}

	for i := 0; i < *e.infoCount; i++ {
		ch <- prometheus.MustNewConstMetric(
			e.cdnFlowHitRate,
			prometheus.GaugeValue,
			flowHitRateAverage,
			*e.infoList[i],
		)

	}

	//	*e.infoCount = 10
	//	log.Println(*e.infoCount)
}
