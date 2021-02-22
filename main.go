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
	Author    string
}

type Logic struct {
	bars *collections.OrderedDict
	min  int
	max  int
}

var block = "\u2580"
var count = 4

func normalize(x, xMin, xMax int) float32 {
	return float32(x-xMin) / float32(xMax-xMin)
}

func getScore(items Logic) {
	maxDigits := (items.max / 10) + 1
	for value := range items.bars.Iterate() {
		switch val := value.(type) {
		case Bars:
			{
				value := int(normalize(val.Commits, items.min, items.max) * float32(count))
				fmt.Print(val.Timestamp)
				n := 0
				fmt.Print(" ")
				commits := val.Commits
				fmt.Print(commits)
				spaces := maxDigits - ((commits / 10) + 1) + 2
				for i := 0; i < spaces; i += 1 {
					fmt.Print(" ")
				}
				for n <= value {
					n += 1
					fmt.Print(strings.Repeat(block, count))
				}
				fmt.Println()
			}
		}

	}
}

func getCommitLog(after, before string, author string) (Logic, error) {
	args := []string{"log", "--pretty=format:%ai|%ae"}
	if after != "" {
		args = append(args, "--after="+after)
	}
	if before != "" {
		args = append(args, "--before="+before)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return Logic{}, err
	}

	bars := collections.New()
	logicStruct := Logic{bars: bars}
	res := strings.Split(string(output), "\n")
	if len(res) == 0 {
		fmt.Println("No commits to plot")
		return logicStruct, nil
	}
	min := 0
	max := 0
	for _, val := range res {
		splitted := strings.Split(val, "|")
		authorName := splitted[1]
		if !strings.Contains(authorName, author) {
			continue
		}
		c := splitted[0][:10]
		value := bars.Get(c)
		commitsForTs := 0
		if value != nil {
			switch val := value.(type) {
			case Bars:
				commitsForTs = val.Commits + 1
			}
		} else {
			commitsForTs = 1
		}
		bars.Set(c, Bars{Timestamp: c, Commits: commitsForTs, Author: authorName})
		if commitsForTs < min {
			min = commitsForTs
		}
		if max < commitsForTs {
			max = commitsForTs
		}
	}
	logicStruct = Logic{bars: bars, min: min, max: max}
	return logicStruct, nil
}

func main() {
	after := flag.String("a", "", "after date (yyyy-mm-dd hh:mm)")
	before := flag.String("b", "", "before date (yyyy-mm-dd hh:mm)")
	author := flag.String("u", "", "author order")
	flag.Parse()

	logicStruct, err := getCommitLog(*after, *before, *author)
	if err != nil {
		fmt.Println("Issue happened")
		os.Exit(0)
	}
	if logicStruct.bars.Length() == 0 {
		fmt.Println("No commits to plot")
	}
	getScore(logicStruct)
}
