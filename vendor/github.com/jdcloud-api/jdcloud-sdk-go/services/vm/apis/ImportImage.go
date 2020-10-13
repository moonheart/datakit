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

package apis

import (
    "github.com/jdcloud-api/jdcloud-sdk-go/core"
)

type ImportImageRequest struct {

    core.JDCloudRequest

    /* 地域ID  */
    RegionId string `json:"regionId"`

    /* 系统架构，可选值：x86_64,i386  */
    Architecture string `json:"architecture"`

    /* 操作系统，可选值：windows,linux  */
    OsType string `json:"osType"`

    /* 平台名称，可选值：CentOS,Ubuntu,Windows Server,Other Linux,Other Windows  */
    Platform string `json:"platform"`

    /* 磁盘格式，可选值：qcow2,vhd,vmdk,raw  */
    DiskFormat string `json:"diskFormat"`

    /* 以此镜像需要制作的系统盘的默认大小，单位GB。最小值40，最大值500，要求值是10的整数倍  */
    SystemDiskSizeGB int `json:"systemDiskSizeGB"`

    /* 要导入镜像的对象存储外链地址  */
    ImageUrl string `json:"imageUrl"`

    /* 镜像的操作系统版本 (Optional) */
    OsVersion *string `json:"osVersion"`

    /* 导入镜像的自定义名称  */
    ImageName string `json:"imageName"`

    /* 导入镜像的描述信息 (Optional) */
    Description *string `json:"description"`

    /* 是否强制导入。强制导入则忽略镜像的合规性检测 (Optional) */
    ForceImport *bool `json:"forceImport"`

    /* 用户导入镜像的幂等性保证。每次创建请传入不同的值，如果传值与某次的clientToken相同，则返还该次的请求结果 (Optional) */
    ClientToken *string `json:"clientToken"`
}

/*
 * param regionId: 地域ID (Required)
 * param architecture: 系统架构，可选值：x86_64,i386 (Required)
 * param osType: 操作系统，可选值：windows,linux (Required)
 * param platform: 平台名称，可选值：CentOS,Ubuntu,Windows Server,Other Linux,Other Windows (Required)
 * param diskFormat: 磁盘格式，可选值：qcow2,vhd,vmdk,raw (Required)
 * param systemDiskSizeGB: 以此镜像需要制作的系统盘的默认大小，单位GB。最小值40，最大值500，要求值是10的整数倍 (Required)
 * param imageUrl: 要导入镜像的对象存储外链地址 (Required)
 * param imageName: 导入镜像的自定义名称 (Required)
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewImportImageRequest(
    regionId string,
    architecture string,
    osType string,
    platform string,
    diskFormat string,
    systemDiskSizeGB int,
    imageUrl string,
    imageName string,
) *ImportImageRequest {

	return &ImportImageRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/regions/{regionId}/images:import",
			Method:  "POST",
			Header:  nil,
			Version: "v1",
		},
        RegionId: regionId,
        Architecture: architecture,
        OsType: osType,
        Platform: platform,
        DiskFormat: diskFormat,
        SystemDiskSizeGB: systemDiskSizeGB,
        ImageUrl: imageUrl,
        ImageName: imageName,
	}
}

/*
 * param regionId: 地域ID (Required)
 * param architecture: 系统架构，可选值：x86_64,i386 (Required)
 * param osType: 操作系统，可选值：windows,linux (Required)
 * param platform: 平台名称，可选值：CentOS,Ubuntu,Windows Server,Other Linux,Other Windows (Required)
 * param diskFormat: 磁盘格式，可选值：qcow2,vhd,vmdk,raw (Required)
 * param systemDiskSizeGB: 以此镜像需要制作的系统盘的默认大小，单位GB。最小值40，最大值500，要求值是10的整数倍 (Required)
 * param imageUrl: 要导入镜像的对象存储外链地址 (Required)
 * param osVersion: 镜像的操作系统版本 (Optional)
 * param imageName: 导入镜像的自定义名称 (Required)
 * param description: 导入镜像的描述信息 (Optional)
 * param forceImport: 是否强制导入。强制导入则忽略镜像的合规性检测 (Optional)
 * param clientToken: 用户导入镜像的幂等性保证。每次创建请传入不同的值，如果传值与某次的clientToken相同，则返还该次的请求结果 (Optional)
 */
