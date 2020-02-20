package main

import (
	"os"

	"github.com/relunctance/grafalog"
)

func main() {
	f, err := os.Open("../mysql/test.logs")
	if err != nil {
		panic("open test.logs is faild")
	}
	defer f.Close()
	g := grafalog.New(f)
	err = g.Run() // default output os.Stdout
	if err != nil {
		panic(err)
	}
}
