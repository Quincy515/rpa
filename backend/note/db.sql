CREATE DATABASE IF NOT EXISTS `baby` DEFAULT CHARACTER SET utf8mb4;

USE `baby`;

-- 创建用户表
DROP TABLE IF EXISTS `bt_user`;
CREATE TABLE `bt_user` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_sn` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户编号',
    `email` varchar(63) NOT NULL DEFAULT '' COMMENT '注册邮箱',
    `store_passwd` varchar(127) NOT NULL DEFAULT '' COMMENT 'bcrypt加密密码',
    `nickname` varchar(31) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar` varchar(127) NOT NULL DEFAULT '' COMMENT '头像',
    `gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '性别:0-男;1-女;2-保密',
    `introduce` varchar(1022) NOT NULL DEFAULT '' COMMENT '个人简介',
    `state` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态。0-第三方登录绑定用户;1-正常;2-禁发文;3-冻结',
    `is_root` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否超级用户，不限制权限，1-是',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_sn` (`user_sn`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建用户统计表
DROP TABLE IF EXISTS `bt_user_count`;
CREATE TABLE `bt_user_count` (
    `user_sn` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户编号',
    `fans_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝数',
    `follow_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '关注数（关注其他用户）',
    `plan_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '计划数',
    `plan_done` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '已完成计划数',
    `zan_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '被赞数',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户计数表';

-- 第三方登录绑定关系
DROP TABLE IF EXISTS `oauth_member_bind`;
CREATE TABLE `oauth_member_bind` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_sn` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户编号',
    `client_type` varchar(20) NOT NULL DEFAULT '' COMMENT '客户端来源类型:qq,weibo,weixin等',
    `type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '类型 type 1:wechat ',
    `openid` varchar(80) NOT NULL DEFAULT '' COMMENT '第三方id',
    `unionid` varchar(100) NOT NULL DEFAULT '',
    `extra` text NOT NULL COMMENT '额外字段',
    `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
    `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_type_openid` (`type`,`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方登录绑定关系';

-- 创建文件表
DROP TABLE IF EXISTS `bt_file`;
CREATE TABLE `bt_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `create_at` datetime default NOW() COMMENT '创建日期',
    `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
    `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
    `ext2` text COMMENT '备用字段2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 创建用户文件表
CREATE TABLE `bt_user_file` (
    `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增ID',
    `user_sn` bigint unsigned NOT NULL DEFAULT 0 COMMENT '文件所有者用户编号',
    `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    UNIQUE KEY `idx_user_file` (`user_sn`, `file_sha1`),
    KEY `idx_status` (`status`),
    KEY `idx_user_id` (`user_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户文件表';