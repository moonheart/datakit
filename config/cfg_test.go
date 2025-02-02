// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	bstoml "github.com/BurntSushi/toml"
	tu "gitlab.jiagouyun.com/cloudcare-tools/cliutils/testutil"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io/dataway"
)

func TestInitCfg(t *testing.T) {
	c := DefaultConfig()

	tomlfile := ".main.toml"
	defer os.Remove(tomlfile) //nolint:errcheck
	tu.Equals(t, nil, c.InitCfg(tomlfile))
}

func TestEnableDefaultsInputs(t *testing.T) {
	cases := []struct {
		list   string
		expect []string
	}{
		{
			list:   "a,a,b,c,d",
			expect: []string{"a", "b", "c", "d"},
		},

		{
			list:   "a,b,c,d",
			expect: []string{"a", "b", "c", "d"},
		},
	}

	c := DefaultConfig()
	for _, tc := range cases {
		c.EnableDefaultsInputs(tc.list)
		tu.Equals(t, len(c.DefaultEnabledInputs), len(tc.expect))
	}
}

func TestSetupGlobalTags(t *testing.T) {
	localIP, err := datakit.LocalIP()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		tags   map[string]string
		expect map[string]string
		fail   bool
	}{
		{
			tags: map[string]string{
				"host": "__datakit_hostname", // ENV host dropped during setup
				"ip":   "__datakit_ip",
				"id":   "__datakit_id",
				"uuid": "__datakit_uuid",
			},
			expect: map[string]string{
				"ip": localIP,
			},
		},

		{
			tags: map[string]string{
				"host": "$datakit_hostname", // ENV host dropped during setup
				"ip":   "$datakit_ip",
				"id":   "$datakit_id",
				"uuid": "$datakit_uuid",
			},
			expect: map[string]string{
				"ip": localIP,
			},
		},

		{
			tags: map[string]string{
				"uuid": "some-uuid",
				"host": "some-host", // ENV host dropped during setup
			},

			expect: map[string]string{
				"uuid": "some-uuid",
			},
		},
	}

	for idx, tc := range cases {
		c := DefaultConfig()
		for k, v := range tc.tags {
			c.GlobalTags[k] = v
		}

		err := c.setupGlobalTags()
		if tc.fail {
			tu.NotOk(t, err, "")
		} else {
			tu.Ok(t, err)
		}

		for k, v := range c.GlobalTags {
			if tc.expect == nil {
				tu.Assert(t, v != tc.expect[k], "[case %d] `%s' != `%s', global tags: %+#v", idx, v, tc.expect[k], c.GlobalTags)
			} else {
				tu.Assert(t, v == tc.expect[k], "[case %d] `%s' != `%s', global tags: %+#v", idx, v, tc.expect[k], c.GlobalTags)
			}
		}
	}
}

func TestProtectedInterval(t *testing.T) {
	cases := []struct {
		enabled              bool
		min, max, in, expect time.Duration
	}{
		{
			enabled: true,
			min:     time.Minute,
			max:     5 * time.Minute,
			in:      time.Second,
			expect:  time.Minute,
		},

		{
			enabled: true,
			min:     time.Minute,
			max:     5 * time.Minute,
			in:      10 * time.Minute,
			expect:  5 * time.Minute,
		},

		{
			enabled: false,
			min:     time.Minute,
			max:     5 * time.Minute,
			in:      time.Second,
			expect:  time.Second,
		},

		{
			enabled: false,
			min:     time.Minute,
			max:     5 * time.Minute,
			in:      10 * time.Minute,
			expect:  10 * time.Minute,
		},
	}

	for _, tc := range cases {
		Cfg.ProtectMode = tc.enabled
		x := ProtectedInterval(tc.min, tc.max, tc.in)
		tu.Equals(t, x, tc.expect)
	}
}

func TestDefaultToml(t *testing.T) {
	c := DefaultConfig()

	buf := new(bytes.Buffer)
	if err := bstoml.NewEncoder(buf).Encode(c); err != nil {
		l.Fatalf("encode main configure failed: %s", err.Error())
	}

	t.Logf("%s", buf.String())
}

