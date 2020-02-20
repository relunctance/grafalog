package mysql

import (
	"github.com/relunctance/grafalog"
)

// the example you can see:  example/mysql/main.go
type Model struct {
	m *Morm
}

func NewModel(dbconfig map[string]string) (*Model, error) {
	m, err := NewMorm(dbconfig)
	m.DB, err = m.NewDb()
	return &Model{m: m}, err
}

func (m *Model) Close() {
	m.m.Close()
}

func (m *Model) Push(vals []grafalog.Dataer) error {
	// TODO 支持批量写入更好
	for _, val := range vals {
		m.m.DB.Create(val) //  利用gorm写入的方法, 直接写入数据库
	}
	return nil
}
