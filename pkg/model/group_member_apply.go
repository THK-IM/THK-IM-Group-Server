package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
	"time"
)

const (
	ApplyStatusInit = iota
	ApplyStatusRejected
	ApplyStatusPassed
)

const (
	TypeInvite = 1
	TypeApply  = 2
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
		Channel      int    `gorm:"channel"`
		Content      string `gorm:"content"`
		Type         int8   `gorm:"type"`
		UIds         string `json:"u_ids"`
		Status       int    `gorm:"status"`
		ReviewUserId *int64 `gorm:"review_user_id"`
		CreateTime   int64  `gorm:"create_time"`
		UpdateTime   int64  `gorm:"update_time"`
	}

	GroupMemberApplyModel interface {
		FindGroupApplies(groupId int64, status *int, count, offset int) ([]*GroupMemberApply, int, error)
		FindOneById(id, groupId int64) (*GroupMemberApply, error)
		FindOneByUIdAndGroupId(uId, groupId int64) (*GroupMemberApply, error)
		InsertApply(groupId, applyUserId int64, reviewUserId *int64, channel, status int, uIds, content string, applyType int8) (*GroupMemberApply, error)
		ReviewApply(id, groupId int64, status int) error
	}

	defaultGroupMemberApplyModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultGroupMemberApplyModel) FindGroupApplies(groupId int64, status *int, count, offset int) ([]*GroupMemberApply, int, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	total := int(0)
	sql := fmt.Sprintf("select count(0) from %s where group_id = ? ", tableName)
	if status != nil {
		sql += fmt.Sprintf("and status = %d", *status)
	}
	err := d.db.Raw(sql, groupId).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	applies := make([]*GroupMemberApply, 0)
	sql = fmt.Sprintf("select * from %s where group_id = ? ", tableName)
	if status != nil {
		sql += fmt.Sprintf("and status = %d ", *status)
	}
	sql += fmt.Sprintf("order by update_time desc limit %d, %d ", count, offset)
	err = d.db.Raw(sql, groupId).Scan(&applies).Error
	return applies, total, err
}

func (d defaultGroupMemberApplyModel) FindOneById(id, groupId int64) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	apply := &GroupMemberApply{}
	err := d.db.Raw(sql, id).Scan(apply).Error
	return apply, err
}

func (d defaultGroupMemberApplyModel) FindOneByUIdAndGroupId(uId, groupId int64) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	sql := fmt.Sprintf("select * from %s where group_id = ? and apply_user_id = ? order by create_time desc limit 0, 1", tableName)
	apply := &GroupMemberApply{}
	err := d.db.Raw(sql, groupId, uId).Scan(apply).Error
	return apply, err
}

func (d defaultGroupMemberApplyModel) InsertApply(groupId, applyUserId int64, reviewUserId *int64, channel, status int, uIds, content string, applyType int8) (*GroupMemberApply, error) {
	tableName := d.genGroupMemberApplyTableName(groupId)
	now := time.Now().UnixMilli()
	apply := &GroupMemberApply{
		Id:           d.snowflakeNode.Generate().Int64(),
		GroupId:      groupId,
		ApplyUserId:  applyUserId,
		UIds:         uIds,
		Type:         applyType,
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
	return fmt.Sprintf("group_join_apply_%d", id%(d.shards))
}

func NewGroupMemberApplyModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) GroupMemberApplyModel {
	return defaultGroupMemberApplyModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
