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

// EmailVerified invokes the domain.EmailVerified API synchronously
// api document: https://help.aliyun.com/api/domain/emailverified.html
func (client *Client) EmailVerified(request *EmailVerifiedRequest) (response *EmailVerifiedResponse, err error) {
	response = CreateEmailVerifiedResponse()
	err = client.DoAction(request, response)
	return
}

// EmailVerifiedWithChan invokes the domain.EmailVerified API asynchronously
// api document: https://help.aliyun.com/api/domain/emailverified.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) EmailVerifiedWithChan(request *EmailVerifiedRequest) (<-chan *EmailVerifiedResponse, <-chan error) {
	responseChan := make(chan *EmailVerifiedResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.EmailVerified(request)
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

// EmailVerifiedWithCallback invokes the domain.EmailVerified API asynchronously
// api document: https://help.aliyun.com/api/domain/emailverified.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) EmailVerifiedWithCallback(request *EmailVerifiedRequest, callback func(response *EmailVerifiedResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *EmailVerifiedResponse
		var err error
		defer close(result)
		response, err = client.EmailVerified(request)
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

// EmailVerifiedRequest is the request struct for api EmailVerified
type EmailVerifiedRequest struct {
	*requests.RpcRequest
	UserClientIp string `position:"Query" name:"UserClientIp"`
	Lang         string `position:"Query" name:"Lang"`
	Email        string `position:"Query" name:"Email"`
}

// EmailVerifiedResponse is the response struct for api EmailVerified
type EmailVerifiedResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateEmailVerifiedRequest creates a request to invoke EmailVerified API
func CreateEmailVerifiedRequest() (request *EmailVerifiedRequest) {
	request = &EmailVerifiedRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "EmailVerified", "domain", "openAPI")
	request.Method = requests.POST
	return
}

// CreateEmailVerifiedResponse creates a response to parse from EmailVerified response
func CreateEmailVerifiedResponse() (response *EmailVerifiedResponse) {
	response = &EmailVerifiedResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
