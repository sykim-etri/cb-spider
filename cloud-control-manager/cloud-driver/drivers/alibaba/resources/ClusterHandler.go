// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by ETRI, 2022.08.

package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/alibaba/main/pmks_test/util"
	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/alibaba/utils/alibaba"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/jeremywohl/flatten"
	"github.com/sirupsen/logrus"
	// call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
)

// tempCalllogger
// 공통로거 만들기 이전까지 사용
var once sync.Once
var tempCalllogger *logrus.Logger

func init() {
	once.Do(func() {
		tempCalllogger = call.GetLogger("HISCALL")
	})
}

type AlibabaClusterHandler struct {
	RegionInfo     idrv.RegionInfo
	CredentialInfo idrv.CredentialInfo
}

// connectionInfo.CredentialInfo.AccessKey
// connectionInfo.CredentialInfo.AccessSecret
// connectionInfo.RegionInfo.Region = "region-1"

func (clusterHandler *AlibabaClusterHandler) CreateCluster(clusterReqInfo irs.ClusterInfo) (irs.ClusterInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called CreateCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "CreateCluster()", "CreateCluster()")

	// 클러스터 생성 요청을 JSON 요청으로 변환
	payload, err := getClusterInfoJSON(clusterReqInfo, clusterHandler.RegionInfo.Region)
	if err != nil {
		cblogger.Error(err)
		return irs.ClusterInfo{}, err
	}

	start := call.Start()
	response_json_str, err := alibaba.CreateCluster(clusterHandler.CredentialInfo.AccessKey, clusterHandler.CredentialInfo.AccessSecret, clusterHandler.RegionInfo.Region, payload)
	loggingInfo(callLogInfo, start)
	if err != nil {
		cblogger.Error(err)
		loggingError(callLogInfo, err)
		return irs.ClusterInfo{}, err
	}

	println(response_json_str)
	// {"cluster_id":"c913aebba53eb40f3978495d92b8da57f","request_id":"2C0836DA-ED3B-5B1E-94C9-5B7E355E2E44","task_id":"T-63185224055a0b07c6000083","instanceId":"c913aebba53eb40f3978495d92b8da57f"}

	var response_json_obj map[string]interface{}
	json.Unmarshal([]byte(response_json_str), &response_json_obj)
	cluster_id := response_json_obj["cluster_id"].(string)
	cluster_info, err := getClusterInfo(clusterHandler.CredentialInfo.AccessKey, clusterHandler.CredentialInfo.AccessSecret, clusterHandler.RegionInfo.Region, cluster_id)
	if err != nil {
		return irs.ClusterInfo{}, err
	}

	// 리턴할 ClusterInfo 만들기
	// 일단은 단순하게 만들어서 반환한다.
	// 추후에 정보 추가 필요

	// NodeGroup 생성 정보가 있는경우 생성을 시도한다.
	// 문제는 Cluster 생성이 완료되어야 NodeGroup 생성이 가능하다.
	// Cluster 생성이 완료되려면 최소 10분 이상 걸린다.
	// 성공할때까지 반복하면서 생성을 시도해야 하는가?
	// for _, node_group := range clusterReqInfo.NodeGroupList {
	// 	res, err := clusterHandler.AddNodeGroup(clusterReqInfo.IId, node_group)
	// 	if err != nil {
	// 		cblogger.Error(err)
	// 		return irs.ClusterInfo{}, err
	// 	}
	// }

	return *cluster_info, nil
}

func (clusterHandler *AlibabaClusterHandler) ListCluster() ([]*irs.ClusterInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called ListCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "ListCluster()", "ListCluster()")

	start := call.Start()
	clusters_json_str, err := alibaba.GetClusters(clusterHandler.CredentialInfo.AccessKey, clusterHandler.CredentialInfo.AccessSecret, clusterHandler.RegionInfo.Region)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return nil, err
	}

	var clusters_json_obj map[string]interface{}
	json.Unmarshal([]byte(clusters_json_str), &clusters_json_obj)
	clusters := clusters_json_obj["clusters"].([]interface{})
	cluster_info_list := make([]*irs.ClusterInfo, len(clusters))
	for i, cluster := range clusters {
		println(i, cluster)
		cluster_id := cluster.(map[string]interface{})["cluster_id"].(string)
		cluster_info_list[i], err = getClusterInfo(clusterHandler.CredentialInfo.AccessKey, clusterHandler.CredentialInfo.AccessSecret, clusterHandler.RegionInfo.Region, cluster_id)
		if err != nil {
			return nil, err
		}
	}

	// return cluster_info_list, nil

	return nil, nil
}

