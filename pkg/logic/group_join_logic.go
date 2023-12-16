package logic

import (
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
)

type GroupApplyLogic struct {
	appCtx *app.Context
}

func NewGroupApplyLogic(appCtx *app.Context) *GroupApplyLogic {
	return &GroupApplyLogic{appCtx: appCtx}
}

func (l *GroupApplyLogic) CreateJoinGroupApply(req *dto.JoinGroupApplyReq) (*dto.CreateGroupRes, error) {
	return nil, nil
}

func (l *GroupApplyLogic) ReviewJoinGroupApply(req *dto.ReviewJoinGroupReq) error {
	return nil
}

func (l *GroupApplyLogic) InviteJoinGroup(req *dto.InviteGroupReq) error {
	return nil
}

func (l *GroupApplyLogic) CancelInviteJoinGroup(req *dto.CancelInviteGroupReq) error {
	return nil
}

func (l *GroupApplyLogic) QueryJoinGroupApplyList(req *dto.QueryJoinGroupApplyListResp) error {
	return nil
}
