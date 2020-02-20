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
