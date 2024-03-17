package logic

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	msgDto "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	msgModel "github.com/thk-im/thk-im-msgapi-server/pkg/model"
	userDto "github.com/thk-im/thk-im-user-server/pkg/dto"
	"image/color"
	"os"
	"strconv"
)

type GroupLogic struct {
	appCtx *app.Context
}

func NewGroupLogic(appCtx *app.Context) *GroupLogic {
	return &GroupLogic{appCtx: appCtx}
}

func (l *GroupLogic) CreatGroup(req *dto.CreateGroupReq, claims baseDto.ThkClaims) (*dto.CreateGroupRes, error) {
	groupId := l.appCtx.GroupModel().NewGroupId()
	displayId := strconv.FormatInt(groupId, 36)
	memberCount := len(req.Members) + 1

	var qrcodeUrl *string = nil
	qrFileName := fmt.Sprintf("%d-qrcode.png", groupId)
	qrFilePath := fmt.Sprintf("tmp/%s", qrFileName)
	url := fmt.Sprintf("https://api.thkim.com/group/%s", displayId)
	errQrcode := qrcode.WriteColorFile(url, qrcode.Medium, 256, color.Black, color.White, qrFilePath)
	if errQrcode != nil {
		l.appCtx.Logger().Error(errQrcode)
	} else {
		qrCodeKey := fmt.Sprintf("group/%d/qrcode/%s", groupId, qrFileName)
		qrcodeUrl, errQrcode = l.appCtx.ObjectStorage().UploadObject(qrCodeKey, qrFilePath)
		if errQrcode != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).WithFields(logrus.Fields(claims)).Error("upload object file error: ", errQrcode)
		}
		errRemove := os.Remove(qrFilePath)
		if errRemove != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).WithFields(logrus.Fields(claims)).Error("remove file error: ", errRemove)
		}
	}
	if qrcodeUrl == nil {
		emptyStr := ""
		qrcodeUrl = &emptyStr
	}

	avatar := req.GroupAvatar
	if avatar == "" {
		ids := []int64{req.UId}
		ids = append(ids, req.Members...)
		if len(ids) > 9 {
			ids = ids[:9]
		}
		avatarFileName := fmt.Sprintf("%d-out.png", groupId)
		avtarPath, errAvatar := l.generateGroupAvatar(groupId, ids, avatarFileName, claims)
		if errAvatar != nil {
			return nil, errAvatar
		}
		avatarKey := fmt.Sprintf("group/%d/avatar/%s", groupId, avatarFileName)
		avatarUrl, errUpload := l.appCtx.ObjectStorage().UploadObject(avatarKey, avtarPath)
		if errUpload != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).WithFields(logrus.Fields(claims)).Error("upload object file error: ", errUpload)
		}
		errRemove := os.Remove(avtarPath)
		if errRemove != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).WithFields(logrus.Fields(claims)).Error("remove file error: ", errRemove)
		}
		avatar = *avatarUrl
	}

	group, err := l.appCtx.GroupModel().CreateGroup(groupId, 0, req.UId, displayId, req.GroupName,
		avatar, req.GroupAnnounce, *qrcodeUrl, nil, memberCount, model.EnterFlagNoReview,
	)
	if err != nil {
		return nil, err
	}
	createSessionReq := &msgDto.CreateSessionReq{
		UId:          req.UId,
		Type:         req.GroupType,
		EntityId:     groupId,
		Members:      req.Members,
		ExtData:      nil,
		Name:         req.GroupName,
		Remark:       "",
		FunctionFlag: msgDto.FuncTextFlag | msgDto.FuncAudioFlag | msgDto.ImageFlag | msgDto.VideoFlag | msgDto.ForwardFlag | msgDto.ForwardFlag,
	}
	createSessionResp, createErr := l.appCtx.MsgApi().CreateSession(createSessionReq, claims)
	if createErr != nil {
		return nil, createErr
	}
	group.SessionId = createSessionResp.SId
	errReset := l.appCtx.GroupModel().ResetGroupSessionId(groupId, group.SessionId)
	if errReset != nil {
		return nil, errReset
	}
	createGroupRes := &dto.CreateGroupRes{
		Group: l.groupModel2Dto(group),
	}
	return createGroupRes, nil
}

func (l *GroupLogic) generateGroupAvatar(groupId int64, members []int64, outName string, claims baseDto.ThkClaims) (string, error) {
	req := &userDto.QueryUsers{Ids: members}
	userMap, err := l.appCtx.UserApi().QueryUsers(req, claims)
	if err != nil {
		return "", err
	}
	urls := make([]string, 0)
	prefix := fmt.Sprintf("%d", groupId)
	for _, v := range userMap {
		if v.Avatar != nil {
			urls = append(urls, *v.Avatar)
		}
	}
	groupAvatarGenerator := NewGroupAvatarGenerator("tmp", prefix, outName)
	return groupAvatarGenerator.Generate(urls)
}