func TestLoadEnv(t *testing.T) {
	cases := []struct {
		name   string
		envs   map[string]string
		expect *Config
	}{
		{
			name: "normal",
			envs: map[string]string{
				"ENV_GLOBAL_TAGS":                     "a=b,c=d",
				"ENV_LOG_LEVEL":                       "debug",
				"ENV_DATAWAY":                         "http://host1.org,http://host2.com",
				"ENV_HOSTNAME":                        "1024.coding",
				"ENV_NAME":                            "testing-datakit",
				"ENV_HTTP_LISTEN":                     "localhost:9559",
				"ENV_RUM_ORIGIN_IP_HEADER":            "not-set",
				"ENV_ENABLE_PPROF":                    "true",
				"ENV_DISABLE_PROTECT_MODE":            "true",
				"ENV_DEFAULT_ENABLED_INPUTS":          "cpu,mem,disk",
				"ENV_ENABLE_ELECTION":                 "1",
				"ENV_DISABLE_404PAGE":                 "on",
				"ENV_DATAWAY_MAX_IDLE_CONNS_PER_HOST": "123",
				"ENV_REQUEST_RATE_LIMIT":              "1234",
				"ENV_DATAWAY_ENABLE_HTTPTRACE":        "any",
				"ENV_DATAWAY_HTTP_PROXY":              "http://1.2.3.4:1234",
			},
			expect: func() *Config {
				cfg := DefaultConfig()

				cfg.Name = "testing-datakit"
				cfg.DataWayCfg = &dataway.DataWayCfg{
					URLs:                []string{"http://host1.org", "http://host2.com"},
					MaxIdleConnsPerHost: 123,
					HTTPProxy:           "http://1.2.3.4:1234",
					Proxy:               true,
					EnableHTTPTrace:     true,
				}

				cfg.HTTPAPI.RUMOriginIPHeader = "not-set"
				cfg.HTTPAPI.Listen = "localhost:9559"
				cfg.HTTPAPI.Disable404Page = true
				cfg.HTTPAPI.RequestRateLimit = 1234.0

				cfg.Logging.Level = "debug"

				cfg.EnablePProf = true
				cfg.Hostname = "1024.coding"
				cfg.ProtectMode = false
				cfg.DefaultEnabledInputs = []string{"cpu", "mem", "disk"}
				cfg.EnableElection = true
				cfg.GlobalTags = map[string]string{
					"a": "b", "c": "d",
				}
				return cfg
			}(),
		},
		{
			name: "test-ENV_REQUEST_RATE_LIMIT",
			envs: map[string]string{
				"ENV_REQUEST_RATE_LIMIT": "1234.0",
			},
			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.HTTPAPI.RequestRateLimit = 1234.0
				return cfg
			}(),
		},

		{
			name: "bad-ENV_REQUEST_RATE_LIMIT",
			envs: map[string]string{
				"ENV_REQUEST_RATE_LIMIT": "0.1234.0",
			},
			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.HTTPAPI.RequestRateLimit = 0
				return cfg
			}(),
		},

		{
			name: "test-ENV_IPDB",
			envs: map[string]string{
				"ENV_IPDB": "iploc",
			},
			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.Pipeline.IPdbType = "iploc"
				return cfg
			}(),
		},

		{
			name: "test-unknown-ENV_IPDB",
			envs: map[string]string{
				"ENV_IPDB": "unknown-ipdb",
			},
			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.Pipeline.IPdbType = "-"
				return cfg
			}(),
		},

		{
			name: "test-ENV_ENABLE_INPUTS",
			envs: map[string]string{
				"ENV_ENABLE_INPUTS": "cpu,mem,disk",
			},
			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.DefaultEnabledInputs = []string{"cpu", "mem", "disk"}
				return cfg
			}(),
		},

		{
			name: "test-ENV_GLOBAL_TAGS",
			envs: map[string]string{
				"ENV_GLOBAL_TAGS": "cpu,mem,disk=sda",
			},

			expect: func() *Config {
				cfg := DefaultConfig()
				cfg.GlobalTags = map[string]string{"disk": "sda"}
				return cfg
			}(),
		},

		{
			name: "test-ENV_DATAWAY_MAX_IDLE_CONNS_PER_HOST",
			envs: map[string]string{
				"ENV_DATAWAY":                         "http://host1.org,http://host2.com",
				"ENV_DATAWAY_MAX_IDLE_CONNS_PER_HOST": "-1",
			},

			expect: func() *Config {
				cfg := DefaultConfig()

				cfg.DataWayCfg = &dataway.DataWayCfg{
					URLs:                []string{"http://host1.org", "http://host2.com"},
					MaxIdleConnsPerHost: 0,
				}

				return cfg
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := DefaultConfig()
			os.Clearenv()
			for k, v := range tc.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}
			if err := c.LoadEnvs(); err != nil {
				t.Error(err)
			}

			a := tc.expect.String()
			b := c.String()
			tu.Equals(t, a, b)
		})
	}
}

