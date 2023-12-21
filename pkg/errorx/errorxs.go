package errorx

import "github.com/thk-im/thk-im-base-server/errorx"

var ErrGroupNotExisted = errorx.NewErrorX(4003001, "group not existed")
var ErrGroupJoinNeedApply = errorx.NewErrorX(4003002, "need apply to join group")
var ErrGroupJoinNeedAdminInvite = errorx.NewErrorX(4003003, "need administrator invite to join group")
var ErrGroupPermission = errorx.NewErrorX(4003004, "group permission error")
var ErrGroupNoAdminOrOwner = errorx.NewErrorX(4003005, "group had no admin or owner ")
