package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/gluster/gluster-block-restapi/glusterblockrestd/apiserver"
	blockhandlers "github.com/gluster/gluster-block-restapi/glusterblockrestd/handlers"
	"github.com/gluster/gluster-block-restapi/pkg/utils"

	log "github.com/sirupsen/logrus"
)

// Below variables are set as flags during build time. The current
// values are just placeholders
var (
	version         = ""
	defaultConfFile = ""
)

var (
	showVersion = flag.Bool("version", false, "Show the version information")
	configFile  = flag.String("config", defaultConfFile, "Config file path")
)

func dumpVersionInfo() {
	fmt.Printf("version   : %s\n", version)
	fmt.Printf("go version: %s\n", runtime.Version())
	fmt.Printf("go OS/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func main() {
	// Init logger with stderr
	err := initLogger("", "-", "info")
	if err != nil {
		log.Fatal("Init logging failed for stderr")
	}

	flag.Parse()

	if *showVersion {
		dumpVersionInfo()
		return
	}
	conf, err := loadConfig(*configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load config file")
	}

	err = validateAddress(conf.Addr)
	if err != nil {
		log.WithError(err).Fatal("Failed to start glusterblockrestd server")
	}
	// Create Log dir
	err = os.MkdirAll(conf.LogDir, 0750)
	if err != nil {
		log.WithError(err).WithField("logdir", conf.LogDir).
			Fatal("Failed to create log directory")
	}

	err = initLogger(conf.LogDir, conf.LogFile, conf.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize logger")
	}

	blockhandlers.SetGlusterBlockCLI(conf.GlusterBlockCLIPath)

	serverRunOpts := &apiserver.ServerRunOptions{
		Addr:                  conf.Addr,
		EnableTLS:             conf.EnableTLS,
		CertFile:              conf.CertFile,
		KeyFile:               conf.KeyFile,
		CorsAllowedOriginList: []string{"*"},
		CorsAllowedMethodList: []string{"GET", "POST", "DELETE", "PUT"},
	}
	server := apiserver.NewServer(serverRunOpts)

	//Add required middlewares here
	server.AddMiddleware(middleware.WithReqLogger)
	errChan := make(chan error)
	closeChan := utils.SetSignalHandler()

	go server.Run(errChan)

	select {
	case err := <-errChan:
		log.WithError(err).Error("encounter an error in starting glusterblockrestd")
	case <-closeChan:
		if err := server.Stop(); err != nil {
			log.WithError(err).Error("error in stopping glusterblock api server")
		}
	}
}
