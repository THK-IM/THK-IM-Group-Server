package errorx

import "github.com/thk-im/thk-im-base-server/errorx"

var ErrGroupNotExisted = errorx.NewErrorX(4003001, "group not existed")
var ErrGroupJoinNeedApply = errorx.NewErrorX(4003002, "you need apply to join group")
