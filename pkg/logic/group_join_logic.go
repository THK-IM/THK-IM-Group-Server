package logic

import (
	"fmt"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
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

func (l *GroupApplyLogic) CreateJoinGroupApply(req *dto.JoinGroupApplyReq, claims baseDto.ThkClaims) error {
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
		errSend := SendReviewGroupJoinMessage(l.appCtx, apply, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("SendGroupApplyJoinMessage: %v %v", apply, errSend)
		}
	}

	return errApply
}

func (l *GroupApplyLogic) ReviewJoinGroupApply(req *dto.ReviewJoinGroupReq, claims baseDto.ThkClaims) error {
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

	sessionUser, errSu := l.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId, claims)
	if errSu != nil {
		return errSu
	}
	if sessionUser.SId == 0 || sessionUser.Role == msgModel.SessionMember {
		return errorx.ErrGroupPermission
	}

	if req.Status == model.ApplyStatusPassed {
		addSessionUserReq := &msgDto.SessionAddUserReq{
			EntityId: group.Id,
			UIds:     []int64{apply.ApplyUserId},
			Role:     msgModel.SessionMember,
		}
		err = l.appCtx.MsgApi().SysAddSessionUser(group.SessionId, addSessionUserReq, claims)
		if err != nil {
			return err
		}
	}
	err = l.appCtx.GroupMemberApplyModel().ReviewApply(apply.Id, apply.GroupId, req.Status)
	if err == nil {
		if req.Status == model.ApplyStatusPassed {
			errSend := SendGroupJoinedMessage(l.appCtx, apply, group.SessionId, claims)
			if errSend != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("SendGroupJoinedMessage: %v %v", apply, errSend)
			}
		} else {
			errSend := SendRejectGroupJoinMessage(l.appCtx, apply, claims)
			if errSend != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("SendReviewGroupJoinMessage: %v %v", apply, errSend)
			}
		}
	}
	return err
}

func (l *GroupApplyLogic) InviteJoinGroup(req *dto.InviteJoinGroupReq, claims baseDto.ThkClaims) error {
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
			errSend := SendReviewGroupJoinMessage(l.appCtx, apply, group.SessionId, claims)
			if errSend != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("SendGroupApplyJoinMessage: %v %v", apply, errSend)
			}
		}
	} else {
		// 无需要审核，直接加入群
		addSessionReq := &msgDto.SessionAddUserReq{
			EntityId: group.Id,
			UIds:     uIds,
			Role:     msgModel.SessionMember,
		}
		err = l.appCtx.MsgApi().SysAddSessionUser(group.SessionId, addSessionReq, claims)
		if err != nil {
			return err
		}
		errSend := SendGroupJoinedMessage(l.appCtx, apply, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupJoinedMessage: %v %v", apply, errSend)
		}
	}
	return nil
}

func (l *GroupApplyLogic) QueryJoinGroupApplyList(req *dto.QueryJoinGroupApplyListReq, claims baseDto.ThkClaims) (*dto.QueryJoinGroupApplyListResp, error) {
	group, errGroup := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if errGroup != nil {
		return nil, errGroup
	}
	if group == nil || group.Id == 0 {
		return nil, errorx.ErrGroupNotExisted
	}

	sessionUser, errSu := l.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId, claims)
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

func (l *GroupApplyLogic) CancelInviteJoinGroup(req *dto.CancelInviteJoinGroupReq, claims baseDto.ThkClaims) error {
	apply, err := l.appCtx.GroupMemberApplyModel().FindOneById(req.ApplyId, req.GroupId)
	if err != nil {
		return err
	}
	if apply.Id == 0 || apply.Status != model.ApplyStatusPassed || apply.Type != model.TypeInvite || apply.ApplyUserId != req.UId {
		return baseErrorx.ErrParamsError
	}

	group, errGroup := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if errGroup != nil {
		return errGroup
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}

	strUIds := strings.Split(apply.UIds, "#")
	uIds := make([]int64, 0)
	for _, strUId := range strUIds {
		uId, errId := strconv.ParseInt(strUId, 10, 64)
		if errId != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("CancelInviteJoinGroup %v %v", apply.UIds, errId)
		}
		uIds = append(uIds, uId)
	}
	if len(uIds) <= 0 {
		return baseErrorx.ErrParamsError
	}
	delSessionUserReq := &msgDto.SessionDelUserReq{
		UIds: uIds,
	}
	return l.appCtx.MsgApi().SysDelSessionUser(group.SessionId, delSessionUserReq, claims)
}

func (l *GroupApplyLogic) DeleteGroupMember(req *dto.DeleteGroupMemberReq, claims baseDto.ThkClaims) error {
	group, errGroup := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if errGroup != nil {
		return errGroup
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}

	delSessionUserReq := &msgDto.SessionDelUserReq{
		UIds: req.DelUIds,
	}
	errDel := l.appCtx.MsgApi().DelSessionUser(group.SessionId, delSessionUserReq, claims)
	if errDel == nil {
		strUIds := make([]string, 0)
		for _, uId := range req.DelUIds {
			strUIds = append(strUIds, fmt.Sprintf("%d", uId))
		}
		errSend := SendGroupQuitMessage(l.appCtx, strings.Join(strUIds, "#"), dto.QuitTypeBeKicked, req.UId, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupQuitMessage: %v ", errSend)
		}
	}
	return errDel
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
