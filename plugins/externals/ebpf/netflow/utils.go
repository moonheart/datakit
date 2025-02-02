//go:build (linux && amd64 && ebpf) || (linux && arm64 && ebpf)
// +build linux,amd64,ebpf linux,arm64,ebpf

package netflow

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/DataDog/ebpf/manager"
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	dkebpf "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/c"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/dnsflow"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/k8sinfo"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	"golang.org/x/sys/unix"
)

var l = logger.DefaultSLogger("ebpf")

var dnsRecord *dnsflow.DNSAnswerRecord

var k8sNetInfo *k8sinfo.K8sNetInfo

func SetDNSRecord(r *dnsflow.DNSAnswerRecord) {
	dnsRecord = r
}

func SetLogger(nl *logger.Logger) {
	l = nl
}

func SetK8sNetInfo(n *k8sinfo.K8sNetInfo) {
	k8sNetInfo = n
}

var SrcIPPortRecorder = func() *srcIPPortRecorder {
	ptr := &srcIPPortRecorder{
		Record: map[[4]uint32]IPPortRecord{},
	}
	go ptr.AutoClean()
	return ptr
}()

type IPPortRecord struct {
	IP [4]uint32
	TS time.Time
}

// 辅助 httpflow 判断 server ip
type srcIPPortRecorder struct {
	sync.RWMutex
	Record map[[4]uint32]IPPortRecord
}

func (record *srcIPPortRecorder) InsertAndUpdate(ip [4]uint32) {
	record.Lock()
	defer record.Unlock()
	record.Record[ip] = IPPortRecord{
		IP: ip,
		TS: time.Now(),
	}
}

