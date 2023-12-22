package logic

import (
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/errorx"
	"github.com/thk-im/thk-im-group-server/pkg/model"
	msgDto "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	msgModel "github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"time"
)

func SendReviewGroupJoinMessage(appCtx *app.Context, apply *model.GroupMemberApply, sessionId int64, claims baseDto.ThkClaims) error {
	role := msgModel.SessionAdmin
	querySessionUsersReq := &msgDto.QuerySessionUsersReq{
		SId:   sessionId,
		Role:  &role,
		MTime: 0,
		Count: 100,
	}
	sessionUsersResp, errSu := appCtx.MsgApi().QuerySessionUsers(sessionId, querySessionUsersReq, claims)
	if errSu != nil {
		return errSu
	}
	if sessionUsersResp == nil || len(sessionUsersResp.Data) == 0 {
		return errorx.ErrGroupNoAdminOrOwner
	}
	body, errBody := dto.NewReviewGroupJoinMsgBody(apply.Id, apply.GroupId, apply.UIds, apply.Type, apply.ApplyUserId).ToJson()
	if errBody != nil {
		return errBody
	}
	admins := make([]int64, 0)
	for _, su := range sessionUsersResp.Data {
		admins = append(admins, su.UId)
	}
	sendMsgReq := &msgDto.SendSysMessageReq{
		Type:      dto.SysMsgTypeReviewJoinGroup,
		CTime:     time.Now().UnixMilli(),
		Body:      body,
		Receivers: admins,
	}
	_, errSend := appCtx.MsgApi().SendSysMessage(sendMsgReq, claims)
	return errSend
}

func SendRejectGroupJoinMessage(appCtx *app.Context, apply *model.GroupMemberApply, claims baseDto.ThkClaims) error {
	body, errBody := dto.NewRejectGroupJoinMsgBody(apply.Id, apply.GroupId, apply.UIds, apply.Type, apply.ApplyUserId).ToJson()
	if errBody != nil {
		return errBody
	}

	sendMsgReq := &msgDto.SendSysMessageReq{
		Type:      dto.SysMsgTypeRejectJoinGroup,
		CTime:     time.Now().UnixMilli(),
		Body:      body,
		Receivers: []int64{apply.ApplyUserId},
	}
	_, errSend := appCtx.MsgApi().SendSysMessage(sendMsgReq, claims)
	return errSend
}

func SendGroupJoinedMessage(appCtx *app.Context, apply *model.GroupMemberApply, sessionId int64, claims baseDto.ThkClaims) error {
	body, errBody := dto.NewGroupJoinMsgBody(apply.UIds, apply.Type, apply.ApplyUserId).ToJson()
	if errBody != nil {
		return errBody
	}
	sendMsgReq := &msgDto.SendMessageReq{
		CId:   appCtx.SnowflakeNode().Generate().Int64(),
		SId:   sessionId,
		Type:  dto.MsgTypeJoinGroup,
		CTime: time.Now().UnixMilli(),
		Body:  body,
		FUid:  0,
	}
	_, errSend := appCtx.MsgApi().SendSessionMessage(sendMsgReq, claims)
	return errSend
}

func SendGroupQuitMessage(appCtx *app.Context, uIds string, quitType int8, oprUId, sessionId int64, claims baseDto.ThkClaims) error {
	body, errBody := dto.NewGroupQuitMsgBody(uIds, quitType, oprUId).ToJson()
	if errBody != nil {
		return errBody
	}
	sendMsgReq := &msgDto.SendMessageReq{
		CId:   appCtx.SnowflakeNode().Generate().Int64(),
		SId:   sessionId,
		Type:  dto.MsgTypeQuitGroup,
		CTime: time.Now().UnixMilli(),
		Body:  body,
		FUid:  0,
	}
	_, errSend := appCtx.MsgApi().SendSessionMessage(sendMsgReq, claims)
	return errSend
}

func SendGroupDisbandMessage(appCtx *app.Context, oprUId, sessionId int64, claims baseDto.ThkClaims) error {
	body, errBody := dto.NewGroupDisbandMsgBody(oprUId).ToJson()
	if errBody != nil {
		return errBody
	}
	sendMsgReq := &msgDto.SendMessageReq{
		CId:   appCtx.SnowflakeNode().Generate().Int64(),
		SId:   sessionId,
		Type:  dto.MsgTypeDisbandGroup,
		CTime: time.Now().UnixMilli(),
		Body:  body,
		FUid:  0,
	}
	_, errSend := appCtx.MsgApi().SendSessionMessage(sendMsgReq, claims)
	return errSend
}

func SendGroupTransferMessage(appCtx *app.Context, oldOwnerId, newOwnerId, sessionId int64, claims baseDto.ThkClaims) error {
	body, errBody := dto.NewGroupTransferMsgBody(oldOwnerId, newOwnerId).ToJson()
	if errBody != nil {
		return errBody
	}
	sendMsgReq := &msgDto.SendMessageReq{
		CId:   appCtx.SnowflakeNode().Generate().Int64(),
		SId:   sessionId,
		Type:  dto.MsgTypeTransferGroup,
		CTime: time.Now().UnixMilli(),
		Body:  body,
		FUid:  0,
	}
	_, errSend := appCtx.MsgApi().SendSessionMessage(sendMsgReq, claims)
	return errSend
}
