//go:build (linux && amd64 && ebpf) || (linux && arm64 && ebpf)
// +build linux,amd64,ebpf linux,arm64,ebpf

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/DataDog/ebpf/manager"
	"github.com/jessevdk/go-flags"
	"github.com/shirou/gopsutil/process"
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	dkbash "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/bashhistory"
	dkdns "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/dnsflow"
	dkfeed "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/feed"
	dkhttpflow "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/httpflow"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/k8sinfo"
	dknetflow "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/netflow"

	dkoffset "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/externals/ebpf/offset"
)

const (
	minInterval = time.Second * 10
	maxInterval = time.Minute * 30
)

var (
	enableEbpfBash = false
	enableEbpfNet  = false
)

var pidFile = filepath.Join(datakit.InstallDir, "externals", "datakit-ebpf.pid")

type Option struct {
	DataKitAPIServer string `long:"datakit-apiserver" description:"DataKit API server" default:"0.0.0.0:9529"`

	HostName string `long:"hostname" description:"host name"`

	Interval string `long:"interval" description:"gather interval" default:"60s"`

	PIDFile string `long:"pidfile" description:"pid file"`

	Log      string `long:"log" description:"log path"`
	LogLevel string `long:"log-level" description:"log file" default:"info"`

	Tags    string `long:"tags" description:"additional tags in 'a=b,c=d,...' format"`
	Enabled string `long:"enabled" description:"enabled plugins list in 'a,b,...' format"`

	Service string `long:"service" description:"service" default:"ebpf"`
}

//  Envs:
// 		K8sURL               string
// 		K8sBearerTokenPath   string
// 		K8sBearerTokenString string

const (
	inputName     = "ebpf"
	inputNameNet  = "ebpf-net"
	inputNameBash = "ebpf-bash"
)

var (
	interval = time.Second * 30
	l        = logger.DefaultSLogger(inputName)
)

var signaIterrrupt = make(chan os.Signal, 1)

