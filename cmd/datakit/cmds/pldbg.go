// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package cmds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/config"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/pipeline"
)

func runPLFlags() error {
	var txt string

	if *flagPLTxtFile != "" {
		txtBytes, err := ioutil.ReadFile(*flagPLTxtFile)
		if err != nil {
			return fmt.Errorf("ioutil.ReadFile: %w", err)
		}
		txt = string(txtBytes)
		txt = strings.TrimSuffix(txt, "\n")
	}

	if txt == "" {
		if *flagPLTxtData != "" {
			txt = *flagPLTxtData
		}
	}

	if txt == "" {
		return fmt.Errorf("no testing string")
	}

	if strings.HasSuffix(txt, "\n") {
		warnf("[E] txt has suffix EOL\n")
	}

	// TODO
	return pipelineDebugger(datakit.Logging, debugPipelineName, txt)
}

func pipelineDebugger(category, plname, txt string) error {
	if err := pipeline.Init(config.Cfg.Pipeline); err != nil {
		return err
	}

	plPath, err := config.GetPipelinePath(category, plname)
	if err != nil {
		return fmt.Errorf("get pipeline failed: %w", err)
	}
	pl, err := pipeline.NewPipelineFromFile(category, plPath)
	if err != nil {
		return fmt.Errorf("new pipeline failed: %w", err)
	}

	start := time.Now()
	opt := &io.PointOption{
		Category: category,
		Time:     time.Now(),
	}
	measurementName := "default"
	pt, err := io.NewPoint(measurementName, nil, map[string]interface{}{pipeline.PipelineMessageField: txt}, opt)
	if err != nil {
		return err
	}
	res, dropFlag, err := pl.Run(pt, nil, *opt)
	if err != nil {
		return fmt.Errorf("run pipeline failed: %w", err)
	}
	cost := time.Since(start)

	if res == nil {
		errorf("[E] No data extracted from pipeline\n")
		return nil
	}

	fields, _ := res.Fields()
	tags := res.Tags()
	if len(fields) == 0 && len(tags) == 0 {
		errorf("[E] No data extracted from pipeline\n")
		return nil
	}

	result := map[string]interface{}{}
	maxWidth := 0

	if *flagPLDate {
		result["time"] = res.Time()
	} else {
		result["time"] = res.Time().UnixNano()
	}

	for k, v := range fields {
		if len(k) > maxWidth {
			maxWidth = len(k)
		}
		result[k] = v
	}

	for k, v := range tags {
		result[k+"#"] = v
		if len(k)+1 > maxWidth {
			maxWidth = len(k) + 1
		}
	}

	measurementName = res.Name()

	if *flagPLTable {
		fmtStr := fmt.Sprintf("%% %ds: %%v", maxWidth)
		lines := []string{}
		for k, v := range result {
			lines = append(lines, fmt.Sprintf(fmtStr, k, v))
		}

		sort.Strings(lines)
		for _, l := range lines {
			fmt.Println(l)
		}
	} else {
		j, err := json.MarshalIndent(result, "", defaultJSONIndent)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(j))
	}

	infof("---------------\n")
	infof("Extracted %d fields, %d tags; measurement(M)<source(L),class(O)...>: %s, drop: %v, cost: %v\n",
		len(fields), len(tags), measurementName, dropFlag, cost)

	return nil
}
