package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
	"hash/crc32"
	"time"
)

const (
	EnterFlagNoReview    = iota // 任意进群
	EnterFlagNeedReview         // 需要审核后进群
	EnterFlagAdminInvite        // 需要管理员邀请才能加入群
)

type (
	Group struct {
		Id          int64   `gorm:"id"`
		DisplayId   string  `gorm:"display_id"`
		SessionId   int64   `gorm:"session_id"`
		OwnerId     int64   `gorm:"owner_id"`
		Name        string  `gorm:"name"`
		Avatar      string  `gorm:"avatar"`
		Announce    string  `gorm:"announce"`
		Qrcode      string  `gorm:"qrcode"`
		ExtData     *string `gorm:"ext_data"`
		MemberCount int     `gorm:"member_count"`
		EnterFlag   int     `gorm:"enter_flag"`
		CreateTime  int64   `gorm:"create_time"`
		UpdateTime  int64   `gorm:"update_time"`
	}

	GroupDisplayId struct {
		DisplayId string `gorm:"display_id"`
		Id        int64  `gorm:"id"`
	}

	GroupModel interface {
		AddGroupMember(groupId int64, count int) error
		NewGroupId() int64
		FindGroup(groupId int64) (*Group, error)
		ResetGroupSessionId(id, sessionId int64) error
		CreateGroup(id, sessionId, ownerId int64, displayId, name, avatar, announce, qrcode string, extData *string, memberCount, enterFlag int) (*Group, error)
		UpdateGroup(groupId int64, name, avatar, announce, qrcode, extData *string, enterFlag, memberCount *int) error
	}

	defaultGroupModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultGroupModel) AddGroupMember(groupId int64, count int) error {
	tableName := d.genGroupTableName(groupId)
	sql := fmt.Sprintf("update %s set member_count = member_count + ? where id = ? ", tableName)
	return d.db.Raw(sql, count, groupId).Error
}

func (d defaultGroupModel) NewGroupId() int64 {
	return d.snowflakeNode.Generate().Int64()
}

func (d defaultGroupModel) FindGroup(groupId int64) (*Group, error) {
	tableName := d.genGroupTableName(groupId)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	group := &Group{}
	err := d.db.Raw(sql, groupId).Scan(group).Error
	return group, err
}

func (d defaultGroupModel) ResetGroupSessionId(id, sessionId int64) error {
	tableName := d.genGroupTableName(id)
	updateMap := make(map[string]interface{})
	updateMap["session_id"] = sessionId
	return d.db.Table(tableName).Where("id=?", id).Updates(updateMap).Error
}

func (d defaultGroupModel) CreateGroup(id, sessionId, ownerId int64, displayId, name, avatar, announce, qrcode string, extData *string, memberCount, enterFlag int) (group *Group, err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	displayIdTableName := d.genGroupDisplayIdTableName(displayId)
	groupDisplay := &GroupDisplayId{
		DisplayId: displayId,
		Id:        id,
	}
	err = tx.Table(displayIdTableName).Create(groupDisplay).Error
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	group = &Group{
		Id:          id,
		DisplayId:   displayId,
		SessionId:   sessionId,
		OwnerId:     ownerId,
		Name:        name,
		Avatar:      avatar,
		Announce:    announce,
		Qrcode:      qrcode,
		ExtData:     extData,
		MemberCount: memberCount,
		EnterFlag:   enterFlag,
		CreateTime:  now,
		UpdateTime:  now,
	}
	tableName := d.genGroupTableName(id)
	err = tx.Table(tableName).Create(group).Error
	return group, err
}

func (d defaultGroupModel) UpdateGroup(groupId int64, name, avatar, announce, qrcode, extData *string, enterFlag, memberCount *int) error {
	if name == nil && avatar == nil && announce == nil && qrcode == nil && extData == nil && enterFlag == nil {
		return nil
	}
	updateMap := make(map[string]interface{})
	if name != nil {
		updateMap["name"] = *name
	}
	if avatar != nil {
		updateMap["avatar"] = *avatar
	}
	if announce != nil {
		updateMap["announce"] = *announce
	}
	if qrcode != nil {
		updateMap["qrcode"] = *qrcode
	}
	if extData != nil {
		updateMap["ext_data"] = *extData
	}
	if enterFlag != nil {
		updateMap["enter_flag"] = *enterFlag
	}
	updateMap["update_time"] = time.Now().UnixMilli()
	return d.db.Table(d.genGroupTableName(groupId)).Where("id = ?", groupId).Updates(updateMap).Error
}

func (d defaultGroupModel) genGroupTableName(id int64) string {
	return fmt.Sprintf("group_%d", id%(d.shards))
}

func (d defaultGroupModel) genGroupDisplayIdTableName(displayId string) string {
	sum := int64(crc32.ChecksumIEEE([]byte(displayId)))
	return fmt.Sprintf("group_display_%d", sum%d.shards)
}

func NewGroupModelModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) GroupModel {
	return defaultGroupModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
