package main

import (
	"flag"
	"fmt"
)

func main() {
	periodicity := flag.String("p", "month", "peridocity definition to day, week, month, year")
	after := flag.String("a", "", "after date (yyyy-mm-dd hh:mm)")
	before := flag.String("b", "", "before date (yyyy-mm-dd hh:mm)")

	flag.Parse()

	fmt.Println(*periodicity, *after, *before)
}
