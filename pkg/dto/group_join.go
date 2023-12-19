package dto

type InviteJoinGroupReq struct {
	UId     int64  `json:"u_id"` // 邀请人id
	GroupId int64  `json:"group_id"`
	ToUId   int64  `json:"to_u_id"` // 被邀请人id
	Channel int    `json:"channel"`
	Content string `json:"content"`
}

type CancelInviteJoinGroupReq struct {
	UId       int64 `json:"u_id"`
	GroupId   int64 `json:"group_id"`
	InviteUId int64 `json:"invite_u_id"`
}

type JoinGroupApplyReq struct {
	UId     int64  `json:"u_id"`
	GroupId int64  `json:"group_id"`
	Channel int    `json:"channel"`
	Content string `json:"content"`
}

type ReviewJoinGroupReq struct {
	Id      int64 `json:"id"`
	GroupId int64 `json:"group_id"`
	UId     int64 `json:"u_id"`
	Status  int   `json:"status"`
}

type JoinGroupApply struct {
	Id         int64  `json:"id"`
	GroupId    int64  `json:"group_id"`
	Status     int    `json:"status"`
	Channel    int    `json:"channel"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type QueryJoinGroupApplyListReq struct {
	UId     int64 `json:"u_id" form:"u_id"`
	GroupId int64 `json:"group_id" form:"group_id"`
	Count   int   `json:"count" form:"count"`
	Offset  int   `json:"offset" form:"offset"`
}

type QueryJoinGroupApplyListResp struct {
	Total int               `json:"total"`
	Data  []*JoinGroupApply `json:"data"`
}
