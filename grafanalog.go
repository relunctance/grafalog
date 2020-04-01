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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/nxadm/tail"
)

const (
	DefaultReadSize = 20

	// 单位: 秒
	FlushTick = 1
)

const ()

type GrafanaLog struct {
	fpath string
	fm    Formater
	db    DBer

	// 单位: 秒 , 每次满多少s开始刷数据
	flushTick int
	// 每次多取多少条, 这个参数会控制刷数据的窗口
	ChunkSize int

	// 临时存储需要刷盘数据
	flushdata []Dataer
}

func New(filepath string) *GrafanaLog {
	return newGrafanaLog(filepath)
}
func newGrafanaLog(f string) *GrafanaLog {
	fmt.Println("newGrafanaLog start")
	g := &GrafanaLog{
		fpath:     f,
		fm:        &DefaultFormat{},
		db:        &DefaultDb{},
		flushTick: FlushTick,
	}
	g.SetReadSize(DefaultReadSize)
	return g
}

// 设置定期刷数据时间, 单位: 秒
func (g *GrafanaLog) SetFlushTick(v int) {
	g.flushTick = v
}

func (g *GrafanaLog) SetReadSize(v int) {
	g.ChunkSize = v
	g.flushdata = make([]Dataer, 0, g.ChunkSize)
}

func (g *GrafanaLog) RegisterFormater(fm Formater) {
	g.fm = fm
}

func (g *GrafanaLog) RegisterDBer(db DBer) {
	g.db = db
}

func (g *GrafanaLog) check() error {
	if g.ChunkSize == 0 {
		return fmt.Errorf("g.ChunkSize should be > 0")
	}
	if g.fpath == "" {
		return fmt.Errorf("file path should not be empty ")
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
	g.startPusher()
	return g.TailLine()
}

func (g *GrafanaLog) TailLine() error {
	t, err := tail.TailFile(g.fpath, tail.Config{Follow: true})
	if err != nil {
		return err
	}
	// 每次满ChunkSize 或者 满3s 都会开始写数据
	for line := range t.Lines { // channel 阻塞兼容日志中新增内容
		data, err := g.fm.Parse([]byte(line.Text))
		if err != nil {
			log.Printf("parse log [%s] ,  err: %v", line.Text, err) // 如果解析日志有问题, 打印操作日志
			continue
		}

		// 推送数据库失败
		if err := g.addData(data); err != nil {
			log.Fatalf("push data faild err: %v", err)
		}
	}
	return nil
}

/*
// 直接一次性读取所有的内容
func (g *GrafanaLog) readAll() (err error) {
	f, err := g.openFile()
	if err != nil {
		return err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	items := make([]Dataer, 0, g.ChunkSize)
	for scan.Scan() {
		data, err := g.fm.Parse(scan.Bytes())
		if err != nil { // 如果解析日志失败会中断执行, 请保证日志格式统一
			return err
		}
		if len(items) == g.ChunkSize {
			err = g.db.Push(items) // 如果推送数据失败也会中断执行
			if err != nil {
				return err
			}
			items = make([]Dataer, 0, g.ChunkSize)
		}
		items = append(items, data)
	}
	if len(items) > 0 {
		g.db.Push(items)
	}
	return nil
}
*/

// 按照大小刷
func (g *GrafanaLog) flushWithFullSize() error {
	if len(g.flushdata) == g.ChunkSize {
		return g.pushItems()
	}
	return nil
}

// 按照时间刷
func (g *GrafanaLog) flushWithFinishTime() error {
	if err := g.pushItems(); err != nil {
		log.Fatalf("time ticker push data faild err: %v", err)
		return err
	}
	return nil
}

func (g *GrafanaLog) pushItems() (err error) {
	if len(g.flushdata) == 0 {
		return
	}
	err = g.db.Push(g.flushdata)
	g.flushdata = make([]Dataer, 0, g.ChunkSize)
	return
}

// 开启一个push协程, 如果不够ChunkSize , 定时也推送
func (g *GrafanaLog) startPusher() {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(g.flushTick) * time.Second):
				g.flushWithFinishTime()
			}
		}
	}()
}

// 每次满ChunkSize 或者 满FlushTime(秒) 都会开始写数据
func (g *GrafanaLog) addData(data Dataer) error {
	g.flushdata = append(g.flushdata, data)
	return g.flushWithFullSize()
}

func (g *GrafanaLog) openFile() (*os.File, error) {
	if !gfile.IsFile(g.fpath) {
		return nil, fmt.Errorf("not exists file: [%s]", g.fpath)
	}
	if !gfile.IsReadable(g.fpath) {
		return nil, fmt.Errorf("[%s] not readable", g.fpath)
	}
	return os.Open(g.fpath)
}
