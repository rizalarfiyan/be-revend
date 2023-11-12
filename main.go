package main

import (
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/logger"
)

func init() {
	config.Init()
	conf := config.Get()
	logger.Init(conf)
	logger.UpdateLogLevel(conf.Logger.Level)
}

func main() {
	logs := logger.Get("main")
	logs.Debug("Hello, World!")
}
