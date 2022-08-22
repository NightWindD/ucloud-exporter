package exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"log"
	"strconv"
	"ucloud-exporter/collector"
)

const cdnNameSpace = "uCloud"

type CdnExporter struct {
	client               *ucdn.UCDNClient
	domainList           *[]ucdn.DomainBaseInfo
	rangeTime            int64
	delayTime            int64
	projectId            string
	cdnRequestHitRate    *prometheus.Desc
	cdnFlowHitRate       *prometheus.Desc
	cdnBandWidth         *prometheus.Desc
	cdnOriginHttpCode4xx *prometheus.Desc
	cdnOriginHttpCode5xx *prometheus.Desc
	cdn95bandwidth       *prometheus.Desc
	cdnResourceRequest   *prometheus.Desc

}

func CdnCloudExporter(domainList *[]ucdn.DomainBaseInfo, projectId string, rangeTime int64, delayTime int64, c *ucdn.UCDNClient) *CdnExporter {
	return &CdnExporter{
		client:     c,
		domainList: domainList,
		rangeTime:  rangeTime,
		delayTime:  delayTime,
		projectId:  projectId,
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

		cdnResourceRequest: prometheus.NewDesc(
			prometheus.BuildFQName(cdnNameSpace, "cdn", "resource_request"),
			"cdn回源请求数",
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
	ch <- e.cdnResourceRequest

}

func (e *CdnExporter) Collect(ch chan<- prometheus.Metric) {

	for _, domain := range *e.domainList {

		var requestHitRateSum float64
		var flowHitRateSum float64
		var bandWidthSum float64
		var bandWidthAverage float64

		var resourceCdnRequestSum float64
		var resourceCdnRequestAverage float64
		var http1xxSum int
		var http2xxSum int
		var http3xxSum int
		var http4xxSum int
		var http5xxSum int
		var http4xxAverage int
		var http5xxAverage int

		hitRateData := collector.RetrieveHitRate(domain.DomainId, e.projectId, e.rangeTime, e.delayTime, e.client).HitRateList
		for _, point := range hitRateData {
			flowHitRateSum += point.FlowHitRate
			requestHitRateSum += point.RequestHitRate
		}
		flowHitRateAverage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", flowHitRateSum/float64(len(hitRateData))), 64)
		if err != nil {
			log.Fatal(err)
		}
		requestHitRateAverage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", requestHitRateSum/float64(len(hitRateData))), 64)
		if err != nil {
			log.Fatal(err)
		}
		bandWidthData := collector.RetrieveBandWidth(domain.DomainId, e.projectId, e.rangeTime, e.delayTime, e.client).BandwidthList
		for _, point := range bandWidthData {
			bandWidthSum += point.CdnBandwidth
		}
		bandWidthAverage, err = strconv.ParseFloat(fmt.Sprintf("%.2f", bandWidthSum/float64(len(bandWidthData))), 64)
		if err != nil {
			log.Fatal(err)
		}

		httpData := collector.RetrieveOriginHttpCode4xx(domain.DomainId, e.projectId, e.rangeTime, e.delayTime, e.client).HttpCodeDetail
		for _, point := range httpData {

			http1xxSum += point.Http1XX.Total
			http2xxSum += point.Http2XX.Total
			http3xxSum += point.Http3XX.Total
			http4xxSum += point.Http4XX.Total
			http5xxSum += point.Http5XX.Total
		}

		http4xxAverage = http4xxSum / len(httpData)
		http5xxAverage = http5xxSum / len(httpData)

		resourceCdnRequestData := collector.RetrieveDomainOriginRequestNum(domain.DomainId, e.projectId, e.rangeTime, e.delayTime, e.client).RequestList
		for _, point := range resourceCdnRequestData {
			resourceCdnRequestSum += point.CdnRequest
		}

		resourceCdnRequestAverage, err = strconv.ParseFloat(fmt.Sprintf("%.2f", resourceCdnRequestSum/float64(len(resourceCdnRequestData))), 64)
		if err != nil {
			log.Fatal(err)
		}

		ch <- prometheus.MustNewConstMetric(
			e.cdnRequestHitRate,
			prometheus.GaugeValue,
			requestHitRateAverage,
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdnBandWidth,
			prometheus.GaugeValue,
			bandWidthAverage,
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdnOriginHttpCode4xx,
			prometheus.GaugeValue,
			float64(http4xxAverage),
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdnOriginHttpCode5xx,
			prometheus.GaugeValue,
			float64(http5xxAverage),
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdn95bandwidth,
			prometheus.GaugeValue,
			collector.Retrieve95BandWidth(domain.DomainId, e.projectId, e.rangeTime, e.delayTime, e.client).CdnBandwidth,
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdnFlowHitRate,
			prometheus.GaugeValue,
			flowHitRateAverage,
			domain.Domain,
		)

		ch <- prometheus.MustNewConstMetric(
			e.cdnResourceRequest,
			prometheus.GaugeValue,
			resourceCdnRequestAverage,
			domain.Domain,
		)
	}

}

