package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getLog(after, before string) ([]map[string]string, error) {
	args := []string{"log", "--pretty=format:%ai|%ae"}
	if after != "" {
		args = append(args, "--after="+after)
	}
	if before != "" {
		args = append(args, "--before="+before)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()

	var items []map[string]string
	if err != nil {
		return items, err
	}
	res := strings.Split(string(output), "\n")
	for _, val := range res {
		c := strings.Split(val, "|")
		m := make(map[string]string)
		m["timestamp"] = c[0]
		m["author"] = c[1]
		items = append(items, m)
	}
	return items, nil
}

func main() {
	//periodicity := flag.String("p", "month", "peridocity definition to day, week, month, year")
	after := flag.String("a", "", "after date (yyyy-mm-dd hh:mm)")
	before := flag.String("b", "", "before date (yyyy-mm-dd hh:mm)")

	flag.Parse()

	items, err := getLog(*after, *before)

	if err != nil {
		fmt.Println("Issue happened")
		os.Exit(0)
	}

	fmt.Println(items)
}
