package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kardianos/service"

	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/cmd/installer/install"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/configtemplate"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/git"
)

var (
	DataKitBaseURL = ""
	DataKitVersion = ""

	datakitUrl = "https://" + path.Join(DataKitBaseURL,
		fmt.Sprintf("datakit-%s-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH, DataKitVersion))

	telegrafUrl = "https://" + path.Join(DataKitBaseURL,
		"telegraf",
		fmt.Sprintf("agent-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH))

	dataUrl = "https://" + path.Join(DataKitBaseURL, "data.tar.gz")

	l = logger.DefaultSLogger("installer")
)

var (
	flagUpgrade     = flag.Bool("upgrade", false, ``)
	flagDatawayHTTP = flag.String("dataway", "", `address of dataway(http://IP:Port?token=xxx), port default 9528`)

	flagInfo         = flag.Bool("info", false, "show installer info")
	flagDownloadOnly = flag.Bool("download-only", false, `download datakit only, not install`)
	flagOTA          = flag.Bool("ota", false, "upgraded by OTA")
	flagInstallOnly  = flag.Bool("install-only", false, "install only, not start")

	flagEnableInputs = flag.String("enable-inputs", "", `default enable inputs(comma splited, example: cpu,mem,disk)`)
	flagDatakitName  = flag.String("name", "", `specify DataKit name, example: prod-env-datakit`)
	flagGlobalTags   = flag.String("global-tags", "", `enable global tags, example: host=__datakit_hostname,ip=__datakit_ip`)
	flagPort         = flag.Int("port", 9529, "datakit HTTP port")

	flagCfgTemplate = flag.String("conf-tmpl", "", `specify input config templates, can be file path or url, e.g, http://res.dataflux.cn/datakit/conf`)

	flagOffline = flag.Bool("offline", false, "offline install mode")
	flagSrcs    = flag.String("srcs", fmt.Sprintf("./datakit-%s-%s-%s.tar.gz,./agent-%s-%s.tar.gz",
		runtime.GOOS, runtime.GOARCH, DataKitVersion, runtime.GOOS, runtime.GOARCH),
		`local path of datakit and agent install files`)
)

const (
	datakitBin = "datakit"

	dlDatakit = "datakit"
	dlAgent   = "agent"
	dlData    = "data"
)

func main() {
	lopt := logger.OPT_DEFAULT | logger.OPT_STDOUT
	if runtime.GOOS != "windows" { // disable color on windows(some color not working under windows)
		lopt |= logger.OPT_COLOR
	}

	logger.SetGlobalRootLogger("", logger.DEBUG, lopt)

	flag.Parse()
	datakit.InitDirs()
	applyFlags()

	// create install dir if not exists
	if err := os.MkdirAll(datakit.InstallDir, 0775); err != nil {
		l.Fatal(err)
	}

	datakit.ServiceExecutable = filepath.Join(datakit.InstallDir, datakitBin)
	if runtime.GOOS == datakit.OSWindows {
		datakit.ServiceExecutable += ".exe"
	}

	svc, err := datakit.NewService()
	if err != nil {
		l.Errorf("new %s service failed: %s", runtime.GOOS, err.Error())
		return
	}

	l.Info("stoping datakit...")
	if err := service.Control(svc, "stop"); err != nil {
		l.Warnf("stop service: %s, ignored", err.Error())
	}

	if *flagOffline && *flagSrcs != "" {
		for _, f := range strings.Split(*flagSrcs, ",") {
			_ = install.ExtractDatakit(f, datakit.InstallDir)
		}
	} else {
		install.CurDownloading = dlDatakit
		install.Download(datakitUrl, datakit.InstallDir, true)
		fmt.Printf("\n")
		install.CurDownloading = dlAgent
		install.Download(telegrafUrl, datakit.InstallDir, true)
		fmt.Printf("\n")
		install.CurDownloading = dlData
		install.Download(dataUrl, datakit.InstallDir, true)
		fmt.Printf("\n")
	}

	if *flagUpgrade { // upgrade new version
		l.Infof("Upgrading to version %s...", DataKitVersion)
		if err := install.UpgradeDatakit(svc); err != nil {
			l.Fatalf("upgrade datakit failed: %s", err.Error())
		}
	} else { // install new datakit
		l.Infof("Installing version %s...", DataKitVersion)
		install.InstallNewDatakit(svc)
	}

	ct := configtemplate.NewCfgTemplate(datakit.InstallDir)
	if err = ct.InstallConfigs(*flagCfgTemplate); err != nil {
		l.Fatalf("fail to intsall config template, %s", err)
	}

	if !*flagInstallOnly {
		l.Infof("starting service %s...", datakit.ServiceName)
		if err = service.Control(svc, "start"); err != nil {
			l.Fatalf("star service: %s", err.Error())
		}
	}

	if *flagUpgrade { // upgrade new version
		l.Info(":) Upgrade Success!")
	} else {
		l.Info(":) Install Success!")
	}

	localIP, err := datakit.LocalIP()
	if err != nil {
		l.Info("get local IP failed: %s", err.Error())
	} else {
		fmt.Printf("\n\tVisit http://%s:%d/stats to see DataKit running status.\n", localIP, *flagPort)
		fmt.Printf("\tVisit http://%s:%d/man to see DataKit manuals.\n\n", localIP, *flagPort)
	}
}

func applyFlags() {

	if *flagInfo {
		fmt.Printf(`
       Version: %s
      Build At: %s
Golang Version: %s
       BaseUrl: %s
       DataKit: %s
      Telegraf: %s
`, git.Version, git.BuildAt, git.Golang, DataKitBaseURL, datakitUrl, telegrafUrl)
		os.Exit(0)
	}

	if *flagDownloadOnly {
		install.DownloadOnly = true

		install.CurDownloading = dlDatakit
		install.Download(datakitUrl,
			fmt.Sprintf("datakit-%s-%s-%s.tar.gz",
				runtime.GOOS, runtime.GOARCH, DataKitVersion), true)
		fmt.Printf("\n")

		install.CurDownloading = dlAgent
		install.Download(telegrafUrl,
			fmt.Sprintf("agent-%s-%s.tar.gz",
				runtime.GOOS, runtime.GOARCH), true)
		fmt.Printf("\n")

		install.CurDownloading = dlData
		install.Download(dataUrl, "data.tar.gz", true)
		fmt.Printf("\n")

		os.Exit(0)
	}

	install.DataWayHTTP = *flagDatawayHTTP
	install.GlobalTags = *flagGlobalTags
	install.Port = *flagPort
	install.DatakitName = *flagDatakitName
	install.EnableInputs = *flagEnableInputs
}
