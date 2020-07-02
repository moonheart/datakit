package prometheus

// Parser inspired from
// https://github.com/prometheus/prom2json/blob/master/main.go

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"strings"
	"time"

	influxdb "github.com/influxdata/influxdb1-client/v2"

	"github.com/influxdata/telegraf"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func convertName(mf *dto.MetricFamily) (metricname string, valuename string) {
	rawname := mf.GetName()

	if promAgent.MetricNameMap != nil {
		if mapname, ok := promAgent.MetricNameMap[rawname]; ok && mapname != "" {
			metricname = mapname
			if strings.HasPrefix(rawname, mapname) {
				valuename = rawname[len(mapname):]
				if valuename[0] == '_' {
					valuename = valuename[1:]
				}
			} else {
				valuename = "value"
			}
			return
		}
	}

	parts := strings.Split(rawname, "_")

	if len(parts) > 2 {
		metricname = strings.Join(parts[:2], "_")
		valuename = strings.Join(parts[2:], "_")
	} else {
		metricname = rawname
		valuename = "value"
	}
	return
}

// Parse returns a slice of Metrics from a text representation of a
// metrics
func ParseV2(buf []byte, header http.Header) ([]*influxdb.Point, error) {
	var metrics []*influxdb.Point
	var parser expfmt.TextParser
	// parse even if the buffer begins with a newline
	buf = bytes.TrimPrefix(buf, []byte("\n"))
	// Read raw data
	buffer := bytes.NewBuffer(buf)
	reader := bufio.NewReader(buffer)

	mediatype, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	// Prepare output
	metricFamilies := make(map[string]*dto.MetricFamily)

	if err == nil && mediatype == "application/vnd.google.protobuf" &&
		params["encoding"] == "delimited" &&
		params["proto"] == "io.prometheus.client.MetricFamily" {
		for {
			mf := &dto.MetricFamily{}
			if _, ierr := pbutil.ReadDelimited(reader, mf); ierr != nil {
				if ierr == io.EOF {
					break
				}
				return nil, fmt.Errorf("reading metric family protocol buffer failed: %s", ierr)
			}
			metricFamilies[mf.GetName()] = mf
		}
	} else {
		metricFamilies, err = parser.TextToMetricFamilies(reader)
		if err != nil {
			return nil, fmt.Errorf("reading text format failed: %s", err)
		}
	}

	// read metrics
	for metricName, mf := range metricFamilies {
		for _, m := range mf.Metric {
			// reading tags
			tags := makeLabels(m)

			if mf.GetType() == dto.MetricType_SUMMARY {
				// summary metric
				telegrafMetrics := makeQuantilesV2(m, tags, metricName, mf.GetType())
				metrics = append(metrics, telegrafMetrics...)
			} else if mf.GetType() == dto.MetricType_HISTOGRAM {
				// histogram metric
				telegrafMetrics := makeBucketsV2(m, tags, metricName, mf.GetType())
				metrics = append(metrics, telegrafMetrics...)
			} else {
				// standard metric
				// reading fields
				fields := make(map[string]interface{})
				fields = getNameAndValueV2(m, metricName)
				// converting to telegraf metric
				if len(fields) > 0 {
					var t time.Time
					if m.TimestampMs != nil && *m.TimestampMs > 0 {
						t = time.Unix(0, *m.TimestampMs*1000000)
					} else {
						t = time.Now()
					}

					pt, err := influxdb.NewPoint(`prometheus`, tags, fields, t)
					if err == nil {
						metrics = append(metrics, pt)
					}

				}
			}
		}
	}

	return metrics, err
}

