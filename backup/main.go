package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/juju/errors"
	"github.com/zssky/log"

	"github.com/mysql-binlog/common/client"
	"github.com/mysql-binlog/common/db"
	"github.com/mysql-binlog/common/final"
	"github.com/mysql-binlog/common/inter"

	"github.com/mysql-binlog/backup/handler"
)

/**
@author: pengan
作用： 合并binlog 并生成binlog 备份不参与合并数据
*/

var (
	// merge config
	mc *handler.MergeConfig

	// dump information
	dumpHost   = flag.String("dumphost", "127.0.0.1", "dump MySQL 域名")
	dumpPort   = flag.Int("dumpport", 3306, "dump MySQL 端口")
	dumpUser   = flag.String("dumpuser", "root", "dump MySQL 用户名")
	dumpPasswd = flag.String("dumppasswd", "secret", "dump MySQL 密码")

	// clusterID 用于记录是属于那个集群的binlog 用来唯一标识 文件路径
	cid = flag.Int64("cid", 1, "集群id")

	// cfs storage path for binlog data
	cfsPath = flag.String("cfspath", "/export/backup/", "cfs 数据存储目录")

	// etcd url
	etcd = flag.String("etcd", "http://localhost:2379", "etcd 请求地址")

	// log level
	level = flag.String("level", "debug", "日志级别log level {debug/info/warn/error}")
)

// 初始化
func initiate() {
	// print input parameter
	log.Info("dump host ", *dumpHost, ", dump port ", *dumpPort, ", dump user ", *dumpUser, ", dump password *****")

	dump := &db.MetaConf{
		Host:     *dumpHost,
		Port:     *dumpPort,
		Db:       "test",
		User:     *dumpUser,
		Password: *dumpPasswd,
	}

	if has, _ := dump.HasGTID(); !has {
		log.Fatal(errors.New("only support gtid opened MySQL instances"))
	}

	// data storage path clusterID
	sp := fmt.Sprintf("%s%d", inter.StdPath(*cfsPath), *cid)

	// 创建目录
	inter.CreateLocalDir(sp)

	etc, err := client.NewEtcdMeta(*etcd, "v1")
	if err != nil {
		log.Fatal(err)
	}

	// newly offset
	o, err := etc.Read(*cid)
	if err != nil {
		log.Fatal(err)
	}

	// get master status
	off := o
	if o == nil {
		pos, err := dump.MasterStatus()
		if err != nil || pos.TrxGtid == "" {
			log.Fatal(err, " or gtid is empty")
		}
		log.Info("start binlog position ", string(pos.TrxGtid))
		off = pos

		off.CID = *cid
		// save newly get offset to etcd as well
		if err := etc.Save(off); err != nil {
			log.Fatalf("save offset{%v} to etcd error %v", off, err)
		}
	}

	log.Debugf("start binlog gtid{%s}, binlog file{%s}, binlog position{%d}", string(off.ExedGtid), off.BinFile, off.BinPos)

	// init merge config
	mc, err = handler.NewMergeConfig(sp, off, dump)
	if err != nil {
		log.Fatal(err)
	}

	// init after math
	errs := make(chan interface{}, 4)
	ctx, cancel := context.WithCancel(context.Background())
	am := &final.After{
		Errs:   errs,
		Ctx:    ctx,
		Cancel: cancel,
	}

	mc.After = am
}

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

	logger()

	initiate()

	f, _ := os.Create("/tmp/cpu.prof")
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	go func() {
		signal.Notify(c)
		s := <-c
		pprof.StopCPUProfile()

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("退出信号", s)
		os.Exit(0)
	}()

	mc.Start()
}
