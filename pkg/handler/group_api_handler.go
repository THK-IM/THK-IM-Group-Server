package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	"github.com/thk-im/thk-im-group-server/pkg/dto"
	"github.com/thk-im/thk-im-group-server/pkg/logic"
)

func createGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.CreateGroupReq{}
		err := context.Bind(req)
		if err != nil {
			appCtx.Logger().Error("createGroup", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errCreate := groupLogic.CreatGroup(req)
		if errCreate != nil {
			appCtx.Logger().Error("createGroup", err)
			baseDto.ResponseInternalServerError(context, errCreate)
		} else {
			appCtx.Logger().Info("createGroup", "success")
			baseDto.ResponseSuccess(context, resp)
		}
	}
}

func joinGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.JoinGroupReq{}
		err := context.Bind(req)
		if err != nil {
			appCtx.Logger().Error("joinGroup", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errJoin := groupLogic.JoinGroup(req)
		if errJoin != nil {
			appCtx.Logger().Error("joinGroup", err)
			baseDto.ResponseInternalServerError(context, errJoin)
		} else {
			appCtx.Logger().Info("joinGroup", "success")
			baseDto.ResponseSuccess(context, resp)
		}
	}
}

func deleteGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.DeleteGroupReq{}
		err := context.Bind(req)
		if err != nil {
			appCtx.Logger().Error("deleteGroup", err)
			baseDto.ResponseBadRequest(context)
		}

		errDelete := groupLogic.DeleteGroup(req)
		if errDelete != nil {
			appCtx.Logger().Error("deleteGroup", err)
			baseDto.ResponseInternalServerError(context, errDelete)
		} else {
			appCtx.Logger().Info("deleteGroup", "success")
			baseDto.ResponseSuccess(context, nil)
		}
	}
}

func queryGroup(appCtx *app.Context) gin.HandlerFunc {
	//groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
	}
}

func transferGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.TransferGroupReq{}
		err := context.Bind(req)
		if err != nil {
			appCtx.Logger().Error("transferGroup", err)
			baseDto.ResponseBadRequest(context)
		}

		errTransfer := groupLogic.TransferGroup(req)
		if errTransfer != nil {
			appCtx.Logger().Error("transferGroup", err)
			baseDto.ResponseInternalServerError(context, errTransfer)
		} else {
			appCtx.Logger().Info("transferGroup", "success")
			baseDto.ResponseSuccess(context, nil)
		}
	}
}

func updateGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(context *gin.Context) {
		req := &dto.UpdateGroupReq{}
		err := context.Bind(req)
		if err != nil {
			appCtx.Logger().Error("updateGroup", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errUpdate := groupLogic.UpdateGroup(req)
		if errUpdate != nil {
			appCtx.Logger().Error("updateGroup", err)
			baseDto.ResponseInternalServerError(context, errUpdate)
		} else {
			appCtx.Logger().Info("updateGroup", "success")
			baseDto.ResponseSuccess(context, resp)
		}
	}
}
