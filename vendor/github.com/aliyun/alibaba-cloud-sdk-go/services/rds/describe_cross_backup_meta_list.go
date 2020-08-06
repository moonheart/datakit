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

// DescribeCrossBackupMetaList invokes the rds.DescribeCrossBackupMetaList API synchronously
// api document: https://help.aliyun.com/api/rds/describecrossbackupmetalist.html
func (client *Client) DescribeCrossBackupMetaList(request *DescribeCrossBackupMetaListRequest) (response *DescribeCrossBackupMetaListResponse, err error) {
	response = CreateDescribeCrossBackupMetaListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeCrossBackupMetaListWithChan invokes the rds.DescribeCrossBackupMetaList API asynchronously
// api document: https://help.aliyun.com/api/rds/describecrossbackupmetalist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCrossBackupMetaListWithChan(request *DescribeCrossBackupMetaListRequest) (<-chan *DescribeCrossBackupMetaListResponse, <-chan error) {
	responseChan := make(chan *DescribeCrossBackupMetaListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeCrossBackupMetaList(request)
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

// DescribeCrossBackupMetaListWithCallback invokes the rds.DescribeCrossBackupMetaList API asynchronously
// api document: https://help.aliyun.com/api/rds/describecrossbackupmetalist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCrossBackupMetaListWithCallback(request *DescribeCrossBackupMetaListRequest, callback func(response *DescribeCrossBackupMetaListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeCrossBackupMetaListResponse
		var err error
		defer close(result)
		response, err = client.DescribeCrossBackupMetaList(request)
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

// DescribeCrossBackupMetaListRequest is the request struct for api DescribeCrossBackupMetaList
type DescribeCrossBackupMetaListRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Pattern              string           `position:"Query" name:"Pattern"`
	PageSize             string           `position:"Query" name:"PageSize"`
	PageIndex            string           `position:"Query" name:"PageIndex"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	BackupSetId          string           `position:"Query" name:"BackupSetId"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	GetDbName            string           `position:"Query" name:"GetDbName"`
	Region               string           `position:"Query" name:"Region"`
}

// DescribeCrossBackupMetaListResponse is the response struct for api DescribeCrossBackupMetaList
type DescribeCrossBackupMetaListResponse struct {
	*responses.BaseResponse
	RequestId        string                             `json:"RequestId" xml:"RequestId"`
	DBInstanceName   string                             `json:"DBInstanceName" xml:"DBInstanceName"`
	PageNumber       int                                `json:"PageNumber" xml:"PageNumber"`
	PageRecordCount  int                                `json:"PageRecordCount" xml:"PageRecordCount"`
	TotalRecordCount int                                `json:"TotalRecordCount" xml:"TotalRecordCount"`
	TotalPageCount   int                                `json:"TotalPageCount" xml:"TotalPageCount"`
	Items            ItemsInDescribeCrossBackupMetaList `json:"Items" xml:"Items"`
}

// CreateDescribeCrossBackupMetaListRequest creates a request to invoke DescribeCrossBackupMetaList API
func CreateDescribeCrossBackupMetaListRequest() (request *DescribeCrossBackupMetaListRequest) {
	request = &DescribeCrossBackupMetaListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeCrossBackupMetaList", "rds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeCrossBackupMetaListResponse creates a response to parse from DescribeCrossBackupMetaList response
func CreateDescribeCrossBackupMetaListResponse() (response *DescribeCrossBackupMetaListResponse) {
	response = &DescribeCrossBackupMetaListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
