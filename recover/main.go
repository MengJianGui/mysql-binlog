package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/zssky/log"

	"github.com/mysql-binlog/common/inter"
	"github.com/mysql-binlog/siddontang/go-mysql/mysql"

	"github.com/mysql-binlog/recover/bpct"
	"github.com/mysql-binlog/recover/res"
	"github.com/mysql-binlog/recover/ss"
)

var (
	// base cfs存储路径
	base = flag.String("base", "/export/backup/127.0.0.1", "cfs 远程存储路径")

	// clusterid
	clusterID = flag.Int64("clusterid", 0, "集群ID")

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

	log.Infof("base path{%s}, cluster id {%d}, db {%s}, table{%s}, user{%s}, time{%s}, log level{%s}", *base, *clusterID, *db, *tb, *user, *time, *level)
	t := inter.ParseTime(*time)

	c := ss.NewCluster(*base, *clusterID)
	tbs, err := c.SelectTables(*db, *tb)
	if err != nil {
		os.Exit(1)
	}
	log.Infof("tables {%v} match reg{%s.%s}", tbs, *db, *tb)

	// take the 1st offset
	s, err := ss.NewSnapshot(*base, *clusterID, t)
	if err != nil {
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
		os.Exit(1)
	}
	log.Infof("snapshot copy ")

	// copy conf
	if err := s.CopyBin(); err != nil {
		os.Exit(1)
	}
	log.Infof("copy conf")

	// auth
	if err := s.Auth(); err != nil {
		os.Exit(1)
	}
	log.Infof("auth file accessory")

	// start MySQL
	if err := s.StartMySQL(); err != nil {
		os.Exit(1)
	}

	// New local MySQL connection POOl
	i, err := bpct.NewInstance(*user, *passwd, 3358)
	if err != nil {
		os.Exit(1)
	}
	defer i.Close()

	// MySQL check
	if err := i.Check(); err != nil {
		os.Exit(1)
	}

	// newly context
	ctx, cancel := context.WithCancel(context.Background())

	// init wait group
	size := len(tbs)
	wg := &sync.WaitGroup{}
	wg.Add(size)

	// init error channels
	errs := make(chan error, 64)
	defer close(errs)

	var trs []*res.TableRecover
	for _, tb := range tbs {
		tr, err := res.NewTable(tb, c.GetClusterPath(), t, ctx, o, i, wg, errs)
		if err != nil {
			// error occur then exit
			os.Exit(1)
		}

		go tr.Recover()

		trs = append(trs, tr)

		log.Infof("start table {%s} recover ", tr.ID())
	}

	go func() {
		select {
		case err := <-errs:
			log.Errorf("get error from routine{%v}", err)
			cancel()
		}
	}()

	log.Infof("wait for all to finish")
	wg.Wait()

	// if just recover then here to return
	switch *rt {
	case "recover":
		log.Infof("recover for {%d} to timestamp{%s} success", *clusterID, *time)
		return
	}

	log.Infof("flush data on MySQL")
	// flush tables with read lock; flush logs;
	if err := i.Flush(); err != nil {
		log.Errorf("flush MySQL data for cluster id{%d} error {%v}", *clusterID, err)
		os.Exit(1)
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
	for _, t := range trs {
		if err := og.Update(t.ExecutedGTIDSet()); err != nil {
			log.Errorf("merge gtid {%s} into original gtid set{%s} error{%v}", t.ExecutedGTIDSet(), o.ExedGtid, err)
			os.Exit(1)
		}
	}

	o.ExedGtid = og.String()
	o.Time = uint32(t)
	o.CID = *clusterID

	if err := s.FlushOffset(o); err != nil {
		log.Errorf("flush offset{%v} to snapshot{%s} index file error{%v}", o, s.ID(), err)
		os.Exit(1)
	}
}