func TestUnmarshalCfg(t *testing.T) {
	cases := []struct {
		raw  string
		fail bool
	}{
		{
			raw: `
	name = "not-set"
http_listen="0.0.0.0:9529"
log = "log"
log_level = "debug"
gin_log = "gin.log"
interval = "10s"
output_file = "out.data"
hostname = "iZb.1024"
default_enabled_inputs = ["cpu", "disk", "diskio", "mem", "swap", "system", "net", "hostobject"]
install_date = 2021-03-25T11:00:19Z

[dataway]
  urls = ["http://testing-openway.cloudcare.cn?token=tkn_2dc4xxxxxxxxxxxxxxxxxxxxxxxxxxxx"]
  timeout = "30s"

[global_tags]
  cluster = ""
  global_test_tag = "global_test_tag_value"
  host = "__datakit_hostname"
  project = ""
  site = ""
  lg= "tl"

[[black_lists]]
  hosts = []
  inputs = []

[[white_lists]]
  hosts = []
  inputs = []
	`,
		},

		{
			raw:  `abc = def`, // invalid toml
			fail: true,
		},

		{
			raw: `
name = "not-set"
http_listen=123  # invalid type
log = "log"`,
			fail: true,
		},

		{
			raw: `
name = "not-set"
log = "log"`,
			fail: false,
		},

		{
			raw: `
hostname = "should-not-set"`,
		},
	}

	tomlfile := ".main.toml"

	defer func() {
		os.Remove(tomlfile) //nolint:errcheck
	}()

	for _, tc := range cases {
		c := DefaultConfig()

		if err := ioutil.WriteFile(tomlfile, []byte(tc.raw), 0o600); err != nil {
			t.Fatal(err)
		}

		err := c.LoadMainTOML(tomlfile)
		if tc.fail {
			tu.NotOk(t, err, "")
			continue
		} else {
			tu.Ok(t, err)
		}

		t.Logf("hostname: %s", c.Hostname)

		if err := os.Remove(tomlfile); err != nil {
			t.Error(err)
		}
	}
}

// go test -v -timeout 30s -run ^TestWriteConfigFile$ gitlab.jiagouyun.com/cloudcare-tools/datakit/config
/*
[sinks]

  [[sinks.sink]]
    categories = ["M", "N", "K", "O", "CO", "L", "T", "R", "S"]
    database = "db0"
    host = "1.1.1.1:8086"
    precision = "ns"
    protocol = "http"
    target = "influxdb"
    timeout = "6s"

  [[sinks.sink]]
    categories = ["M", "N", "K", "O", "CO", "L", "T", "R", "S"]
    database = "db1"
    host = "1.1.1.1:8087"
    precision = "ns"
    protocol = "http"
    target = "influxdb"
    timeout = "6s"

[sinks]

  [[sinks.sink]]
*/
func TestWriteConfigFile(t *testing.T) {
	c := DefaultConfig()

	cases := []struct {
		name string
		in   []map[string]interface{}
	}{
		{
			name: "has_data",
			in: []map[string]interface{}{
				{
					"target":     "influxdb",
					"categories": []string{"M", "N", "K", "O", "CO", "L", "T", "R", "S"},
					"host":       "1.1.1.1:8086",
					"protocol":   "http",
					"precision":  "ns",
					"database":   "db0",
					"timeout":    "5s",
				},
				{
					"target":       "logstash",
					"categories":   []string{"L"},
					"host":         "1.1.1.1:8080",
					"protocol":     "http",
					"request_path": "/twitter/tweet/1",
					"timeout":      "5s",
				},
			},
		},
		{
			name: "no_data",
			in: []map[string]interface{}{
				{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.Sinks.Sink = tc.in
			mcdata, err := datakit.TomlMarshal(c)
			if err != nil {
				panic(err)
			}
			fmt.Println("=====================================================")
			fmt.Println(string(mcdata))
		})
	}
}

func TestGetElectionNamespace(t *testing.T) {
	Cfg = DefaultConfig()
	tu.Equals(t, GetElectionNamespace(), "")
	Cfg.Namespace = "test"
	tu.Equals(t, GetElectionNamespace(), "")

	Cfg.EnableElection = true
	tu.Equals(t, GetElectionNamespace(), "test")

	Cfg.Namespace = ""
	tu.Equals(t, GetElectionNamespace(), "default")
}
