// Copyright 2018 JDCLOUD.COM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// NOTE: This class is auto generated by the jdcloud code generator program.

package models


type UpdateWebHookReq struct {

    /* 是否启用, 1表示启用webHook，0表示禁用webHook，默认为1 (Optional) */
    Active int64 `json:"active"`

    /* webHook content (Optional) */
    Content string `json:"content"`

    /* webHook 协议,http或者https  */
    Protocol string `json:"protocol"`

    /* webHook secret，用户请求签名，防伪造 (Optional) */
    Secret string `json:"secret"`

    /* webHook url  */
    Url string `json:"url"`
}
