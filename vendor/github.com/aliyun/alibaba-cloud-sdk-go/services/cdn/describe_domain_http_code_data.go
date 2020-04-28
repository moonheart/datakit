package cdn

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

// DescribeDomainHttpCodeData invokes the cdn.DescribeDomainHttpCodeData API synchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedata.html
func (client *Client) DescribeDomainHttpCodeData(request *DescribeDomainHttpCodeDataRequest) (response *DescribeDomainHttpCodeDataResponse, err error) {
	response = CreateDescribeDomainHttpCodeDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainHttpCodeDataWithChan invokes the cdn.DescribeDomainHttpCodeData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainHttpCodeDataWithChan(request *DescribeDomainHttpCodeDataRequest) (<-chan *DescribeDomainHttpCodeDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainHttpCodeDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainHttpCodeData(request)
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

// DescribeDomainHttpCodeDataWithCallback invokes the cdn.DescribeDomainHttpCodeData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainHttpCodeDataWithCallback(request *DescribeDomainHttpCodeDataRequest, callback func(response *DescribeDomainHttpCodeDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainHttpCodeDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainHttpCodeData(request)
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

// DescribeDomainHttpCodeDataRequest is the request struct for api DescribeDomainHttpCodeData
type DescribeDomainHttpCodeDataRequest struct {
	*requests.RpcRequest
	StartTime  string           `position:"Query" name:"StartTime"`
	DomainName string           `position:"Query" name:"DomainName"`
	EndTime    string           `position:"Query" name:"EndTime"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
	Interval   string           `position:"Query" name:"Interval"`
}

// DescribeDomainHttpCodeDataResponse is the response struct for api DescribeDomainHttpCodeData
type DescribeDomainHttpCodeDataResponse struct {
	*responses.BaseResponse
	RequestId    string                                   `json:"RequestId" xml:"RequestId"`
	DomainName   string                                   `json:"DomainName" xml:"DomainName"`
	StartTime    string                                   `json:"StartTime" xml:"StartTime"`
	EndTime      string                                   `json:"EndTime" xml:"EndTime"`
	DataInterval string                                   `json:"DataInterval" xml:"DataInterval"`
	HttpCodeData HttpCodeDataInDescribeDomainHttpCodeData `json:"HttpCodeData" xml:"HttpCodeData"`
}

// CreateDescribeDomainHttpCodeDataRequest creates a request to invoke DescribeDomainHttpCodeData API
func CreateDescribeDomainHttpCodeDataRequest() (request *DescribeDomainHttpCodeDataRequest) {
	request = &DescribeDomainHttpCodeDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeDomainHttpCodeData", "", "")
	return
}

// CreateDescribeDomainHttpCodeDataResponse creates a response to parse from DescribeDomainHttpCodeData response
func CreateDescribeDomainHttpCodeDataResponse() (response *DescribeDomainHttpCodeDataResponse) {
	response = &DescribeDomainHttpCodeDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
