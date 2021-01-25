/*
 * ELB
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 */

package model

import (
	"encoding/json"

	"strings"
)

// 健康检查对象。
type UpdateHealthMonitorOption struct {
	// 功能说明：管理状态true/false。使用说明：默认为true，true表示开启健康检查，false表示关闭健康检查。
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
	// 健康检查间隔。
	Delay *int32 `json:"delay,omitempty"`
	// 功能说明：健康检查测试member健康状态时，发送的http请求的域名。仅当type为HTTP时生效。使用说明：默认为空，表示使用负载均衡器的vip作为http请求的目的地址。以数字或字母开头，只能包含数字、字母、’-’、’.’。
	DomainName *string `json:"domain_name,omitempty"`
	// 期望HTTP响应状态码，指定下列值：单值，例如200；列表，例如200，202；区间，例如200-204。仅当type为HTTP时生效。该字段为预留字段，暂未启用。
	ExpectedCodes *string `json:"expected_codes,omitempty"`
	// HTTP方法，可以为GET、HEAD、POST、PUT、DELETE、TRACE、OPTIONS、CONNECT、PATCH。仅当type为HTTP时生效。该字段为预留字段，暂未启用。
	HttpMethod *string `json:"http_method,omitempty"`
	// 最大重试次数
	MaxRetries *int32 `json:"max_retries,omitempty"`
	// 健康检查连续成功多少次后，将后端服务器的健康检查状态由ONLIEN判定为OFFLINE
	MaxRetriesDown *int32 `json:"max_retries_down,omitempty"`
	// 健康检查端口号。默认为空，表示使用后端云服务器组的端口。
	MonitorPort *int32 `json:"monitor_port,omitempty"`
	// 健康检查名称。
	Name *string `json:"name,omitempty"`
	// 健康检查的超时时间。建议该值小于delay的值。
	Timeout *int32 `json:"timeout,omitempty"`
	// 功能说明：健康检查测试member健康时发送的http请求路径。默认为“/”。使用说明：以“/”开头。仅当type为HTTP时生效。
	UrlPath *string `json:"url_path,omitempty"`
	// 描述：健康检查类型。   取值：TCP,UDP_CONNECT,HTTP,HTTPS,PING   约束：   1、若pool的protocol为QUIC，则type只能是UDP
	Type *string `json:"type,omitempty"`
}

func (o UpdateHealthMonitorOption) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return "UpdateHealthMonitorOption struct{}"
	}

	return strings.Join([]string{"UpdateHealthMonitorOption", string(data)}, " ")
}
