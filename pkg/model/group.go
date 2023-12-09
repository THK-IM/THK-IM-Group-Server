package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	EnterFlagNoReview = iota
	EnterFlagNeedReview
	EnterFlagAdminInvite
)

type (
	Group struct {
		Id          int64   `gorm:"id"`
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

	GroupModel interface {
		FindGroup(groupId int64) (*Group, error)
		CreateGroup(sessionId, ownerId int64, name, avatar, announce, qrcode string, extData *string, memberCount, enterFlag int) (*Group, error)
		UpdateGroup(groupId int64, name, avatar, announce, qrcode, extData *string, enterFlag, memberCount *int) error
	}

	defaultGroupModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultGroupModel) FindGroup(groupId int64) (*Group, error) {
	tableName := d.genGroupTableName(groupId)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	group := &Group{}
	err := d.db.Raw(sql, groupId).Scan(group).Error
	return group, err
}

func (d defaultGroupModel) CreateGroup(sessionId, ownerId int64, name, avatar, announce, qrcode string, extData *string, memberCount, enterFlag int) (*Group, error) {
	id := d.snowflakeNode.Generate().Int64()
	now := time.Now().UnixMilli()
	group := &Group{
		Id:          id,
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
	err := d.db.Table(tableName).Create(group).Error
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

func NewGroupModelModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) GroupModel {
	return defaultGroupModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
