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
	"bytes"
	"fmt"

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

// buildInsertSql 创建批量插入语句
// isUpdate =true:
//		当包含的列存在唯一索引的情况下会同时进行更新,不会出现重复的数据
//      如果没有唯一索引的情况下则批量写入
// isUpdate =false:
//      默认情况下批量写入
func (m *Model) insertSqlBuild(item grafalog.Dataer, isUpdate bool) (*bytes.Buffer, string) {
	s := m.m.DB.NewScope(item)
	fields := s.GetStructFields()
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(s.QuotedTableName())
	buf.WriteString("(")
	var duplicateSql string = " ON DUPLICATE KEY UPDATE "
	for i, field := range fields {
		fname := field.Tag.Get("gorm")
		buf.WriteString("`" + fname + "`")
		if isUpdate {
			duplicateSql += "`" + fname + "` = VALUES(`" + fname + "`)"
		}
		if i != len(fields)-1 {
			if isUpdate {
				duplicateSql += ","
			}
			buf.WriteString(",")
		}
	}
	buf.WriteString(") VALUES")
	return buf, duplicateSql
}

func (m *Model) buildInsertSql(vals []grafalog.Dataer, isUpdate bool) *bytes.Buffer {
	if len(vals) == 0 {
		return nil
	}
	buf, duplicateSql := m.insertSqlBuild(vals[0], isUpdate)
	for j, val := range vals {
		buf.WriteString("(")
		s := m.m.DB.NewScope(val)
		fs := s.Fields()
		for i, field := range fs {
			v := m.Addslashes(field.Field.Interface())
			buf.WriteString("'" + v + "'")
			if i != len(fs)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")
		if j != len(vals)-1 {
			buf.WriteString(",")
		}
	}
	if isUpdate {
		buf.WriteString(duplicateSql)
	}
	return buf
}

func (m *Model) Addslashes(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// implement grafalog.DBer
func (m *Model) Push(vals []grafalog.Dataer) error {
	return m.BatchInsert(vals, true)
}

func (m *Model) BatchInsert(vals []grafalog.Dataer, isUpdate bool) error {
	if len(vals) == 0 {
		return nil
	}
	return m.m.DB.Exec(m.buildInsertSql(vals, isUpdate).String()).Error
}