func (clusterHandler *AlibabaClusterHandler) GetCluster(clusterIID irs.IID) (irs.ClusterInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called GetCluster()")
	// callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "GetCluster()")

	// start := call.Start()
	// cluster_info, err := getClusterInfo(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId)
	// LoggingInfo(callLogInfo, start)
	// if err != nil {
	// 	return irs.ClusterInfo{}, err
	// }

	// return *cluster_info, nil

	return irs.ClusterInfo{}, nil
}

func (clusterHandler *AlibabaClusterHandler) DeleteCluster(clusterIID irs.IID) (bool, error) {
	cblogger.Info("Alibaba Cloud Driver: called DeleteCluster()")
	// callLogInfo := GetCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "DeleteCluster()")

	// start := call.Start()
	// res, err := alibaba.DeleteCluster(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId)
	// LoggingInfo(callLogInfo, start)
	// if err != nil {
	// 	return false, err
	// }
	// if res != "" {
	// 	// 삭제 처리를 성공하면 ""를 리턴한다.
	// 	// 삭제 처리를 실패하면 에러 메시지를 리턴한다.
	// 	return false, errors.New(res)
	// }

	// return true, nil

	return false, nil
}

func (clusterHandler *AlibabaClusterHandler) AddNodeGroup(clusterIID irs.IID, nodeGroupReqInfo irs.NodeGroupInfo) (irs.NodeGroupInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called AddNodeGroup()")

	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "AddNodeGroup()")

	// 노드 그룹 생성 요청을 JSON 요청으로 변환
	payload, err := getNodeGroupJSONString(nodeGroupReqInfo)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	start := call.Start()
	result_json_str, err := alibaba.CreateNodeGroup(clusterHandler.CredentialInfo.AccessKey, clusterHandler.CredentialInfo.AccessSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, payload)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	var result_json_obj map[string]interface{}
	json.Unmarshal([]byte(result_json_str), &result_json_obj)
	if result_json_obj["errors"] != nil {
		return irs.NodeGroupInfo{}, errors.New(result_json_str)
	}
	// uuid := result_json_obj["uuid"].(string)
	// node_group_info, err := getNodeGroupInfo(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId, uuid)
	// if err != nil {
	// 	return irs.NodeGroupInfo{}, err
	// }

	// return *node_group_info, nil

	return irs.NodeGroupInfo{}, nil
}

func (clusterHandler *AlibabaClusterHandler) ListNodeGroup(clusterIID irs.IID) ([]*irs.NodeGroupInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called ListNodeGroup()")
	// callLogInfo := GetCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "ListNodeGroup()")

	// node_group_info_list := []*irs.NodeGroupInfo{}

	// start := call.Start()
	// node_groups_json_str, err := alibaba.GetNodeGroups(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId)
	// LoggingInfo(callLogInfo, start)
	// if err != nil {
	// 	return node_group_info_list, err
	// }
	// var node_groups_json_obj map[string]interface{}
	// json.Unmarshal([]byte(node_groups_json_str), &node_groups_json_obj)
	// node_groups := node_groups_json_obj["nodegroups"].([]interface{})
	// for _, node_group := range node_groups {
	// 	node_group_info, err := getNodeGroupInfo(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId, node_group.(map[string]interface{})["uuid"].(string))
	// 	if err != nil {
	// 		return node_group_info_list, err
	// 	}
	// 	node_group_info_list = append(node_group_info_list, node_group_info)
	// }

	// return node_group_info_list, nil

	return nil, nil
}

func (clusterHandler *AlibabaClusterHandler) GetNodeGroup(clusterIID irs.IID, nodeGroupIID irs.IID) (irs.NodeGroupInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called GetNodeGroup()")
	// callLogInfo := GetCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "GetNodeGroup()")

	// start := call.Start()
	// temp, err := getNodeGroupInfo(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId, nodeGroupIID.SystemId)
	// LoggingInfo(callLogInfo, start)
	// if err != nil {
	// 	return irs.NodeGroupInfo{}, err
	// }

	// return *temp, nil

	return irs.NodeGroupInfo{}, nil
}

