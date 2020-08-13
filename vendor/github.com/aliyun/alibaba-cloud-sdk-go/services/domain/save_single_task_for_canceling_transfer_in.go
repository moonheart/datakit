package domain

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

// SaveSingleTaskForCancelingTransferIn invokes the domain.SaveSingleTaskForCancelingTransferIn API synchronously
// api document: https://help.aliyun.com/api/domain/savesingletaskforcancelingtransferin.html
func (client *Client) SaveSingleTaskForCancelingTransferIn(request *SaveSingleTaskForCancelingTransferInRequest) (response *SaveSingleTaskForCancelingTransferInResponse, err error) {
	response = CreateSaveSingleTaskForCancelingTransferInResponse()
	err = client.DoAction(request, response)
	return
}

// SaveSingleTaskForCancelingTransferInWithChan invokes the domain.SaveSingleTaskForCancelingTransferIn API asynchronously
// api document: https://help.aliyun.com/api/domain/savesingletaskforcancelingtransferin.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCancelingTransferInWithChan(request *SaveSingleTaskForCancelingTransferInRequest) (<-chan *SaveSingleTaskForCancelingTransferInResponse, <-chan error) {
	responseChan := make(chan *SaveSingleTaskForCancelingTransferInResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveSingleTaskForCancelingTransferIn(request)
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

// SaveSingleTaskForCancelingTransferInWithCallback invokes the domain.SaveSingleTaskForCancelingTransferIn API asynchronously
// api document: https://help.aliyun.com/api/domain/savesingletaskforcancelingtransferin.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCancelingTransferInWithCallback(request *SaveSingleTaskForCancelingTransferInRequest, callback func(response *SaveSingleTaskForCancelingTransferInResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveSingleTaskForCancelingTransferInResponse
		var err error
		defer close(result)
		response, err = client.SaveSingleTaskForCancelingTransferIn(request)
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

// SaveSingleTaskForCancelingTransferInRequest is the request struct for api SaveSingleTaskForCancelingTransferIn
type SaveSingleTaskForCancelingTransferInRequest struct {
	*requests.RpcRequest
	DomainName   string `position:"Query" name:"DomainName"`
	UserClientIp string `position:"Query" name:"UserClientIp"`
	Lang         string `position:"Query" name:"Lang"`
}

// SaveSingleTaskForCancelingTransferInResponse is the response struct for api SaveSingleTaskForCancelingTransferIn
type SaveSingleTaskForCancelingTransferInResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveSingleTaskForCancelingTransferInRequest creates a request to invoke SaveSingleTaskForCancelingTransferIn API
func CreateSaveSingleTaskForCancelingTransferInRequest() (request *SaveSingleTaskForCancelingTransferInRequest) {
	request = &SaveSingleTaskForCancelingTransferInRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "SaveSingleTaskForCancelingTransferIn", "domain", "openAPI")
	request.Method = requests.POST
	return
}

// CreateSaveSingleTaskForCancelingTransferInResponse creates a response to parse from SaveSingleTaskForCancelingTransferIn response
func CreateSaveSingleTaskForCancelingTransferInResponse() (response *SaveSingleTaskForCancelingTransferInResponse) {
	response = &SaveSingleTaskForCancelingTransferInResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
