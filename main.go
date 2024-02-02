package main

import (
	"os"
	"time"
	"github.com/dbubel/vision/config"
	"github.com/dbubel/vision/cmd/api"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
)

var (
	// BUILD_TAG is populated in the CircleCI build process.
	BUILD_TAG      = "unknown"
	BUILD_DATE     = "unknown"
	BUILD_GIT_HASH = "unknown"
)

type TimeFormatter struct {
	logrus.Formatter
}

func (u TimeFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.In(time.Local)
	return u.Formatter.Format(e)
}

func main() {
	log := logrus.New()
	log.SetFormatter(TimeFormatter{Formatter: &logrus.JSONFormatter{}})

	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		logrus.WithError(err).Error("Error parsing config")
		return
	}

	if cfg.Environment == config.ENV_PROD {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	cfg.BuildTag = BUILD_TAG
	cfg.BuildDate = BUILD_DATE
	cfg.GitHash = BUILD_GIT_HASH

	c := cli.NewCLI("vision server", BUILD_TAG)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"serve": func() (cli.Command, error) {
			return &api.ServeCommand{
				Config: cfg,
				Log:    log,
			}, nil
		},
	}
	_, err := c.Run()
	if err != nil {
		logrus.WithError(err).Fatalln("Error running serve command")
	}
}

