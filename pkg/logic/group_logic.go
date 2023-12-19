package logic

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	msgDto "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	msgModel "github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"image/color"
	"strconv"
	"time"
)

type GroupLogic struct {
	appCtx *app.Context
}

func NewGroupLogic(appCtx *app.Context) *GroupLogic {
	return &GroupLogic{appCtx: appCtx}
}

func (g *GroupLogic) CreatGroup(req *dto.CreateGroupReq) (*dto.CreateGroupRes, error) {
	groupId := g.appCtx.GroupModel().NewGroupId()
	displayId := strconv.FormatInt(groupId, 36)
	memberCount := len(req.Members) + 1

	var qrcodeUrl *string = nil
	qrFileName := fmt.Sprintf("%s-%d-qrcode.png", req.GroupName, time.Now().UnixMilli()/1000)
	qrFilePath := fmt.Sprintf("tmp/%s", qrFileName)
	url := fmt.Sprintf("https://api.thkim.com/group/%s", displayId)
	errQrcode := qrcode.WriteColorFile(url, qrcode.Medium, 256, color.Black, color.White, qrFilePath)
	if errQrcode != nil {
		g.appCtx.Logger().Error(errQrcode)
	} else {
		qrCodeKey := fmt.Sprintf("group/avatar/%d/%s", groupId, qrFileName)
		qrcodeUrl, errQrcode = g.appCtx.ObjectStorage().UploadObject(qrCodeKey, qrFilePath)
		if errQrcode != nil {
			g.appCtx.Logger().Error("upload object file error: ", errQrcode)
		}
	}
	if qrcodeUrl == nil {
		emptyStr := ""
		qrcodeUrl = &emptyStr
	}

	group, err := g.appCtx.GroupModel().CreateGroup(groupId, 0, req.UId, displayId, req.GroupName,
		"", req.GroupAnnounce, *qrcodeUrl, nil, memberCount, model.EnterFlagNoReview,
	)
	if err != nil {
		return nil, err
	}
	createSessionReq := &msgDto.CreateSessionReq{
		UId:      req.UId,
		Type:     req.GroupType,
		EntityId: groupId,
		Members:  req.Members,
		ExtData:  nil,
		Name:     req.GroupName,
		Remark:   "",
	}
	createSessionResp, createErr := g.appCtx.MsgApi().CreateSession(createSessionReq)
	if createErr != nil {
		return nil, createErr
	}
	group.SessionId = createSessionResp.SId
	errReset := g.appCtx.GroupModel().ResetGroupSessionId(groupId, group.SessionId)
	if errReset != nil {
		return nil, errReset
	}
	createGroupRes := &dto.CreateGroupRes{
		Group: g.groupModel2Dto(group),
	}
	return createGroupRes, nil
}

func (g *GroupLogic) UpdateGroup(req *dto.UpdateGroupReq) (*dto.UpdateGroupRes, error) {
	group, err := g.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return nil, err
	}
	if group == nil || group.Id == 0 {
		return nil, errorx.ErrGroupNotExisted
	}
	sessionUser, errSu := g.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId)
	if errSu != nil {
		return nil, errSu
	}
	if sessionUser.SId == 0 || sessionUser.Role == msgModel.SessionMember {
		return nil, errorx.ErrGroupPermission
	}

	err = g.appCtx.GroupModel().UpdateGroup(group.Id, req.Name, req.Avatar, req.Announce, nil, req.ExtData, req.EnterFlag)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Avatar != nil {
		group.Avatar = *req.Avatar
	}
	if req.Announce != nil {
		group.Announce = *req.Announce
	}
	group.ExtData = req.ExtData
	if req.EnterFlag != nil {
		group.EnterFlag = *req.EnterFlag
	}
	res := &dto.UpdateGroupRes{
		Group: g.groupModel2Dto(group),
	}
	return res, nil
}

func (g *GroupLogic) JoinGroup(req *dto.JoinGroupReq) (*dto.JoinGroupRes, error) {
	group, err := g.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errorx.ErrGroupNotExisted
	}
	if group.EnterFlag == model.EnterFlagNoReview {
		return nil, errorx.ErrGroupJoinNeedApply
	}
	if group.EnterFlag == model.EnterFlagAdminInvite {
		return nil, errorx.ErrGroupJoinNeedAdminInvite
	}
	addSessionReq := &msgDto.SessionAddUserReq{
		EntityId: group.Id,
		UIds:     []int64{req.UId},
		Role:     msgModel.SessionMember,
	}
	err = g.appCtx.MsgApi().AddSessionUser(group.SessionId, addSessionReq)
	if err != nil {
		return nil, err
	}
	_ = g.appCtx.GroupModel().AddGroupMember(group.Id, 1)
	group.MemberCount += 1

	joinGroupRes := &dto.JoinGroupRes{
		Group: g.groupModel2Dto(group),
	}
	return joinGroupRes, nil
}

func (g *GroupLogic) DeleteGroup(req *dto.DeleteGroupReq) error {
	group, err := g.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	sessionUser, errSu := g.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId)
	if errSu != nil {
		return errSu
	}
	if sessionUser.SId == 0 {
		return errorx.ErrGroupPermission
	}
	if sessionUser.Role == msgModel.SessionOwner {
		// 群主解散群
		delReq := &msgDto.DelSessionReq{
			Id: group.SessionId,
		}
		errDel := g.appCtx.MsgApi().DelSession(group.SessionId, delReq)
		if errDel != nil {
			return errDel
		} else {
			return g.appCtx.GroupModel().DelGroup(group.Id)
		}
	} else {
		// 群成员退出群
		delReq := &msgDto.SessionDelUserReq{
			UIds: []int64{req.UId},
		}
		errDel := g.appCtx.MsgApi().DelSessionUser(group.SessionId, delReq)
		if errDel != nil {
			return errDel
		} else {
			return g.appCtx.GroupModel().AddGroupMember(group.Id, -1)
		}
	}
}

func (g *GroupLogic) TransferGroup(req *dto.TransferGroupReq) error {
	group, err := g.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	if group.OwnerId != req.UId {
		return errorx.ErrGroupPermission
	}

	return g.appCtx.GroupModel().UpdateGroupOwner(req.GroupId, req.ToUId)
}

func (g *GroupLogic) groupModel2Dto(group *model.Group) *dto.Group {
	return &dto.Group{
		Id:          group.Id,
		DisplayId:   group.DisplayId,
		OwnerId:     group.OwnerId,
		SessionId:   group.SessionId,
		Qrcode:      group.Qrcode,
		MemberCount: group.MemberCount,
		Name:        group.Name,
		Avatar:      group.Avatar,
		Announce:    group.Announce,
		ExtData:     group.ExtData,
		EnterFlag:   group.EnterFlag,
	}
}
