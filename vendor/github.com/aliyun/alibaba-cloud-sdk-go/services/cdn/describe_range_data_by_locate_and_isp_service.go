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

// DescribeRangeDataByLocateAndIspService invokes the cdn.DescribeRangeDataByLocateAndIspService API synchronously
// api document: https://help.aliyun.com/api/cdn/describerangedatabylocateandispservice.html
func (client *Client) DescribeRangeDataByLocateAndIspService(request *DescribeRangeDataByLocateAndIspServiceRequest) (response *DescribeRangeDataByLocateAndIspServiceResponse, err error) {
	response = CreateDescribeRangeDataByLocateAndIspServiceResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeRangeDataByLocateAndIspServiceWithChan invokes the cdn.DescribeRangeDataByLocateAndIspService API asynchronously
// api document: https://help.aliyun.com/api/cdn/describerangedatabylocateandispservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeRangeDataByLocateAndIspServiceWithChan(request *DescribeRangeDataByLocateAndIspServiceRequest) (<-chan *DescribeRangeDataByLocateAndIspServiceResponse, <-chan error) {
	responseChan := make(chan *DescribeRangeDataByLocateAndIspServiceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeRangeDataByLocateAndIspService(request)
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

// DescribeRangeDataByLocateAndIspServiceWithCallback invokes the cdn.DescribeRangeDataByLocateAndIspService API asynchronously
// api document: https://help.aliyun.com/api/cdn/describerangedatabylocateandispservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeRangeDataByLocateAndIspServiceWithCallback(request *DescribeRangeDataByLocateAndIspServiceRequest, callback func(response *DescribeRangeDataByLocateAndIspServiceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeRangeDataByLocateAndIspServiceResponse
		var err error
		defer close(result)
		response, err = client.DescribeRangeDataByLocateAndIspService(request)
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

// DescribeRangeDataByLocateAndIspServiceRequest is the request struct for api DescribeRangeDataByLocateAndIspService
type DescribeRangeDataByLocateAndIspServiceRequest struct {
	*requests.RpcRequest
	DomainNames   string           `position:"Query" name:"DomainNames"`
	LocationNames string           `position:"Query" name:"LocationNames"`
	StartTime     string           `position:"Query" name:"StartTime"`
	IspNames      string           `position:"Query" name:"IspNames"`
	EndTime       string           `position:"Query" name:"EndTime"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeRangeDataByLocateAndIspServiceResponse is the response struct for api DescribeRangeDataByLocateAndIspService
type DescribeRangeDataByLocateAndIspServiceResponse struct {
	*responses.BaseResponse
	RequestId  string `json:"RequestId" xml:"RequestId"`
	JsonResult string `json:"JsonResult" xml:"JsonResult"`
}

// CreateDescribeRangeDataByLocateAndIspServiceRequest creates a request to invoke DescribeRangeDataByLocateAndIspService API
func CreateDescribeRangeDataByLocateAndIspServiceRequest() (request *DescribeRangeDataByLocateAndIspServiceRequest) {
	request = &DescribeRangeDataByLocateAndIspServiceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeRangeDataByLocateAndIspService", "", "")
	return
}

// CreateDescribeRangeDataByLocateAndIspServiceResponse creates a response to parse from DescribeRangeDataByLocateAndIspService response
func CreateDescribeRangeDataByLocateAndIspServiceResponse() (response *DescribeRangeDataByLocateAndIspServiceResponse) {
	response = &DescribeRangeDataByLocateAndIspServiceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
