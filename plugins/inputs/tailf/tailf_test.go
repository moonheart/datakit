// +build !solaris

package tailf

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
)

var (
	root = "/tmp/tailf_test"

	dir = "/tmp/tailf_test/1/2"

	paths = []string{
		"/tmp/tailf_test/zero.txt",
		"/tmp/tailf_test/1/one.txt",
		"/tmp/tailf_test/1/2/two.txt",
	}
)

func __init() {
	logger.SetGlobalRootLogger("", logger.DEBUG, logger.OPT_DEFAULT)
	l = logger.SLogger(inputName)
}

func TestWrite(t *testing.T) {
	defer func() {
		os.RemoveAll(root)
	}()

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	var files []*os.File
	for _, path := range paths {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	count := 0
	for {
		if count > 1000 {
			return
		}
		for index, file := range files {
			file.WriteString(time.Now().Format(time.RFC3339Nano) +
				fmt.Sprintf(" -- index: %d -- count: %d\n", index, count))
			time.Sleep(100 * time.Millisecond)
		}
		count++
	}
}

func TestMain(t *testing.T) {
	__init()
	testAssert = true

	var tailer = Tailf{
		Regexs:           []string{".txt"},
		Paths:            []string{root},
		Source:           "NAXXRAMAS",
		FormBeginning:    false,
		UpdateFiles:      true,
		UpdateFilesCycle: 3 * time.Second,
	}

	go tailer.Run()

	time.Sleep(90 * time.Second)
}

func TestFileList(t *testing.T) {

	for {
		fmt.Println(getFileList([]string{root}))
		time.Sleep(500 * time.Millisecond)
	}
}
