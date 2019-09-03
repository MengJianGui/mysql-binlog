package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/zssky/log"

	"github.com/mysql-binlog/common/inter"
	"github.com/mysql-binlog/siddontang/go-mysql/mysql"

	"github.com/mysql-binlog/recover/bpct"
	"github.com/mysql-binlog/recover/res"
	"github.com/mysql-binlog/recover/ss"
)

var (
	// path cfs存储路径
	path = flag.String("path", "/export/backup/127.0.0.1", "cfs 远程存储路径")

	// clusterid
	clusterID = flag.Int64("cluster_id", 0, "集群ID")

	// time
	time = flag.String("time", "2999-12-30 23:59:59", "截止时间")

	// db
	db = flag.String("dbreg", "", "需要恢复的库名正则")

	// tb
	tb = flag.String("tbreg", "", "需要恢复的表名正则")

	// user
	user = flag.String("user", "root", "恢复目标 MySQL user")

	// password
	passwd = flag.String("password", "secret", "恢复目标 MySQL password")

	// type
	rt = flag.String("rt", "snapshot", "恢复类型recover type including{recover, snapshot} two kinds")

	// dump for MySQL using separated mode or integrated mode
	mode = flag.String("mode", "integrated", "separated or integrated 表示是否将每个表的binlog事件独立而不往一个binlog文件写")

	// log level
	level = flag.String("level", "debug", "日志级别log level {debug/info/warn/error}")
)

// logger 初始化logger
func logger() {
	// 日志输出到标准输出
	log.SetOutput(os.Stdout)

	// 设置日志级别
	log.SetLevelByString(*level)
}

func main() {
	// 解析导入参数
	flag.Parse()

	// init logger
	logger()

	log.Infof("path path{%s}, cluster id {%d}, db {%s}, table{%s}, user{%s}, time{%s}, log level{%s}, mode {%s}", *path, *clusterID, *db, *tb, *user, *time, *level, *mode)
	t := inter.ParseTime(*time)

	c := ss.NewCluster(*path, *mode, *clusterID)
	// check timestamp
	b, err := c.CheckTime(t)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Infof("all table timestamp > parameter timestamp{%d} is %v", t, b)
	if !b {
		// all table timestamp < parameter timestamp
		log.Warnf("latest timestamp for increment data < para.timestamp", t)
		os.Exit(0)
	}

	// recover tables
	rtbs, err := c.SelectTables(*db, *tb)
	if err != nil {
		os.Exit(1)
	}
	log.Infof("tables {%v} match reg{%s.%s}", rtbs, *db, *tb)

	// total tables
	ttbs, err := c.SelectTables(".*", ".*")
	if err != nil {
		os.Exit(1)
	}
	log.Debugf("take total {%v} tables {%s.%s}", ttbs, ".*", ".*")

	// take the 1st offset
	s, err := ss.NewSnapshot(*path, *clusterID, t)
	if err != nil {
		log.Errorf("func NewSnapshot error{%v}", err)
		os.Exit(1)
	}

	// take newly offset
	o, err := s.Offset()
	if err != nil {
		os.Exit(1)
	}
	log.Infof("init newly offset{%v}", o)

	// copy
	if err := s.CopyData(); err != nil {
		log.Errorf("func CopyData from cfs error{%v}", err)
		os.Exit(1)
	}
	log.Infof("snapshot copy ")

	// copy conf
	if err := s.CopyBin(); err != nil {
		log.Errorf("func CopyBin error{%v}", err)
		os.Exit(1)
	}
	log.Infof("copy conf")

	// auth
	if err := s.Auth(); err != nil {
		log.Errorf("grant all auth to file to mysql error{%v}", err)
		os.Exit(1)
	}
	log.Infof("auth file accessory")

	// start MySQL
	if err := s.StartMySQL(); err != nil {
		log.Errorf("func StartMySQL error{%v}", err)
		os.Exit(1)
	}

	// New local MySQL connection POOl
	oi, err := bpct.NewInstance(*user, *passwd, 3358)
	if err != nil {
		os.Exit(1)
	}
	defer oi.Close()

	// MySQL check
	if err := oi.Check(); err != nil {
		os.Exit(1)
	}

	// newly context
	ctx, cancel := context.WithCancel(context.Background())

	// init error channels
	errs := make(chan error, 64)
	defer close(errs)

	rs, err := res.Recovering(*mode, rtbs, ttbs, c.GetClusterPath(), t, ctx, o, *user, *passwd, 3358, errs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case e := <-errs:
				cancel()
				log.Error(e)
				os.Exit(1)
			}
		}
	}()

	log.Infof("flush data on MySQL")
	// flush tables with read lock; flush logs;
	if err := oi.Flush(); err != nil {
		log.Errorf("flush MySQL data for cluster id{%d} error {%v}", *clusterID, err)
		os.Exit(1)
	}

	// if just recover then here to return
	switch *rt {
	case "recover":
		log.Infof("recover for {%d} to timestamp{%s} success", *clusterID, *time)
		for range errs {
			// block forever
		}
	default:
		log.Debugf("keep continue for cluster id {%d}", *clusterID)
		// keep going
	}

	log.Infof("to stop MySQL server")
	if err := s.StopMySQL(*user, *passwd); err != nil {
		log.Errorf("stop MySQL using user{%s} and password{*******} error{%v}", *user, err)
		os.Exit(1)
	}

	log.Infof("to copy data to cfs")
	if err := s.Copy2Cfs(); err != nil {
		log.Errorf("copy data to cfs error {%v}", err)
		os.Exit(1)
	}

	// take gtid
	og, err := mysql.ParseMysqlGTIDSet(o.ExedGtid)
	if err != nil {
		log.Errorf("parse mysql gtid{%s} error{%v}", o.ExedGtid, err)
		os.Exit(1)
	}
	// write newly offset to snapshot directory
	for _, t := range rs {
		for _, g := range strings.Split(t.ExecutedGTID(), ",") {
			g := strings.TrimSpace(g)
			if strings.EqualFold(g, "") {
				// skip empty string
				continue
			}
			if err := og.Update(g); err != nil {
				log.Errorf("merge gtid {%s} into original gtid set{%s} error{%v}", g, o.ExedGtid, err)
				os.Exit(1)
			}
		}
	}

	o.ExedGtid = og.String()
	o.Time = uint32(t)
	o.CID = *clusterID

	if err := s.FlushOffset(o); err != nil {
		log.Errorf("flush offset{%v} to snapshot{%s} index file error{%v}", o, s.ID(), err)
		os.Exit(1)
	}

	// return success
	os.Exit(0)
}
