package collector

import (
	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"log"
	"time"
)

func RetrieveHitRate(projectId string, rangeTime int64, delayTime int64, client *ucdn.UCDNClient) (response *ucdn.GetUcdnDomainHitRateResponse) {
	req := client.NewGetUcdnDomainHitRateRequest()
	req.ProjectId = ucloud.String(projectId)
	req.Type = ucloud.Int(0)
	req.BeginTime = ucloud.Int(int(time.Now().Unix() - rangeTime))
	req.EndTime = ucloud.Int(int(time.Now().Unix() - delayTime))
	req.Areacode = ucloud.String("cn")

	newUCdnHitRate, err := client.GetUcdnDomainHitRate(req)
	if err != nil {
		log.Println(err)
	}
	return newUCdnHitRate
}

func RetrieveBandWidth(projectId string, rangeTime int64, delayTime int64, client *ucdn.UCDNClient) (response *ucdn.GetNewUcdnDomainBandwidthResponse) {
	req := client.NewGetNewUcdnDomainBandwidthRequest()
	req.ProjectId = ucloud.String(projectId)
	req.Type = ucloud.Int(3)
	req.BeginTime = ucloud.Int(int(time.Now().Unix() - rangeTime))
	req.EndTime = ucloud.Int(int(time.Now().Unix() - delayTime))
	req.Areacode = ucloud.String("cn")

	newUCdnBandWidth, err := client.GetNewUcdnDomainBandwidth(req)
	if err != nil {
		log.Println(err)
	}

	return newUCdnBandWidth
}

func RetrieveOriginHttpCode4xx(projectId string, rangeTime int64, delayTime int64, client *ucdn.UCDNClient) (response *ucdn.GetUcdnDomainHttpCodeV2Response) {
	req := client.NewGetUcdnDomainHttpCodeV2Request()
	req.ProjectId = ucloud.String(projectId)
	req.Type = ucloud.Int(0)
	req.BeginTime = ucloud.Int(int(time.Now().Unix() - rangeTime))
	req.EndTime = ucloud.Int(int(time.Now().Unix() - delayTime))
	req.Areacode = ucloud.String("cn")

	newUCdnRequestNum, err := client.GetUcdnDomainHttpCodeV2(req)

	if err != nil {
		log.Println(err)
	}

	return newUCdnRequestNum
}

func RetrieveOriginHttpCode5xx(projectId string, rangeTime int64, delayTime int64, client *ucdn.UCDNClient) (response *ucdn.GetUcdnDomainHttpCodeV2Response) {
	req := client.NewGetUcdnDomainHttpCodeV2Request()
	req.ProjectId = ucloud.String(projectId)
	req.Type = ucloud.Int(0)
	req.BeginTime = ucloud.Int(int(time.Now().Unix() - rangeTime))
	req.EndTime = ucloud.Int(int(time.Now().Unix() - delayTime))
	req.Areacode = ucloud.String("cn")

	newUCdnRequestNum, err := client.GetUcdnDomainHttpCodeV2(req)

	if err != nil {
		log.Println(err)
	}

	return newUCdnRequestNum
}

func Retrieve95BandWidth(projectId string, rangeTime int64, delayTime int64, client *ucdn.UCDNClient) (response *ucdn.GetUcdnDomain95BandwidthV2Response) {
	req := client.NewGetUcdnDomain95BandwidthV2Request()
	req.ProjectId = ucloud.String(projectId)
	req.BeginTime = ucloud.Int(int(time.Now().Unix() - rangeTime))
	req.EndTime = ucloud.Int(int(time.Now().Unix() - delayTime))
	req.Areacode = ucloud.String("cn")

	newUCdn95BandWidth, err := client.GetUcdnDomain95BandwidthV2(req)

	if err != nil {
		log.Println(err)
	}

	return newUCdn95BandWidth
}

func RetrieveInfoList(projectId string, client *ucdn.UCDNClient) (response *ucdn.GetUcdnDomainInfoListResponse) {
	req := client.NewGetUcdnDomainInfoListRequest()
	req.ProjectId = ucloud.String(projectId)

	newUCdnInfoList, err := client.GetUcdnDomainInfoList(req)

	if err != nil {
		log.Println(err)
	}

	return newUCdnInfoList
}
