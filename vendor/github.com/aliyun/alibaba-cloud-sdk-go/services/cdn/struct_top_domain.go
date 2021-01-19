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

// TopDomain is a nested struct in cdn response
type TopDomain struct {
	DomainName     string  `json:"DomainName" xml:"DomainName"`
	Rank           int64   `json:"Rank" xml:"Rank"`
	TotalTraffic   string  `json:"TotalTraffic" xml:"TotalTraffic"`
	TrafficPercent string  `json:"TrafficPercent" xml:"TrafficPercent"`
	MaxBps         float64 `json:"MaxBps" xml:"MaxBps"`
	MaxBpsTime     string  `json:"MaxBpsTime" xml:"MaxBpsTime"`
	TotalAccess    int64   `json:"TotalAccess" xml:"TotalAccess"`
}
