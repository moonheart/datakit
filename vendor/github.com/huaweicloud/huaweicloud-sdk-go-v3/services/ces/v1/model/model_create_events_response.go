/*
 * CES
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 */

package model

import (
	"encoding/json"

	"strings"
)

// Response Object
type CreateEventsResponse struct {
	// 响应参数。
	Body           *[]CreateEventsResponseBody `json:"body,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o CreateEventsResponse) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return "CreateEventsResponse struct{}"
	}

	return strings.Join([]string{"CreateEventsResponse", string(data)}, " ")
}
