package dialtesting

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	//	"github.com/jinzhu/copier"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	uhttp "gitlab.jiagouyun.com/cloudcare-tools/cliutils/network/http"
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/system/rtpanic"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	dt "gitlab.jiagouyun.com/cloudcare-tools/kodo/dialtesting"
)

var (
	AuthorizationType = `DIAL_TESTING`
	SignHeaders       = []string{
		`Content-MD5`,
		`Content-Type`,
		`Date`,
	}

	inputName = "dialtesting"
	l         = logger.DefaultSLogger(inputName)

	MaxFails = 100

	chromeCtxs chan *context.Context
)

const (
	maxCrashCnt = 6
	RegionInfo  = "region"
)

var (
	apiTasksNum      int
	headlessTasksNum int
	chromeCurCount   int
)

type Input struct {
	Region       string            `toml:"region,omitempty"`
	RegionId     string            `toml:"region_id"`
	Server       string            `toml:"server,omitempty"`
	AK           string            `toml:"ak"`
	SK           string            `toml:"sk"`
	PullInterval string            `toml:"pull_interval,omitempty"`
	TimeOut      *datakit.Duration `toml:"time_out,omitempty"` //单位为秒
	Workers      int               `toml:"workers,omitempty"`
	Tags         map[string]string

	cli *http.Client
	//class string

	curTasks map[string]*dialer
	wg       sync.WaitGroup
	pos      int64 // current largest-task-update-time
}

const sample = `
[[inputs.dialtesting]]
  # 中心任务存储的服务地址，即df_dialtesting center service。
  # 此处同时可配置成本地json 文件全路径 "files:///your/dir/json-file-name", 为task任务的json字符串。
  server = "https://dflux-dial.dataflux.cn"

  # require，节点惟一标识ID
  region_id = "default"

  # 若server配为中心任务服务地址时，需要配置相应的ak或者sk
  ak = ""
  sk = ""

  pull_interval = "1m"

  time_out = "1m"
  workers = 6
  [inputs.dialtesting.tags]
  # some_tag = "some_value"
  # more_tag = "some_other_value"
  # ...`

func (dt *Input) SampleConfig() string {
	return sample
}

func (dt *Input) Catalog() string {
	return "network"
}

func (dt *Input) SampleMeasurement() []inputs.Measurement {
	return []inputs.Measurement{
		&httpMeasurement{},
	}

}

func (i *Input) AvailableArchs() []string {
	return datakit.AllArch
}

func (i *Input) NewChromePool(total int) chan *context.Context {
	if total == 0 {
		total = 1
	} //默认为1

	p := make(chan *context.Context, total)
	for i := 0; i < total; i++ {
		p <- dt.NewChromedpCtx(false, ``)
	}

	return p
}

func (d *Input) Run() {

	l = logger.SLogger(inputName)

	// 根据Server配置，若为服务地址则定时拉取任务数据；
	// 若为本地json文件，则读取任务

	if d.Workers == 0 {
		d.Workers = 6
	}

	chromeCtxs = make(chan *context.Context, d.Workers)

	reqURL, err := url.Parse(d.Server)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return
	}

	l.Debugf(`%+#v, %+#v`, d.cli, d.TimeOut)

	if d.TimeOut == nil {
		d.cli.Timeout = 60 * time.Second
	} else {
		d.cli.Timeout = d.TimeOut.Duration
	}

	switch reqURL.Scheme {
	case "http", "https":
		d.doServerTask() // task server

	case "file":
		d.doLocalTask(reqURL.Path)

	case "":
		d.doLocalTask(reqURL.String())

	default:
		l.Warnf(`no invalid scheme`)
	}
}

