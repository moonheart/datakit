// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package io

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb1-client/models"
	lp "gitlab.jiagouyun.com/cloudcare-tools/cliutils/lineproto"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
)

var (
	// For logging, we use measurement-name as source value
	// in kodo, so there should not be any tag/field named
	// with `source`.
	// For object, we use measurement-name as class value
	// in kodo, so there should not be any tag/field named
	// with `class`.
	DisabledTagKeys = map[string][]string{
		datakit.Logging: {"source"},
		datakit.Object:  {"class"},
		// others data type not set...
	}

	DisabledFieldKeys = map[string][]string{
		datakit.Logging: {"source"},
		datakit.Object:  {"class"},
		// others data type not set...
	}

	Callback func(models.Point) (models.Point, error) = nil

	Strict        = true
	MaxTags   int = 256  // limit tag count
	MaxFields int = 1024 // limit field count

	// limit tag/field key/value length.
	MaxTagKeyLen     int = 256
	MaxFieldKeyLen   int = 256
	MaxTagValueLen   int = 1024
	MaxFieldValueLen int = 32 * 1024 // if field value is string,limit to 32K

	Precision string = "n"
)

func SetExtraTags(k, v string) {
	extraTags[k] = v
}

func doMakePoint(name string,
	tags map[string]string,
	fields map[string]interface{},
	opt *lp.Option) (*Point, error) {
	p, warnings, err := lp.MakeLineProtoPointWithWarnings(name, tags, fields, opt)

	if err != nil {
		return nil, err
	} else if len(warnings) > 0 {
		warningsStr := ""
		for _, warn := range warnings {
			warningsStr += warn.Message + ";"
		}
		l.Warnf("make metric(%s) point successfully but with warnings: %s", name, warningsStr)
	}

	return &Point{Point: p}, nil
}

type PointOption struct {
	Time              time.Time
	Category          string
	DisableGlobalTags bool
	Strict            bool
	MaxFieldValueLen  int
}

func defaultPointOption() *PointOption {
	return &PointOption{
		Time:     time.Now(),
		Category: datakit.Metric,
		Strict:   true,
	}
}

func NewPoint(name string,
	tags map[string]string,
	fields map[string]interface{},
	opt *PointOption) (*Point, error) {
	if opt == nil {
		opt = defaultPointOption()
	}

	lpOpt := &lp.Option{
		Time:      opt.Time,
		Strict:    opt.Strict,
		Precision: "n",

		MaxTags:   MaxTags,
		MaxFields: MaxFields,
		ExtraTags: extraTags,

		MaxTagKeyLen:     MaxTagKeyLen,
		MaxFieldKeyLen:   MaxFieldKeyLen,
		MaxTagValueLen:   MaxTagValueLen,
		MaxFieldValueLen: MaxFieldValueLen,

		// not set
		DisabledTagKeys:   nil,
		DisabledFieldKeys: nil,
		Callback:          nil,
	}

	if opt.DisableGlobalTags {
		lpOpt.ExtraTags = nil
	}
	if opt.MaxFieldValueLen > 0 {
		lpOpt.MaxFieldValueLen = opt.MaxFieldValueLen
	}
	switch opt.Category {
	case datakit.Metric:
		lpOpt.EnablePointInKey = true
		lpOpt.DisabledTagKeys = DisabledTagKeys[opt.Category]
		lpOpt.DisabledFieldKeys = DisabledFieldKeys[opt.Category]
		lpOpt.DisableStringField = true // ingore string field value in metric point
	case datakit.Network,
		datakit.KeyEvent,
		datakit.Object,
		datakit.CustomObject,
		datakit.Logging,
		datakit.Tracing,
		datakit.RUM,
		datakit.Security:
		lpOpt.DisabledTagKeys = DisabledTagKeys[opt.Category]
		lpOpt.DisabledFieldKeys = DisabledFieldKeys[opt.Category]
	default:
		return nil, fmt.Errorf("invalid point category: %s", opt.Category)
	}
	return doMakePoint(name, tags, fields, lpOpt)
}

// MakePoint Deprecated.
func MakePoint(name string,
	tags map[string]string,
	fields map[string]interface{},
	t ...time.Time) (*Point, error) {
	lpOpt := &lp.Option{
		Strict:    true,
		Precision: "n",

		MaxTags:   MaxTags,
		MaxFields: MaxFields,
		ExtraTags: extraTags,

		MaxTagKeyLen:     MaxTagKeyLen,
		MaxFieldKeyLen:   MaxFieldKeyLen,
		MaxTagValueLen:   MaxTagValueLen,
		MaxFieldValueLen: MaxFieldValueLen,

		// not set
		DisabledTagKeys:   nil,
		DisabledFieldKeys: nil,
		Callback:          nil,
	}

	if len(t) > 0 {
		lpOpt.Time = t[0]
	} else {
		lpOpt.Time = time.Now().UTC()
	}

	return doMakePoint(name, tags, fields, lpOpt)
}

// MakeMetric Deprecated.
func MakeMetric(name string,
	tags map[string]string,
	fields map[string]interface{},
	t ...time.Time) ([]byte, error) {
	p, err := MakePoint(name, tags, fields, t...)
	if err != nil {
		return nil, err
	}

	return []byte(p.Point.String()), nil
}
