package res

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/zssky/log"

	"github.com/mysql-binlog/siddontang/go-mysql/mysql"
	"github.com/mysql-binlog/siddontang/go-mysql/replication"

	"github.com/mysql-binlog/common/inter"
	"github.com/mysql-binlog/common/meta"
	"github.com/mysql-binlog/common/utils"

	"github.com/mysql-binlog/recover/bpct"
)

/***
* recover on table each table have one routine
*
*/
const (
	// 	StmtEndFlag        = 1
	StmtEndFlag = 1

	logSuffix = ".log"
)

type int64s []int64

func (s int64s) Len() int           { return len(s) }
func (s int64s) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64s) Less(i, j int) bool { return s[i] < s[j] }

// TableRecover
type TableRecover struct {
	table string          // table name
	path  string          // table binlog path
	time  int64           // final time for binlog recover
	ctx   context.Context // context
	off   *meta.Offset    // start offset
	mg    mysql.GTIDSet   // merged gtid
	cg    mysql.GTIDSet   // current gtid set
	og    mysql.GTIDSet   // origin gtid offset
	local *bpct.Instance  // MySQL local connection pool
	wg    *sync.WaitGroup // wait group for outside
	errs  chan error      // error channel

	// followings are for event type
	done   bool                      // table finish dump
	parser *replication.BinlogParser // binlog parser
	desc   []byte                    // desc format event
	buffer bytes.Buffer              // byte buffer
}

// NewTable for outside use
func NewTable(table, clusterPath string, time int64, ctx context.Context, o *meta.Offset, l *bpct.Instance, wg *sync.WaitGroup, errs chan error) (*TableRecover, error) {
	if strings.HasSuffix(clusterPath, "/") {
		clusterPath = strings.TrimSuffix(clusterPath, "/")
	}

	og, err := mysql.ParseMysqlGTIDSet(o.ExedGtid)
	if err != nil {
		log.Errorf("parse executed gtid{%s} error{%v}", o.ExedGtid, err)
		return nil, err
	}

	mg, err := mysql.ParseMysqlGTIDSet(o.ExedGtid)
	if err != nil {
		log.Errorf("parse executed gtid{%s} error{%v}", o.ExedGtid, err)
		return nil, err
	}

	return &TableRecover{
		table:  table,
		path:   fmt.Sprintf("%s/%s", clusterPath, table),
		time:   time,
		ctx:    ctx,
		off:    o,
		mg:     mg,
		og:     og,
		local:  l,
		wg:     wg,
		errs:   errs,
		parser: replication.NewBinlogParser(),
	}, nil
}

// ID for routine
func (t *TableRecover) ID() string {
	return fmt.Sprintf("%s/%d", t.path, t.time)
}

// latestTime find the latestTime log file
func latestTime(time int64, path string) (int64, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("read dir %s error {%v}", path, err)
		return 0, err
	}

	mx := time
	// range files
	for _, f := range fs {
		n := f.Name()
		if strings.HasSuffix(n, logSuffix) {
			ts := strings.TrimSuffix(n, logSuffix)
			t, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				log.Errorf("parse int{%s} error {%v}", ts, err)
				return 0, err
			}

			if t > time {
				// continue for timestamp is bigger than
				continue
			}
			// using max timestamp distance
			dist := time - t
			if dist > 0 && mx > dist {
				mx = dist
			}
		}
	}

	// error for max value not have changed then means no snapshot get error
	if mx == time {
		err := fmt.Errorf("no snapshot get from directory{%s}", path)
		log.Error(err)
		return 0, err
	}

	// return absolute file path with no error
	return time - mx, nil
}

// rangeLogs using start, end
func rangeLogs(start, end int64, p string) ([]int64, error) {
	// find all binlog file between start and end
	fs, err := ioutil.ReadDir(p)
	if err != nil {
		log.Errorf("read dir %s error {%v}", p, err)
		return nil, err
	}

	// result
	var rst int64s
	for _, f := range fs {
		n := f.Name()
		if strings.HasPrefix(n, logSuffix) {
			// totally log files

			ts := strings.TrimSuffix(n, logSuffix)
			t, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				log.Errorf("parse int value{%s} failed {%v}", ts, err)
				return nil, err
			}

			if t >= start && t <= end {
				rst = append(rst, t)
			}
		}
	}

	if len(rst) == 0 {
		log.Warnf("no suitable log file found on path{%s} between{%d, %d}", p, start, end)
		return []int64{}, nil
	}

	sort.Sort(rst)

	return rst, nil
}

// selectLogs to apply MySQL binlog
func (t *TableRecover) selectLogs(start, end int64) (int64s, error) {
	ss, err := latestTime(start, t.path)
	if err != nil {
		log.Errorf("get latestTime log file according timestamp{%d} error{%v}", start, err)
		return nil, err
	}

	return rangeLogs(ss, end, t.path)
}