func (d *Input) doServerTask() {

	var f rtpanic.RecoverCallback

	f = func(stack []byte, err error) {

		defer rtpanic.Recover(f, nil)

		du, err := time.ParseDuration(d.PullInterval)
		if err != nil {
			l.Warnf("invalid frequency: %s, use default", d.PullInterval)
			du = time.Minute
		}
		if du > 24*time.Hour || du < time.Minute {
			l.Warnf("invalid frequency: %s, use default", d.PullInterval)
			du = time.Minute
		}

		tick := time.NewTicker(du)
		defer tick.Stop()

		for {
			select {
			case <-tick.C:
				j, err := d.pullTask()
				if err != nil {
					l.Warnf(`%s,ignore`, err.Error())
					continue
				}
				l.Debugf(`task: %s %v`, string(j), d.pos)
				d.dispatchTasks(j)

			case <-datakit.Exit.Wait():
				l.Info("exit")
				return

				// TODO: 调接口发送每个任务的执行情况，便于中心对任务的管理
			}
		}
	}

	f(nil, nil)

}

func (d *Input) doLocalTask(path string) {

	j, err := d.getLocalJsonTasks(path)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return
	}

	d.dispatchTasks(j)

	<-datakit.Exit.Wait()
}

func (d *Input) newTaskRun(t dt.Task) (*dialer, error) {

	//
	if err := t.Init(); err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, err
	}

	switch t.Class() {
	case dt.ClassHTTP:
		apiTasksNum++
	case dt.ClassHeadless: // chromedp 缓慢增加
		headlessTasksNum++
		if headlessTasksNum/3+1 > chromeCurCount && chromeCurCount < d.Workers {
			chromeCtxs <- dt.NewChromedpCtx(false, ``)
			chromeCurCount++
			l.Debugf(`worker:%d, chromeCurCount:%d, tasks:%d`, d.Workers, chromeCurCount, headlessTasksNum)
		}
	case dt.ClassDNS:
		// TODO
	case dt.ClassTCP:
		// TODO
	case dt.ClassOther:
		// TODO
	case RegionInfo:
		break
		//no need dealwith
	default:
		l.Errorf("unknown task type")
		break
	}

	dialer, err := newDialer(t, d.Tags)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, err
	}

	d.wg.Add(1)
	go func(id string) {
		defer d.wg.Done()
		protectedRun(dialer)
		l.Infof("input %s exited", id)
	}(t.ID())

	return dialer, nil

}

func protectedRun(d *dialer) {

	crashcnt := 0
	var f rtpanic.RecoverCallback
	f = func(trace []byte, err error) {
		defer rtpanic.Recover(f, nil)
		if trace != nil {
			l.Warnf("task %s panic: %+#v", d.task.ID(), err)
			crashcnt++
			if crashcnt > maxCrashCnt {
				l.Warnf("task %s crashed %d times, exit now", d.task.ID(), crashcnt)
				return
			}
		}
		d.run()
	}

	f(nil, nil)
}

type taskPullResp struct {
	Content map[string]interface{} `json:"content"`
}

func (d *Input) dispatchTasks(j []byte) error {
	var resp taskPullResp

	if err := json.Unmarshal(j, &resp); err != nil {
		l.Error(err)
		return err
	}

	for k, arr := range resp.Content {

		switch k {
		case RegionInfo:
			for k, v := range arr.(map[string]interface{}) {
				switch v.(type) {
				case bool:
					if v.(bool) {
						d.Tags[k] = `true`
					} else {
						d.Tags[k] = `false`
					}
				default:
					d.Tags[k] = v.(string)
				}
			}

		default:
		}
	}

	for k, arr := range resp.Content {

		l.Debugf(`class: %s`, k)

		if k == RegionInfo {
			continue
		}

		for _, j := range arr.([]interface{}) {
			var t dt.Task

			switch k {
			case dt.ClassHTTP:
				t = &dt.HTTPTask{}
			case dt.ClassHeadless:
				t = &dt.HeadlessTask{}
			case dt.ClassDNS:
				// TODO
			case dt.ClassTCP:
				// TODO
			case dt.ClassOther:
				// TODO
			case RegionInfo:
				break
				//no need dealwith
			default:
				l.Errorf("unknown task type: %s", k)
				break
			}
			if err := json.Unmarshal([]byte(j.(string)), &t); err != nil {
				l.Errorf(`%s`, err.Error())
				return err
			}

			// update dialer pos
			ts := t.UpdateTimeUs()
			if d.pos < ts {
				d.pos = ts
			}

			l.Debugf(`%+#v id: %s`, d.curTasks[t.ID()], t.ID())

			if dialer, ok := d.curTasks[t.ID()]; ok { // update task

				if dialer.failCnt >= MaxFails {
					l.Warnf(`failed %d times,ignore`, dialer.failCnt)
					delete(d.curTasks, t.ID())
					continue
				}

				if err := dialer.updateTask(t); err != nil {
					l.Warnf(`%s,ignore`, err.Error())
				}

				if strings.ToLower(t.Status()) == dt.StatusStop {
					delete(d.curTasks, t.ID())
				}

			} else { // create new task

				if strings.ToLower(t.Status()) == dt.StatusStop {
					l.Warnf(`%s status is stop, exit ignore`, t.ID())
					continue
				}

				l.Debugf(`create new task %+#v`, t)
				dialer, err := d.newTaskRun(t)
				if err != nil {
					l.Errorf(`%s, ignore`, err.Error())
				} else {
					d.curTasks[t.ID()] = dialer
				}

			}
		}

		//case dt.ClassHeadless:
	}
	return nil
}