func NewImportImageRequestWithAllParams(
    regionId string,
    architecture string,
    osType string,
    platform string,
    diskFormat string,
    systemDiskSizeGB int,
    imageUrl string,
    osVersion *string,
    imageName string,
    description *string,
    forceImport *bool,
    clientToken *string,
) *ImportImageRequest {

    return &ImportImageRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/images:import",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
        RegionId: regionId,
        Architecture: architecture,
        OsType: osType,
        Platform: platform,
        DiskFormat: diskFormat,
        SystemDiskSizeGB: systemDiskSizeGB,
        ImageUrl: imageUrl,
        OsVersion: osVersion,
        ImageName: imageName,
        Description: description,
        ForceImport: forceImport,
        ClientToken: clientToken,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewImportImageRequestWithoutParam() *ImportImageRequest {

    return &ImportImageRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/images:import",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
    }
}

/* param regionId: 地域ID(Required) */
func (r *ImportImageRequest) SetRegionId(regionId string) {
    r.RegionId = regionId
}

/* param architecture: 系统架构，可选值：x86_64,i386(Required) */
func (r *ImportImageRequest) SetArchitecture(architecture string) {
    r.Architecture = architecture
}

/* param osType: 操作系统，可选值：windows,linux(Required) */
func (r *ImportImageRequest) SetOsType(osType string) {
    r.OsType = osType
}

/* param platform: 平台名称，可选值：CentOS,Ubuntu,Windows Server,Other Linux,Other Windows(Required) */
func (r *ImportImageRequest) SetPlatform(platform string) {
    r.Platform = platform
}

/* param diskFormat: 磁盘格式，可选值：qcow2,vhd,vmdk,raw(Required) */
func (r *ImportImageRequest) SetDiskFormat(diskFormat string) {
    r.DiskFormat = diskFormat
}

/* param systemDiskSizeGB: 以此镜像需要制作的系统盘的默认大小，单位GB。最小值40，最大值500，要求值是10的整数倍(Required) */
func (r *ImportImageRequest) SetSystemDiskSizeGB(systemDiskSizeGB int) {
    r.SystemDiskSizeGB = systemDiskSizeGB
}

/* param imageUrl: 要导入镜像的对象存储外链地址(Required) */
func (r *ImportImageRequest) SetImageUrl(imageUrl string) {
    r.ImageUrl = imageUrl
}

/* param osVersion: 镜像的操作系统版本(Optional) */
func (r *ImportImageRequest) SetOsVersion(osVersion string) {
    r.OsVersion = &osVersion
}

/* param imageName: 导入镜像的自定义名称(Required) */
func (r *ImportImageRequest) SetImageName(imageName string) {
    r.ImageName = imageName
}

/* param description: 导入镜像的描述信息(Optional) */
func (r *ImportImageRequest) SetDescription(description string) {
    r.Description = &description
}

/* param forceImport: 是否强制导入。强制导入则忽略镜像的合规性检测(Optional) */
func (r *ImportImageRequest) SetForceImport(forceImport bool) {
    r.ForceImport = &forceImport
}

/* param clientToken: 用户导入镜像的幂等性保证。每次创建请传入不同的值，如果传值与某次的clientToken相同，则返还该次的请求结果(Optional) */
func (r *ImportImageRequest) SetClientToken(clientToken string) {
    r.ClientToken = &clientToken
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r ImportImageRequest) GetRegionId() string {
    return r.RegionId
}

type ImportImageResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result ImportImageResult `json:"result"`
}

type ImportImageResult struct {
    ImageId string `json:"imageId"`
    ImportTaskId int `json:"importTaskId"`
}