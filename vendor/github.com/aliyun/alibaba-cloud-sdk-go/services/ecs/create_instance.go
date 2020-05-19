package ecs

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateInstance invokes the ecs.CreateInstance API synchronously
// api document: https://help.aliyun.com/api/ecs/createinstance.html
func (client *Client) CreateInstance(request *CreateInstanceRequest) (response *CreateInstanceResponse, err error) {
	response = CreateCreateInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// CreateInstanceWithChan invokes the ecs.CreateInstance API asynchronously
// api document: https://help.aliyun.com/api/ecs/createinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateInstanceWithChan(request *CreateInstanceRequest) (<-chan *CreateInstanceResponse, <-chan error) {
	responseChan := make(chan *CreateInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateInstance(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateInstanceWithCallback invokes the ecs.CreateInstance API asynchronously
// api document: https://help.aliyun.com/api/ecs/createinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateInstanceWithCallback(request *CreateInstanceRequest, callback func(response *CreateInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateInstanceResponse
		var err error
		defer close(result)
		response, err = client.CreateInstance(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateInstanceRequest is the request struct for api CreateInstance
type CreateInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId               requests.Integer          `position:"Query" name:"ResourceOwnerId"`
	HpcClusterId                  string                    `position:"Query" name:"HpcClusterId"`
	HttpPutResponseHopLimit       requests.Integer          `position:"Query" name:"HttpPutResponseHopLimit"`
	SecurityEnhancementStrategy   string                    `position:"Query" name:"SecurityEnhancementStrategy"`
	KeyPairName                   string                    `position:"Query" name:"KeyPairName"`
	SpotPriceLimit                requests.Float            `position:"Query" name:"SpotPriceLimit"`
	DeletionProtection            requests.Boolean          `position:"Query" name:"DeletionProtection"`
	ResourceGroupId               string                    `position:"Query" name:"ResourceGroupId"`
	HostName                      string                    `position:"Query" name:"HostName"`
	Password                      string                    `position:"Query" name:"Password"`
	StorageSetPartitionNumber     requests.Integer          `position:"Query" name:"StorageSetPartitionNumber"`
	Tag                           *[]CreateInstanceTag      `position:"Query" name:"Tag"  type:"Repeated"`
	AutoRenewPeriod               requests.Integer          `position:"Query" name:"AutoRenewPeriod"`
	NodeControllerId              string                    `position:"Query" name:"NodeControllerId"`
	Period                        requests.Integer          `position:"Query" name:"Period"`
	DryRun                        requests.Boolean          `position:"Query" name:"DryRun"`
	OwnerId                       requests.Integer          `position:"Query" name:"OwnerId"`
	CapacityReservationPreference string                    `position:"Query" name:"CapacityReservationPreference"`
	VSwitchId                     string                    `position:"Query" name:"VSwitchId"`
	PrivateIpAddress              string                    `position:"Query" name:"PrivateIpAddress"`
	SpotStrategy                  string                    `position:"Query" name:"SpotStrategy"`
	PeriodUnit                    string                    `position:"Query" name:"PeriodUnit"`
	InstanceName                  string                    `position:"Query" name:"InstanceName"`
	AutoRenew                     requests.Boolean          `position:"Query" name:"AutoRenew"`
	InternetChargeType            string                    `position:"Query" name:"InternetChargeType"`
	ZoneId                        string                    `position:"Query" name:"ZoneId"`
	InternetMaxBandwidthIn        requests.Integer          `position:"Query" name:"InternetMaxBandwidthIn"`
	UseAdditionalService          requests.Boolean          `position:"Query" name:"UseAdditionalService"`
	Affinity                      string                    `position:"Query" name:"Affinity"`
	ImageId                       string                    `position:"Query" name:"ImageId"`
	ClientToken                   string                    `position:"Query" name:"ClientToken"`
	VlanId                        string                    `position:"Query" name:"VlanId"`
	SpotInterruptionBehavior      string                    `position:"Query" name:"SpotInterruptionBehavior"`
	IoOptimized                   string                    `position:"Query" name:"IoOptimized"`
	SecurityGroupId               string                    `position:"Query" name:"SecurityGroupId"`
	InternetMaxBandwidthOut       requests.Integer          `position:"Query" name:"InternetMaxBandwidthOut"`
	Description                   string                    `position:"Query" name:"Description"`
	SystemDiskCategory            string                    `position:"Query" name:"SystemDisk.Category"`
	CapacityReservationId         string                    `position:"Query" name:"CapacityReservationId"`
	SystemDiskPerformanceLevel    string                    `position:"Query" name:"SystemDisk.PerformanceLevel"`
	UserData                      string                    `position:"Query" name:"UserData"`
	PasswordInherit               requests.Boolean          `position:"Query" name:"PasswordInherit"`
	HttpEndpoint                  string                    `position:"Query" name:"HttpEndpoint"`
	InstanceType                  string                    `position:"Query" name:"InstanceType"`
	Arn                           *[]CreateInstanceArn      `position:"Query" name:"Arn"  type:"Repeated"`
	InstanceChargeType            string                    `position:"Query" name:"InstanceChargeType"`
	DeploymentSetId               string                    `position:"Query" name:"DeploymentSetId"`
	InnerIpAddress                string                    `position:"Query" name:"InnerIpAddress"`
	ResourceOwnerAccount          string                    `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                  string                    `position:"Query" name:"OwnerAccount"`
	Tenancy                       string                    `position:"Query" name:"Tenancy"`
	SystemDiskDiskName            string                    `position:"Query" name:"SystemDisk.DiskName"`
	RamRoleName                   string                    `position:"Query" name:"RamRoleName"`
	DedicatedHostId               string                    `position:"Query" name:"DedicatedHostId"`
	ClusterId                     string                    `position:"Query" name:"ClusterId"`
	CreditSpecification           string                    `position:"Query" name:"CreditSpecification"`
	SpotDuration                  requests.Integer          `position:"Query" name:"SpotDuration"`
	DataDisk                      *[]CreateInstanceDataDisk `position:"Query" name:"DataDisk"  type:"Repeated"`
	StorageSetId                  string                    `position:"Query" name:"StorageSetId"`
	SystemDiskSize                requests.Integer          `position:"Query" name:"SystemDisk.Size"`
	ImageFamily                   string                    `position:"Query" name:"ImageFamily"`
	HttpTokens                    string                    `position:"Query" name:"HttpTokens"`
	SystemDiskDescription         string                    `position:"Query" name:"SystemDisk.Description"`
}

// CreateInstanceTag is a repeated param struct in CreateInstanceRequest
type CreateInstanceTag struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// CreateInstanceArn is a repeated param struct in CreateInstanceRequest
type CreateInstanceArn struct {
	Rolearn       string `name:"Rolearn"`
	RoleType      string `name:"RoleType"`
	AssumeRoleFor string `name:"AssumeRoleFor"`
}

// CreateInstanceDataDisk is a repeated param struct in CreateInstanceRequest
type CreateInstanceDataDisk struct {
	DiskName           string `name:"DiskName"`
	SnapshotId         string `name:"SnapshotId"`
	Size               string `name:"Size"`
	Encrypted          string `name:"Encrypted"`
	PerformanceLevel   string `name:"PerformanceLevel"`
	EncryptAlgorithm   string `name:"EncryptAlgorithm"`
	Description        string `name:"Description"`
	Category           string `name:"Category"`
	KMSKeyId           string `name:"KMSKeyId"`
	Device             string `name:"Device"`
	DeleteWithInstance string `name:"DeleteWithInstance"`
}

// CreateInstanceResponse is the response struct for api CreateInstance
type CreateInstanceResponse struct {
	*responses.BaseResponse
	RequestId  string  `json:"RequestId" xml:"RequestId"`
	InstanceId string  `json:"InstanceId" xml:"InstanceId"`
	TradePrice float64 `json:"TradePrice" xml:"TradePrice"`
}

// CreateCreateInstanceRequest creates a request to invoke CreateInstance API
func CreateCreateInstanceRequest() (request *CreateInstanceRequest) {
	request = &CreateInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "CreateInstance", "ecs", "openAPI")
	return
}

// CreateCreateInstanceResponse creates a response to parse from CreateInstance response
func CreateCreateInstanceResponse() (response *CreateInstanceResponse) {
	response = &CreateInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
