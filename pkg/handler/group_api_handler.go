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
	"strconv"
)

// curl -i -X POST -d '{"user_id": 1, "members": [5, 6], "group_name": "1-group", "group_announce": "12143242423", "group_type": 2}' "http://192.168.1.9:16000/group"
func createGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.CreateGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		resp, errCreate := groupLogic.CreatGroup(req, claims)
		if errCreate != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createGroup %v %v", req, errCreate)
			baseDto.ResponseInternalServerError(ctx, errCreate)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("createGroup %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

// curl -i -X POST -d '{"user_id": 90}' "http://192.168.1.9:16000/group/1736282436255359732/join"
func joinGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		groupId, errGroupId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errGroupId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("joinGroup %v", errGroupId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req := &dto.JoinGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("joinGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}
		req.GroupId = groupId

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("joinGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		resp, errJoin := groupLogic.JoinGroup(req, claims)
		if errJoin != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("joinGroup %v %v", req, errJoin)
			baseDto.ResponseInternalServerError(ctx, errJoin)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("joinGroup %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func deleteGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.DeleteGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errDelete := groupLogic.DeleteGroup(req, claims)
		if errDelete != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteGroup %v %v", req, errDelete)
			baseDto.ResponseInternalServerError(ctx, errDelete)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteGroup %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func queryGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		strId := ctx.Param("id")
		var groupId int64 = 0
		if strId != "" {
			id, errId := strconv.ParseInt(strId, 10, 64)
			if errId != nil {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryGroup %v", errId)
				baseDto.ResponseBadRequest(ctx)
			}
			groupId = id
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)

		resp, errQuery := groupLogic.QueryGroup(groupId)
		if errQuery != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryGroup %v %v", requestUid, errQuery)
			baseDto.ResponseInternalServerError(ctx, errQuery)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("queryGroup %v %v", requestUid, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func searchGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		displayId := ctx.Query("display_id")

		requestUid := ctx.GetInt64(userSdk.UidKey)

		resp, errQuery := groupLogic.SearchGroup(displayId)
		if errQuery != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryGroup %v %v", requestUid, errQuery)
			baseDto.ResponseInternalServerError(ctx, errQuery)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("queryGroup %v %v", requestUid, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func transferGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.TransferGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("transferGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("transferGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		errTransfer := groupLogic.TransferGroup(req, claims)
		if errTransfer != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("transferGroup %v %v", req, errTransfer)
			baseDto.ResponseInternalServerError(ctx, errTransfer)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("transferGroup %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateGroup(appCtx *app.Context) gin.HandlerFunc {
	groupLogic := logic.NewGroupLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		req := &dto.UpdateGroupReq{}
		err := ctx.BindJSON(req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateGroup %v", err)
			baseDto.ResponseBadRequest(ctx)
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateGroup %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		resp, errUpdate := groupLogic.UpdateGroup(req, claims)
		if errUpdate != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateGroup %v %v", req, errUpdate)
			baseDto.ResponseInternalServerError(ctx, errUpdate)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("updateGroup %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}
