package dto

import "encoding/json"

const (
	SysMsgTypeReviewJoinGroup = -30 // 审核进群消息
	SysMsgTypeRejectJoinGroup = -31 // 拒绝通过群加入申请消息

	MsgTypeJoinGroup     = 10 // 加入群消息类型
	MsgTypeQuitGroup     = 11 // 退出群消息类型
	MsgTypeDisbandGroup  = 12 // 群解散消息类型
	MsgTypeTransferGroup = 13 // 转让群消息类型
)

const (
	QuitTypeLeave    = 1
	QuitTypeBeKicked = 2
)

// ReviewGroupJoinMsgBody 审核进群消息body
type ReviewGroupJoinMsgBody struct {
	ApplyId  int64  `json:"apply_id"`  // 申请id
	GroupId  int64  `json:"group_id"`  // 群id
	UIds     string `json:"u_ids"`     // 多个uid #号隔开
	JoinType int8   `json:"join_type"` // 1申请自己加入，2申请邀请别人加入
	OprUId   int64  `json:"opr_u_id"`  // 操作人id
}

func (g *ReviewGroupJoinMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewReviewGroupJoinMsgBody(applyId, groupId int64, uIds string, joinType int8, oprUId int64) *ReviewGroupJoinMsgBody {
	return &ReviewGroupJoinMsgBody{
		ApplyId:  applyId,
		GroupId:  groupId,
		UIds:     uIds,
		JoinType: joinType,
		OprUId:   oprUId,
	}
}

// RejectGroupJoinMsgBody 拒绝进群消息body
type RejectGroupJoinMsgBody struct {
	ApplyId  int64  `json:"apply_id"`  // 申请id
	GroupId  int64  `json:"group_id"`  // 群id
	UIds     string `json:"u_ids"`     // 多个uid #号隔开
	JoinType int8   `json:"join_type"` // 1申请自己加入，2申请邀请别人加入
	OprUId   int64  `json:"opr_u_id"`  // 操作人id
}

func (g *RejectGroupJoinMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewRejectGroupJoinMsgBody(applyId, groupId int64, uIds string, joinType int8, oprUId int64) *RejectGroupJoinMsgBody {
	return &RejectGroupJoinMsgBody{
		ApplyId:  applyId,
		GroupId:  groupId,
		UIds:     uIds,
		JoinType: joinType,
		OprUId:   oprUId,
	}
}

// GroupJoinMsgBody 进群消息body
type GroupJoinMsgBody struct {
	UIds     string `json:"u_ids"`     // 多个uid #号隔开
	JoinType int8   `json:"join_type"` // 1自己加入，2别人邀请加入
	OprUId   int64  `json:"opr_u_id"`  // 操作人id
}

func (g *GroupJoinMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewGroupJoinMsgBody(uIds string, joinType int8, oprUId int64) *GroupJoinMsgBody {
	return &GroupJoinMsgBody{
		UIds:     uIds,
		JoinType: joinType,
		OprUId:   oprUId,
	}
}

// GroupQuitMsgBody 退群消息body
type GroupQuitMsgBody struct {
	UIds     string `json:"u_ids"`       // 多个uid #号隔开
	QuitType int8   `json:"join_type"`   // 1自己退出，2其他人踢出
	OprUId   int64  `json:"invite_u_id"` // 操作人id
}

func (g *GroupQuitMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewGroupQuitMsgBody(uIds string, quitType int8, oprUId int64) *GroupQuitMsgBody {
	return &GroupQuitMsgBody{
		UIds:     uIds,
		QuitType: quitType,
		OprUId:   oprUId,
	}
}

// GroupDisbandMsgBody 群解散消息body
type GroupDisbandMsgBody struct {
	OprUId int64 `json:"invite_u_id"` // 操作人id
}

func (g *GroupDisbandMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewGroupDisbandMsgBody(oprUId int64) *GroupDisbandMsgBody {
	return &GroupDisbandMsgBody{
		OprUId: oprUId,
	}
}

// GroupTransferMsgBody 群转让消息body
type GroupTransferMsgBody struct {
	OlderOwnerId int64 `json:"older_owner_id"`
	NewOwnerId   int64 `json:"new_owner_id"`
}

func (g *GroupTransferMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewGroupTransferMsgBody(olderOwnerId, newOwnerId int64) *GroupTransferMsgBody {
	return &GroupTransferMsgBody{
		OlderOwnerId: olderOwnerId,
		NewOwnerId:   newOwnerId,
	}
}
