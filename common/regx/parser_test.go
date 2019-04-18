package regx

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	ddls := []string{
		"ALTER TABLE sbtest1 ADD COLUMN name VARCHAR(100) DEFAULT NULL COMMENT '姓名'",

		//"CREATE INDEX k_1 ON sbtest1(k)",
		//"CREATE INDEX k_1 ON sysbench.sbtest1(k)",
		//"CREATE INDEX k_1 ON `sysbench`.`sbtest1`(k)",
		//
		//"ALTER TABLE 		`mydb.mytable` ADD `field2` DATE  NULL  AFTER `field1`;",
		//"ALTER TABLE `mydb`.`mytable` ADD `field2` DATE  NULL  AFTER `field1`;",
		//"ALTER TABLE `myTable` ADD `field2` DATE  NULL  AFTER `field1`;",
		//"ALTER TABLE mydb.mytable ADD `field2` DATE  NULL  AFTER `field1`;",
		//"ALTER TABLE mytable ADD `field2` DATE  NULL  AFTER `field1`;",
		//"ALTER TABLE mydb.mytable ADD field2 DATE  NULL  AFTER `field1`;",
		//
		//"rename 			table `mydb`.`mytable` to `mydb`.`mytable1`",
		//"rename table `mytable` to `mytable1`",
		//"rename table mydb.mytable to mydb.mytable1",
		//"rename table mytable to mytable1",
		//"rename table `mydb`.`mytable` to `mydb`.`mytable2`, `mydb`.`mytable3` to `mydb`.`mytable1`",
		//"rename table `mytable` to `mytable2`, `mytable3` to `mytable1`",
		//"rename table mydb.mytable to mydb.mytable2, mydb.mytable3 to mydb.mytable1",
		//"rename table mytable to mytable2, mytable3 to mytable1",

		"drop table test1",
		"DROP			 TABLE test1",
		"DROP TABLE test1",
		"DROP table IF EXISTS test.test1",
		"drop table `test1`",
		"DROP TABLE `test1`",
		"DROP table IF EXISTS `test`.`test1`",
		"DROP TABLE `test1` /* generated by server */",
		"DROP table if exists test1",
		"DROP table if exists `test1`",
		"DROP table if exists test.test1",
		"DROP table if exists `test`.test1",
		"DROP table if exists `test`.`test1`",
		"DROP table if exists test.`test1`",
		"DROP table if exists test.`test1`",

		//"CREATE TABLE `position` (`id` bigint(20)) primary key NOT NULL AUTO_INCREMENT COMMENT '主键id'",
		//"CREATE TABLE `position`(`id` bigint(20)) NOT NULL AUTO_INCREMENT COMMENT '主键id'",
		//"CREATE				 TABLE `mydb.mytable` (`id` int(10)) ENGINE=InnoDB",
		//"CREATE TABLE `mytable` (`id` int(10)) ENGINE=InnoDB",
		//"CREATE TABLE IF NOT EXISTS `mytable` (`id` int(10)) ENGINE=InnoDB",
		//"CREATE TABLE IF NOT EXISTS mytable (`id` int(10)) ENGINE=InnoDB",
		//"CREATE TABLE position(`id` bigint(20)) NOT NULL AUTO_INCREMENT COMMENT '主键id'",
	}

	for _, ddl := range ddls {
		tbs, matched := Parse([]byte(strings.TrimSpace(ddl)), []byte("mydb"))
		if matched {
			for _, tb := range tbs {
				fmt.Println(string(tb))
			}
		}
	}
}
