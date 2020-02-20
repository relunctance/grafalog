package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/relunctance/grafalog"
	"github.com/relunctance/grafalog/driver/mysql"
)

var _ grafalog.Dataer = new(Logtest)
var _ grafalog.Formater = new(Formate)

type Formate struct {
}

func (f *Formate) Parse(line []byte) (grafalog.Dataer, error) {
	vals := bytes.Split(line, []byte("\t"))
	d := &Logtest{}
	d.CostTime, _ = strconv.Atoi(string(vals[0]))
	d.Ukey = string(vals[1])
	d.Api = string(vals[2])
	d.Value = string(vals[3])
	d.Ctime = string(vals[4])
	d.Msg = string(vals[5])
	return d, nil
}

// 注意这里表名: logtests , gorm会自动拼接上一个s
type Logtest struct {
	CostTime int    `gorm:"cost_time"`
	Ukey     string `gorm:"ukey"`
	Api      string `gorm:"api"`
	Value    string `gorm:"value"`
	Ctime    string `gorm:"ctime"`
	Msg      string `gorm:"msg"`
}

func (d *Logtest) Item() []byte {
	s := fmt.Sprintf("ukey:%s, cost_time:%d, api:%s, value:%s, msg:%s, ctime:%s",
		d.Ukey,
		d.CostTime,
		d.Api,
		d.Value,
		d.Msg,
		d.Ctime,
	)
	return []byte(s)

}

func main() {
	f, err := os.Open("test.logs")
	if err != nil {
		panic("open test.logs is faild")
	}
	defer f.Close()
	m, err := mysql.NewModel(map[string]string{
		"username": "root",
		"password": "123456QWER",
		"host":     "127.0.0.1",
		"port":     "3306",
		"dbname":   "test",
		"charset":  "utf8",
	})
	defer m.Close()
	if err != nil {
		panic(err)
	}
	g := grafalog.New(f)
	g.SetReadSize(20) // set push size
	g.RegisterDBer(m)
	g.RegisterFormater(&Formate{})
	err = g.Run()
	if err != nil {
		panic(err)
	}
}
