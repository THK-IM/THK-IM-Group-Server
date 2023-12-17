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
		err := context.BindJSON(req)
		if err != nil {
			appCtx.Logger().Errorf("createGroup %v", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errCreate := groupLogic.CreatGroup(req)
		if errCreate != nil {
			appCtx.Logger().Errorf("createGroup %v %v", req, errCreate)
			baseDto.ResponseInternalServerError(context, errCreate)
		} else {
			appCtx.Logger().Infof("createGroup %v %v", req, resp)
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
			appCtx.Logger().Errorf("joinGroup %v", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errJoin := groupLogic.JoinGroup(req)
		if errJoin != nil {
			appCtx.Logger().Errorf("joinGroup %v %v", req, errJoin)
			baseDto.ResponseInternalServerError(context, errJoin)
		} else {
			appCtx.Logger().Infof("joinGroup %v %v", req, resp)
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
			appCtx.Logger().Errorf("deleteGroup %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errDelete := groupLogic.DeleteGroup(req)
		if errDelete != nil {
			appCtx.Logger().Errorf("deleteGroup %v %v", req, errDelete)
			baseDto.ResponseInternalServerError(context, errDelete)
		} else {
			appCtx.Logger().Infof("deleteGroup %v", req)
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
			appCtx.Logger().Errorf("transferGroup %v", err)
			baseDto.ResponseBadRequest(context)
		}

		errTransfer := groupLogic.TransferGroup(req)
		if errTransfer != nil {
			appCtx.Logger().Errorf("transferGroup %v %v", req, errTransfer)
			baseDto.ResponseInternalServerError(context, errTransfer)
		} else {
			appCtx.Logger().Infof("transferGroup %v", req)
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
			appCtx.Logger().Errorf("updateGroup %v", err)
			baseDto.ResponseBadRequest(context)
		}

		resp, errUpdate := groupLogic.UpdateGroup(req)
		if errUpdate != nil {
			appCtx.Logger().Errorf("updateGroup %v %v", req, errUpdate)
			baseDto.ResponseInternalServerError(context, errUpdate)
		} else {
			appCtx.Logger().Infof("updateGroup %v %v", req, resp)
			baseDto.ResponseSuccess(context, resp)
		}
	}
}
