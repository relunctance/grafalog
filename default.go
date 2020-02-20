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
