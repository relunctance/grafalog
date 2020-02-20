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

import "fmt"

type DefaultDb struct{}

func (db *DefaultDb) outputStdout(val Dataer) error {
	if val == nil {
		return fmt.Errorf("val:Dataer is nil")
	}
	fmt.Println(string(val.Item()))
	return nil

}

func (db *DefaultDb) Push(vals []Dataer) error {
	for _, val := range vals {
		err := db.outputStdout(val)
		if err != nil {
			return err
		}
	}
	return nil
}

type DefaultFormat struct{}

func (f *DefaultFormat) Parse(val []byte) (Dataer, error) {
	return &DefaultData{data: val}, nil
}

type DefaultData struct {
	data []byte
}

func (d *DefaultData) Item() []byte {
	return d.data
}
