package inter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// FileType 定义的文件类型 {day, hour, min, sec}
type FileType string

// FileSuffix 文件后缀类型
type FileSuffix string

// CharStd 字符标准化 用作格式化表名 统一处理
var CharStd = strings.ToLower

// StdPath 标准化路径
var PathStd = StdPath

// DbPath 数据文件路径
var DataPath = "data/"

// SnapPath 快照目录
var SnapPath = "snap/"

// StorePath 存储文件路径
var StorePath = "store/"

// LogPath 日志文件路径
var LogPath = "logs/"

//// Delimiter binlog 分隔符
//var Delimiter = "/*!*/;\n"

// Delimiter binlog 分隔符
var Delimiter = ""

// DaySeconds
var DaySeconds = uint32(24 * 3600)

// HourSeconds
var HourSeconds = uint32(3600)

// MinSeconds
var MinSeconds = uint32(60)

// Seconds
var Second = uint32(1)

// FileMode
var FileMode = os.FileMode(0666)

type FileName struct {
	Path   string     `json:"path"`   // Path 文件路径
	Name   string     `json:"name"`   // Name 文件标准名称		example : sec_1551661307_1551661308.log
	Prefix FileType   `json:"prefix"` // Prefix file type {sec, min, hour, day}
	Suffix FileSuffix `json:"suffix"` // suffix 后缀名
}

type FileNames []*FileName

const (
	// LogEventSuppressUseF use db flag
	LogEventSuppressUseF = uint16(0x8)

	// DAY 日期字符串常量
	DAY FileType = "day"

	// HOUR 小时字符串常量
	HOUR FileType = "hour"

	// MINUTE 分钟字符串常量
	MINUTE FileType = "min"

	// SECOND 秒 字符串常量
	SECOND FileType = "sec"

	// SPLIT 分隔符
	SPLIT = "_"

	// LogSuffix 结束符
	LogSuffix FileSuffix = ".log"

	// TarSuffix
	TarSuffix FileSuffix = ".tar"

	// StatusSuccess
	StatusSuccess = "success"

	// StatusFailure
	StatusFailure = "failure"

	// public folder means: that all data should use the public ddl when restore data from storage
	Public = "public"
)

// AbsolutePath 绝对路径
type AbsolutePath struct {
	TmpSrc   string   // TmpSrc 临时source 根路径
	TmpDst   string   // TmpDst 临时dest 根路径
	FileType FileType // file type 文件类型 {day, hour, minute, second}
	Host     string   // Host mysql 域名/ip
	Table    string   // schema.table 表名全称
	Start    int64    // Start 起始时间
	End      int64    // End 终止时间
}

// SourcePath 文件源路径
func (p *AbsolutePath) TmpSourcePath() string {
	if strings.HasSuffix(p.TmpSrc, "/") {
		p.TmpSrc = strings.TrimSuffix(p.TmpSrc, "/")
	}

	// dst + host + "/" + inter.SecondTable(table)
	return fmt.Sprintf("%s/%s/%s%s%s", p.TmpSrc, p.Host, p.Table, SPLIT, p.FileType)
}

// DstPath 目标文件路径
func (p *AbsolutePath) TmpDstPath() string {
	CreateLocalDir(fmt.Sprintf("%s/%s", p.TmpDst, p.Host))

	if strings.HasSuffix(p.TmpDst, "/") {
		p.TmpDst = strings.TrimSuffix(p.TmpDst, "/")
	}

	// dst + host + "/" + inter.SecondTable(table)
	return fmt.Sprintf("%s/%s/%s%s%s", p.TmpDst, p.Host, p.Table, SPLIT, p.FileType)
}

// DstLogFileName 目标文件名称 不包括路径
func (p *AbsolutePath) DstLogFileName() string {
	return fmt.Sprintf("%s%s%010d%s%010d%s", p.FileType, SPLIT, p.Start, SPLIT, p.End, LogSuffix)
}

func (p *AbsolutePath) NextPrefix() string {
	return fmt.Sprintf("%s%s%d", p.FileType, SPLIT, p.End)
}

func (p *AbsolutePath) DstTarFileName() string {
	return fmt.Sprintf("%s%s%010d%s%010d%s", p.FileType, SPLIT, p.Start, SPLIT, p.End, TarSuffix)
}

func (l FileNames) Len() int { return len(l) }

func (l FileNames) Less(i, j int) bool {
	ei, _ := strconv.ParseInt(strings.TrimSuffix(strings.Split(l[i].Name, SPLIT)[2], string(l[i].Suffix)), 10, 64)
	ej, _ := strconv.ParseInt(strings.TrimSuffix(strings.Split(l[j].Name, SPLIT)[2], string(l[j].Suffix)), 10, 64)
	return ei < ej
}

func (l FileNames) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (f *FileName) GetRange() (int64, int64) {
	ps := strings.Split(f.Name, SPLIT)
	e, _ := strconv.ParseInt(strings.TrimSuffix(ps[2], string(f.Suffix)), 10, 64)
	s, _ := strconv.ParseInt(ps[1], 10, 64)

	return s, e
}

// CreateLocalDir create directory
func CreateLocalDir(path string) {
	// create dir
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// ParseTime 解析字符串時間
func ParseTime(end string) int64 {
	stopTime, _ := time.Parse("2006-01-02 15:04:05", end)
	return stopTime.Unix()
}

func StdPath(p string) string {
	if !strings.HasSuffix(p, "/") {
		return fmt.Sprintf("%s/", p)
	}
	return p
}
