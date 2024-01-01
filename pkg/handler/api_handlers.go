package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-group-server/pkg/app"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func RegisterGroupApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()
	ipAuth := middleware.WhiteIpAuth(appCtx.Config().IpWhiteList, appCtx.Logger())
	userApi := appCtx.UserApi()
	userTokenAuth := userSdk.UserTokenAuth(userApi, appCtx.Logger())

	var authMiddleware gin.HandlerFunc
	if appCtx.Config().DeployMode == conf.DeployExposed {
		authMiddleware = userTokenAuth
	} else if appCtx.Config().DeployMode == conf.DeployBackend {
		authMiddleware = ipAuth
	} else {
		panic(errors.New("check deployMode conf"))
	}

	groupRoute := httpEngine.Group("group")
	groupRoute.Use(authMiddleware)
	{
		groupRoute.POST("", createGroup(appCtx))                               // 创建群
		groupRoute.GET("", searchGroup(appCtx))                                // 搜索群资料
		groupRoute.GET("/:id", queryGroup(appCtx))                             // 查询群资料
		groupRoute.PUT("/:id", updateGroup(appCtx))                            // 修改群信息
		groupRoute.POST("/:id/join", joinGroup(appCtx))                        // 加入群
		groupRoute.DELETE("/:id", deleteGroup(appCtx))                         // 删除群
		groupRoute.POST("/:id/transfer", transferGroup(appCtx))                // 转让群
		groupRoute.POST("/:id/member/apply", createJoinGroupApply(appCtx))     // 申请加入群
		groupRoute.GET("/:id/member/apply", queryJoinGroupApply(appCtx))       // 查询群申请列表
		groupRoute.POST("/:id/member/review", reviewJoinGroupApply(appCtx))    // 审核群加入申请
		groupRoute.POST("/:id/member/invite", inviteJoinGroup(appCtx))         // 邀请加入
		groupRoute.DELETE("/:id/member/invite", cancelInviteJoinGroup(appCtx)) // 取消邀请加入
		groupRoute.DELETE("/:id/member", deleteGroupMember(appCtx))            // 删除群成员
	}
}
