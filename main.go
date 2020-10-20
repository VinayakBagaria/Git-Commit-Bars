package main

import (
	"flag"
	"fmt"
	"git-bars/collections"
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

func normalize(x, xMin, xMax int) int {
	return int(float32(x-xMin) / float32(xMax-xMin))
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

func getCommitLog(after, before string, reverse string) ([]map[string]string, error) {
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
	if reverse != "" {
		i := 0
		j := len(items)
		for i < j {
			items[i], items[j] = items[j], items[i]
			i += 1
			j -= 1
		}
	}
	return items, nil
}

func filterAsPerPeriodicity(items []map[string]string, periodicity string) *collections.OrderedDict {
	bars := collections.New()
	for i := 0; i < len(items); i++ {
		timestamp := items[i]["timestamp"][:10]
		value := bars.Get(timestamp)
		if value != nil {
			switch val := value.(type) {
			case *Bars:
				bars.Set(timestamp, Bars{Timestamp: val.Timestamp, Commits: 1})
			}
		} else {
			bars.Set(timestamp, Bars{Timestamp: timestamp, Commits: 1})
		}

		//data := bars.Get(timestamp)
		//if data != nil {
		//	bars.Set(timestamp, Bars{Timestamp: timestamp, Commits: 1})
		//} else {
		//	//bars.Set(timestamp, Bars{Timestamp: timestamp, Commits: data.Commits + 1})
		//}
		//if data, found := bars[timestamp]; found {
		//	bars[timestamp] = Bars{Timestamp: data.Timestamp, Commits: data.Commits + 1}
		//} else {
		//	bars[timestamp] = Bars{Timestamp: timestamp, Commits: 0}
		//}
	}
	return bars
}

func main() {
	periodicity := flag.String("p", "month", "periodicity definition to day, week, month, year")
	after := flag.String("a", "", "after date (yyyy-mm-dd hh:mm)")
	before := flag.String("b", "", "before date (yyyy-mm-dd hh:mm)")
	reverse := flag.String("r", "", "reverse date order")
	flag.Parse()

	items, err := getCommitLog(*after, *before, *reverse)
	if err != nil {
		fmt.Println("Issue happened")
		os.Exit(0)
	}

	if len(items) > 0 {
		data := filterAsPerPeriodicity(items, *periodicity)
		fmt.Printf("%d commits over %d days\n", len(items), data.Length())
		//printData(getScore(data))
	} else {
		fmt.Println("No commits to plot")
	}
}
