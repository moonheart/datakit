/*
 * IAM
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 */

package model

import (
	"encoding/json"

	"strings"
)

//
type KeystoneUpdateUserPasswordRequestBody struct {
	User *KeystoneUpdatePasswordOption `json:"user"`
}

func (o KeystoneUpdateUserPasswordRequestBody) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return "KeystoneUpdateUserPasswordRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateUserPasswordRequestBody", string(data)}, " ")
}
