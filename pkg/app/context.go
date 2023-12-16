package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-group-server/pkg/loader"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

type Context struct {
	*server.Context
}

func (c *Context) GroupModel() model.GroupModel {
	return c.Context.ModelMap["group"].(model.GroupModel)
}

func (c *Context) GroupMemberApplyModel() model.GroupMemberApplyModel {
	return c.Context.ModelMap["group_apply_member"].(model.GroupMemberApplyModel)
}

func (c *Context) UserApi() userSdk.UserApi {
	return c.Context.SdkMap["user_api"].(userSdk.UserApi)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	c.Context.SdkMap = loader.LoadSdks(c.Config().Sdks, c.Logger())
	c.Context.ModelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
