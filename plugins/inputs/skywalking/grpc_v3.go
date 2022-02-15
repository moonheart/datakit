package skywalking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	itrace "gitlab.jiagouyun.com/cloudcare-tools/datakit/io/trace"
	skyimpl "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs/skywalking/v3/compile"
	"google.golang.org/grpc"
)

func runServerV3(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("start skywalking V3 grpc server %s failed: %v", addr, err)

		return
	}
	log.Infof("skywalking v3 listening on: %s", addr)

	srv := grpc.NewServer()
	skyimpl.RegisterTraceSegmentReportServiceServer(srv, &TraceReportServerV3{})
	skyimpl.RegisterEventServiceServer(srv, &EventServerV3{})
	skyimpl.RegisterJVMMetricReportServiceServer(srv, &JVMMetricReportServerV3{})
	skyimpl.RegisterManagementServiceServer(srv, &ManagementServerV3{})
	skyimpl.RegisterConfigurationDiscoveryServiceServer(srv, &DiscoveryServerV3{})
	if err = srv.Serve(listener); err != nil {
		log.Error(err)
	}
	log.Info("skywalking v3 exits")
}

type TraceReportServerV3 struct {
	skyimpl.UnimplementedTraceSegmentReportServiceServer
}

func (s *TraceReportServerV3) Collect(tsc skyimpl.TraceSegmentReportService_CollectServer) (err error) {
	for {
		segobj, err := tsc.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return tsc.SendAndClose(&skyimpl.Commands{})
			}
			log.Error(err.Error())

			return err
		}

		log.Debug("v3 segment received")

		dktrace, err := segobjToDkTrace(segobj)
		if err != nil {
			log.Error(err.Error())

			return err
		}

		if len(dktrace) == 0 {
			log.Warn("empty datakit trace")
		} else {
			afterGather.Run(inputName, dktrace, false)
		}
	}
}

func (*TraceReportServerV3) CollectInSync(
	ctx context.Context,
	seg *skyimpl.SegmentCollection) (*skyimpl.Commands, error) {
	log.Debugf("reveived collect insync: %s", seg.String())

	return &skyimpl.Commands{}, nil
}

func segobjToDkTrace(segment *skyimpl.SegmentObject) (itrace.DatakitTrace, error) {
	var dktrace itrace.DatakitTrace
	for _, span := range segment.Spans {
		if span == nil {
			continue
		}

		dkspan := &itrace.DatakitSpan{
			TraceID:   segment.TraceId,
			SpanID:    fmt.Sprintf("%s%d", segment.TraceSegmentId, span.SpanId),
			ParentID:  "0",
			EndPoint:  span.Peer,
			Operation: span.OperationName,
			Service:   segment.Service,
			Source:    inputName,
			Start:     span.StartTime * int64(time.Millisecond),
			Duration:  (span.EndTime - span.StartTime) * int64(time.Millisecond),
		}

		if span.SpanType == skyimpl.SpanType_Entry {
			if len(span.Refs) > 0 {
				dkspan.ParentID = fmt.Sprintf("%s%d", span.Refs[0].ParentTraceSegmentId, span.Refs[0].ParentSpanId)
			}
		} else {
			dkspan.ParentID = fmt.Sprintf("%s%d", segment.TraceSegmentId, span.ParentSpanId)
		}

		dkspan.Status = itrace.STATUS_OK
		if span.IsError {
			dkspan.Status = itrace.STATUS_ERR
		}

		switch span.SpanType {
		case skyimpl.SpanType_Entry:
			dkspan.SpanType = itrace.SPAN_TYPE_ENTRY
		case skyimpl.SpanType_Local:
			dkspan.SpanType = itrace.SPAN_TYPE_LOCAL
		case skyimpl.SpanType_Exit:
			dkspan.SpanType = itrace.SPAN_TYPE_EXIT
		}

		sourceTags := make(map[string]string)
		for _, tag := range span.Tags {
			sourceTags[tag.Key] = tag.Value
		}
		dkspan.Tags = itrace.MergeInToCustomerTags(customerKeys, tags, sourceTags)

		if defSampler != nil {
			dkspan.Priority = defSampler.Priority
			dkspan.SamplingRateGlobal = defSampler.SamplingRateGlobal
		}

		if buf, err := json.Marshal(span); err != nil {
			log.Warn(err.Error())
		} else {
			dkspan.Content = string(buf)
		}

		dktrace = append(dktrace, dkspan)
	}

	return dktrace, nil
}

type EventServerV3 struct {
	skyimpl.UnimplementedEventServiceServer
}

func (*EventServerV3) Collect(esrv skyimpl.EventService_CollectServer) error {
	for {
		event, err := esrv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return esrv.SendAndClose(&skyimpl.Commands{})
			}
			log.Debug(err.Error())

			return err
		}

		log.Debugf("reveived service event: %s", event.String())
	}
}

type ManagementServerV3 struct {
	skyimpl.UnimplementedManagementServiceServer
}

func (*ManagementServerV3) ReportInstanceProperties(ctx context.Context,
	mng *skyimpl.InstanceProperties) (*skyimpl.Commands, error) {
	var kvpStr string
	for _, kvp := range mng.Properties {
		kvpStr += fmt.Sprintf("[%v:%v]", kvp.Key, kvp.Value)
	}
	log.Debugf("ReportInstanceProperties service:%v instance:%v properties:%v", mng.Service, mng.ServiceInstance, kvpStr)

	return &skyimpl.Commands{}, nil
}

func (*ManagementServerV3) KeepAlive(
	ctx context.Context,
	ping *skyimpl.InstancePingPkg) (*skyimpl.Commands, error) {
	log.Debugf("KeepAlive service:%v instance:%v", ping.Service, ping.ServiceInstance)

	return &skyimpl.Commands{}, nil
}

type JVMMetricReportServerV3 struct {
	skyimpl.UnimplementedJVMMetricReportServiceServer
}

func (*JVMMetricReportServerV3) Collect(ctx context.Context,
	jvm *skyimpl.JVMMetricCollection) (*skyimpl.Commands, error) {
	log.Debugf("JVMMetricReportService service:%v instance:%v", jvm.Service, jvm.ServiceInstance)

	return &skyimpl.Commands{}, nil
}

type DiscoveryServerV3 struct {
	skyimpl.UnimplementedConfigurationDiscoveryServiceServer
}

func (*DiscoveryServerV3) FetchConfigurations(ctx context.Context,
	cfgReq *skyimpl.ConfigurationSyncRequest) (*skyimpl.Commands, error) {
	log.Debugf("DiscoveryServerV3 service: %s", cfgReq.String())

	return &skyimpl.Commands{}, nil
}
