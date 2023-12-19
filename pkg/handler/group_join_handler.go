package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/logic"
)

func createJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.JoinGroupApplyReq{}
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().Errorf("createJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errCreate := applyLogic.CreateJoinGroupApply(req)
		if errCreate != nil {
			appCtx.Logger().Errorf("createJoinGroupApply %v %v", req, errCreate)
			baseDto.ResponseInternalServerError(context, errCreate)
		} else {
			appCtx.Logger().Infof("createJoinGroupApply %v ", req)
			baseDto.ResponseSuccess(context, nil)
		}
	}
}

func reviewJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.ReviewJoinGroupReq{}
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().Errorf("reviewJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errReview := applyLogic.ReviewJoinGroupApply(req)
		if errReview != nil {
			appCtx.Logger().Errorf("reviewJoinGroupApply %v %v", req, errReview)
			baseDto.ResponseInternalServerError(context, errReview)
		} else {
			appCtx.Logger().Infof("reviewJoinGroupApply %v", req)
			baseDto.ResponseSuccess(context, nil)
		}
	}
}

func queryJoinGroupApply(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.QueryJoinGroupApplyListReq{}
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().Errorf("queryJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errQuery := applyLogic.QueryJoinGroupApplyList(req)
		if errQuery != nil {
			appCtx.Logger().Errorf("queryJoinGroupApply %v %v", req, errQuery)
			baseDto.ResponseInternalServerError(context, errQuery)
		} else {
			appCtx.Logger().Infof("queryJoinGroupApply %v %v", req, resp)
			baseDto.ResponseSuccess(context, resp)
		}
	}
}

func inviteJoinGroup(appCtx *app.Context) gin.HandlerFunc {
	applyLogic := logic.NewGroupApplyLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.InviteJoinGroupReq{}
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().Errorf("createJoinGroupApply %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errInivte := applyLogic.InviteJoinGroup(req)
		if errInivte != nil {
			appCtx.Logger().Errorf("createJoinGroupApply %v %v", req, errInivte)
			baseDto.ResponseInternalServerError(context, errInivte)
		} else {
			appCtx.Logger().Infof("createJoinGroupApply %v", req)
			baseDto.ResponseSuccess(context, nil)
		}
	}
}
