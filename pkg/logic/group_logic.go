package logic

import (
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
)

type GroupLogic struct {
	appCtx *app.Context
}

func NewGroupLogic(appCtx *app.Context) *GroupLogic {
	return &GroupLogic{appCtx: appCtx}
}

func (g *GroupLogic) CreatGroup(req *dto.CreateGroupReq) (*dto.CreateGroupRes, error) {
	return nil, nil
}

func (g *GroupLogic) UpdateGroup(req *dto.UpdateGroupReq) (*dto.CreateGroupRes, error) {
	return nil, nil
}

func (g *GroupLogic) JoinGroup(req *dto.JoinGroupReq) (*dto.JoinGroupRes, error) {
	return nil, nil
}

func (g *GroupLogic) DeleteGroup(req *dto.DeleteGroupReq) error {
	return nil
}

func (g *GroupLogic) TransferGroup(req *dto.TransferGroupReq) error {
	return nil
}

func (g *GroupLogic) QueryGroupList(req *dto.QueryGroupListReq) (*dto.QueryGroupListResp, error) {
	return nil, nil
}
