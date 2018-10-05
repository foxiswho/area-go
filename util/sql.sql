CREATE TABLE `area` (
  `id` bigint(28) NOT NULL AUTO_INCREMENT,
  `name` CHAR(50) DEFAULT NULL COMMENT '名称',
  `name_traditional` VARCHAR(50) DEFAULT NULL COMMENT '繁体名称',
  `name_en` VARCHAR(100) DEFAULT NULL COMMENT '英文名称',
  `parent_id` INT(11) DEFAULT '0' COMMENT '上级栏目ID',
  `type` TINYINT(4) DEFAULT '0' COMMENT '类别;0默认;1又名;2;3属于;11已合并到;12已更名为',
  `sort` INT(11) DEFAULT '0' COMMENT '排序',
  `type_name` VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
  `other_name` VARCHAR(50) DEFAULT NULL COMMENT '根据类别名称填写',
  `name_format` CHAR(80) DEFAULT NULL COMMENT '格式化全称',
  PRIMARY KEY (`id`),
  KEY `id` (`id`,`parent_id`,`sort`) USING BTREE,
  KEY `name` (`name`),
  KEY `name_format` (`name_format`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='地区表';