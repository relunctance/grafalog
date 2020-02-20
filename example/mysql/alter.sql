CREATE database test;
use test;
CREATE TABLE `logtests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `cost_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'cost seconds time ',
  `ukey` varchar(255) NOT NULL DEFAULT '' COMMENT 'unique key',
  `api` varchar(10) NOT NULL DEFAULT '',
  `value` varchar(255) NOT NULL DEFAULT '',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'create time',
  `msg` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='test-log';
