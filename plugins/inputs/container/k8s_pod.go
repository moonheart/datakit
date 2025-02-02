// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	v1 "k8s.io/api/core/v1"
)

var (
	_ k8sResourceMetricInterface = (*pod)(nil)
	_ k8sResourceObjectInterface = (*pod)(nil)
)

type pod struct {
	client    k8sClientX
	extraTags map[string]string
	items     []v1.Pod
}

func newPod(client k8sClientX, extraTags map[string]string) *pod {
	return &pod{
		client:    client,
		extraTags: extraTags,
	}
}

func (p *pod) name() string {
	return "pod"
}

func (p *pod) pullItems() error {
	if len(p.items) != 0 {
		return nil
	}

	list, err := p.client.getPods().List(context.Background(), metaV1ListOption)
	if err != nil {
		return fmt.Errorf("failed to get pods resource: %w", err)
	}

	p.items = list.Items
	return nil
}

func (p *pod) metric() (inputsMeas, error) {
	if err := p.pullItems(); err != nil {
		return nil, err
	}
	var res inputsMeas

	for _, item := range p.items {
		met := &podMetric{
			tags: map[string]string{
				"pod":       item.Name,
				"namespace": defaultNamespace(item.Namespace),
				// "condition":  "",
				// "deployment": "",
				// "daemonset":  "",
			},
			fields: map[string]interface{}{
				"ready": 0,
				// "scheduled": 0,
				// "volumes_persistentvolumeclaims_readonly": 0,
				// "unschedulable": 0,
			},
			time: time.Now(),
		}

		containerReadyCount := 0
		for _, cs := range item.Status.ContainerStatuses {
			if cs.State.Running != nil {
				containerReadyCount++
			}
		}
		met.fields["ready"] = containerReadyCount

		met.tags.append(p.extraTags)
		res = append(res, met)
	}

	count, _ := p.count()
	for ns, c := range count {
		met := &podMetric{
			tags:   map[string]string{"namespace": ns},
			fields: map[string]interface{}{"count": c},
			time:   time.Now(),
		}
		met.tags.append(p.extraTags)
		res = append(res, met)
	}

	return res, nil
}

func (p *pod) count() (map[string]int, error) {
	if err := p.pullItems(); err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, item := range p.items {
		m[defaultNamespace(item.Namespace)]++
	}

	if len(m) == 0 {
		m["default"] = 0
	}

	return m, nil
}

