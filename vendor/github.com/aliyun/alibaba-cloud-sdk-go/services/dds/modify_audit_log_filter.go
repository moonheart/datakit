package dds

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

// ModifyAuditLogFilter invokes the dds.ModifyAuditLogFilter API synchronously
// api document: https://help.aliyun.com/api/dds/modifyauditlogfilter.html
func (client *Client) ModifyAuditLogFilter(request *ModifyAuditLogFilterRequest) (response *ModifyAuditLogFilterResponse, err error) {
	response = CreateModifyAuditLogFilterResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyAuditLogFilterWithChan invokes the dds.ModifyAuditLogFilter API asynchronously
// api document: https://help.aliyun.com/api/dds/modifyauditlogfilter.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAuditLogFilterWithChan(request *ModifyAuditLogFilterRequest) (<-chan *ModifyAuditLogFilterResponse, <-chan error) {
	responseChan := make(chan *ModifyAuditLogFilterResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyAuditLogFilter(request)
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

// ModifyAuditLogFilterWithCallback invokes the dds.ModifyAuditLogFilter API asynchronously
// api document: https://help.aliyun.com/api/dds/modifyauditlogfilter.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAuditLogFilterWithCallback(request *ModifyAuditLogFilterRequest, callback func(response *ModifyAuditLogFilterResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyAuditLogFilterResponse
		var err error
		defer close(result)
		response, err = client.ModifyAuditLogFilter(request)
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

// ModifyAuditLogFilterRequest is the request struct for api ModifyAuditLogFilter
type ModifyAuditLogFilterRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	RoleType             string           `position:"Query" name:"RoleType"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Filter               string           `position:"Query" name:"Filter"`
}

// ModifyAuditLogFilterResponse is the response struct for api ModifyAuditLogFilter
type ModifyAuditLogFilterResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyAuditLogFilterRequest creates a request to invoke ModifyAuditLogFilter API
func CreateModifyAuditLogFilterRequest() (request *ModifyAuditLogFilterRequest) {
	request = &ModifyAuditLogFilterRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "ModifyAuditLogFilter", "Dds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyAuditLogFilterResponse creates a response to parse from ModifyAuditLogFilter response
func CreateModifyAuditLogFilterResponse() (response *ModifyAuditLogFilterResponse) {
	response = &ModifyAuditLogFilterResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