func (l *GroupLogic) QueryGroup(id int64) (*dto.QueryGroupRes, error) {
	group, err := l.appCtx.GroupModel().FindGroup(id)
	if err != nil {
		return nil, err
	}
	if group.Id == 0 {
		return nil, baseErrorx.ErrParamsError
	}
	return &dto.QueryGroupRes{Group: l.groupModel2Dto(group)}, nil
}

func (l *GroupLogic) SearchGroup(displayId string) (*dto.QueryGroupRes, error) {
	group, err := l.appCtx.GroupModel().FindGroupByDisplayId(displayId)
	if err != nil {
		return nil, err
	}
	if group.Id == 0 {
		return nil, baseErrorx.ErrParamsError
	}
	return &dto.QueryGroupRes{Group: l.groupModel2Dto(group)}, nil
}

func (l *GroupLogic) UpdateGroup(req *dto.UpdateGroupReq, claims baseDto.ThkClaims) (*dto.UpdateGroupRes, error) {
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return nil, err
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

	err = l.appCtx.GroupModel().UpdateGroup(group.Id, req.Name, req.Avatar, req.Announce, nil, req.ExtData, req.EnterFlag)
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
		Group: l.groupModel2Dto(group),
	}
	return res, nil
}

func (l *GroupLogic) JoinGroup(req *dto.JoinGroupReq, claims baseDto.ThkClaims) (*dto.JoinGroupRes, error) {
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
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
	uIds := fmt.Sprintf("%d", req.UId)
	apply, errApply := l.appCtx.GroupMemberApplyModel().InsertApply(group.Id, req.UId, nil, req.Channel, model.ApplyStatusPassed, uIds, req.Content, model.TypeApply)
	if errApply != nil {
		return nil, errApply
	}
	err = l.appCtx.MsgApi().SysAddSessionUser(group.SessionId, addSessionReq, claims)
	if err != nil {
		return nil, err
	}
	_ = l.appCtx.GroupModel().AddGroupMember(group.Id, 1)
	group.MemberCount += 1

	errSend := SendGroupJoinedMessage(l.appCtx, apply, group.SessionId, claims)
	if errSend != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupJoinedMessage: %v %v", apply, err)
	}

	joinGroupRes := &dto.JoinGroupRes{
		Group: l.groupModel2Dto(group),
	}
	return joinGroupRes, nil
}

func (l *GroupLogic) DeleteGroup(req *dto.DeleteGroupReq, claims baseDto.ThkClaims) error {
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	sessionUser, errSu := l.appCtx.MsgApi().QuerySessionUser(group.SessionId, req.UId, claims)
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
		errDel := l.appCtx.MsgApi().DelSession(group.SessionId, delReq, claims)
		if errDel != nil {
			return errDel
		}
		errSend := SendGroupDisbandMessage(l.appCtx, req.UId, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupDisbandMessage: %v %v", req.UId, err)
		}
		return l.appCtx.GroupModel().DelGroup(group.Id)
	} else {
		// 群成员退出群
		delReq := &msgDto.SessionDelUserReq{
			UIds: []int64{req.UId},
		}
		errDel := l.appCtx.MsgApi().SysDelSessionUser(group.SessionId, delReq, claims)
		if errDel != nil {
			return errDel
		}
		uIds := fmt.Sprintf("%d", req.UId)
		errSend := SendGroupQuitMessage(l.appCtx, uIds, dto.QuitTypeLeave, req.UId, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupDisbandMessage: %v %v", req.UId, err)
		}
		return l.appCtx.GroupModel().AddGroupMember(group.Id, -1)
	}
}

func (l *GroupLogic) TransferGroup(req *dto.TransferGroupReq, claims baseDto.ThkClaims) error {
	group, err := l.appCtx.GroupModel().FindGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil || group.Id == 0 {
		return errorx.ErrGroupNotExisted
	}
	if group.OwnerId != req.UId {
		return errorx.ErrGroupPermission
	}
	err = l.appCtx.GroupModel().UpdateGroupOwner(req.GroupId, req.ToUId)
	if err != nil {
		errSend := SendGroupTransferMessage(l.appCtx, req.UId, req.ToUId, group.SessionId, claims)
		if errSend != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendGroupTransferMessage: %v %v", req.UId, err)
		}
	}
	return err
}

func (l *GroupLogic) groupModel2Dto(group *model.Group) *dto.Group {
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
		CreateTime:  group.CreateTime,
		UpdateTime:  group.UpdateTime,
	}
}