func (p *pod) object() (inputsMeas, error) {
	if err := p.pullItems(); err != nil {
		return nil, err
	}
	var res inputsMeas

	podIDs := make(map[string]interface{})

	for _, item := range p.items {
		obj := &podObject{
			tags: map[string]string{
				"name":         fmt.Sprintf("%v", item.UID),
				"pod_name":     item.Name,
				"node_name":    item.Spec.NodeName,
				"phase":        fmt.Sprintf("%v", item.Status.Phase),
				"qos_class":    fmt.Sprintf("%v", item.Status.QOSClass),
				"state":        fmt.Sprintf("%v", item.Status.Phase), // Depercated
				"status":       fmt.Sprintf("%v", item.Status.Phase),
				"cluster_name": defaultClusterName(item.ClusterName),
				"namespace":    defaultNamespace(item.Namespace),
			},
			fields: map[string]interface{}{
				"age":         int64(time.Since(item.CreationTimestamp.Time).Seconds()),
				"availale":    len(item.Status.ContainerStatuses),
				"create_time": item.CreationTimestamp.Time.Unix(),
			},
			time: time.Now(),
		}

		if n := getHostname(); n != "" {
			obj.tags["host"] = n // 指定 pod 所在的 host
		}

		for _, ref := range item.OwnerReferences {
			if ref.Kind == "ReplicaSet" {
				obj.tags["replica_set"] = ref.Name
				break
			}
		}
		if deployment := getDeployment(item.Labels["app"], item.Namespace); deployment != "" {
			obj.tags["deployment"] = deployment
		}

		for _, containerStatus := range item.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil {
				obj.tags["state"] = containerStatus.State.Waiting.Reason // Depercated
				obj.tags["status"] = containerStatus.State.Waiting.Reason
				break
			}
		}
		obj.tags.append(p.extraTags)

		containerReadyCount := 0
		for _, cs := range item.Status.ContainerStatuses {
			if cs.State.Running != nil {
				containerReadyCount++
			}
		}
		obj.fields["ready"] = containerReadyCount

		restartCount := 0
		for _, containerStatus := range item.Status.InitContainerStatuses {
			restartCount += int(containerStatus.RestartCount)
		}
		for _, containerStatus := range item.Status.ContainerStatuses {
			restartCount += int(containerStatus.RestartCount)
		}
		for _, containerStatus := range item.Status.EphemeralContainerStatuses {
			restartCount += int(containerStatus.RestartCount)
		}
		obj.fields["restart"] = restartCount
		obj.fields["restarts"] = restartCount

		obj.fields.addMapWithJSON("annotations", item.Annotations)
		obj.fields.addLabel(item.Labels)
		obj.fields.mergeToMessage(obj.tags)
		obj.fields.delete("annotations")

		if cli, ok := p.client.(*k8sClient); ok && cli.metricsClient != nil {
			met, err := getPodSrvMetric(cli.metricsClient, item.Namespace, item.Name)
			if err != nil {
				l.Debugf("unable get pod metric %s, namespace %s, name %s", err, defaultNamespace(item.Namespace), item.Name)
			} else {
				obj.fields["cpu_usage"] = met.fields["cpu_usage"]
				obj.fields["memory_usage_bytes"] = met.fields["memory_usage_bytes"]
			}
		}

		res = append(res, obj)

		podIDs[string(item.UID)] = nil

		tempItem := item
		if err := tryRunInput(&tempItem); err != nil {
			l.Warnf("failed to run input(discovery), %s", err)
		}
	}

	for id, inputList := range discoveryInputsMap {
		if _, ok := podIDs[id]; ok {
			continue
		}
		for _, ii := range inputList {
			if ii == nil {
				continue
			}
			if inp, ok := ii.(inputs.InputV2); ok {
				inp.Terminate()
			}
		}
	}

	return res, nil
}

//nolint:deadcode,unused
func getPodLables(k8sClient k8sClientX, podname, podnamespace string) (map[string]string, error) {
	pod, err := queryPodMetaData(k8sClient, podname, podnamespace)
	if err != nil {
		return nil, err
	}
	return pod.labels(), nil
}

func getPodAnnotations(k8sClient k8sClientX, podname, podnamespace string) (map[string]string, error) {
	pod, err := queryPodMetaData(k8sClient, podname, podnamespace)
	if err != nil {
		return nil, err
	}
	return pod.annotations(), nil
}

type podMeta struct{ *v1.Pod }

func queryPodMetaData(k8sClient k8sClientX, podname, podnamespace string) (*podMeta, error) {
	pod, err := k8sClient.getPodsForNamespace(podnamespace).Get(context.Background(), podname, metaV1GetOption)
	if err != nil {
		return nil, err
	}
	return &podMeta{pod}, nil
}

func (item *podMeta) labels() map[string]string { return item.Labels }

func (item *podMeta) annotations() map[string]string { return item.Annotations }

func (item *podMeta) containerImage(name string) string {
	for _, container := range item.Spec.Containers {
		if container.Name == name {
			return container.Image
		}
	}
	return ""
}

func (item *podMeta) replicaSet() string {
	for _, ref := range item.OwnerReferences {
		if ref.Kind == "ReplicaSet" {
			return ref.Name
		}
	}
	return ""
}

type podMetric struct {
	tags   tagsType
	fields fieldsType
	time   time.Time
}

