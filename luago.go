package main

import (
	"io/ioutil"
	"luago/binchunk"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		proto := binchunk.Undump(data)
		list(proto)
	}
}