func (record *srcIPPortRecorder) Query(ip [4]uint32) (*IPPortRecord, error) {
	record.RLock()
	defer record.RUnlock()
	if v, ok := record.Record[ip]; ok {
		return &v, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

const (
	cleanTickerIPPortDur = time.Minute * 3
	cleanIPPortDur       = time.Minute * 5
)

func (record *srcIPPortRecorder) CleanOutdateData() {
	record.Lock()
	defer record.Unlock()
	ts := time.Now()
	needDelete := [][4]uint32{}
	for k, v := range record.Record {
		if ts.Sub(v.TS) > cleanIPPortDur {
			needDelete = append(needDelete, k)
		}
	}
	for _, v := range needDelete {
		delete(record.Record, v)
	}
}

func (record *srcIPPortRecorder) AutoClean() {
	ticker := time.NewTicker(cleanTickerIPPortDur)
	for {
		<-ticker.C
		record.CleanOutdateData()
	}
}

func NewNetFlowManger(constEditor []manager.ConstantEditor, closedEventHandler func(cpu int, data []byte,
	perfmap *manager.PerfMap, manager *manager.Manager)) (*manager.Manager, error) {
	// 部分 kretprobe 类型程序需设置 maxactive， https://www.kernel.org/doc/Documentation/kprobes.txt.
	m := &manager.Manager{
		Probes: []*manager.Probe{
			{
				Section: "kprobe/sockfd_lookup_light", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/sockfd_lookup_light", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/do_sendfile", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/do_sendfile", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/tcp_set_state", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/inet_csk_accept", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/inet_csk_listen_stop", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/tcp_close", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/tcp_retransmit_skb", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/tcp_sendmsg", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/tcp_cleanup_rbuf", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/ip_make_skb", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/udp_recvmsg", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/udp_recvmsg", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/inet_bind", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/inet_bind", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/inet6_bind", KProbeMaxActive: 128,
			}, {
				Section: "kretprobe/inet6_bind", KProbeMaxActive: 128,
			}, {
				Section: "kprobe/udp_destroy_sock", KProbeMaxActive: 128,
			},
		},
		PerfMaps: []*manager.PerfMap{
			{
				Map: manager.Map{
					Name: "bpfmap_closed_event",
				},
				PerfMapOptions: manager.PerfMapOptions{
					// sizeof(connection_closed_info) > 112 Byte, pagesize ~= 4k,
					// if cpus = 8, 5 conn/per connection_closed_info
					PerfRingBufferSize: 32 * os.Getpagesize(),
					DataHandler:        closedEventHandler,
				},
			},
		},
	}
	mOpts := manager.Options{
		RLimit: &unix.Rlimit{
			Cur: math.MaxUint64,
			Max: math.MaxUint64,
		},
		ConstantEditors: constEditor,
	}
	if buf, err := dkebpf.Asset("netflow.o"); err != nil {
		return nil, err
	} else if err := m.InitWithOptions((bytes.NewReader(buf)), mOpts); err != nil {
		return nil, err
	}

	return m, nil
}

func ConvertConn2Measurement(connR *ConnResult, name string) []inputs.Measurement {
	collectCache := []inputs.Measurement{}

	for k, v := range connR.result {
		if ConnNotNeedToFilter(k, v) {
			m := ConvConn2M(k, v, name, connR.tags, connR.ts)
			collectCache = append(collectCache, m)
		}
	}
	return collectCache
}

func ConvConn2M(k ConnectionInfo, v ConnFullStats, name string,
	tags map[string]string, ts time.Time) inputs.Measurement {
	m := measurement{
		name:   name,
		tags:   map[string]string{},
		fields: map[string]interface{}{},
		ts:     ts,
	}
	for k, v := range tags {
		m.tags[k] = v
	}

	m.tags["status"] = "info"
	m.tags["pid"] = fmt.Sprint(k.Pid)

	isV6 := !ConnAddrIsIPv4(k.Meta)
	if k.Saddr[0] == 0 && k.Saddr[1] == 0 && k.Daddr[0] == 0 && k.Daddr[1] == 0 {
		if k.Saddr[2] == 0xffff0000 && k.Daddr[2] == 0xffff0000 {
			isV6 = false
		} else if k.Saddr[2] == 0 && k.Daddr[2] == 0 && k.Saddr[3] > 1 && k.Daddr[3] > 1 {
			isV6 = false
		}
	}

	if !isV6 {
		m.tags["src_ip_type"] = ConnIPv4Type(k.Saddr[3])
		m.tags["dst_ip_type"] = ConnIPv4Type(k.Daddr[3])
		m.tags["family"] = "IPv4"
	} else {
		m.tags["src_ip_type"] = ConnIPv6Type(k.Saddr)

		m.tags["dst_ip_type"] = ConnIPv6Type(k.Daddr)
		m.tags["family"] = "IPv6"
	}

	m.tags["src_ip"] = U32BEToIP(k.Saddr, isV6).String()

	dstIP := U32BEToIP(k.Daddr, isV6)
	m.tags["dst_ip"] = dstIP.String()

	if dnsRecord != nil {
		m.tags["dst_domain"] = dnsRecord.LookupAddr(dstIP)
	}

	if k.Sport == math.MaxUint32 {
		m.tags["src_port"] = "*"
	} else {
		m.tags["src_port"] = fmt.Sprintf("%d", k.Sport)
	}

	m.tags["dst_port"] = fmt.Sprintf("%d", k.Dport)

	m.fields["bytes_read"] = int64(v.Stats.RecvBytes)
	m.fields["bytes_written"] = int64(v.Stats.SentBytes)

	if ConnProtocolIsTCP(k.Meta) {
		m.tags["transport"] = "tcp"
		m.fields["retransmits"] = int64(v.TCPStats.Retransmits)
		m.fields["rtt"] = int64(v.TCPStats.Rtt)
		m.fields["rtt_var"] = int64(v.TCPStats.RttVar)
		m.fields["tcp_closed"] = v.TotalClosed
		m.fields["tcp_established"] = v.TotalEstablished
	} else {
		m.tags["transport"] = "udp"
	}
	m.tags["direction"] = connDirection2Str(v.Stats.Direction)

	if k8sNetInfo != nil {
		srcK8sFlag := false
		dstK8sFlag := false
		_, srcPoName, srcSvcName, ns, svcP, err := k8sNetInfo.QueryPodInfo(m.tags["src_ip"], k.Sport, m.tags["transport"])
		if err == nil {
			srcK8sFlag = true
			m.tags["src_k8s_namespace"] = ns
			m.tags["src_k8s_pod_name"] = srcPoName
			m.tags["src_k8s_service_name"] = srcSvcName
			if svcP == k.Sport {
				m.tags["direction"] = "incoming"
			}
		}

		_, dstPodName, dstSvcName, ns, svcP, err := k8sNetInfo.QueryPodInfo(m.tags["dst_ip"], k.Dport, m.tags["transport"])
		if err == nil {
			dstK8sFlag = true
			m.tags["dst_k8s_namespace"] = ns
			m.tags["dst_k8s_pod_name"] = dstPodName
			m.tags["dst_k8s_service_name"] = dstSvcName
			if svcP == k.Dport {
				m.tags["direction"] = "outgoing"
			}

		} else {
			dstSvcName, ns, err := k8sNetInfo.QuerySvcInfo(m.tags["dst_ip"])
			if err == nil {
				dstK8sFlag = true
				m.tags["dst_k8s_namespace"] = ns
				m.tags["dst_k8s_pod_name"] = "N/A"
				m.tags["dst_k8s_service_name"] = dstSvcName
				m.tags["direction"] = "outgoing"
			}
		}

		if srcK8sFlag || dstK8sFlag {
			m.tags["sub_source"] = "K8s"
			if !srcK8sFlag {
				m.tags["src_k8s_namespace"] = "N/A"
				m.tags["src_k8s_pod_name"] = "N/A"
				m.tags["src_k8s_service_name"] = "N/A"
			}
			if !dstK8sFlag {
				m.tags["dst_k8s_namespace"] = "N/A"
				m.tags["dst_k8s_pod_name"] = "N/A"
				m.tags["dst_k8s_service_name"] = "N/A"
			}
		}
	}

	if ConnProtocolIsTCP(k.Meta) {
		l.Debug(fmt.Sprintf("pid %s: %s:%s->%s(%s):%s r/w: %d/%d e/c: %d/%d "+
			"re: %d rtt/rttvar: %.2fms/%.2fms (%s, %s)",
			m.tags["pid"], m.tags["src_ip"], m.tags["src_port"], m.tags["dst_ip"], m.tags["dst_domain"],
			m.tags["dst_port"], m.fields["bytes_read"], m.fields["bytes_written"],
			m.fields["tcp_established"], m.fields["tcp_closed"], m.fields["retransmits"],
			float64(v.TCPStats.Rtt)/1000., float64(v.TCPStats.RttVar)/1000, m.tags["transport"], m.tags["direction"]))
	} else {
		l.Debug(fmt.Sprintf("pid %s: %s:%s->%s(%s):%s r/w: %d/%d (%s, %s)",
			m.tags["pid"], m.tags["src_ip"], m.tags["src_port"], m.tags["dst_ip"],
			m.tags["dst_domain"], m.tags["dst_port"], m.fields["bytes_read"],
			m.fields["bytes_written"], m.tags["transport"], m.tags["direction"]))
	}
	return &m
}

func U32BEToIPv4Array(addr uint32) [4]uint8 {
	var ip [4]uint8
	for x := 0; x < 4; x++ {
		ip[x] = uint8(addr & 0xff)
		addr >>= 8
	}
	return ip
}

func SwapU16(v uint16) uint16 {
	return ((v & 0x00ff) << 8) | ((v & 0xff00) >> 8)
}

func U32BEToIPv6Array(addr [4]uint32) [8]uint16 {
	var ip [8]uint16
	for x := 0; x < 4; x++ {
		ip[(x * 2)] = SwapU16(uint16(addr[x] & 0xffff))         // uint32 低16位
		ip[(x*2)+1] = SwapU16(uint16((addr[x] >> 16) & 0xffff)) //	高16位
	}
	return ip
}

func U32BEToIP(addr [4]uint32, isIPv6 bool) net.IP {
	ip := net.IP{}
	if !isIPv6 {
		v4 := U32BEToIPv4Array(addr[3])
		for _, v := range v4 {
			ip = append(ip, v)
		}
	} else {
		v6 := U32BEToIPv6Array(addr)
		for _, v := range v6 {
			ip = append(ip, byte((v&0xff00)>>8), byte(v&0x00ff)) // SwapU16(v)
		}
	}
	return ip
}

// ConnNotNeedToFilter 规则: 1. 过滤源 IP 和目标 IP 相同的连接;
// 2. 过滤 loopback ip 的连接;
// 3. 过滤一个采集周期内的无数据收发的连接;
// 4. 过滤端口 为 0 或 ip address 为 :: or 0.0.0.0 的连接;
// 需过滤，函数返回 False.
func ConnNotNeedToFilter(conn ConnectionInfo, connStats ConnFullStats) bool {
	if (conn.Saddr[0]|conn.Saddr[1]|conn.Saddr[2]|conn.Saddr[3]) == 0 ||
		(conn.Daddr[0]|conn.Daddr[1]|conn.Daddr[2]|conn.Daddr[3]) == 0 ||
		conn.Sport == 0 || conn.Dport == 0 {
		return false
	}
	if ConnAddrIsIPv4(conn.Meta) { // IPv4
		if (conn.Saddr[3]&0xff) == 127 && (conn.Daddr[3]&0xff) == 127 {
			return false
		}
	} else { // IPv6
		if conn.Saddr[2] == 0xffff0000 && conn.Daddr[2] == 0xffff0000 {
			if (conn.Saddr[3]&0xff) == 127 && (conn.Daddr[3]&0xff) == 127 {
				return false
			}
		} else if (conn.Saddr[0]|conn.Saddr[1]|conn.Saddr[2]) == 0 && conn.Saddr[3] == 1 &&
			(conn.Daddr[0]|conn.Daddr[1]|conn.Daddr[2]) == 0 && conn.Daddr[3] == 1 {
			return false
		}
	}

	// 过滤上一周期的无变化的连接
	if connStats.Stats.RecvBytes == 0 && connStats.Stats.SentBytes == 0 &&
		connStats.TotalClosed == 0 && connStats.TotalEstablished == 0 {
		return false
	}

	return true
}

type ConnTCPWithoutPidStats struct {
	TCPStats         ConnectionTCPStats
	TotalEstablished int64
	TotalClosed      int64
	Pids             map[uint32]bool
}

type ConnTCPWithoutPid struct {
	Conns map[ConnectionInfo]ConnTCPWithoutPidStats
}

func (cwp *ConnTCPWithoutPid) CleanupConns() {
	needCleanup := []ConnectionInfo{}
	for k, v := range cwp.Conns {
		if v.TotalClosed == v.TotalEstablished {
			needCleanup = append(needCleanup, k)
		}
	}
	for _, k := range needCleanup {
		delete(cwp.Conns, k)
	}
}

func (cwp *ConnTCPWithoutPid) Update(conn ConnectionInfo, fullStats ConnFullStats) ConnFullStats {
	pid := conn.Pid
	conn.Pid = 0
	result := fullStats
	if v, ok := cwp.Conns[conn]; ok {
		if _, ok := v.Pids[pid]; !ok {
			v.Pids[pid] = false // 首次记录
		}
		if len(v.Pids) > 1 && result.TotalEstablished > 0 { // 复数个 pid
			if !v.Pids[pid] { // 对首次记录的 pid 对应的连接次数进行处理
				if result.TotalEstablished > 0 {
					result.TotalEstablished -= 1
				}
			}
		}
		v.Pids[pid] = true
		v.TotalClosed += result.TotalClosed
		v.TotalEstablished += result.TotalEstablished
		cwp.Conns[conn] = v
	} else {
		cwp.Conns[conn] = ConnTCPWithoutPidStats{
			Pids:             map[uint32]bool{pid: true},
			TotalEstablished: result.TotalEstablished,
			TotalClosed:      result.TotalClosed,
		}
	}
	return result
}

func newConnTCPWithoutPid() *ConnTCPWithoutPid {
	return &ConnTCPWithoutPid{
		Conns: make(map[ConnectionInfo]ConnTCPWithoutPidStats),
	}
}

var connTCPWithoutPid = newConnTCPWithoutPid()

// MergeConns 聚合 src port 为临时端口(32768 ~ 60999)的连接,
// 被聚合的端口号被设置为
// cat /proc/sys/net/ipv4/ip_local_port_range.
func MergeConns(preResult *ConnResult) {
	resultTmpConn := map[ConnectionInfo]ConnFullStats{}
	if len(preResult.result) < 1 {
		return
	}

	connInfoList := ConnInfoList{}

	for k, v := range preResult.result {
		connInfoList = append(connInfoList, k)
		r := connTCPWithoutPid.Update(k, v)
		preResult.result[k] = r
	}
	connTCPWithoutPid.CleanupConns()
	sort.Sort(connInfoList)
	connRecord := map[ConnectionInfo]bool{}
	lastIndex := -1
	for k := 0; k < len(connInfoList); k++ {
		if !isEphemeralPort(connInfoList[k].Sport) {
			continue
		}

		switch {
		case lastIndex < 0:
			lastIndex = k
			resultTmpConn[connInfoList[k]] = preResult.result[connInfoList[k]]
			delete(preResult.result, connInfoList[k])
		case ConnCmpNoSPort(connInfoList[lastIndex], connInfoList[k]):
			connRecord[connInfoList[lastIndex]] = true
			resultTmpConn[connInfoList[lastIndex]] = StatsTCPOp("+", resultTmpConn[connInfoList[lastIndex]],
				preResult.result[connInfoList[k]].Stats, preResult.result[connInfoList[k]].TCPStats)

			connfull := resultTmpConn[connInfoList[lastIndex]]
			connfull.TotalEstablished += preResult.result[connInfoList[k]].TotalEstablished
			connfull.TotalClosed += preResult.result[connInfoList[k]].TotalClosed
			resultTmpConn[connInfoList[lastIndex]] = connfull

			delete(preResult.result, connInfoList[k])
		default:
			k--
			lastIndex = -1
		}
	}

	for k, v := range resultTmpConn {
		if _, ok := connRecord[k]; ok {
			k.Sport = math.MaxUint32
		}
		preResult.result[k] = v
	}
}

func ConnCmpNoSPort(expected, actual ConnectionInfo) bool {
	expected.Sport = 0
	actual.Sport = 0
	return expected == actual
}

func ConnCmpNoPid(expected, actual ConnectionInfo) bool {
	expected.Pid = 0
	actual.Pid = 0
	return expected == actual
}

type ConnInfoList []ConnectionInfo

func (l ConnInfoList) Len() int {
	return len(l)
}

func (l ConnInfoList) Less(i, j int) bool {
	metaI := l[i].Meta
	metaJ := l[j].Meta

	// family (ipv4)
	if metaI&ConnL3Mask != metaJ&ConnL3Mask {
		return metaI&ConnL3Mask == ConnL3IPv4
	}

	// transport (tcp)
	if metaI&ConnL4Mask != metaJ&ConnL4Mask {
		return metaI&ConnL4Mask == ConnL4TCP
	}

	// src ip, dst ip
	if metaI&ConnL3Mask == ConnL3IPv4 { // ipv4
		if l[i].Saddr[3] != l[j].Saddr[3] {
			return l[i].Saddr[3] < l[j].Saddr[3]
		}
		if l[i].Daddr[3] != l[j].Daddr[3] { // dst ip
			return l[i].Daddr[3] < l[j].Daddr[3]
		}
	} else { // ipv6
		if l[i].Saddr != l[j].Saddr {
			for x := 0; x < 4; x++ {
				if l[i].Saddr[x] > l[j].Saddr[x] {
					return false
				}
			}
			return true
		}
		if l[i].Daddr != l[j].Daddr {
			for x := 0; x < 4; x++ {
				if l[i].Daddr[x] > l[j].Daddr[x] {
					return false
				}
			}
			return true
		}
	}

	// dst port
	if l[i].Dport != l[j].Dport {
		return l[i].Dport < l[j].Dport
	}

	// src port
	if l[i].Sport != l[j].Sport {
		return l[i].Sport < l[j].Sport
	}

	// pid
	if l[i].Pid != l[j].Pid {
		return l[i].Pid < l[j].Pid
	}

	// all equal
	return false
}

func (l ConnInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

const (
	EphemeralPortMin = 32768
	EphemeralPortMax = 60999
)

func isEphemeralPort(port uint32) bool {
	return port >= EphemeralPortMin && port <= EphemeralPortMax
}

func IPPortFilterIn(conn *ConnectionInfo) bool {
	if conn.Sport == 0 || conn.Dport == 0 {
		return false
	}

	if ConnAddrIsIPv4(conn.Meta) {
		if (conn.Saddr[3]&0xFF == 0x7F) || (conn.Daddr[3]&0xFF == 0x7F) {
			return false
		}
	} else if (conn.Saddr[0]|conn.Saddr[1]) == 0x00 || (conn.Daddr[0]|conn.Daddr[1]) == 0x00 {
		if (conn.Saddr[2] == 0xffff0000 && conn.Saddr[3]&0xFF == 0x7F) ||
			(conn.Daddr[2] == 0xffff0000 && conn.Daddr[3]&0xFF == 0x7F) {
			return false
		} else if (conn.Saddr[2] == 0x0 && conn.Saddr[3] == 0x01000000) ||
			(conn.Daddr[2] == 0x0 && conn.Daddr[3] == 0x01000000) {
			return false
		}
	}
	return true
}
