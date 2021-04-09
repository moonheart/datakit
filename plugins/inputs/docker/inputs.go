package docker

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/docker/docker/api/types"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
)

func init() {
	inputs.Add(inputName, func() inputs.Input {
		return &Inputs{
			newEnvClient: NewEnvClient,
			newClient:    NewClient,
			Tags:         make(map[string]string),
		}
	})
}

type Inputs struct {
	Endpoint              string            `toml:"endpoint"`
	CollectMetricInterval string            `toml:"collect_metric_interval"`
	CollectObjectInterval string            `toml:"collect_object_interval"`
	CollectLogging        bool              `toml:"collect_logging"`
	IncludeExited         bool              `toml:"include_exited"`
	ClientConfig                            // tls config
	LogOption             []*LogOption      `toml:"log_option"`
	Tags                  map[string]string `toml:"tags"`

	collectMetricDuration time.Duration
	collectObjectDuration time.Duration
	timeoutDuration       time.Duration

	newEnvClient         func() (Client, error)
	newClient            func(string, *tls.Config) (Client, error)
	containerLogsOptions types.ContainerLogsOptions

	client Client

	opts types.ContainerListOptions
	wg   sync.WaitGroup
}

func (*Inputs) SampleConfig() string {
	return sampleCfg
}

func (*Inputs) Catalog() string {
	return "docker"
}

func (*Inputs) PipelineConfig() map[string]string {
	return nil
}

func (this *Inputs) Run() {
	l = logger.SLogger(inputName)
	if this.initCfg() {
		return
	}
	l.Info("docker input start")

	gatherTick := time.NewTicker(this.collectMetricDuration)

	for {
		select {
		case <-datakit.Exit.Wait():
			return

		case <-gatherTick.C:
			data, err := this.gather()
			if err != nil {
			}
			if err := io.NamedFeed(data, io.Metric, inputName); err != nil {
				l.Error(err)
			}
			this.gatherLog()

		case <-time.After(this.collectObjectDuration):
			data, err := this.gather()
			if err != nil {
			}
			if err := io.NamedFeed(data, io.Object, inputName); err != nil {
				l.Error(err)
			}
		}
	}
}

func (this *Inputs) initCfg() bool {
	for {
		select {
		case <-datakit.Exit.Wait():
			l.Info("exit")
			return true
		default:
			// nil
		}

		if err := this.loadCfg(); err != nil {
			l.Error(err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	return false
}
