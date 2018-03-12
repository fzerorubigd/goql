package main

import "github.com/fzerorubigd/goql/internal/parse"

func main() {
	parse.Test("select * from xx where id in (1,3434343,4.3) AND test = 'sss'")
}
