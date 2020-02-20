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
package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/relunctance/goutils/fc"
)

type Morm struct {
	conf map[string]string
	DB   *gorm.DB
}

func NewMorm(conf map[string]string) (*Morm, error) {
	m := &Morm{}
	err := m.setConfig(conf)
	if err != nil {
		return m, err
	}
	return m, err
}

// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (m *Morm) buildDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", m.conf["username"], m.conf["password"], m.conf["host"], m.conf["port"], m.conf["dbname"], m.conf["charset"])
}

func (m *Morm) setConfig(conf map[string]string) error {
	shouldBeKeys := []string{
		"username",
		"password",
		"host",
		"port",
		"dbname",
		"charset",
	}
	for key, _ := range conf {
		if !fc.InStringArray(key, shouldBeKeys) {
			return fmt.Errorf("key:[%s] unknow , should be exist: %v", key, shouldBeKeys)
		}
	}
	m.conf = conf
	return nil
}

// 单独开启连接查询 , 外面请记得关闭
func (m *Morm) NewDb() (*gorm.DB, error) {
	return m.createDb()
}

func (m *Morm) createDb() (*gorm.DB, error) {
	return gorm.Open("mysql", m.buildDsn())
}

func (m *Morm) Db() *gorm.DB {
	if m == nil {
		return nil
	}
	if m.DB == nil {
		m.DB, _ = m.createDb()
	}
	return m.DB
}

func (m *Morm) Close() error {
	if m == nil || m.DB == nil {
		return nil
	}
	return m.DB.Close()
}
