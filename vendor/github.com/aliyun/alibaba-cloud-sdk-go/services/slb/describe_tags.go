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

// DescribeTags invokes the slb.DescribeTags API synchronously
// api document: https://help.aliyun.com/api/slb/describetags.html
func (client *Client) DescribeTags(request *DescribeTagsRequest) (response *DescribeTagsResponse, err error) {
	response = CreateDescribeTagsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeTagsWithChan invokes the slb.DescribeTags API asynchronously
// api document: https://help.aliyun.com/api/slb/describetags.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeTagsWithChan(request *DescribeTagsRequest) (<-chan *DescribeTagsResponse, <-chan error) {
	responseChan := make(chan *DescribeTagsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeTags(request)
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

// DescribeTagsWithCallback invokes the slb.DescribeTags API asynchronously
// api document: https://help.aliyun.com/api/slb/describetags.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeTagsWithCallback(request *DescribeTagsRequest, callback func(response *DescribeTagsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeTagsResponse
		var err error
		defer close(result)
		response, err = client.DescribeTags(request)
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

// DescribeTagsRequest is the request struct for api DescribeTags
type DescribeTagsRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	DistinctKey          requests.Boolean `position:"Query" name:"DistinctKey"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Tags                 string           `position:"Query" name:"Tags"`
	LoadBalancerId       string           `position:"Query" name:"LoadBalancerId"`
}

// DescribeTagsResponse is the response struct for api DescribeTags
type DescribeTagsResponse struct {
	*responses.BaseResponse
	RequestId  string  `json:"RequestId" xml:"RequestId"`
	PageSize   int     `json:"PageSize" xml:"PageSize"`
	PageNumber int     `json:"PageNumber" xml:"PageNumber"`
	TotalCount int     `json:"TotalCount" xml:"TotalCount"`
	TagSets    TagSets `json:"TagSets" xml:"TagSets"`
}

// CreateDescribeTagsRequest creates a request to invoke DescribeTags API
func CreateDescribeTagsRequest() (request *DescribeTagsRequest) {
	request = &DescribeTagsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "DescribeTags", "slb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeTagsResponse creates a response to parse from DescribeTags response
func CreateDescribeTagsResponse() (response *DescribeTagsResponse) {
	response = &DescribeTagsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
