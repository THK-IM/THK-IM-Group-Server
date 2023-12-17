package logic

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	msgDto "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
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

	group, err := g.appCtx.GroupModel().CreateGroup(groupId, 0, req.UserId, displayId, req.GroupName,
		"", req.GroupAnnounce, *qrcodeUrl, nil, memberCount, model.EnterFlagNoReview,
	)
	if err != nil {
		return nil, err
	}
	createSessionReq := &msgDto.CreateSessionReq{
		UId:      req.UserId,
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
		Group: &dto.Group{
			Id:          groupId,
			DisplayId:   group.DisplayId,
			OwnerId:     req.UserId,
			SessionId:   createSessionResp.SId,
			Qrcode:      group.Qrcode,
			MemberCount: group.MemberCount,
			Name:        group.Name,
			Avatar:      group.Avatar,
			Announce:    group.Announce,
			ExtData:     group.ExtData,
			EnterFlag:   group.EnterFlag,
		},
	}
	return createGroupRes, nil
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
