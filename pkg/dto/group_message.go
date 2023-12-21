package dto

import "encoding/json"

const (
	SysMsgTypeApplyJoinGroup       = -104 // 申请加入群消息
	SysMsgTypeRejectApplyJoinGroup = -105 // 拒绝通过群加入申请消息

	MsgTypeJoinGroup    = 40
	MsgTypeQuitGroup    = 41
	MsgTypeDisbandGroup = 42
)

const (
	QuitTypeLeave    = 1
	QuitTypeBeKicked = 2
)

type GroupApplyJoinMsgBody struct {
	UIds     string `json:"u_ids"`     // 多个uid #号隔开
	JoinType int8   `json:"join_type"` // 1申请自己加入，2申请邀请别人加入
	OprUId   int64  `json:"opr_u_id"`  // 操作人id
}

func (g *GroupApplyJoinMsgBody) ToJson() (string, error) {
	d, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func NewGroupApplyJoinMsgBody(uIds string, joinType int8, oprUId int64) *GroupApplyJoinMsgBody {
	return &GroupApplyJoinMsgBody{
		UIds:     uIds,
		JoinType: joinType,
		OprUId:   oprUId,
	}
}

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
