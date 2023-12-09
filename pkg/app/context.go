package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-group-server/pkg/loader"
	"github.com/thk-im/thk-im-group-server/pkg/model"
)

type Context struct {
	*server.Context
	modelMap map[string]interface{}
}

func (c *Context) GroupModel() model.GroupModel {
	return c.modelMap["group"].(model.GroupModel)
}

func (c *Context) GroupMemberApplyModel() model.GroupMemberApplyModel {
	return c.modelMap["group_apply_member"].(model.GroupMemberApplyModel)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	c.modelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
