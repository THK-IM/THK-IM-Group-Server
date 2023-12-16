package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-group-server/pkg/app"
)

func RegisterGroupApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()
	userAuth := middleware.UserTokenAuth(appCtx.Context)
	ipAuth := middleware.WhiteIpAuth(appCtx.Context)
	var authMiddleware gin.HandlerFunc
	if appCtx.Config().DeployMode == conf.DeployExposed {
		authMiddleware = userAuth
	} else if appCtx.Config().DeployMode == conf.DeployBackend {
		authMiddleware = ipAuth
	} else {
		panic(errors.New("check deployMode conf"))
	}

	groupRoute := httpEngine.Group("group")
	groupRoute.Use(authMiddleware)
	{
		groupRoute.POST("", createGroup(appCtx))                // 创建群
		groupRoute.GET("/:id", queryGroup(appCtx))              // 查询群资料
		groupRoute.PUT("/:id", updateGroup(appCtx))             // 修改群信息
		groupRoute.DELETE("/:id", deleteGroup(appCtx))          // 删除群
		groupRoute.POST("/:id/transfer", transferGroup(appCtx)) // 转让群

		groupRoute.GET("/apply", queryJoinGroupApply(appCtx))              // 查询群申请列表
		groupRoute.POST("/apply", createJoinGroupApply(appCtx))            // 申请加入群
		groupRoute.POST("/apply/:id/review", reviewJoinGroupApply(appCtx)) // 审核群加入申请
		groupRoute.POST("/invite", inviteJoinGroup(appCtx))                // 邀请加入
	}
}
