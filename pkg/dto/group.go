package dto

type Group struct {
	Id          int64   `json:"id"`
	DisplayId   string  `json:"display_id"`
	OwnerId     int64   `json:"owner_id"`
	SessionId   int64   `json:"session_id"`
	Qrcode      string  `json:"qrcode"`
	MemberCount int     `json:"member_count"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Announce    string  `json:"announce"`
	ExtData     *string `json:"ext_data"`
	EnterFlag   int     `json:"enter_flag"`
}

type CreateGroupReq struct {
	UserId        int64   `json:"user_id"`
	Members       []int64 `json:"members"`
	GroupName     string  `json:"group_name"`
	GroupAnnounce string  `json:"group_announce"` // 群公告
	GroupType     int     `json:"group_type"`     // 2普通群，3 超级群
}

type CreateGroupRes struct {
	*Group
}

type UpdateGroupReq struct {
	UserId    int64   `json:"user_id"`
	Name      *string `json:"name"`
	Avatar    *string `json:"avatar"`
	Announce  *string `json:"announce"`
	ExtData   *string `json:"ext_data"`
	EnterFlag *int    `json:"enter_flag"`
}

type UpdateGroupRes struct {
	*Group
}

type DeleteGroupReq struct {
	UserId  int64 `json:"user_id"`
	GroupId int64 `json:"group_id"`
}

type JoinGroupReq struct {
	UserId  int64 `json:"user_id"`
	GroupId int64
}

type JoinGroupRes struct {
	*Group
}

type TransferGroupReq struct {
	UId     int64 `json:"u_id"`
	ToUId   int64 `json:"to_u_id"`
	GroupId int64 `json:"group_id"`
}

type QueryGroupListReq struct {
	UserId int64 `json:"user_id" form:"user_id"`
	Count  int   `json:"count" form:"count"`
	Offset int   `json:"offset" form:"offset"`
}

type QueryGroupListResp struct {
	Total int64    `json:"total"`
	Data  []*Group `json:"data"`
}
