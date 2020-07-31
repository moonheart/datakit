package slb

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

// StopLoadBalancerListener invokes the slb.StopLoadBalancerListener API synchronously
// api document: https://help.aliyun.com/api/slb/stoploadbalancerlistener.html
func (client *Client) StopLoadBalancerListener(request *StopLoadBalancerListenerRequest) (response *StopLoadBalancerListenerResponse, err error) {
	response = CreateStopLoadBalancerListenerResponse()
	err = client.DoAction(request, response)
	return
}

// StopLoadBalancerListenerWithChan invokes the slb.StopLoadBalancerListener API asynchronously
// api document: https://help.aliyun.com/api/slb/stoploadbalancerlistener.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) StopLoadBalancerListenerWithChan(request *StopLoadBalancerListenerRequest) (<-chan *StopLoadBalancerListenerResponse, <-chan error) {
	responseChan := make(chan *StopLoadBalancerListenerResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.StopLoadBalancerListener(request)
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

// StopLoadBalancerListenerWithCallback invokes the slb.StopLoadBalancerListener API asynchronously
// api document: https://help.aliyun.com/api/slb/stoploadbalancerlistener.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) StopLoadBalancerListenerWithCallback(request *StopLoadBalancerListenerRequest, callback func(response *StopLoadBalancerListenerResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *StopLoadBalancerListenerResponse
		var err error
		defer close(result)
		response, err = client.StopLoadBalancerListener(request)
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

// StopLoadBalancerListenerRequest is the request struct for api StopLoadBalancerListener
type StopLoadBalancerListenerRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ListenerPort         requests.Integer `position:"Query" name:"ListenerPort"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ListenerProtocol     string           `position:"Query" name:"ListenerProtocol"`
	Tags                 string           `position:"Query" name:"Tags"`
	LoadBalancerId       string           `position:"Query" name:"LoadBalancerId"`
}

// StopLoadBalancerListenerResponse is the response struct for api StopLoadBalancerListener
type StopLoadBalancerListenerResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateStopLoadBalancerListenerRequest creates a request to invoke StopLoadBalancerListener API
func CreateStopLoadBalancerListenerRequest() (request *StopLoadBalancerListenerRequest) {
	request = &StopLoadBalancerListenerRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "StopLoadBalancerListener", "slb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateStopLoadBalancerListenerResponse creates a response to parse from StopLoadBalancerListener response
func CreateStopLoadBalancerListenerResponse() (response *StopLoadBalancerListenerResponse) {
	response = &StopLoadBalancerListenerResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