func main() {
	signal.Notify(signaIterrrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	opt, gTags, err := parseFlags()
	if err != nil {
		feedLastErrorLoop(err, signaIterrrupt)
		return
	}
	if opt.PIDFile != "" {
		pidFile = opt.PIDFile
	}

	if err := savePid(); err != nil {
		l.Fatal(err)
	}

	dumpStderr2File()

	dkfeed.DataKitAPIServer = opt.DataKitAPIServer

	ebpfNetPostURL := fmt.Sprintf("http://%s%s?input="+inputNameNet, dkfeed.DataKitAPIServer, datakit.Network)
	ebpfBashPostURL := fmt.Sprintf("http://%s%s?input="+inputNameBash, dkfeed.DataKitAPIServer, datakit.Logging)

	logOpt := logger.Option{
		Path:  opt.Log,
		Level: opt.LogLevel,
		Flags: logger.OPT_DEFAULT,
	}

	if err := logger.InitRoot(&logOpt); err != nil {
		l.Errorf("set root log faile: %s", err.Error())
	}

	l = logger.SLogger(inputName)

	dknetflow.SetLogger(l)
	dkdns.SetLogger(l)
	dkoffset.SetLogger(l)
	dkbash.SetLogger(l)
	dkhttpflow.SetLogger(l)

	// duration 介于 10s ～ 30min，若非，默认设为 30s.
	if tmp, err := time.ParseDuration(opt.Interval); err == nil {
		if tmp < maxInterval || tmp > maxInterval {
			tmp = time.Second * 30
		}
		interval = tmp
		l.Debug("interval: ", opt.Interval)
	} else {
		interval = time.Second * 60
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l.Info("datakit-ebpf starting ...")

	// ebpf-net
	if enableEbpfNet {
		offset, err := getOffset()
		if err != nil {
			feedLastErrorLoop(err, signaIterrrupt)
			return
		}
		l.Debugf("%+v", offset)

		k8sinfo, err := newK8sInfoFromENV()
		if err != nil {
			l.Warn(err)
		} else {
			go k8sinfo.AutoUpdate(ctx)
			dknetflow.SetK8sNetInfo(k8sinfo)
			dkhttpflow.SetK8sNetInfo(k8sinfo)
		}
		constEditor := dkoffset.NewConstEditor(offset)

		netflowTracer := dknetflow.NewNetFlowTracer()
		ebpfNetManger, err := dknetflow.NewNetFlowManger(constEditor, netflowTracer.ClosedEventHandler)
		if err != nil {
			feedLastErrorLoop(err, signaIterrrupt)
			return
		}
		// Start the manager
		if err := ebpfNetManger.Start(); err != nil {
			feedLastErrorLoop(err, signaIterrrupt)
			return
		} else {
			l.Info(" >>> datakit ebpf-net tracer(ebpf) starting ...")
		}
		defer ebpfNetManger.Stop(manager.CleanAll) //nolint:errcheck

		// used for dns reverse
		dnsRecord := dkdns.NewDNSRecord()
		dknetflow.SetDNSRecord(dnsRecord)

		// run dnsflow
		if tp, err := dkdns.NewTPacketDNS(); err != nil {
			l.Error(err)
		} else {
			dnsTracer := dkdns.NewDNSFlowTracer()
			go dnsTracer.Run(ctx, tp, gTags, dnsRecord, ebpfNetPostURL)
		}

		// run netflow
		err = netflowTracer.Run(ctx, ebpfNetManger, ebpfNetPostURL, gTags, interval)
		if err != nil {
			feedLastErrorLoop(err, signaIterrrupt)
			return
		}

		tracer := dkhttpflow.NewHTTPFlowTracer(gTags, ebpfNetPostURL)
		if err := tracer.Run(ctx); err != nil {
			l.Error(err)
		}
	}

	// ebpf-bash
	if enableEbpfBash {
		l.Info(" >>> datakit ebpf-bash tracer(ebpf) starting ...")
		bashTracer := dkbash.NewBashTracer()
		err := bashTracer.Run(ctx, gTags, ebpfBashPostURL, interval)
		if err != nil {
			feedLastErrorLoop(err, signaIterrrupt)
			return
		}
	}

	if enableEbpfBash || enableEbpfNet {
		<-signaIterrrupt
	}

	l.Info("datakit-ebpf exit")
	quit()
}

func getOffset() (*dkoffset.OffsetGuessC, error) {
	bpfManger, err := dkoffset.NewGuessManger()
	if err != nil {
		return nil, err
	}
	// Start the manager
	if err := bpfManger.Start(); err != nil {
		return nil, err
	}
	loopCount := 5
	defer bpfManger.Stop(manager.CleanAll) //nolint:errcheck
	for i := 0; i < loopCount; i++ {
		mapG, err := dkoffset.BpfMapGuessInit(bpfManger)
		if err != nil {
			return nil, err
		}
		status, err := dkoffset.GuessOffset(mapG, nil)
		if err != nil {
			if i == loopCount-1 {
				return nil, err
			}
			l.Error(err)
			continue
		} else {
			return status, nil
		}
	}
	return nil, err
}

func feedLastErrorLoop(err error, ch chan os.Signal) {
	l.Error(err)

	extLastErr := dkfeed.ExternalLastErr{
		Input:      inputName,
		ErrContent: err.Error(),
	}
	if err := dkfeed.FeedLastError(extLastErr); err != nil {
		l.Error(err)
	}

	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			if err := dkfeed.FeedLastError(extLastErr); err != nil {
				l.Error(err)
			}
		case <-ch:
			return
		}
	}
}