func (clusterHandler *AlibabaClusterHandler) SetNodeGroupAutoScaling(clusterIID irs.IID, nodeGroupIID irs.IID, on bool) (bool, error) {
	cblogger.Info("Alibaba Cloud Driver: called SetNodeGroupAutoScaling()")
	return false, errors.New("SetNodeGroupAutoScaling is not supported")
}

func (clusterHandler *AlibabaClusterHandler) ChangeNodeGroupScaling(clusterIID irs.IID, nodeGroupIID irs.IID, desiredNodeSize int, minNodeSize int, maxNodeSize int) (irs.NodeGroupInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called ChangeNodeGroupScaling()")
	return irs.NodeGroupInfo{}, errors.New("ChangeNodeGroupScaling is not supported")
}

func (clusterHandler *AlibabaClusterHandler) RemoveNodeGroup(clusterIID irs.IID, nodeGroupIID irs.IID) (bool, error) {
	cblogger.Info("Alibaba Cloud Driver: called RemoveNodeGroup()")
	// callLogInfo := GetCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "RemoveNodeGroup()")

	// start := call.Start()
	// res, err := alibaba.DeleteNodeGroup(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId, nodeGroupIID.SystemId)
	// LoggingInfo(callLogInfo, start)
	// if err != nil {
	// 	return false, err
	// }
	// if res != "" {
	// 	// 삭제 처리를 성공하면 ""를 리턴한다.
	// 	// 삭제 처리를 실패하면 에러 메시지를 리턴한다.
	// 	return false, errors.New(res)
	// }

	// return true, nil

	return false, nil
}

func (clusterHandler *AlibabaClusterHandler) UpgradeCluster(clusterIID irs.IID, newVersion string) (irs.ClusterInfo, error) {
	cblogger.Info("Alibaba Cloud Driver: called UpgradeCluster()")
	// callLogInfo := GetCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "UpgradeCluster()")

	// node_groups_json_str, err := alibaba.GetNodeGroups(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId)
	// if err != nil {
	// 	return irs.ClusterInfo{}, err
	// }
	// var node_groups_json_obj map[string]interface{}
	// json.Unmarshal([]byte(node_groups_json_str), &node_groups_json_obj)
	// node_groups := node_groups_json_obj["nodegroups"].([]interface{})
	// for _, node_group := range node_groups {
	// 	start := call.Start()
	// 	res, err := alibaba.UpgradeCluster(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId, node_group.(map[string]interface{})["uuid"].(string), newVersion)
	// 	LoggingInfo(callLogInfo, start)
	// 	if res != "" {
	// 		return irs.ClusterInfo{}, err
	// 	}
	// }

	// clusterInfo, err := getClusterInfo(clusterHandler.ClusterClient.Endpoint, clusterHandler.ClusterClient.TokenID, clusterIID.SystemId)
	// if err != nil {
	// 	return irs.ClusterInfo{}, err
	// }

	// return *clusterInfo, nil

	return irs.ClusterInfo{}, nil
}

