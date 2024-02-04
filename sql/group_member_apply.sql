CREATE TABLE IF NOT EXISTS `group_member_apply_%s`
(
    `id`             BIGINT  PRIMARY KEY NOT NULL COMMENT '申请id',
    `group_id`       BIGINT  NOT NULL COMMENT '群id',
    `apply_user_id`  BIGINT  NOT NULL COMMENT '申请用户id',
    `channel`        INT     NOT NULL COMMENT '渠道:1群号,2群二维码,3邀请,4分享',
    `content`        TEXT    NOT NULL COMMENT '申请内容',
    `type`           TINYINT NOT NULL COMMENT '类型 1 申请进群，2邀请进群',
    `u_ids`          TEXT    NOT NULL COMMENT '进群用户id多个id#号隔开',
    `review_user_id` BIGINT COMMENT '审核用户id',
    `status`         INT     NOT NULL DEFAULT 0 COMMENT '申请状态 0 申请中，1 拒绝，2 通过',
    `update_time`    BIGINT  NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time`    BIGINT  NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    INDEX `GroupMemberApply_GroupId_IDX` (`group_id`)
);