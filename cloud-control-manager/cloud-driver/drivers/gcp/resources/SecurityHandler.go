// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// program by ysjeon@mz.co.kr, 2019.07.
// modify by devunet@mz.co.kr, 2019.11.

package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"
	compute "google.golang.org/api/compute/v1"
)

type GCPSecurityHandler struct {
	Region     idrv.RegionInfo
	Ctx        context.Context
	Client     *compute.Service
	Credential idrv.CredentialInfo
}

//@TODO : 이슈
func (securityHandler *GCPSecurityHandler) CreateSecurity(securityReqInfo irs.SecurityReqInfo) (irs.SecurityInfo, error) {
	cblogger.Info(securityReqInfo)

	vNetworkHandler := GCPVPCHandler{
		Client:     securityHandler.Client,
		Region:     securityHandler.Region,
		Ctx:        securityHandler.Ctx,
		Credential: securityHandler.Credential,
	}

	vNetInfo, errVnet := vNetworkHandler.GetVPC(securityReqInfo.VpcIID)
	spew.Dump(vNetInfo)
	if errVnet != nil {
		cblogger.Error(errVnet)
		return irs.SecurityInfo{}, errVnet
	}

	projectID := securityHandler.Credential.ProjectID
	// @TODO: SecurityGroup 생성 요청 파라미터 정의 필요
	ports := *securityReqInfo.SecurityRules
	var firewallAllowed []*compute.FirewallAllowed

	for _, item := range ports {
		var port string
		fp := item.FromPort
		tp := item.ToPort

		if tp != "" && fp != "" {
			port = fp + "-" + tp
		}
		if tp != "" && fp == "" {
			port = tp
		}
		if tp == "" && fp != "" {
			port = fp
		}

		firewallAllowed = append(firewallAllowed, &compute.FirewallAllowed{
			IPProtocol: item.IPProtocol,
			Ports: []string{
				port,
			},
		})
	}

	var sgDirection string
	if strings.EqualFold(securityReqInfo.Direction, "inbound") {
		sgDirection = "INGRESS"
	} else if strings.EqualFold(securityReqInfo.Direction, "outbound") {
		sgDirection = "EGRESS"
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + projectID
	//networkURL := prefix + "/global/networks/" + securityReqInfo.VpcIID.NameId
	networkURL := prefix + "/global/networks/" + securityReqInfo.VpcIID.SystemId

	fireWall := &compute.Firewall{
		Allowed:   firewallAllowed,
		Direction: sgDirection, //INGRESS(inbound), EGRESS(outbound)
		SourceRanges: []string{
			"0.0.0.0/0",
		},
		Name: securityReqInfo.IId.NameId,
		TargetTags: []string{
			securityReqInfo.IId.NameId,
		},
		Network: networkURL,
	}

	res, err := securityHandler.Client.Firewalls.Insert(projectID, fireWall).Do()
	if err != nil {
		cblogger.Error(err)
		return irs.SecurityInfo{}, err
	}
	fmt.Println("create result : ", res)
	time.Sleep(time.Second * 3)
	//secInfo, _ := securityHandler.GetSecurity(securityReqInfo.IId)
	secInfo, _ := securityHandler.GetSecurity(irs.IID{SystemId: securityReqInfo.IId.NameId})
	return secInfo, nil
}

func (securityHandler *GCPSecurityHandler) ListSecurity() ([]*irs.SecurityInfo, error) {
	//result, err := securityHandler.Client.ListAll(securityHandler.Ctx)
	projectID := securityHandler.Credential.ProjectID
	result, err := securityHandler.Client.Firewalls.List(projectID).Do()
	if err != nil {
		cblogger.Error(err)
		return nil, err
	}

	var securityInfo []*irs.SecurityInfo
	for _, item := range result.Items {
		name := item.Name
		//systemId := strconv.FormatUint(item.Id, 10)
		//secInfo, _ := securityHandler.GetSecurity(irs.IID{NameId: name, SystemId: systemId})
		secInfo, _ := securityHandler.GetSecurity(irs.IID{SystemId: name})

		securityInfo = append(securityInfo, &secInfo)
	}

	return securityInfo, nil
}

func (securityHandler *GCPSecurityHandler) GetSecurity(securityIID irs.IID) (irs.SecurityInfo, error) {
	projectID := securityHandler.Credential.ProjectID

	security, err := securityHandler.Client.Firewalls.Get(projectID, securityIID.SystemId).Do()
	if err != nil {
		cblogger.Error(err)
		return irs.SecurityInfo{}, err
	}
	var securityRules []irs.SecurityRuleInfo
	for _, item := range security.Allowed {
		var portArr []string
		var fromPort string
		var toPort string
		if ports := item.Ports; ports != nil {
			portArr = strings.Split(item.Ports[0], "-")
			fromPort = portArr[0]
			if len(portArr) > 1 {
				toPort = portArr[len(portArr)-1]
			} else {
				toPort = ""
			}

		} else {
			fromPort = ""
			toPort = ""
		}

		securityRules = append(securityRules, irs.SecurityRuleInfo{
			FromPort:   fromPort,
			ToPort:     toPort,
			IPProtocol: item.IPProtocol,
			Direction:  security.Direction,
		})
	}
	vpcArr := strings.Split(security.Network, "/")
	vpcName := vpcArr[len(vpcArr)-1]
	securityInfo := irs.SecurityInfo{
		IId: irs.IID{
			NameId: security.Name,
			//SystemId: strconv.FormatUint(security.Id, 10),
			SystemId: security.Name,
		},
		VpcIID: irs.IID{
			NameId:   vpcName,
			SystemId: vpcName,
		},

		Direction: security.Direction,
		KeyValueList: []irs.KeyValue{
			{"Priority", strconv.FormatInt(security.Priority, 10)},
			{"SourceRanges", security.SourceRanges[0]},
			{"Allowed", security.Allowed[0].IPProtocol},
			{"Vpc", vpcName},
		},
		SecurityRules: &securityRules,
	}

	return securityInfo, nil
}

func (securityHandler *GCPSecurityHandler) DeleteSecurity(securityIID irs.IID) (bool, error) {
	projectID := securityHandler.Credential.ProjectID

	res, err := securityHandler.Client.Firewalls.Delete(projectID, securityIID.SystemId).Do()
	if err != nil {
		cblogger.Error(err)
		return false, err
	}
	fmt.Println(res)
	return true, nil
}