func getClusterInfo(access_key string, access_secret string, region_id string, cluster_id string) (*irs.ClusterInfo, error) {
	defer func() {
		if r := recover(); r != nil {
			cblogger.Error("getClusterInfo() failed!", r)
		}
	}()

	cluster_json_str, err := alibaba.GetCluster(access_key, access_secret, region_id, cluster_id)
	if err != nil {
		return nil, err
	}
	println(cluster_json_str)
	flat, err := flatten.FlattenString(cluster_json_str, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}
	println(flat)

	// k,v 추출
	// k,v 변환 규칙 작성 [k,v]:[ClusterInfo.k, ClusterInfo.v]
	// 변환 규칙에 따라 k,v 변환
	var cluster_json_obj map[string]interface{}
	json.Unmarshal([]byte(cluster_json_str), &cluster_json_obj)

	// https://www.alibabacloud.com/help/doc-detail/86987.html
	// Initializing	Creating the cloud resources that are used by the cluster.
	// Creation Failed	Failed to create the cloud resources that are used by the cluster.
	// Running	The cloud resources used by the cluster are created.
	// Updating	Updating the metadata of the cluster.
	// Scaling	Adding nodes to the cluster.
	// Removing	Removing nodes from the cluster.
	// Upgrading	Upgrading the cluster.
	// Draining	Evicting pods from a node to other nodes. After all pods are evicted from the node, the node becomes unschudulable.
	// Deleting	Deleting the cluster.
	// Deletion Failed	Failed to delete the cluster.
	// Deleted (invisible to users)	The cluster is deleted.

	// ClusterCreating ClusterStatus = "Creating"
	// ClusterActive   ClusterStatus = "Active"
	// ClusterInactive ClusterStatus = "Inactive"
	// ClusterUpdating ClusterStatus = "Updating"
	// ClusterDeleting ClusterStatus = "Deleting"

	health_status := cluster_json_obj["state"].(string)
	cluster_status := irs.ClusterActive
	if strings.EqualFold(health_status, "Initializing") {
		cluster_status = irs.ClusterCreating
	} else if strings.EqualFold(health_status, "Updating") {
		cluster_status = irs.ClusterUpdating
	} else if strings.EqualFold(health_status, "Creation Failed") {
		cluster_status = irs.ClusterInactive
	} else if strings.EqualFold(health_status, "Deleting") {
		cluster_status = irs.ClusterDeleting
	} else if strings.EqualFold(health_status, "Running") {
		cluster_status = irs.ClusterActive
	}

	println(cluster_status)

	created_at := cluster_json_obj["created"].(string) // 2022-09-08T09:02:16+08:00,
	datetime, err := time.Parse(time.RFC3339, created_at)
	if err != nil {
		panic(err)
	}

	// name
	// cluster_id
	// current_version
	// security_group_id
	// vpc_id
	// state
	// created
	cluster_info := &irs.ClusterInfo{
		IId: irs.IID{
			NameId:   cluster_json_obj["name"].(string),
			SystemId: cluster_json_obj["cluster_id"].(string),
		},
		Version: cluster_json_obj["current_version"].(string),
		Network: irs.NetworkInfo{
			VpcIID: irs.IID{
				NameId:   "",
				SystemId: cluster_json_obj["vpc_id"].(string),
			},
			SecurityGroupIIDs: []irs.IID{
				{
					NameId:   "",
					SystemId: cluster_json_obj["security_group_id"].(string),
				},
			},
		},
		Status:      cluster_status,
		CreatedTime: datetime,
		// KeyValueList: []irs.KeyValue{}, // flatten data 입력하기
	}
	println(cluster_info)

	// NodeGroups
	node_groups_json_str, err := alibaba.ListNodeGroup(access_key, access_secret, region_id, cluster_id)
	if err != nil {
		return nil, err
	}
	print(node_groups_json_str)
	// {"NextToken":"","TotalCount":0,"nodepools":[],"request_id":"4529A823-F344-5EA6-8E60-47FC30117668"}

	// k,v 추출
	// k,v 변환 규칙 작성 [k,v]:[NodeGroup.k, NodeGroup.v]
	// 변환 규칙에 따라 k,v 변환
	flat, err = flatten.FlattenString(node_groups_json_str, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}
	println(flat)

	var node_groups_json_obj map[string]interface{}
	json.Unmarshal([]byte(node_groups_json_str), &node_groups_json_obj)
	node_groups := node_groups_json_obj["nodepools"].([]interface{})
	for _, node_group := range node_groups {
		// printFlattenJSON(node_group)

		// "nodepool_info.nodepool_id": "np02b049a03b8141858697497e12a61aa1",
		node_group_id := node_group.(map[string]interface{})["nodepool_info"].(map[string]interface{})["nodepool_id"].(string)

		// TODO
		node_group_info, err := getNodeGroupInfo(access_key, access_secret, region_id, cluster_id, node_group_id)
		if err != nil {
			return nil, err
		}
		cluster_info.NodeGroupList = append(cluster_info.NodeGroupList, *node_group_info)
	}

	return cluster_info, nil
}

func printFlattenJSON(json_obj interface{}) {
	temp, err := json.MarshalIndent(json_obj, "", "  ")
	if err != nil {
		println(err)
	} else {
		flat, err := flatten.FlattenString(string(temp), "", flatten.DotStyle)
		if err != nil {
			println(err)
		} else {
			println(flat)
		}
	}
}

