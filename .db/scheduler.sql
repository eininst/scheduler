CREATE TABLE `scheduler_user`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `name`        varchar(64)  DEFAULT '',
    `password`    varchar(100) DEFAULT '',
    `real_name`   varchar(32)  default '',
    `role`        varchar(20)  DEFAULT '',
    `head`        varchar(200) default '',
    `mail`        varchar(200) default '',
    `create_time` varchar(20)  DEFAULT '',
    `status`      varchar(20)  DEFAULT '',
    PRIMARY KEY (`id`),
    KEY           `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

insert into scheduler_user(name,password, real_name, role,status) values('admin','96e79218965eb72c92a549dd5a330112', 'admin','admin','ok');

CREATE TABLE `scheduler_task`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id`      bigint(20) DEFAULT 0,
    `name`         varchar(100)  DEFAULT '',
    `group`        varchar(64)   default '',
    `spec`         varchar(64)   DEFAULT '',
    `url`          varchar(128)  DEFAULT 0,
    `method`       varchar(128)  DEFAULT '',
    `content_type` varchar(128)  DEFAULT '',
    `body`         varchar(1024) DEFAULT '',
    `timeout`      int(10) DEFAULT 0,
    `max_retries`  int(10) DEFAULT 0,
    `desc`        varchar(2048)   default '',
    `status`       varchar(32)   DEFAULT '',
    `create_time`  varchar(20)   DEFAULT '',
    PRIMARY KEY (`id`),
    KEY            `idx_user_id` (`user_id`),
    UNIQUE KEY `idx_name` (`name`),
    KEY            `idx_group` (`group`),
    KEY            `idx_status` (`status`),
    KEY            `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `scheduler_task_excute`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20) DEFAULT '',
    `task_id`     bigint(20) DEFAULT 0,
    `task_name`   varchar(100)  DEFAULT '',
    `task_url`    varchar(128)  DEFAULT '',
    `task_obj`    varchar(2000) DEFAULT '',
    `code`        int(10) DEFAULT '',
    `response`    varchar(2000) DEFAULT '',
    `start_time`  varchar(32)   DEFAULT '',
    `end_time`    varchar(32)   DEFAULT '',
    `duration`    int(10) DEFAULT 0,
    `create_time` varchar(20)   DEFAULT '',
    PRIMARY KEY (`id`),
    KEY           `idx_user_id` (`user_id`),
    KEY           `idx_task_id` (`task_id`),
    KEY           `idx_code` (`code`),
    KEY           `idx_duration` (`duration`),
    KEY           `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
