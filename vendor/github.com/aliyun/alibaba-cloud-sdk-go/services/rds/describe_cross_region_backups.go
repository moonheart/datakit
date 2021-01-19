package rds

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

// DescribeCrossRegionBackups invokes the rds.DescribeCrossRegionBackups API synchronously
// api document: https://help.aliyun.com/api/rds/describecrossregionbackups.html
func (client *Client) DescribeCrossRegionBackups(request *DescribeCrossRegionBackupsRequest) (response *DescribeCrossRegionBackupsResponse, err error) {
	response = CreateDescribeCrossRegionBackupsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeCrossRegionBackupsWithChan invokes the rds.DescribeCrossRegionBackups API asynchronously
// api document: https://help.aliyun.com/api/rds/describecrossregionbackups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCrossRegionBackupsWithChan(request *DescribeCrossRegionBackupsRequest) (<-chan *DescribeCrossRegionBackupsResponse, <-chan error) {
	responseChan := make(chan *DescribeCrossRegionBackupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeCrossRegionBackups(request)
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

// DescribeCrossRegionBackupsWithCallback invokes the rds.DescribeCrossRegionBackups API asynchronously
// api document: https://help.aliyun.com/api/rds/describecrossregionbackups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCrossRegionBackupsWithCallback(request *DescribeCrossRegionBackupsRequest, callback func(response *DescribeCrossRegionBackupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeCrossRegionBackupsResponse
		var err error
		defer close(result)
		response, err = client.DescribeCrossRegionBackups(request)
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

// DescribeCrossRegionBackupsRequest is the request struct for api DescribeCrossRegionBackups
type DescribeCrossRegionBackupsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	StartTime            string           `position:"Query" name:"StartTime"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	EndTime              string           `position:"Query" name:"EndTime"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	CrossBackupRegion    string           `position:"Query" name:"CrossBackupRegion"`
	CrossBackupId        requests.Integer `position:"Query" name:"CrossBackupId"`
}

// DescribeCrossRegionBackupsResponse is the response struct for api DescribeCrossRegionBackups
type DescribeCrossRegionBackupsResponse struct {
	*responses.BaseResponse
	RequestId        string                            `json:"RequestId" xml:"RequestId"`
	RegionId         string                            `json:"RegionId" xml:"RegionId"`
	StartTime        string                            `json:"StartTime" xml:"StartTime"`
	EndTime          string                            `json:"EndTime" xml:"EndTime"`
	TotalRecordCount int                               `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageRecordCount  int                               `json:"PageRecordCount" xml:"PageRecordCount"`
	PageNumber       int                               `json:"PageNumber" xml:"PageNumber"`
	Items            ItemsInDescribeCrossRegionBackups `json:"Items" xml:"Items"`
}

// CreateDescribeCrossRegionBackupsRequest creates a request to invoke DescribeCrossRegionBackups API
func CreateDescribeCrossRegionBackupsRequest() (request *DescribeCrossRegionBackupsRequest) {
	request = &DescribeCrossRegionBackupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeCrossRegionBackups", "rds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeCrossRegionBackupsResponse creates a response to parse from DescribeCrossRegionBackups response
func CreateDescribeCrossRegionBackupsResponse() (response *DescribeCrossRegionBackupsResponse) {
	response = &DescribeCrossRegionBackupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