func getNodeGroupInfo(host string, token string, cluster_id string, node_group_id string) (*irs.NodeGroupInfo, error) {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			cblogger.Error("getNodeGroupInfo() failed!", r)
	// 		}
	// 	}()

	// 	node_group_json_str, err := alibaba.GetNodeGroup(host, token, cluster_id, node_group_id)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if node_group_json_str == "" {
	// 		return nil, errors.New("not found")
	// 	}

	// 	var node_group_json_obj map[string]interface{}
	// 	json.Unmarshal([]byte(node_group_json_str), &node_group_json_obj)

	// 	auto_scaling, _ := strconv.ParseBool(node_group_json_obj["labels"].(map[string]interface{})["ca_enable"].(string))
	// 	ca_min_node_count, _ := strconv.ParseInt(node_group_json_obj["labels"].(map[string]interface{})["ca_min_node_count"].(string), 10, 32)
	// 	ca_max_node_count, _ := strconv.ParseInt(node_group_json_obj["labels"].(map[string]interface{})["ca_max_node_count"].(string), 10, 32)
	// 	node_count := int(node_group_json_obj["node_count"].(float64))

	// 	node_group_info := irs.NodeGroupInfo{
	// 		IId: irs.IID{
	// 			NameId:   node_group_json_obj["name"].(string),
	// 			SystemId: node_group_json_obj["uuid"].(string),
	// 		},
	// 		ImageIID: irs.IID{
	// 			NameId:   "",
	// 			SystemId: node_group_json_obj["image_id"].(string),
	// 		},
	// 		VMSpecName:      node_group_json_obj["flavor_id"].(string),
	// 		RootDiskType:    node_group_json_obj["labels"].(map[string]interface{})["boot_volume_size"].(string),
	// 		RootDiskSize:    node_group_json_obj["labels"].(map[string]interface{})["boot_volume_type"].(string),
	// 		KeyPairIID:      irs.IID{},
	// 		OnAutoScaling:   auto_scaling,
	// 		MinNodeSize:     int(ca_min_node_count),
	// 		MaxNodeSize:     int(ca_max_node_count),
	// 		DesiredNodeSize: int(node_count),
	// 		NodeList:        []irs.IID{},
	// 		KeyValueList:    []irs.KeyValue{},
	// 	}

	//return &node_group_info, nil

	return nil, nil
}

func getClusterInfoJSON(clusterInfo irs.ClusterInfo, region_id string) (string, error) {

	defer func() {
		if r := recover(); r != nil {
			cblogger.Error("getClusterInfoJSON failed", r)
		}
	}()

	// clusterInfo := irs.ClusterInfo{
	// 	IId: irs.IID{
	// 		NameId:   "cluster-x",
	// 		SystemId: "",
	// 	},
	// 	Version: "1.22.10-aliyun.1",
	// 	Network: irs.NetworkInfo{
	// 		VpcIID: irs.IID{NameId: "", SystemId: "vpc-2zek5slojo5bh621ftnrg"},
	// 	},
	// 	KeyValueList: []irs.KeyValue{
	// 		{
	// 			Key:   "container_cidr",
	// 			Value: "172.31.0.0/16",
	// 		},
	// 		{
	// 			Key:   "service_cidr",
	// 			Value: "172.32.0.0/16",
	// 		},
	// 		{
	// 			Key:   "master_vswitch_id",
	// 			Value: "vsw-2ze0qpwcio7r5bx3nqbp1",
	// 		},
	// 	},
	// }

	//cidr: Valid values: 10.0.0.0/16-24, 172.16-31.0.0/16-24, and 192.168.0.0/16-24.
	container_cidr := ""
	service_cidr := ""
	master_vswitch_id := ""

	for _, v := range clusterInfo.KeyValueList {
		switch v.Key {
		case "container_cidr":
			container_cidr = v.Value
		case "service_cidr":
			service_cidr = v.Value
		case "master_vswitch_id":
			master_vswitch_id = v.Value
		}
	}

	temp := `{
		"name": "%s",
		"region_id": "%s",
		"cluster_type": "ManagedKubernetes",
		"kubernetes_version": "1.22.10-aliyun.1",
		"vpcid": "%s",
		"container_cidr": "%s",
		"service_cidr": "%s",
		"num_of_nodes": 0,
		"master_vswitch_ids": [
			"%s"
		]
	}`

	clusterInfoJSON := fmt.Sprintf(temp, clusterInfo.IId.NameId, region_id, clusterInfo.Network.VpcIID.SystemId, container_cidr, service_cidr, master_vswitch_id)

	return clusterInfoJSON, nil
}

