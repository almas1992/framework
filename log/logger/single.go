package logger

import (
	"path/filepath"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/log/formatter"
	"github.com/goravel/framework/support"
)

type Single struct {
	config config.Config
	json   foundation.Json
}

func NewSingle(config config.Config, json foundation.Json) *Single {
	return &Single{
		config: config,
		json:   json,
	}
}

func (single *Single) Handle(channel string) (logrus.Hook, error) {
	logPath := single.config.GetString(channel + ".path")
	if logPath == "" {
		return nil, errors.LogEmptyLogFilePath
	}

	logPath = filepath.Join(support.RelativePath, logPath)
	levels := getLevels(single.config.GetString(channel + ".level"))
	pathMap := lfshook.PathMap{}
	for _, level := range levels {
		pathMap[level] = logPath
	}

	return lfshook.NewHook(
		pathMap,
		formatter.NewGeneral(single.config, single.json),
	), nil
}

func getLevels(level string) []logrus.Level {
	if level == "panic" {
		return []logrus.Level{
			logrus.PanicLevel,
		}
	}

	if level == "fatal" {
		return []logrus.Level{
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
	}

	if level == "error" {
		return []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
	}

	if level == "warning" {
		return []logrus.Level{
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
	}

	if level == "info" {
		return []logrus.Level{
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
	}

	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
