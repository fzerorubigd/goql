package main

import "github.com/fzerorubigd/goql/internal/parse"

func main() {
	parse.Test("select * from xx where id in (1,2,3,4) AND test = 'sss'")
}