func getNodeGroupJSONString(nodeGroupReqInfo irs.NodeGroupInfo) (string, error) {

	defer func() {
		if r := recover(); r != nil {
			cblogger.Error("getNodeGroupJSONString failed", r)
		}
	}()

	// new_node_group := &irs.NodeGroupInfo{
	// 	IId:             irs.IID{NameId: "nodepoolx100", SystemId: ""},
	// 	ImageIID:        irs.IID{NameId: "", SystemId: "image_id"}, // 이미지 id 선택 추가
	// 	VMSpecName:      "ecs.c6.xlarge",
	// 	RootDiskType:    "cloud_essd",
	// 	RootDiskSize:    "70",
	// 	KeyPairIID:      irs.IID{NameId: "kp1", SystemId: ""},
	// 	OnAutoScaling:   true,
	// 	DesiredNodeSize: 1,
	// 	MinNodeSize:     0,
	// 	MaxNodeSize:     3,
	// 	// KeyValueList: []irs.KeyValue{ // 클러스터 조회해서 처리한다. // //vswitch_id":"vsw-2ze0qpwcio7r5bx3nqbp1"
	// 	// 	{
	// 	// 		Key:   "vswitch_ids",
	// 	// 		Value: "vsw-2ze0qpwcio7r5bx3nqbp1",
	// 	// 	},
	// 	// },
	// }

	name := nodeGroupReqInfo.IId.NameId
	//image_id := nodeGroupReqInfo.ImageIID.SystemId
	enable := nodeGroupReqInfo.OnAutoScaling
	max_instances := nodeGroupReqInfo.MaxNodeSize
	min_instances := nodeGroupReqInfo.MinNodeSize
	// desired_instances := nodeGroupReqInfo.DesiredNodeSize // not supported in alibaba

	instance_type := nodeGroupReqInfo.VMSpecName
	key_pair := nodeGroupReqInfo.KeyPairIID.NameId

	system_disk_category := nodeGroupReqInfo.RootDiskType
	system_disk_size, _ := strconv.ParseInt(nodeGroupReqInfo.RootDiskSize, 10, 32)

	vswitch_id := "vsw-2ze0qpwcio7r5bx3nqbp1" // get vswitch_id, get from cluster info

	temp := `{
		"nodepool_info": {
			"name": "%s"
		},
		"auto_scaling": {
			"enable": %t,
			"max_instances": %d,
			"min_instances": %d
		},
		"scaling_group": {
			"instance_types": ["%s"],
			"key_pair": "%s",
			"system_disk_category": "%s",
			"system_disk_size": %d,
			"vswitch_ids": ["%s"]
		},
		"management": {
			"enable":true
		}
	}`

	payload := fmt.Sprintf(temp, name, enable, max_instances, min_instances, instance_type, key_pair, system_disk_category, system_disk_size, vswitch_id)

	return payload, nil
}

// getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "ListCluster()", "ListCluster()")
func getCallLogScheme(region string, resourceType call.RES_TYPE, resourceName string, apiName string) call.CLOUDLOGSCHEMA {
	cblogger.Info(fmt.Sprintf("Call %s %s", call.ALIBABA, apiName))
	return call.CLOUDLOGSCHEMA{
		CloudOS:      call.ALIBABA,
		RegionZone:   region,
		ResourceType: resourceType,
		ResourceName: resourceName,
		CloudOSAPI:   apiName,
	}
}

func loggingError(hiscallInfo call.CLOUDLOGSCHEMA, err error) {
	hiscallInfo.ErrorMSG = err.Error()
	tempCalllogger.Info(call.String(hiscallInfo))
}

func loggingInfo(hiscallInfo call.CLOUDLOGSCHEMA, start time.Time) {
	hiscallInfo.ElapsedTime = call.Elapsed(start)
	tempCalllogger.Info(call.String(hiscallInfo))
}
