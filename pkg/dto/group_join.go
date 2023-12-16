package dto

type InviteGroupReq struct {
	UId       int64 `json:"u_id"`
	GroupId   int64 `json:"group_id"`
	InviteUId int64 `json:"invite_u_id"`
}

type CancelInviteGroupReq struct {
	UId       int64 `json:"u_id"`
	GroupId   int64 `json:"group_id"`
	InviteUId int64 `json:"invite_u_id"`
}

type JoinGroupApplyReq struct {
	UId     int64 `json:"u_id"`
	GroupId int64 `json:"group_id"`
}

type ReviewJoinGroupReq struct {
	UId int64 `json:"u_id"`
}

type JoinGroupApply struct {
	Id         int64 `json:"id"`
	Status     int   `json:"status"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

type QueryJoinGroupApplyListReq struct {
	UserId  int64 `json:"user_id" form:"user_id"`
	GroupId int64 `json:"group_id" form:"group_id"`
	Count   int   `json:"count" form:"count"`
	Offset  int   `json:"offset" form:"offset"`
}

type QueryJoinGroupApplyListResp struct {
	Total int64             `json:"total"`
	Data  []*JoinGroupApply `json:"data"`
}