// init opt, dkutil.DataKitAPIServer, datakitPostURL.
func parseFlags() (*Option, map[string]string, error) {
	opt := Option{}
	gTags := map[string]string{}
	if _, err := flags.Parse(&opt); err != nil {
		return nil, nil, err
	}

	optEnabled := strings.Split(opt.Enabled, ",")
	for _, item := range optEnabled {
		switch item {
		case inputNameNet:
			enableEbpfNet = true
		case inputNameBash:
			enableEbpfBash = true
		}
	}

	optTags := strings.Split(opt.Tags, ";")
	for _, item := range optTags {
		tagArr := strings.Split(item, "=")

		if len(tagArr) == 2 {
			tagKey := strings.Trim(tagArr[0], " ")
			tagVal := strings.Trim(tagArr[1], " ")
			if tagKey != "" {
				gTags[tagKey] = tagVal
			}
		}
	}

	if gTags["host"] == "" && opt.HostName == "" {
		var err error
		if gTags["host"], err = os.Hostname(); err != nil {
			l.Error(err)
			gTags["host"] = "no-value"
		}
	} else if opt.HostName != "" {
		gTags["host"] = opt.HostName
	}

	gTags["service"] = opt.Service

	if opt.Log == "" {
		opt.Log = filepath.Join(datakit.InstallDir, "externals", "datakit-ebpf.log")
	}

	return &opt, gTags, nil
}

func newK8sInfoFromENV() (*k8sinfo.K8sNetInfo, error) {
	var k8sURL string
	k8sBearerTokenPath := ""
	k8sBearerTokenStr := ""

	if v, ok := os.LookupEnv("K8S_URL"); ok && v != "" {
		k8sURL = v
	} else {
		k8sURL = "https://kubernetes.default:443"
	}

	// net.LookupHost()

	if v, ok := os.LookupEnv("K8S_BEARER_TOKEN_STRING"); ok && v != "" {
		k8sBearerTokenStr = v
	}

	if v, ok := os.LookupEnv("K8S_BEARER_TOKEN_PATH"); ok && v != "" {
		k8sBearerTokenPath = v
	}

	if k8sBearerTokenPath == "" && k8sBearerTokenStr == "" {
		//nolint:gosec
		k8sBearerTokenPath = "/run/secrets/kubernetes.io/serviceaccount/token"
	}

	var cli *k8sinfo.K8sClient
	var err error
	if k8sBearerTokenPath != "" {
		cli, err = k8sinfo.NewK8sClientFromBearerToken(k8sURL,
			k8sBearerTokenPath)
		if err != nil {
			return nil, err
		}
	} else {
		cli, err = k8sinfo.NewK8sClientFromBearerTokenString(k8sURL,
			k8sBearerTokenStr)
		if err != nil {
			return nil, err
		}
	}
	if cli == nil {
		return nil, fmt.Errorf("new k8s client")
	}
	kinfo, err := k8sinfo.NewK8sNetInfo(cli)
	if err != nil {
		return nil, err
	}

	return kinfo, err
}

func quit() {
	_ = os.Remove(pidFile)
}

func savePid() error {
	if isRuning() {
		return fmt.Errorf("ebpf still running, PID: %s", pidFile)
	}

	pid := os.Getpid()
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0o644) //nolint:gosec
}

func isRuning() bool {
	var oidPid int64
	var name string
	var p *process.Process

	cont, err := ioutil.ReadFile(filepath.Clean(pidFile))
	// pid文件不存在
	if err != nil {
		return false
	}

	oidPid, err = strconv.ParseInt(string(cont), 10, 32)
	if err != nil {
		return false
	}

	p, _ = process.NewProcess(int32(oidPid))
	name, _ = p.Name()

	return name == "datakit-ebpf"
}

const (
	FileModeRW = 0o644
	DirModeRW  = 0o755
)

func dumpStderr2File() {
	dirpath := filepath.Join(datakit.InstallDir, "externals")
	filepath := filepath.Join(dirpath, "datakit-ebpf.stderr")
	if err := os.MkdirAll(dirpath, DirModeRW); err != nil {
		l.Warn(err)
		return
	}
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, FileModeRW) //nolint:gosec
	if err != nil {
		l.Error(err)
		return
	}
	_, err = f.WriteString(fmt.Sprintf("\n========= %s =========\n", time.Now().UTC()))
	if err != nil {
		l.Error(err)
	}
	_ = syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0) // for arm64 arch
	if err = f.Close(); err != nil {
		l.Error(err)
	}
}
