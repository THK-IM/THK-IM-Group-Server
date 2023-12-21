package logic

import (
	"fmt"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	msgDto "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	msgModel "github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"strconv"
	"strings"
)

type GroupApplyLogic struct {
	appCtx *app.Context
}

func NewGroupApplyLogic(appCtx *app.Context) *GroupApplyLogic {
	return &GroupApplyLogic{appCtx: appCtx}
}

func (l *GroupApplyLogic) CreateJoinGroupApply(req *dto.JoinGroupApplyReq) error {
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	if group.EnterFlag == model.EnterFlagAdminInvite {
		return errorx.ErrGroupJoinNeedAdminInvite
	}
	apply, errApply := l.appCtx.GroupMemberApplyModel().InsertApply(group.Id, req.UId, nil, req.Channel, model.ApplyStatusInit, fmt.Sprintf("%d", req.UId), req.Content, model.TypeApply)

	if errApply == nil {
		errSend := SendGroupApplyJoinMessage(l.appCtx, apply, group.SessionId)
		if errSend != nil {
			l.appCtx.Logger().Errorf("SendGroupApplyJoinMessage: %v %v", apply, errSend)
		}
	}

	return errApply
}

func (l *GroupApplyLogic) ReviewJoinGroupApply(req *dto.ReviewJoinGroupReq) error {
	apply, err := l.appCtx.GroupMemberApplyModel().FindOneById(req.Id, req.GroupId)
	if err != nil {
		return err
	}
	if apply == nil || apply.Id == 0 {
		return baseErrorx.ErrParamsError
	}
	if req.Status != model.ApplyStatusInit {
		return baseErrorx.ErrParamsError
	}

	group, errGroup := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if errGroup != nil {
		return errGroup
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}

	sessionUser, errSu := l.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId)
	if errSu != nil {
		return errSu
	}
	if sessionUser.SId == 0 || sessionUser.Role == msgModel.SessionMember {
		return errorx.ErrGroupPermission
	}

	if req.Status == model.ApplyStatusPassed {
		addSessionReq := &msgDto.SessionAddUserReq{
			EntityId: group.Id,
			UIds:     []int64{apply.ApplyUserId},
			Role:     msgModel.SessionMember,
		}
		err = l.appCtx.MsgApi().AddSessionUser(group.SessionId, addSessionReq)
		if err != nil {
			return err
		}
	}
	err = l.appCtx.GroupMemberApplyModel().ReviewApply(apply.Id, apply.GroupId, req.Status)
	if err == nil {
		if req.Status == model.ApplyStatusPassed {
			errSend := SendGroupJoinedMessage(l.appCtx, apply, group.SessionId)
			if errSend != nil {
				l.appCtx.Logger().Errorf("SendGroupJoinedMessage: %v %v", apply, errSend)
			}
		} else {
			// TODO 审核不通过结果告知申请用户
		}
	}
	return err
}

func (l *GroupApplyLogic) InviteJoinGroup(req *dto.InviteJoinGroupReq) error {
	strUIds := strings.Split(req.InviteUIds, "#")
	uIds := make([]int64, 0)
	for _, strUId := range strUIds {
		uId, errStr := strconv.ParseInt(strUId, 10, 64)
		if errStr != nil {
			return baseErrorx.ErrParamsError
		}
		uIds = append(uIds, uId)
	}
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	if group.EnterFlag == model.EnterFlagAdminInvite {
		return errorx.ErrGroupJoinNeedAdminInvite
	}
	if group.EnterFlag != model.EnterFlagNeedReview && group.EnterFlag != model.EnterFlagNoReview {
		// 进群flag值错误
		return baseErrorx.ErrInternalServerError
	}

	apply, errApply := l.appCtx.GroupMemberApplyModel().InsertApply(group.Id, req.UId, nil, req.Channel, model.ApplyStatusInit, req.InviteUIds, req.Content, model.TypeInvite)
	if errApply != nil {
		return errApply
	}
	if group.EnterFlag == model.EnterFlagNeedReview {
		if errApply == nil {
			errSend := SendGroupApplyJoinMessage(l.appCtx, apply, group.SessionId)
			if errSend != nil {
				l.appCtx.Logger().Errorf("SendGroupApplyJoinMessage: %v %v", apply, errSend)
			}
		}
	} else {
		// 无需要审核，直接加入群
		addSessionReq := &msgDto.SessionAddUserReq{
			EntityId: group.Id,
			UIds:     uIds,
			Role:     msgModel.SessionMember,
		}
		err = l.appCtx.MsgApi().AddSessionUser(group.SessionId, addSessionReq)
		if err != nil {
			return err
		}
		errSend := SendGroupJoinedMessage(l.appCtx, apply, group.SessionId)
		if errSend != nil {
			l.appCtx.Logger().Error("SendGroupJoinedMessage: %v %v", apply, err)
		}
	}
	return nil
}

func (l *GroupApplyLogic) QueryJoinGroupApplyList(req *dto.QueryJoinGroupApplyListReq) (*dto.QueryJoinGroupApplyListResp, error) {
	group, errGroup := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if errGroup != nil {
		return nil, errGroup
	}
	if group == nil || group.Id == 0 {
		return nil, errorx.ErrGroupNotExisted
	}

	sessionUser, errSu := l.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId)
	if errSu != nil {
		return nil, errSu
	}
	if sessionUser.SId == 0 || sessionUser.Role == msgModel.SessionMember {
		return nil, errorx.ErrGroupPermission
	}

	applies, total, errQuery := l.appCtx.GroupMemberApplyModel().FindGroupApplies(req.GroupId, nil, req.Count, req.Offset)
	if errQuery != nil {
		return nil, errQuery
	}
	dtoApplies := make([]*dto.JoinGroupApply, 0)
	for _, apply := range applies {
		dtoApply := l.applyModel2Dto(apply)
		dtoApplies = append(dtoApplies, dtoApply)
	}
	resp := &dto.QueryJoinGroupApplyListResp{
		Total: total,
		Data:  dtoApplies,
	}
	return resp, nil
}

func (l *GroupApplyLogic) CancelInviteJoinGroup(req *dto.CancelInviteJoinGroupReq) error {
	return nil
}

func (l *GroupApplyLogic) applyModel2Dto(m *model.GroupMemberApply) *dto.JoinGroupApply {
	return &dto.JoinGroupApply{
		Id:         m.Id,
		GroupId:    m.GroupId,
		Status:     m.Status,
		Channel:    m.Channel,
		Content:    m.Content,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}
