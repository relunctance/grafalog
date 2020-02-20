// Copyright (c) 2020 Gao.QiLin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package grafalog

import (
	"bufio"
	"fmt"
	"io"
)

type GrafanaLog struct {
	f  io.Reader
	fm Formater
	db DBer
	// 每次多取多少条
	ReadSize int
}

func New(f io.Reader) *GrafanaLog {
	return newGrafanaLog(f)
}
func newGrafanaLog(f io.Reader) *GrafanaLog {
	g := &GrafanaLog{
		ReadSize: 10,
		f:        f,
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

func (g *GrafanaLog) check() error {
	if g.ReadSize == 0 {
		return fmt.Errorf("g.ReadSize should be > 0")
	}
	if g.f == nil {
		return fmt.Errorf("io.Reader is nil")
	}
	if g.db == nil {
		return fmt.Errorf("g.db is nil , please run grafalog.RegisterDBer(DBer)")
	}
	if g.fm == nil {
		return fmt.Errorf("g.fm is nil , you should run grafalog.RegisterFormater(Formater)")
	}
	return nil
}

func (g *GrafanaLog) Run() error {
	if err := g.check(); err != nil {
		return err
	}
	scan := bufio.NewScanner(g.f)
	items := make([]Dataer, 0, g.ReadSize)
	for scan.Scan() {
		data, err := g.fm.Parse(scan.Bytes())
		if err != nil { // 如果解析日志失败会中断执行, 请保证日志格式统一
			return err
		}
		if len(items) == g.ReadSize {
			err = g.db.Push(items) // 如果推送数据失败也会中断执行
			if err != nil {
				return err
			}
			items = make([]Dataer, 0, g.ReadSize)
		}
		items = append(items, data)
	}
	if len(items) > 0 {
		g.db.Push(items)
	}
	return nil
}
