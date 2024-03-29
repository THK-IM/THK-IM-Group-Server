package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/handler"
	"os"
)

func init() {
	tmpFolder := "tmp"
	if _, err := os.Stat(tmpFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tmpFolder, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	configPath := "etc/group_server.yaml"
	config := &conf.Config{}
	if err := conf.LoadConfig(configPath, config); err != nil {
		panic(err)
	}
	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterGroupApiHandlers(appCtx)

	appCtx.StartServe()
}
