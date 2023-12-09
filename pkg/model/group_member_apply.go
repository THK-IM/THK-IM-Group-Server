package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	ApplyStatusInit = iota
	ApplyStatusRejected
	ApplyStatusPassed
)

const (
	ApplyChannelGroupId = iota + 1
	ApplyChannelQRCode
	ApplyChannelInvite
	ApplyChannelShare
)

type (
	GroupMemberApply struct {
		Id           int64  `gorm:"id"`
		GroupId      int64  `gorm:"group_id"`
		ApplyUserId  int64  `gorm:"apply_user_id"`
		InviteUserId *int64 `gorm:"invite_user_id"`
		ReviewUserId *int64 `gorm:"review_user_id"`
		Channel      int    `gorm:"channel"`
		Status       int    `gorm:"status"`
		Content      string `gorm:"content"`
		CreateTime   int64  `gorm:"create_time"`
		UpdateTime   int64  `gorm:"update_time"`
	}

	GroupMemberApplyModel interface {
		findOneById(id, groupId int64) (*GroupMemberApply, error)
		findOneByUIdAndGroupId(uId, groupId int64) (*GroupMemberApply, error)
		InsertApply(groupId, applyUserId int64, inviteUserId, reviewUserId *int64, channel, status int, content string) (*GroupMemberApply, error)
		ReviewApply(id, groupId int64, status int) error
	}

	defaultGroupMemberApplyModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultGroupMemberApplyModel) findOneById(id, groupId int64) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	apply := &GroupMemberApply{}
	err := d.db.Raw(sql, id).Scan(apply).Error
	return apply, err
}

func (d defaultGroupMemberApplyModel) findOneByUIdAndGroupId(uId, groupId int64) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	sql := fmt.Sprintf("select * from %s where group_id = ? and apply_user_id = ? order by create_time desc limit 0, 1", tableName)
	apply := &GroupMemberApply{}
	err := d.db.Raw(sql, groupId, uId).Scan(apply).Error
	return apply, err
}

func (d defaultGroupMemberApplyModel) InsertApply(groupId, applyUserId int64, inviteUserId, reviewUserId *int64, channel, status int, content string) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	now := time.Now().UnixMilli()
	apply := &GroupMemberApply{
		Id:           d.snowflakeNode.Generate().Int64(),
		GroupId:      groupId,
		ApplyUserId:  applyUserId,
		InviteUserId: inviteUserId,
		ReviewUserId: reviewUserId,
		Channel:      channel,
		Status:       status,
		Content:      content,
		CreateTime:   now,
		UpdateTime:   now,
	}
	err := d.db.Table(tableName).Create(apply).Error
	return apply, err
}

func (d defaultGroupMemberApplyModel) ReviewApply(id, groupId int64, status int) error {
	tableName := d.genGroupMemberApplyTableName(groupId)
	sql := fmt.Sprintf("update %s set status = ?, update_time = ?  where id = ? ", tableName)
	return d.db.Exec(sql, status, time.Now(), id).Error
}

func (d defaultGroupMemberApplyModel) genGroupMemberApplyTableName(id int64) string {
	return fmt.Sprintf("group_member_apply_%d", id%(d.shards))
}

func NewGroupMemberModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) GroupMemberApplyModel {
	return defaultGroupMemberApplyModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
