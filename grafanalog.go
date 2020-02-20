package grafanalog

import (
	"bufio"
	"io"
)

type GrafanaLog struct {
	f  io.Reader
	fm Formater
	db DBer
	// 每次多去多少条 , -1 是所有
	ReadSize int
}

func NewGrafanaLog(f io.Reader) *GrafanaLog {
	g := &GrafanaLog{
		ReadSize: 10,
		fm:       &DefaultFormat{},
		db:       &DefaultDb{},
	}
	return g
}

func (g *GrafanaLog) SetReadSize(v int) {
	g.ReadSize = v
}

func (g *GrafanaLog) RegisterFormater(fm Formater) {
	g.fm = fm
}

func (g *GrafanaLog) RegisterDBer(db DBer) {
	g.db = db
}

// 数据库地址
type DBer interface {
	Push(Dataer) error
	PushMulti([]Dataer) error
}

// 数据存储单元
type Dataer interface {
	Item() []byte
}

// 解析器
type Formater interface {
	Parse([]byte) (Dataer, error)
}

func (g *GrafanaLog) Run() error {
	scan := bufio.NewScanner(g.f)
	items := make([]Dataer, 0, g.ReadSize)
	for scan.Scan() {
		data, err := g.fm.Parse(scan.Bytes())
		if err != nil { // 扫描日志失败
			return err
		}
		if len(items) == g.ReadSize {
			items = make([]Dataer, 0, g.ReadSize)
		}
		items = append(items, data)
		err = g.db.PushMulti(items)
		if err != nil {
			return err
		}
	}
	return nil
}
