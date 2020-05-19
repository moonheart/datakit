package aegis

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

// DeleteDingTalk invokes the aegis.DeleteDingTalk API synchronously
// api document: https://help.aliyun.com/api/aegis/deletedingtalk.html
func (client *Client) DeleteDingTalk(request *DeleteDingTalkRequest) (response *DeleteDingTalkResponse, err error) {
	response = CreateDeleteDingTalkResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteDingTalkWithChan invokes the aegis.DeleteDingTalk API asynchronously
// api document: https://help.aliyun.com/api/aegis/deletedingtalk.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDingTalkWithChan(request *DeleteDingTalkRequest) (<-chan *DeleteDingTalkResponse, <-chan error) {
	responseChan := make(chan *DeleteDingTalkResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteDingTalk(request)
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

// DeleteDingTalkWithCallback invokes the aegis.DeleteDingTalk API asynchronously
// api document: https://help.aliyun.com/api/aegis/deletedingtalk.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDingTalkWithCallback(request *DeleteDingTalkRequest, callback func(response *DeleteDingTalkResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteDingTalkResponse
		var err error
		defer close(result)
		response, err = client.DeleteDingTalk(request)
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

// DeleteDingTalkRequest is the request struct for api DeleteDingTalk
type DeleteDingTalkRequest struct {
	*requests.RpcRequest
	SourceIp string `position:"Query" name:"SourceIp"`
	Ids      string `position:"Query" name:"Ids"`
}

// DeleteDingTalkResponse is the response struct for api DeleteDingTalk
type DeleteDingTalkResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteDingTalkRequest creates a request to invoke DeleteDingTalk API
func CreateDeleteDingTalkRequest() (request *DeleteDingTalkRequest) {
	request = &DeleteDingTalkRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("aegis", "2016-11-11", "DeleteDingTalk", "vipaegis", "openAPI")
	return
}

// CreateDeleteDingTalkResponse creates a response to parse from DeleteDingTalk response
func CreateDeleteDingTalkResponse() (response *DeleteDingTalkResponse) {
	response = &DeleteDingTalkResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
