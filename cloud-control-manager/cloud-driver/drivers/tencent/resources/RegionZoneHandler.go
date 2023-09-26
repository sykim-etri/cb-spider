package resources

import (
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentRegionZoneHandler struct {
	Region idrv.RegionInfo
	Client *cvm.Client
}

func (regionZoneHandler *TencentRegionZoneHandler) ListRegionZone() ([]*irs.RegionZoneInfo, error) {

	responseRegions, err := DescribeRegions(regionZoneHandler.Client)
	if err != nil {
		cblogger.Error("err : DescribeRegions")
		cblogger.Error(err)
		return nil, err
	}

	clientProfile := profile.NewClientProfile()
	clientProfile.Language = "en-US" // lang default set is zh-CN -> set as en-US.

	var regionZoneInfoList []*irs.RegionZoneInfo
	for _, region := range responseRegions.Response.RegionSet {
		tempClient, _ := cvm.NewClient(regionZoneHandler.Client.Client.GetCredential(), *region.Region, clientProfile)
		responseZones, _ := DescribeZones(tempClient)

		var zoneInfoList []irs.ZoneInfo
		for _, zone := range responseZones.Response.ZoneSet {
			zoneInfo := irs.ZoneInfo{}
			zoneInfo.Name = *zone.Zone
			zoneInfo.DisplayName = *zone.ZoneName
			zoneInfo.Status = GetZoneStatus(*zone.ZoneState)

			// keyValueList 삭제 https://github.com/cloud-barista/cb-spider/issues/930#issuecomment-1734817828
			// keyValueList, err := ConvertKeyValueList(zone)
			// if err != nil {
			// 	cblogger.Errorf("err : ConvertKeyValueList [%s]", *zone.ZoneName)
			// 	cblogger.Error(err)
			// 	keyValueList = nil
			// }
			// zoneInfo.KeyValueList = keyValueList

			zoneInfoList = append(zoneInfoList, zoneInfo)
		}

		regionInfo := irs.RegionZoneInfo{}
		regionInfo.Name = *region.Region
		regionInfo.DisplayName = *region.RegionName
		regionInfo.ZoneList = zoneInfoList

		// keyValueList 삭제 https://github.com/cloud-barista/cb-spider/issues/930#issuecomment-1734817828
		// keyValueList, err := ConvertKeyValueList(region)
		// if err != nil {
		// 	cblogger.Errorf("err : ConvertKeyValueList [%s]", *region.Region)
		// 	cblogger.Error(err)
		// 	keyValueList = nil
		// }
		// regionInfo.KeyValueList = keyValueList

		regionZoneInfoList = append(regionZoneInfoList, &regionInfo)
	}

	return regionZoneInfoList, nil
}

func (regionZoneHandler *TencentRegionZoneHandler) GetRegionZone(Name string) (irs.RegionZoneInfo, error) {
	responseRegions, err := DescribeRegions(regionZoneHandler.Client)
	if err != nil {
		cblogger.Error("err : DescribeRegions")
		cblogger.Error(err)
		return irs.RegionZoneInfo{}, err
	}

	var targetRegion cvm.RegionInfo
	for _, region := range responseRegions.Response.RegionSet {
		if *region.Region == Name {
			targetRegion = *region
		}
	}

	clientProfile := profile.NewClientProfile()
	clientProfile.Language = "en-US" // lang default set is zh-CN -> set as en-US.

	// var regionZoneInfo irs.RegionZoneInfo

	tempClient, _ := cvm.NewClient(regionZoneHandler.Client.Client.GetCredential(), Name, clientProfile)
	responseZones, _ := DescribeZones(tempClient)

	var zoneInfoList []irs.ZoneInfo
	for _, zone := range responseZones.Response.ZoneSet {
		zoneInfo := irs.ZoneInfo{}
		zoneInfo.Name = *zone.Zone
		zoneInfo.DisplayName = *zone.ZoneName
		zoneInfo.Status = GetZoneStatus(*zone.ZoneState)

		// keyValueList 삭제 https://github.com/cloud-barista/cb-spider/issues/930#issuecomment-1734817828
		// keyValueList, err := ConvertKeyValueList(zone)
		// if err != nil {
		// 	cblogger.Errorf("err : ConvertKeyValueList [%s]", *zone.ZoneName)
		// 	cblogger.Error(err)
		// 	keyValueList = nil
		// }
		// zoneInfo.KeyValueList = keyValueList

		zoneInfoList = append(zoneInfoList, zoneInfo)
	}

	regionZoneInfo := irs.RegionZoneInfo{}
	regionZoneInfo.Name = *targetRegion.Region
	regionZoneInfo.DisplayName = *targetRegion.RegionName
	regionZoneInfo.ZoneList = zoneInfoList

	// keyValueList 삭제 https://github.com/cloud-barista/cb-spider/issues/930#issuecomment-1734817828
	// keyValueList, err := ConvertKeyValueList(targetRegion)
	// if err != nil {
	// 	cblogger.Errorf("err : ConvertKeyValueList [%s]", Name)
	// 	cblogger.Error(err)
	// 	keyValueList = nil
	// }
	// regionZoneInfo.KeyValueList = keyValueList

	return regionZoneInfo, nil
}

func (regionZoneHandler *TencentRegionZoneHandler) ListOrgRegion() (string, error) {

	responseRegions, err := DescribeRegions(regionZoneHandler.Client)
	if err != nil {
		cblogger.Error(err)
		return "", err
	}

	jsonString, errJson := ConvertJsonString(responseRegions)
	if errJson != nil {
		cblogger.Error(err)
		return "", err
	}

	return jsonString, err
}

func (regionZoneHandler *TencentRegionZoneHandler) ListOrgZone() (string, error) {

	responseRegions, err := DescribeRegions(regionZoneHandler.Client)
	if err != nil {
		cblogger.Error(err)
		return "", err
	}

	clientProfile := profile.NewClientProfile()
	clientProfile.Language = "en-US" // lang default set is zh-CN -> set as en-US.

	var responseZonesList []*cvm.DescribeZonesResponse

	for _, region := range responseRegions.Response.RegionSet {
		tempClient, _ := cvm.NewClient(regionZoneHandler.Client.Client.GetCredential(), *region.Region, clientProfile)
		responseZones, _ := DescribeZones(tempClient)

		responseZonesList = append(responseZonesList, responseZones)
	}

	jsonString, errJson := ConvertJsonString(responseZonesList)
	if errJson != nil {
		cblogger.Error(err)
		return "", err
	}

	return jsonString, nil
}
