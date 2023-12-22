package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/logic"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func createJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.JoinGroupApplyReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createJoinGroupApply %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errCreate := applyLogic.CreateJoinGroupApply(req, claims)
		if errCreate != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createJoinGroupApply %v %v", req, errCreate)
			baseDto.ResponseInternalServerError(ctx, errCreate)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("createJoinGroupApply %v ", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func reviewJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.ReviewJoinGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("reviewJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("reviewJoinGroupApply %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errReview := applyLogic.ReviewJoinGroupApply(req, claims)
		if errReview != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("reviewJoinGroupApply %v %v", req, errReview)
			baseDto.ResponseInternalServerError(ctx, errReview)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("reviewJoinGroupApply %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func queryJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.QueryJoinGroupApplyListReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryJoinGroupApply %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		resp, errQuery := applyLogic.QueryJoinGroupApplyList(req, claims)
		if errQuery != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryJoinGroupApply %v %v", req, errQuery)
			baseDto.ResponseInternalServerError(ctx, errQuery)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("queryJoinGroupApply %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func inviteJoinGroup(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.InviteJoinGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("cancelInviteJoinGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errInvite := applyLogic.InviteJoinGroup(req, claims)
		if errInvite != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createJoinGroupApply %v %v", req, errInvite)
			baseDto.ResponseInternalServerError(ctx, errInvite)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("createJoinGroupApply %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func cancelInviteJoinGroup(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.CancelInviteJoinGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("cancelInviteJoinGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("cancelInviteJoinGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errInvite := applyLogic.CancelInviteJoinGroup(req, claims)
		if errInvite != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("cancelInviteJoinGroup %v %v", req, errInvite)
			baseDto.ResponseInternalServerError(ctx, errInvite)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("cancelInviteJoinGroup %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteGroupMember(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.DeleteGroupMemberReq{}
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteGroupMember %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errDelete := applyLogic.DeleteGroupMember(req, claims)
		if errDelete != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteGroupMember %v %v", req, errDelete)
			baseDto.ResponseInternalServerError(context, errDelete)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteGroupMember %v", req)
			baseDto.ResponseSuccess(context, nil)
		}
	}
}
