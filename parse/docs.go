// Package parse is a simple sql lexer/parser with limited functionality
// it can handle queries like :
// 	select fields from table where "field" = 'string' and (another_field=100 or boolean_field) order by field_1 desc , field_2 asc limit 10, 100
// the result is some sort of abstract source tree :)
// the package is based on net/html and rob pike talk about writing lexer in go.
package parse