func (p *podMetric) LineProto() (*io.Point, error) {
	return io.NewPoint("kube_pod", p.tags, p.fields, &io.PointOption{Time: p.time, Category: datakit.Metric})
}

//nolint:lll
func (*podMetric) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: "kube_pod",
		Desc: "Kubernetes pod 指标数据",
		Type: "metric",
		Tags: map[string]interface{}{
			"pod":       inputs.NewTagInfo("Name must be unique within a namespace."),
			"namespace": inputs.NewTagInfo("Namespace defines the space within each name must be unique."),
		},
		Fields: map[string]interface{}{
			"count": &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.NCount, Desc: "Number of pods"},
			"ready": &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.NCount, Desc: "Describes whether the pod is ready to serve requests."},
		},
	}
}

type podObject struct {
	tags   tagsType
	fields fieldsType
	time   time.Time
}

func (p *podObject) LineProto() (*io.Point, error) {
	return io.NewPoint("kubelet_pod", p.tags, p.fields, &io.PointOption{Time: p.time, Category: datakit.Object})
}

//nolint:lll
func (*podObject) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: "kubelet_pod",
		Desc: "Kubernetes pod 对象数据",
		Type: "object",
		Tags: map[string]interface{}{
			"name":         inputs.NewTagInfo("UID"),
			"pod_name":     inputs.NewTagInfo("Name must be unique within a namespace."),
			"node_name":    inputs.NewTagInfo("NodeName is a request to schedule this pod onto a specific node."),
			"cluster_name": inputs.NewTagInfo("The name of the cluster which the object belongs to."),
			"namespace":    inputs.NewTagInfo("Namespace defines the space within each name must be unique."),
			"phase":        inputs.NewTagInfo("The phase of a Pod is a simple, high-level summary of where the Pod is in its lifecycle.(Pending/Running/Succeeded/Failed/Unknown)"),
			"state":        inputs.NewTagInfo("Reason the container is not yet running. (Depercated, use status)"),
			"status":       inputs.NewTagInfo("Reason the container is not yet running."),
			"qos_class":    inputs.NewTagInfo("The Quality of Service (QOS) classification assigned to the pod based on resource requirements"),
			"deployment":   inputs.NewTagInfo("The name of the deployment which the object belongs to. (Probably empty)"),
			"replica_set":  inputs.NewTagInfo("The name of the replicaSet which the object belongs to. (Probably empty)"),
		},
		Fields: map[string]interface{}{
			"age":                &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.DurationSecond, Desc: "age (seconds)"},
			"create_time":        &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.TimestampSec, Desc: "CreationTimestamp is a timestamp representing the server time when this object was created.(second)"},
			"restart":            &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.NCount, Desc: "The number of times the container has been restarted. (Depercated, use restarts)"},
			"restarts":           &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.NCount, Desc: "The number of times the container has been restarted."},
			"ready":              &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "Describes whether the pod is ready to serve requests."},
			"available":          &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "Number of containers"},
			"cpu_usage":          &inputs.FieldInfo{DataType: inputs.Float, Unit: inputs.Percent, Desc: "The percentage of cpu used"},
			"memory_usage_bytes": &inputs.FieldInfo{DataType: inputs.Float, Unit: inputs.SizeByte, Desc: "The number of memory used in bytes"},
			"message":            &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "object details"},
		},
	}
}

func getHostname() string {
	if e := os.Getenv("ENV_K8S_NODE_NAME"); e != "" {
		return e
	}
	n, _ := os.Hostname()
	return n
}

//nolint:gochecknoinits
func init() {
	registerK8sResourceMetric(func(c k8sClientX, m map[string]string) k8sResourceMetricInterface { return newPod(c, m) })
	registerK8sResourceObject(func(c k8sClientX, m map[string]string) k8sResourceObjectInterface { return newPod(c, m) })
	registerMeasurement(&podMetric{})
	registerMeasurement(&podObject{})
}
