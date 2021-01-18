package pipeline

import (
	"testing"
	"strconv"
	"fmt"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/pipeline/geo"
	"github.com/stretchr/testify/assert"
)

type funcCase struct {
	desc     string
	data     string
	script   string
	expected interface{}
	key      string
	err      error
	fail     bool
}

type EscapeError string

func (e EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

func TestDefaultTimeFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"a":{"time":"06/Jan/2017:16:16:37 +0000","second":2,"thrid":"abc","forth":true},"age":47}`,
			script: `json(_, a.time) default_time(a.time)`,
			expected: "1483719397000000000",
			key: "a.time",
			err: nil,
		},
		{
			data: `{"a":{"time":"2014-12-16 06:20:00 UTC","second":2,"thrid":"abc","forth":true},"age":47}`,
			script: `json(_, a.time) default_time(a.time)`,
			expected: "1418682000000000000",
			key: "a.time",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)

		assert.Equal(t, err, nil)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assert.Equal(t, r, tt.expected)
	}
}

func TestUrlencodeFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"url[0]":"+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B","second":2}`,
			script: "json(_, `url[0]`) url_decode(`url[0]`)",
			expected: " ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;",
			key: "url[0]",
			err: nil,
		},
		{
			data: `{"url":"http%3a%2f%2fwww.baidu.com%2fs%3fwd%3d%e6%b5%8b%e8%af%95","second":2}`,
			script: `json(_, url) url_decode(url)`,
			expected: `http://www.baidu.com/s?wd=测试`,
			key: "url",
			err: nil,
		},
		{
			data: `{"url":"","second":2}`,
			script: `json(_, url) url_decode(url)`,
			expected: ``,
			key: "url",
			err: nil,
		},
		{
			data: `{"url":"+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B","second":2}`,
			script: `json(_, url) url_decode(url)`,
			expected: " ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;",
			key: "url",
			err: nil,
		},
		{
			data: `{"url":"+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B","second":2}`,
			script: `json(_, url) url_decode("url", "aaa")`,
			expected: "+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B",
			key: "url",
			err: nil,
		},
		{
			data: `{"aa":{"url":"http%3a%2f%2fwww.baidu.com%2fs%3fwd%3d%e6%b5%8b%e8%af%95"},"second":2}`,
			script: `json(_, aa.url) url_decode(aa.url)`,
			expected: `http://www.baidu.com/s?wd=测试`,
			key: "aa.url",
			err: nil,
		},
		{
			data: `{"aa":{"aa.url":"http%3a%2f%2fwww.baidu.com%2fs%3fwd%3d%e6%b5%8b%e8%af%95"},"second":2}`,
			script: "json(_, aa.`aa.url`) url_decode(aa.`aa.url`)",
			expected: `http://www.baidu.com/s?wd=测试`,
			key: "aa.aa.url",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)

		assertEqual(t, err, nil)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assertEqual(t, r, tt.expected)
	}
}

func TestGeoIpFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"ip":"116.228.89.206", "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, ip) geoip(ip)`,
			expected: "Shanghai",
			key: "city",
			err: nil,
		},
		{
			data: `{"ip":"192.168.0.1", "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, ip) geoip(ip)`,
			expected: "-",
			key: "city",
			err: nil,
		},
		{
			data: `{"ip":"192.168.0.1", "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, "ip") geoip("ip")`,
			expected: "-",
			key: "city",
			err: nil,
		},
		{
			data: `{"ip":"", "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, "ip") geoip(ip)`,
			expected: "unknown",
			key: "city",
			err: nil,
		},
		{
			data: `{"aa": {"ip":"116.228.89.206"}, "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, aa.ip) geoip("aa.ip")`,
			expected: "Shanghai",
			key: "city",
			err: nil,
		},
	}

	geo.Init()

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assertEqual(t, r, tt.expected)
	}
}

func TestUserAgentFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"userAgent":"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36", "second":2,"thrid":"abc","forth":true}`,
			script: `json(_, userAgent) user_agent(userAgent)`,
			expected: "Windows 7",
			key: "os",
			err: nil,
		},
		{
			data: `{"userAgent":"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"}`,
			script: `json(_, userAgent) user_agent(userAgent)`,
			expected: "Googlebot",
			key: "browser",
			err: nil,
		},
		{
			data: `{"userAgent":"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html"}`,
			script: `json(_, userAgent) user_agent(userAgent)`,
			expected: "",
			key: "engine",
			err: nil,
		},
		{
			data: `{"userAgent":"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"}`,
			script: `json(_, userAgent) user_agent(userAgent)`,
			expected: "bingbot",
			key: "browser",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assertEqual(t, r, tt.expected)
	}
}

func TestDatetimeFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"a":{"timestamp": "1610103765000", "second":2},"age":47}`,
			script: `json(_, a.timestamp) datetime(a.timestamp, 'ms', 'YYYY-MM-dd hh:mm:ss')`,
			expected: "2021-01-08 07:02:45",
			key: "a.timestamp",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assertEqual(t, r, tt.expected)
	}
}

func TestGroupFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"status": 200,"age":47}`,
			script: `json(_, status) group_between(status, [299, 200], "ok")`,
			expected: "ok",
			key: "status",
			err: nil,
		},
		{
			data: `{"status": 200,"age":47}`,
			script: `json(_, status) group_between(status, [299, 200], "ok", newkey)`,
			expected: "",
			key: "newkey",
			err: nil,
		},
		{
			data: `{"status": 200,"age":47}`,
			script: `json(_, status) group_between(status, [200, 299], "ok", newkey)`,
			expected: "ok",
			key: "newkey",
			err: nil,
		},
		{
			data: `{"status": 200,"age":47}`,
			script: `json(_, status) group_between(status, [300, 400], "ok", newkey)`,
			expected: "",
			key: "newkey",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		fmt.Println("=======>", r)

		assertEqual(t, r, tt.expected)
	}
}

func TestGroupInFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"status": "test","age":"47"}`,
			script: `json(_, status) group_in(status, [200, 47, "test"], "ok", newkey)`,
			expected: "ok",
			key: "newkey",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		assertEqual(t, r, tt.expected)
	}
}

func TestNullIfFunc(t *testing.T) {
	var testCase = []*funcCase{
		{
			data: `{"a":{"first": 1,"second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, "1")`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first": "1","second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, 1)`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first": "","second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, "")`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first": true,"second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, true)`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first": 2.3, "second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, 2.3)`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first": 2,"second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, 2)`,
			expected: "",
			key: "a.first",
			err: nil,
		},
		{
			data: `{"a":{"first":"2.3","second":2,"thrid":"aBC","forth":true},"age":47}`,
			script: `json(_, a.first) nullif(a.first, "2.3")`,
			expected: "",
			key: "a.first",
			err: nil,
		},
	}

	for _, tt := range testCase {
		p, err := NewPipeline(tt.script)
		assertEqual(t, err, p.lastErr)

		p.Run(tt.data)

		r, err := p.getContentStr(tt.key)

		fmt.Println("out ======>", p.Output)
		fmt.Println("======>", r)

		assertEqual(t, r, tt.expected)
	}
}