// Get Quantiles for summary metric & Buckets for histogram
func makeQuantilesV2(m *dto.Metric, tags map[string]string, metricName string, metricType dto.MetricType) []*influxdb.Point {
	var metrics []*influxdb.Point
	fields := make(map[string]interface{})
	var t time.Time
	if m.TimestampMs != nil && *m.TimestampMs > 0 {
		t = time.Unix(0, *m.TimestampMs*1000000)
	} else {
		t = time.Now()
	}
	fields[metricName+"_count"] = float64(m.GetSummary().GetSampleCount())
	fields[metricName+"_sum"] = float64(m.GetSummary().GetSampleSum())
	//met, err := metric.New("prometheus", tags, fields, t, valueType(metricType))
	pt, err := influxdb.NewPoint(`prometheus`, tags, fields, t)
	if err == nil {
		metrics = append(metrics, pt)
	}

	for _, q := range m.GetSummary().Quantile {
		newTags := tags
		fields = make(map[string]interface{})

		newTags["quantile"] = fmt.Sprint(q.GetQuantile())
		fields[metricName] = float64(q.GetValue())

		//quantileMetric, err := metric.New("prometheus", newTags, fields, t, valueType(metricType))
		pt, err := influxdb.NewPoint(`prometheus`, newTags, fields, t)
		if err == nil {
			metrics = append(metrics, pt)
		}
	}
	return metrics
}

// Get Buckets  from histogram metric
func makeBucketsV2(m *dto.Metric, tags map[string]string, metricName string, metricType dto.MetricType) []*influxdb.Point {
	var metrics []*influxdb.Point
	fields := make(map[string]interface{})
	var t time.Time
	if m.TimestampMs != nil && *m.TimestampMs > 0 {
		t = time.Unix(0, *m.TimestampMs*1000000)
	} else {
		t = time.Now()
	}
	fields[metricName+"_count"] = float64(m.GetHistogram().GetSampleCount())
	fields[metricName+"_sum"] = float64(m.GetHistogram().GetSampleSum())

	//met, err := metric.New("prometheus", tags, fields, t, valueType(metricType))
	pt, err := influxdb.NewPoint(`prometheus`, tags, fields, t)
	if err == nil {
		metrics = append(metrics, pt)
	}

	for _, b := range m.GetHistogram().Bucket {
		newTags := tags
		fields = make(map[string]interface{})
		newTags["le"] = fmt.Sprint(b.GetUpperBound())
		fields[metricName+"_bucket"] = float64(b.GetCumulativeCount())

		//histogramMetric, err := metric.New("prometheus", newTags, fields, t, valueType(metricType))
		pt, err := influxdb.NewPoint(`prometheus`, newTags, fields, t)
		if err == nil {
			metrics = append(metrics, pt)
		}
	}
	return metrics
}

// Parse returns a slice of Metrics from a text representation of a
// metrics
func Parse(buf []byte, header http.Header) ([]*influxdb.Point, error) {
	var metrics []*influxdb.Point
	var parser expfmt.TextParser
	// parse even if the buffer begins with a newline
	buf = bytes.TrimPrefix(buf, []byte("\n"))
	// Read raw data
	buffer := bytes.NewBuffer(buf)
	reader := bufio.NewReader(buffer)

	mediatype, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	// Prepare output
	metricFamilies := make(map[string]*dto.MetricFamily)

	if err == nil && mediatype == "application/vnd.google.protobuf" &&
		params["encoding"] == "delimited" &&
		params["proto"] == "io.prometheus.client.MetricFamily" {
		for {
			mf := &dto.MetricFamily{}
			if _, ierr := pbutil.ReadDelimited(reader, mf); ierr != nil {
				if ierr == io.EOF {
					break
				}
				return nil, fmt.Errorf("reading metric family protocol buffer failed: %s", ierr)
			}
			metricFamilies[mf.GetName()] = mf
		}
	} else {
		metricFamilies, err = parser.TextToMetricFamilies(reader)
		if err != nil {
			return nil, fmt.Errorf("reading text format failed: %s", err)
		}
	}

	// read metrics
	for metricName, mf := range metricFamilies {
		_ = metricName
		finalMetricName, valuename := convertName(mf)
		for _, m := range mf.Metric {
			// reading tags
			tags := makeLabels(m)
			// reading fields
			fields := make(map[string]interface{})
			if mf.GetType() == dto.MetricType_SUMMARY {
				// summary metric
				fields = makeQuantiles(m)
				fields["count"] = float64(m.GetSummary().GetSampleCount())
				fields["sum"] = float64(m.GetSummary().GetSampleSum())
			} else if mf.GetType() == dto.MetricType_HISTOGRAM {
				// histogram metric
				fields = makeBuckets(m)
				fields["count"] = float64(m.GetHistogram().GetSampleCount())
				fields["sum"] = float64(m.GetHistogram().GetSampleSum())

			} else {
				// standard metric
				fields = getNameAndValue(m, valuename)
			}
			// converting to telegraf metric
			if len(fields) > 0 {
				var t time.Time
				if m.TimestampMs != nil && *m.TimestampMs > 0 {
					t = time.Unix(0, *m.TimestampMs*1000000)
				} else {
					t = time.Now()
				}
				//metric, err := metric.New(finalMetricName, tags, fields, t, valueType(mf.GetType()))
				pt, err := influxdb.NewPoint(finalMetricName, tags, fields, t)
				if err == nil {
					metrics = append(metrics, pt)
				}
			}
		}
	}

	return metrics, err
}

