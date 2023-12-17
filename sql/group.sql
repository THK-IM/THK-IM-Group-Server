CREATE TABLE IF NOT EXISTS `group_%s`
(
    `id`           BIGINT       NOT NULL COMMENT '群id',
    `display_id`   varchar(20)  NOT NULL COMMENT '显示id',
    `session_id`   BIGINT       NOT NULL COMMENT '会话id',
    `owner_id`     BIGINT       NOT NULL COMMENT '群主',
    `name`         varchar(100) NOT NULL COMMENT '群名称',
    `avatar`       TEXT COMMENT '群头像',
    `announce`     TEXT COMMENT '公告',
    `qrcode`       TEXT COMMENT '二维码',
    `member_count` INT COMMENT '群成员数量',
    `ext_data`     Text COMMENT '扩展字段',
    `enter_flag`   INT          NOT NULL default 0 COMMENT '进群条件，0 扫码或通过群id随意进群，1 申请通过后可以进入 2 管理员邀请 ',
    `update_time`  BIGINT       NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time`  BIGINT       NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `Group_IDX` (`id`)
);