func (d *Input) getLocalJsonTasks(path string) ([]byte, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, err
	}

	//转化结构，json结构转成与kodo服务一样的格式
	var resp map[string][]interface{}
	if err := json.Unmarshal(data, &resp); err != nil {
		l.Error(err)
		return nil, err
	}

	res := map[string]interface{}{}
	for k, v := range resp {
		vs := []string{}
		for _, v1 := range v {
			dt, err := json.Marshal(v1)
			if err != nil {
				l.Error(err)
				return nil, err
			}

			vs = append(vs, string(dt))
		}

		res[k] = vs
	}

	tasks := taskPullResp{
		Content: res,
	}
	rs, err := json.Marshal(tasks)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return rs, nil
}

func (d *Input) pullTask() ([]byte, error) {
	reqURL, err := url.Parse(d.Server)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, err
	}

	var res []byte
	for i := 0; i <= 3; i++ {
		statusCode := 0
		res, statusCode, err = d.pullHTTPTask(reqURL, d.pos)
		if statusCode/100 != 5 { //500 err 重试
			break
		}
	}

	return res, err
}

func signReq(req *http.Request, ak, sk string) {

	so := &uhttp.SignOption{
		AuthorizationType: AuthorizationType,
		SignHeaders:       SignHeaders,
		SK:                sk,
	}

	reqSign, err := so.SignReq(req)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("DIAL_TESTING %s:%s", ak, reqSign))
}

func (d *Input) pullHTTPTask(reqURL *url.URL, sinceUs int64) ([]byte, int, error) {

	reqURL.Path = "/v1/task/pull"
	reqURL.RawQuery = fmt.Sprintf("region_id=%s&since=%d", d.RegionId, sinceUs)

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, 5, err
	}

	bodymd5 := fmt.Sprintf("%x", md5.Sum([]byte("")))
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Content-MD5", bodymd5)
	req.Header.Set("Connection", "close")
	signReq(req, d.AK, d.SK)

	resp, err := d.cli.Do(req)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, 5, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Errorf(`%s`, err.Error())
		return nil, 0, err
	}

	defer resp.Body.Close()
	switch resp.StatusCode / 100 {
	case 2: // ok
		return body, resp.StatusCode / 100, nil
	default:
		l.Warn("request %s failed(%s): %s", d.Server, resp.Status, string(body))
		//error_code = kodo.RegionNotFoundOrDisabled, 停止掉所有任务
		if strings.Contains(string(body), `kodo.RegionNotFoundOrDisabled`) {
			//stop all
			d.stopAlltask()
		}
		return nil, resp.StatusCode / 100, fmt.Errorf("pull task failed")
	}

}

func (d *Input) stopAlltask() {
	for tid, dialer := range d.curTasks {
		dialer.stop()
		delete(d.curTasks, tid)
	}
}

func init() {
	inputs.Add(inputName, func() inputs.Input {
		return &Input{
			Tags:     map[string]string{},
			curTasks: map[string]*dialer{},
			wg:       sync.WaitGroup{},
			cli: &http.Client{
				Timeout: 30 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
					TLSHandshakeTimeout: 30 * time.Second,
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 100,
				},
			},
		}
	})
}
