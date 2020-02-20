package grafalog

// what database you want to push
type DBer interface {
	Push([]Dataer) error
}

//
type Dataer interface {
	Item() []byte
}

// 解析器
type Formater interface {
	Parse([]byte) (Dataer, error)
}
