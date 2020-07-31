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

// DescribeLiveSnapshotConfig invokes the cdn.DescribeLiveSnapshotConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/describelivesnapshotconfig.html
func (client *Client) DescribeLiveSnapshotConfig(request *DescribeLiveSnapshotConfigRequest) (response *DescribeLiveSnapshotConfigResponse, err error) {
	response = CreateDescribeLiveSnapshotConfigResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLiveSnapshotConfigWithChan invokes the cdn.DescribeLiveSnapshotConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivesnapshotconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveSnapshotConfigWithChan(request *DescribeLiveSnapshotConfigRequest) (<-chan *DescribeLiveSnapshotConfigResponse, <-chan error) {
	responseChan := make(chan *DescribeLiveSnapshotConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLiveSnapshotConfig(request)
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

// DescribeLiveSnapshotConfigWithCallback invokes the cdn.DescribeLiveSnapshotConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivesnapshotconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveSnapshotConfigWithCallback(request *DescribeLiveSnapshotConfigRequest, callback func(response *DescribeLiveSnapshotConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLiveSnapshotConfigResponse
		var err error
		defer close(result)
		response, err = client.DescribeLiveSnapshotConfig(request)
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

// DescribeLiveSnapshotConfigRequest is the request struct for api DescribeLiveSnapshotConfig
type DescribeLiveSnapshotConfigRequest struct {
	*requests.RpcRequest
	PageNum       requests.Integer `position:"Query" name:"PageNum"`
	AppName       string           `position:"Query" name:"AppName"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	StreamName    string           `position:"Query" name:"StreamName"`
	Order         string           `position:"Query" name:"Order"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeLiveSnapshotConfigResponse is the response struct for api DescribeLiveSnapshotConfig
type DescribeLiveSnapshotConfigResponse struct {
	*responses.BaseResponse
	RequestId                    string                       `json:"RequestId" xml:"RequestId"`
	PageNum                      int                          `json:"PageNum" xml:"PageNum"`
	PageSize                     int                          `json:"PageSize" xml:"PageSize"`
	Order                        string                       `json:"Order" xml:"Order"`
	TotalNum                     int                          `json:"TotalNum" xml:"TotalNum"`
	TotalPage                    int                          `json:"TotalPage" xml:"TotalPage"`
	LiveStreamSnapshotConfigList LiveStreamSnapshotConfigList `json:"LiveStreamSnapshotConfigList" xml:"LiveStreamSnapshotConfigList"`
}

// CreateDescribeLiveSnapshotConfigRequest creates a request to invoke DescribeLiveSnapshotConfig API
func CreateDescribeLiveSnapshotConfigRequest() (request *DescribeLiveSnapshotConfigRequest) {
	request = &DescribeLiveSnapshotConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeLiveSnapshotConfig", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeLiveSnapshotConfigResponse creates a response to parse from DescribeLiveSnapshotConfig response
func CreateDescribeLiveSnapshotConfigResponse() (response *DescribeLiveSnapshotConfigResponse) {
	response = &DescribeLiveSnapshotConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
