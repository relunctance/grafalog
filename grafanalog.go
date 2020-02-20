package grafanalog

import (
	"fmt"
	"os"

	"github.com/relunctance/goutils/fc"
)

type GrafanaLog struct {
	Path string

	// 每次多去多少条 , -1 是所有
	ReadSize int
}

func NewGrafanaLog(path string) *GrafanaLog {
	g := &GrafanaLog{
		ReadSize: 10,
		Path:     path,
	}
	return g
}

// 数据库地址
type DBer interface {
	Conn() (Conner, error)
}

// 数据存储单元
type Dataer interface {
	Item() []byte
	Items() [][]byte
}

/*
// 支持实现推送的
type Conner interface {
	Push(Dataer) error
}
*/

// 解析器
type Formater interface {
	Parse([]byte) Dataer
}

func ParseData(line []byte, fm Formater) Dataer {
	return fm.Parse(line)
}

func (g *GrafanaLog) Run(d DBer) error {
	_, err := g.openfile(g.Path)
	if err != nil {
		return err
	}
	return nil
	/*
		db, err := d.Conn()

		//if err := db.Push(); err != nil {
		//return err
		//}
		return nil
	*/

}

func (g *GrafanaLog) openfile(path string) (*os.File, error) {
	if !fc.IsExist(path) {
		return nil, fmt.Errorf("not exists file : [%s]", path)
	}
	return os.Open(path)
}
