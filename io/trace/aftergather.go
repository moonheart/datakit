package trace

import (
	"sync"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	dkio "gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
)

var once = sync.Once{}

// CalculatorFunc is func type for calculation, statistics, etc
// any data changes in DatakitTraces will be saved and affect the next actions afterwards.
type CalculatorFunc func(dktrace DatakitTrace)

// FilterFunc is func type for data filter.
// Return the DatakitTraces that need to propagate to next action and
// return ture if one want to skip all FilterFunc afterwards, false otherwise.
type FilterFunc func(dktrace DatakitTrace) (DatakitTrace, bool)

type AfterGather struct {
	sync.Mutex
	calculators []CalculatorFunc
	filters     []FilterFunc
}

func NewAfterGather() *AfterGather {
	return &AfterGather{}
}

// AppendCalculator will append new calculators into AfterGather structure,
// and run them as the order they added.
func (aga *AfterGather) AppendCalculator(calc ...CalculatorFunc) {
	aga.Lock()
	defer aga.Unlock()

	aga.calculators = append(aga.calculators, calc...)
}

// AppendFilter will append new filters into AfterGather structure,
// and run them as the order they added if no filter func return false
// to break out the filters loop.
func (aga *AfterGather) AppendFilter(filter ...FilterFunc) {
	aga.Lock()
	defer aga.Unlock()

	aga.filters = append(aga.filters, filter...)
}

func (aga *AfterGather) Run(inputName string, dktrace DatakitTrace, stricktMod bool) {
	once.Do(func() {
		log = logger.SLogger(packageName)
	})

	if inputName == "" || len(dktrace) == 0 {
		log.Warnf("wrong parameters for AfterGather.Run(inputName: %s, dktrace:%v)", inputName, dktrace)

		return
	}

	for i := range aga.calculators {
		aga.calculators[i](dktrace)
	}
	var skip bool
	for i := range aga.filters {
		if dktrace, skip = aga.filters[i](dktrace); skip {
			break
		}
	}
	if dktrace == nil {
		return
	}

	if pts := BuildPointsBatch(inputName, dktrace, stricktMod); len(pts) != 0 {
		if err := dkio.Feed(inputName, datakit.Tracing, pts, &dkio.Option{HighFreq: true}); err != nil {
			log.Errorf("io feed points error: %s", err.Error())
		}
	} else {
		log.Warn("BuildPointsBatch return empty points array")
	}
}

// BuildPointsBatch builds points from whole trace
func BuildPointsBatch(inputName string, dktrace DatakitTrace, strict bool) []*dkio.Point {
	var pts []*dkio.Point
	for i := range dktrace {
		if pt, err := BuildPoint(dktrace[i], strict); err != nil {
			log.Errorf("build point error: %s", err.Error())
		} else {
			pts = append(pts, pt)
		}
	}

	return pts
}

// BuildPoint builds point from DatakitSpan.
func BuildPoint(dkspan *DatakitSpan, strict bool) (*dkio.Point, error) {
	var (
		tags   = make(map[string]string)
		fields = make(map[string]interface{})
	)

	tags[TAG_PROJECT] = dkspan.Project
	tags[TAG_OPERATION] = dkspan.Operation
	tags[TAG_SERVICE] = dkspan.Service
	tags[TAG_VERSION] = dkspan.Version
	tags[TAG_ENV] = dkspan.Env
	tags[TAG_HTTP_METHOD] = dkspan.HTTPMethod
	tags[TAG_HTTP_CODE] = dkspan.HTTPStatusCode

	if dkspan.SourceType != "" {
		tags[TAG_TYPE] = dkspan.SourceType
	} else {
		tags[TAG_TYPE] = SPAN_SERVICE_CUSTOM
	}

	for k, v := range dkspan.Tags {
		tags[k] = v
	}

	tags[TAG_SPAN_STATUS] = dkspan.Status

	if dkspan.EndPoint != "" {
		tags[TAG_ENDPOINT] = dkspan.EndPoint
	} else {
		tags[TAG_ENDPOINT] = "null"
	}

	if dkspan.SpanType != "" {
		tags[TAG_SPAN_TYPE] = dkspan.SpanType
	} else {
		tags[TAG_SPAN_TYPE] = SPAN_TYPE_ENTRY
	}

	if dkspan.ContainerHost != "" {
		tags[TAG_CONTAINER_HOST] = dkspan.ContainerHost
	}

	if dkspan.ParentID == "" {
		dkspan.ParentID = "0"
	}

	fields[FIELD_START] = dkspan.Start / int64(time.Microsecond)
	fields[FIELD_DURATION] = dkspan.Duration / int64(time.Microsecond)
	fields[FIELD_MSG] = dkspan.Content
	fields[FIELD_RESOURCE] = dkspan.Resource
	fields[FIELD_PARENTID] = dkspan.ParentID
	fields[FIELD_TRACEID] = dkspan.TraceID
	fields[FIELD_SPANID] = dkspan.SpanID

	return dkio.NewPoint(dkspan.Source, tags, fields, &dkio.PointOption{
		Time:     time.Unix(0, dkspan.Start),
		Category: datakit.Tracing,
		Strict:   strict,
	})
}
