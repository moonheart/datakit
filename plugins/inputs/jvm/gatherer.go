package jvm

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

const defaultFieldName = "value"

type Gatherer struct {
	metrics  []Metric
	requests []ReadRequest
}

func NewGatherer(metrics []Metric) *Gatherer {
	return &Gatherer{
		metrics:  metrics,
		requests: makeReadRequests(metrics),
	}
}

// Gather adds points to an accumulator from responses returned
// by a Jolokia agent.
func (g *Gatherer) Gather(client *Client, ja *JolokiaAgent) error {
	tags := make(map[string]string)
	if client.config.ProxyConfig != nil {
		tags["jolokia_proxy_url"] = client.URL
	} else {
		tags["jolokia_agent_url"] = client.URL
	}

	requests := makeReadRequests(g.metrics)
	responses, err := client.read(requests)
	if err != nil {
		return err
	}

	g.gatherResponses(responses, mergeTags(tags, ja.Tags), ja)

	return nil
}

// gatherResponses adds points to an accumulator from the ReadResponse objects
// returned by a Jolokia agent.
func (g *Gatherer) gatherResponses(responses []ReadResponse, tags map[string]string, ja *JolokiaAgent) error {
	series := make(map[string][]point)

	for _, metric := range g.metrics {
		points, ok := series[metric.Name]
		if !ok {
			points = make([]point, 0)
		}

		responsePoints, responseErrors := g.generatePoints(metric, responses)
		points = append(points, responsePoints...)
		if len(responseErrors) != 0 {
			return responseErrors[0]
		}

		series[metric.Name] = points
	}

	for measurement, points := range series {
		for _, point := range compactPoints(points) {
			if len(point.Fields) == 0 {
				continue
			}

			field := make(map[string]interface{})
			for k, v := range point.Fields {
				if v == nil {
					continue
				}

				if jn, ok := v.(json.Number); ok {
					field[k] = convertJsonNumber(k, jn, ja.Types)
				} else {
					field[k] = v
				}
			}

			if len(field) == 0 {
				continue
			}

			ja.collectCache = append(ja.collectCache, &JvmMeasurement{
				measurement,
				mergeTags(point.Tags, tags),
				field,
				time.Now(),
			})
		}
	}
	return nil
}

// generatePoints creates points for the supplied metric from the ReadResponse
// objects returned by the Jolokia client.
func (g *Gatherer) generatePoints(metric Metric, responses []ReadResponse) ([]point, []error) {
	points := make([]point, 0)
	errors := make([]error, 0)

	for _, response := range responses {
		switch response.Status {
		case 200:
			break
		case 404:
			continue
		default:
			errors = append(errors, fmt.Errorf("unexpected status in response from target %s (%q): %d",
				response.RequestTarget, response.RequestMbean, response.Status))
			continue
		}

		if !metricMatchesResponse(metric, response) {
			continue
		}

		pb := newPointBuilder(metric, response.RequestAttributes, response.RequestPath)
		for _, point := range pb.Build(metric.Mbean, response.Value) {
			if response.RequestTarget != "" {
				point.Tags["jolokia_agent_url"] = response.RequestTarget
			}

			points = append(points, point)
		}
	}

	return points, errors
}

// mergeTags combines two tag sets into a single tag set.
func mergeTags(metricTags, outerTags map[string]string) map[string]string {
	tags := make(map[string]string)
	for k, v := range outerTags {
		tags[k] = v
	}
	for k, v := range metricTags {
		tags[k] = v
	}

	return tags
}

// metricMatchesResponse returns true when the name, attributes, and path
// of a Metric match the corresponding elements in a ReadResponse object
// returned by a Jolokia agent.
func metricMatchesResponse(metric Metric, response ReadResponse) bool {
	if !metric.MatchObjectName(response.RequestMbean) {
		return false
	}

	if len(metric.Paths) == 0 {
		return len(response.RequestAttributes) == 0
	}

	for _, attribute := range response.RequestAttributes {
		if metric.MatchAttributeAndPath(attribute, response.RequestPath) {
			return true
		}
	}

	return false
}

// compactPoints attempts to remove points by compacting points
// with matching tag sets. When a match is found, the fields from
// one point are moved to another, and the empty point is removed.
func compactPoints(points []point) []point {
	compactedPoints := make([]point, 0)

	for _, sourcePoint := range points {
		keepPoint := true

		for _, compactPoint := range compactedPoints {
			if !tagSetsMatch(sourcePoint.Tags, compactPoint.Tags) {
				continue
			}

			keepPoint = false
			for key, val := range sourcePoint.Fields {
				compactPoint.Fields[key] = val
			}
		}

		if keepPoint {
			compactedPoints = append(compactedPoints, sourcePoint)
		}
	}

	return compactedPoints
}

// tagSetsMatch returns true if two maps are equivalent.
func tagSetsMatch(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for ak, av := range a {
		bv, ok := b[ak]
		if !ok {
			return false
		}
		if av != bv {
			return false
		}
	}

	return true
}

// makeReadRequests creates ReadRequest objects from metrics definitions.
func makeReadRequests(metrics []Metric) []ReadRequest {
	var requests []ReadRequest
	for _, metric := range metrics {
		if len(metric.Paths) == 0 {
			requests = append(requests, ReadRequest{
				Mbean:      metric.Mbean,
				Attributes: []string{},
			})
		} else {
			attributes := make(map[string][]string)

			for _, path := range metric.Paths {
				segments := strings.Split(path, "/")
				attribute := segments[0]

				if _, ok := attributes[attribute]; !ok {
					attributes[attribute] = make([]string, 0)
				}

				if len(segments) > 1 {
					paths := attributes[attribute]
					attributes[attribute] = append(paths, strings.Join(segments[1:], "/"))
				}
			}
			rootAttributes := findRequestAttributesWithoutPaths(attributes)
			if len(rootAttributes) > 0 {
				requests = append(requests, ReadRequest{
					Mbean:      metric.Mbean,
					Attributes: rootAttributes,
				})
			}

			for _, deepAttribute := range findRequestAttributesWithPaths(attributes) {
				for _, path := range attributes[deepAttribute] {
					requests = append(requests, ReadRequest{
						Mbean:      metric.Mbean,
						Attributes: []string{deepAttribute},
						Path:       path,
					})
				}
			}
		}
	}
	return requests
}

func findRequestAttributesWithoutPaths(attributes map[string][]string) []string {
	results := make([]string, 0)
	for attr, paths := range attributes {
		if len(paths) == 0 {
			results = append(results, attr)
		}
	}

	sort.Strings(results)
	return results
}

func findRequestAttributesWithPaths(attributes map[string][]string) []string {
	results := make([]string, 0)
	for attr, paths := range attributes {
		if len(paths) != 0 {
			results = append(results, attr)
		}
	}

	sort.Strings(results)
	return results
}

func convertJsonNumber(fieldName string, jn json.Number, fieldTyp map[string]string) interface{} {
	if fieldTyp != nil {
		return convertSpecifyType(fieldName, jn, fieldTyp)
	}

	if intVal, err := jn.Int64(); err == nil {
		return intVal
	}

	if floatVal, err := jn.Float64(); err == nil {
		return floatVal
	}

	return jn.String()
}

func convertSpecifyType(fieldName string, jn json.Number, fieldTyp map[string]string) interface{} {
	val, _ := jn.Float64()
	if typ, ok := fieldTyp[fieldName]; ok && typ == "int" {
		return int64(val)
	}
	return val
}
