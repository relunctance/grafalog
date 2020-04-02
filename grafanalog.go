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
	"sync"
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

	debug bool
	mx    sync.Mutex
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

func (g *GrafanaLog) SetDebug(v bool) {
	g.debug = v
}

func (g *GrafanaLog) logPrintf(format string, v ...interface{}) {
	if g.debug {
		log.Printf(format, v...)
	}
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
	// 每次满ChunkSize 或者 满flushTick 都会开始写数据
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

// 按照大小刷
func (g *GrafanaLog) flushWithFullSize() error {
	if len(g.flushdata) >= g.ChunkSize {
		return g.pushItems("from=sizeflush")
	}
	return nil
}

// 按照时间刷
func (g *GrafanaLog) flushWithFinishTime() error {
	err := g.pushItems("from=timetickflush")
	if err != nil {
		log.Fatalf("time ticker push data faild err: %v", err)
		return err
	}
	g.logPrintf("wait data , timetick[%d] second\n", g.flushTick)
	return nil
}

func (g *GrafanaLog) copyData() []Dataer {
	tmpdata := make([]Dataer, len(g.flushdata))
	copy(tmpdata, g.flushdata)
	// 写入失败情况下会丢弃掉对应的数据, 再次其他日志
	g.flushdata = make([]Dataer, 0, g.ChunkSize) // 重置
	return tmpdata
}

func (g *GrafanaLog) pushItems(fr string) (err error) {
	if len(g.flushdata) == 0 {
		return
	}
	// 后台去慢慢写入 , 不阻塞下次任务
	go func(fr string) {
		data := g.copyData()
		g.logPrintf("from:[%s] , length:%d \n", fr, len(data))
		err = g.db.Push(data) // 有可能写入数据库超时, 网络超时等异常情况
		if err == nil {
			g.logPrintf("success push data size: [%d]\n", len(data))
		}
	}(fr)

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
	g.mx.Lock()
	defer g.mx.Unlock()
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
