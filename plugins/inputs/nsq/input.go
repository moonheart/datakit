package nsq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/net"
	timex "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/time"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
)

const (
	inputName = "nsq"
	catalog   = "nsq"

	updateEndpointListInterval = time.Second * 30
	minInterval                = time.Second * 3
	defaultInterval            = time.Second * 10

	nsqdStatsPattern = `%s/stats?format=json`
	lookupdPattern   = `%s/nodes`
)

var l = logger.DefaultSLogger(inputName)

type Input struct {
	Lookupd            string            `toml:"lookupd"`
	NSQDs              []string          `toml:"nsqd"`
	Interval           string            `toml:"interval"`
	TLSCA              string            `toml:"tls_ca"`
	TLSCert            string            `toml:"tls_cert"`
	TLSKey             string            `toml:"tls_key"`
	InsecureSkipVerify bool              `toml:"insecure_skip_verify"`
	Tags               map[string]string `toml:"tags"`

	lookupdEndpoint  string
	nsqdEndpointList map[string]interface{}

	httpClient *http.Client
	duration   time.Duration

	pauseCh chan bool
	pause   bool
}

func newInput() *Input {
	return &Input{
		Tags:             make(map[string]string),
		nsqdEndpointList: make(map[string]interface{}),
		pauseCh:          make(chan bool, 1),
		httpClient:       &http.Client{Timeout: 5 * time.Second},
	}
}

func (*Input) SampleConfig() string { return sampleCfg }

func (*Input) Catalog() string { return catalog }

func (*Input) AvailableArchs() []string { return datakit.AllArch }

func (*Input) SampleMeasurement() []inputs.Measurement {
	return []inputs.Measurement{
		&nsqTopicMeasurement{},
	}
}

func (this *Input) Run() {
	l = logger.SLogger(inputName)

	if this.setup() {
		return
	}

	gatherTicker := time.NewTicker(this.duration)
	defer gatherTicker.Stop()

	updateListTicker := time.NewTicker(updateEndpointListInterval)
	defer updateListTicker.Stop()

	for {
		select {
		case <-datakit.Exit.Wait():
			l.Info("exit")
			return

		case <-gatherTicker.C:
			if this.pause {
				l.Debugf("not leader, skipped")
				continue
			}
			this.gather()

		case <-updateListTicker.C:
			if this.pause {
				l.Debugf("not leader, skipped")
				continue
			}
			if this.isLookupd() {
				if err := this.updateEndpointListByLookupd(this.lookupdEndpoint); err != nil {
					l.Error(err)
					continue
				}
				l.Debugf("nsqd endpoint list: %v", this.nsqdEndpointList)
			}

		case this.pause = <-this.pauseCh:
			// nil
		}
	}
}

func (this *Input) setup() bool {
	var err error
	this.duration, err = timex.ParseDuration(this.Interval)
	if err != nil {
		l.Warnf("parse duration error: %s", err)
	}
	if this.duration < minInterval {
		this.duration = defaultInterval
		l.Debugf("interval should large %s, got %s, use default interval %s", minInterval, this.Interval, defaultInterval)
	}

	if this.httpClient == nil {
		this.httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	for {
		select {
		case <-datakit.Exit.Wait():
			l.Info("exit")
			return true
		default:
			// nil
		}
		time.Sleep(time.Second)

		if this.TLSCA != "" {
			tlsconfig := &net.TLSClientConfig{
				CaCerts:            []string{this.TLSCA},
				Cert:               this.TLSCert,
				CertKey:            this.TLSKey,
				InsecureSkipVerify: this.InsecureSkipVerify,
			}

			tc, err := tlsconfig.TLSConfig()
			if err != nil {
				l.Error(err)
				continue
			}
			this.httpClient.Transport = &http.Transport{TLSClientConfig: tc}
		}

		if this.isLookupd() {
			u, err := buildURL(fmt.Sprintf(lookupdPattern, this.Lookupd))
			if err != nil {
				l.Error(err)
				continue
			}
			this.lookupdEndpoint = u.String()
			if err := this.updateEndpointListByLookupd(this.lookupdEndpoint); err != nil {
				l.Error(err)
			}
		} else {
			if len(this.NSQDs) == 0 {
				l.Error("invalid nsqd endpoints")
				continue
			}
			for _, n := range this.NSQDs {
				u, err := buildURL(fmt.Sprintf(nsqdStatsPattern, n))
				if err != nil {
					l.Error(err)
					continue
				}
				this.nsqdEndpointList[u.String()] = nil
			}
		}

		break
	}

	return false
}

func (this *Input) isLookupd() bool {
	return this.Lookupd != ""
}

func (this *Input) gather() {
	if len(this.nsqdEndpointList) == 0 {
		return
	}
	start := time.Now()

	var wg sync.WaitGroup

	target := newStats(this.Tags)
	for endpoint := range this.nsqdEndpointList {
		wg.Add(1)
		go func(e string) {
			defer wg.Done()

			if err := this.gatherEndpoint(e, target); err != nil {
				l.Warn(err)
			}
		}(endpoint)
	}

	wg.Wait()

	pts, err := target.makePoint()
	if err != nil {
		l.Error(err)
	}

	if len(pts) == 0 {
		return
	}
	if err := io.Feed(inputName, datakit.Metric, pts, &io.Option{CollectCost: time.Since(start)}); err != nil {
		l.Error(err)
	}
}

func (this *Input) gatherEndpoint(endpoint string, target *stats) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	body, err := this.httpGet(u.String())
	if err != nil {
		return err
	}

	return target.add(body)
}

func (this *Input) updateEndpointListByLookupd(lookupdEndpoint string) error {
	body, err := this.httpGet(lookupdEndpoint)
	if err != nil {
		return err
	}

	var endpoints []string
	lk := &LookupNodes{}
	if err := json.Unmarshal(body, lk); err != nil {
		return fmt.Errorf("error parsing response: %s", err)
	}

	for idx := range lk.Producers {
		addr := fmt.Sprintf("http://%s:%d", lk.Producers[idx].BroadcastAddress, lk.Producers[idx].HTTPPort)
		endpoint, err := buildURL(fmt.Sprintf(nsqdStatsPattern, addr))
		if err != nil {
			l.Warn(err)
			continue
		}
		endpoints = append(endpoints, endpoint.String())
	}

	for _, endpoint := range endpoints {
		if _, ok := this.nsqdEndpointList[endpoint]; !ok {
			this.nsqdEndpointList[endpoint] = nil
		}
	}

	return nil
}

func (this *Input) httpGet(u string) ([]byte, error) {
	r, err := this.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error while polling %s: %s", u, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %s", u, r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf(`error reading body: %s`, err)
	}

	return body, nil
}

func (this *Input) Pause() error {
	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()
	select {
	case this.pauseCh <- true:
		return nil
	case <-tick.C:
		return fmt.Errorf("pause %s failed", inputName)
	}
}

func (this *Input) Resume() error {
	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()
	select {
	case this.pauseCh <- false:
		return nil
	case <-tick.C:
		return fmt.Errorf("resume %s failed", inputName)
	}
}

func buildURL(u string) (*url.URL, error) {
	addr, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("unable to parse address '%s': %s", u, err)
	}
	return addr, nil
}

func init() {
	inputs.Add(inputName, func() inputs.Input {
		return newInput()
	})
}
