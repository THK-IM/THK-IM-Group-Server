package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/handler"
)

func main() {
	configPath := "etc/group_server.yaml"
	config := conf.LoadConfig(configPath)

	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterGroupApiHandlers(appCtx)

	appCtx.StartServe()
}