// Recover
func (t *TableRecover) Recover() {
	log.Debugf("Snapshot for table %s", t.table)
	defer t.wg.Done()

	// take selected log files
	lfs, err := t.selectLogs(int64(t.off.Time), t.time)
	if err != nil {
		t.errs <- err
		return
	}

	onEventFunc := func(e *replication.BinlogEvent) error {
		select {
		case <-t.ctx.Done():
			return fmt.Errorf("canceled by context")
		default:
			if int64(e.Header.Timestamp) > t.time {
				t.done = true
				return nil
			}

			switch e.Header.EventType {
			case replication.FORMAT_DESCRIPTION_EVENT:
				// sql executor begin
				if err := t.local.Begin(t.table); err != nil {
					log.Error(err)
					return err
				}

				t.desc = utils.Base64Encode(e.RawData)
				// sql executor
				if err := t.local.Execute(t.table, []byte(fmt.Sprintf("BINLOG '\n%s\n'%s", t.desc, inter.Delimiter))); err != nil {
					log.Error(err)
					return err
				}

				// sql executor commit
				if err := t.local.Commit(t.table); err != nil {
					log.Error(err)
					return err
				}
			case replication.QUERY_EVENT:
				if t.cg != nil && t.og.Contain(t.cg) {
					log.Debugf("current gtid{%s} already executed on snapshot {%s}", t.cg.String(), t.off.ExedGtid)
					return nil
				}

				qe := e.Event.(*replication.QueryEvent)
				switch strings.ToUpper(string(qe.Query)) {
				case "BEGIN":
					if err := t.local.Begin(t.table); err != nil {
						log.Error(err)
						return err
					}
				case "COMMIT":
					if err := t.local.Commit(t.table); err != nil {
						log.Error(err)
						return err
					}
				case "ROLLBACK":
				case "SAVEPOINT":
				default:
					if err := t.local.Begin(t.table); err != nil {
						log.Error(err)
						return err
					}

					// first use db
					if e.Header.Flags&inter.LogEventSuppressUseF == 0 && qe.Schema != nil && len(qe.Schema) != 0 {
						use := fmt.Sprintf("use %s", qe.Schema)
						if err := t.local.Execute(t.table, []byte(use)); err != nil {
							log.Error(err)
							return err
						}
					}

					// then execute ddl
					if err := t.local.Execute(t.table, qe.Query); err != nil {
						log.Error(err)
						return err
					}
					if err := t.local.Commit(t.table); err != nil {
						log.Error(err)
						return err
					}
				}
			case replication.XID_EVENT:
				if t.cg != nil && t.og.Contain(t.cg) {
					log.Debugf("current gtid{%s} already executed on snapshot {%s}", t.cg.String(), t.off.ExedGtid)
					return nil
				}
				if err := t.local.Commit(t.table); err != nil {
					log.Error(err)
					return err
				}
			case replication.TABLE_MAP_EVENT:
				if t.cg != nil && t.og.Contain(t.cg) {
					log.Debugf("current gtid{%s} already executed on snapshot {%s}", t.cg.String(), t.off.ExedGtid)
					return nil
				}
				// write to buffer first
				t.buffer.WriteString(fmt.Sprintf("BINLOG '\n%s", utils.Base64Encode(e.RawData)))
			case replication.WRITE_ROWS_EVENTv0,
				replication.WRITE_ROWS_EVENTv1,
				replication.WRITE_ROWS_EVENTv2,
				replication.DELETE_ROWS_EVENTv0,
				replication.DELETE_ROWS_EVENTv1,
				replication.DELETE_ROWS_EVENTv2,
				replication.UPDATE_ROWS_EVENTv0,
				replication.UPDATE_ROWS_EVENTv1,
				replication.UPDATE_ROWS_EVENTv2:
				// check current gtid
				if t.cg != nil && t.og.Contain(t.cg) {
					log.Debugf("current gtid{%s} already executed on snapshot {%s}", t.cg.String(), t.off.ExedGtid)
					return nil
				}
				t.buffer.WriteString(fmt.Sprintf("\n%s", utils.Base64Encode(e.RawData)))

				if e.Event.(*replication.RowsEvent).Flags == StmtEndFlag {
					t.buffer.WriteString(fmt.Sprintf("\n'%s", inter.Delimiter))

					// execute
					if err := t.local.Execute(t.table, t.buffer.Bytes()); err != nil {
						log.Error(err)
						return err
					}

					// reset buffer for reuse
					t.buffer.Reset()
				}
			case replication.GTID_EVENT:
				g := e.Event.(*replication.GTIDEvent)

				s := fmt.Sprintf("%s:%d", string(g.SID), g.GNO)
				c, err := mysql.ParseMysqlGTIDSet(s)
				if err != nil {
					log.Errorf("parse og{%s} error {%v}", s, err)
					return err
				}

				t.cg = c

				if err := t.mg.Update(s); err != nil {
					log.Errorf("merge gtid{%s} into gtid{%s} error{%v}", s, t.mg.String(), err)
					return err
				}
			}
			return nil
		}
	}

	// range log files
	for _, l := range lfs {
		// parse each local binlog files

		lf := fmt.Sprintf("%s/%d", t.path, l)
		log.Debugf("parse binlog file{%s}", lf)

		if err := t.parser.ParseFile(lf, 4, onEventFunc); err != nil {
			// write check to io writer
			log.Error("parse file error ", lf, ", error ", err)
			t.errs <- err
			return
		}

		if t.done {
			log.Infof("table {%s} binlog finish apply", t.table)
			return
		}
	}
}