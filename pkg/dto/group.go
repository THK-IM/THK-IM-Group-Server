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
	CreateTime  int64   `json:"create_time"`
	UpdateTime  int64   `json:"update_time"`
}

type CreateGroupReq struct {
	UId           int64   `json:"u_id"`
	Members       []int64 `json:"members"`
	GroupName     string  `json:"group_name"`
	GroupAvatar   string  `json:"group_avatar"`
	GroupAnnounce string  `json:"group_announce"` // 群公告
	GroupType     int     `json:"group_type"`     // 2普通群，3 超级群
}

type CreateGroupRes struct {
	*Group
}

type UpdateGroupReq struct {
	GroupId   int64   `json:"group_id"`
	UId       int64   `json:"u_id"`
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
	UId     int64 `json:"u_id"`
	GroupId int64 `json:"group_id"`
}

type JoinGroupReq struct {
	UId     int64  `json:"u_id"`
	GroupId int64  `json:"group_id"`
	Channel int    `json:"channel"`
	Content string `json:"content"`
}

type JoinGroupRes struct {
	*Group
}

type TransferGroupReq struct {
	UId     int64 `json:"u_id"`
	ToUId   int64 `json:"to_u_id"`
	GroupId int64 `json:"group_id"`
}

type QueryLatestGroupReq struct {
	UId    int64 `json:"u_id"`
	Count  int   `json:"count" form:"count"`
	Offset int   `json:"offset" form:"offset"`
}

type QueryLatestGroupRes struct {
	Total int64    `json:"total"`
	Data  []*Group `json:"data"`
}

type QueryGroupReq struct {
	UId       int64   `json:"u_id"`
	GroupId   *int64  `json:"group_id"`
	DisplayId *string `json:"display_id"`
}

type QueryGroupRes struct {
	*Group
}
