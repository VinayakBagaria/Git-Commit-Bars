package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Bars struct {
	Timestamp string
	Commits   int
}

func printData(values map[string]int) {
	for key, value := range values {
		block := "\u2580"
		n := 0
		fmt.Print(key)
		fmt.Print(" ")
		for n <= value {
			n += 1
			fmt.Print(block)
		}
		fmt.Println()
	}
}

func normalize(x, xmin, xmax int) int {
	return int(float32(x-xmin) / float32(xmax-xmin))
}

func getScore(items map[string]Bars) map[string]int {
	intVals := []int{}
	for _, value := range items {
		intVals = append(intVals, value.Commits)
	}
	xMax := intVals[0]
	xMin := intVals[0]

	for _, val := range intVals {
		if xMax < val {
			xMax = val
		}
		if xMin > val {
			xMax = val
		}
	}

	dateWiseLength := make(map[string]int)
	for key, value := range items {
		dateWiseLength[key] = normalize(value.Commits, xMin, xMax)
	}
	return dateWiseLength
}

func getCommitLog(after, before string) ([]map[string]string, error) {
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

func filterAsPerPeriodicity(items []map[string]string, periodicity string) map[string]Bars {
	bars := make(map[string]Bars)

	for i := 0; i < len(items); i++ {
		timestamp := items[i]["timestamp"][:10]
		bar := Bars{Timestamp: timestamp, Commits: 0}
		if _, found := bars[timestamp]; !found {
			bar.Commits += 1
		}
		bars[timestamp] = bar
	}

	return bars
}

func main() {
	periodicity := flag.String("p", "month", "peridocity definition to day, week, month, year")
	after := flag.String("a", "", "after date (yyyy-mm-dd hh:mm)")
	before := flag.String("b", "", "before date (yyyy-mm-dd hh:mm)")
	flag.Parse()

	items, err := getCommitLog(*after, *before)
	if err != nil {
		fmt.Println("Issue happened")
		os.Exit(0)
	}

	if len(items) > 0 {
		fmt.Printf("%d commits\n", len(items))
		data := filterAsPerPeriodicity(items, *periodicity)
		printData(getScore(data))
	} else {
		fmt.Println("No commits to plot")
	}
}