func valueType(mt dto.MetricType) telegraf.ValueType {
	switch mt {
	case dto.MetricType_COUNTER:
		return telegraf.Counter
	case dto.MetricType_GAUGE:
		return telegraf.Gauge
	case dto.MetricType_SUMMARY:
		return telegraf.Summary
	case dto.MetricType_HISTOGRAM:
		return telegraf.Histogram
	default:
		return telegraf.Untyped
	}
}

// Get Quantiles from summary metric
func makeQuantiles(m *dto.Metric) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, q := range m.GetSummary().Quantile {
		if !math.IsNaN(q.GetValue()) {
			fields[fmt.Sprint(q.GetQuantile())] = float64(q.GetValue())
		}
	}
	return fields
}

// Get Buckets  from histogram metric
func makeBuckets(m *dto.Metric) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, b := range m.GetHistogram().Bucket {
		fields[fmt.Sprint(b.GetUpperBound())] = float64(b.GetCumulativeCount())
	}
	return fields
}

// Get labels from metric
func makeLabels(m *dto.Metric) map[string]string {
	result := map[string]string{}
	for _, lp := range m.Label {
		result[lp.GetName()] = lp.GetValue()
	}
	return result
}

// Get name and value from metric
// func getNameAndValue(m *dto.Metric) map[string]interface{} {
// 	fields := make(map[string]interface{})
// 	if m.Gauge != nil {
// 		if !math.IsNaN(m.GetGauge().GetValue()) {
// 			fields["gauge"] = float64(m.GetGauge().GetValue())
// 		}
// 	} else if m.Counter != nil {
// 		if !math.IsNaN(m.GetCounter().GetValue()) {
// 			fields["counter"] = float64(m.GetCounter().GetValue())
// 		}
// 	} else if m.Untyped != nil {
// 		if !math.IsNaN(m.GetUntyped().GetValue()) {
// 			fields["value"] = float64(m.GetUntyped().GetValue())
// 		}
// 	}
// 	return fields
// }

// Get name and value from metric
func getNameAndValueV2(m *dto.Metric, metricName string) map[string]interface{} {
	fields := make(map[string]interface{})
	if m.Gauge != nil {
		if !math.IsNaN(m.GetGauge().GetValue()) {
			fields[metricName] = float64(m.GetGauge().GetValue())
		}
	} else if m.Counter != nil {
		if !math.IsNaN(m.GetCounter().GetValue()) {
			fields[metricName] = float64(m.GetCounter().GetValue())
		}
	} else if m.Untyped != nil {
		if !math.IsNaN(m.GetUntyped().GetValue()) {
			fields[metricName] = float64(m.GetUntyped().GetValue())
		}
	}
	return fields
}

func getNameAndValue(m *dto.Metric, fieldname string) map[string]interface{} {
	fields := make(map[string]interface{})

	if m.Gauge != nil {
		if !math.IsNaN(m.GetGauge().GetValue()) {
			fields[fieldname] = float64(m.GetGauge().GetValue())
		}
	} else if m.Counter != nil {
		if !math.IsNaN(m.GetCounter().GetValue()) {
			fields[fieldname] = float64(m.GetCounter().GetValue())
		}
	} else if m.Untyped != nil {
		if !math.IsNaN(m.GetUntyped().GetValue()) {
			fields[fieldname] = float64(m.GetUntyped().GetValue())
		}
	}
	return fields
}
