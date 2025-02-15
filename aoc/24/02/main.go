package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	contents := string(b)
	reports := make([][]int, 1000)
	for i, r := range strings.Split(contents, "\n") {
		if len(r) == 0 {
			continue
		}
		lvls := strings.Split(r, " ")
		rep := make([]int, len(lvls))
		for i, lvl := range lvls {
			ilvl, err := strconv.Atoi(lvl) //not a WoW reference
			if err != nil {
				panic(err)
			}
			rep[i] = ilvl
		}
		reports[i] = rep
	}

	// returns whether the report is safe
	checkReport := func(report []int) bool {
		priorLevel, priorDiff := report[0], 0
		for i := 1; i < len(report); i++ {
			diff := report[i] - priorLevel
			// valid change: ascension from 1 to 3 levels if it wasn't descending previously
			if diff >= 1 && diff <= 3 && priorDiff >= 0 {
				priorLevel, priorDiff = report[i], diff
				continue
			}
			// valid change: descension from 1 to 3 levels if it wasn't descending previously
			if diff <= -1 && diff >= -3 && priorDiff <= 0 {
				priorLevel, priorDiff = report[i], diff
				continue
			}
			// invalid change: unsafe report
			return false
		}
		return true
	}

	numSafe := 0
	numSafeDampened := 0
	for _, report := range reports {
		if checkReport(report) {
			numSafe++
		} else {
			for li := range report {
				if checkReport(append(append([]int{}, report[0:li]...), report[li+1:]...)) {
					numSafeDampened++
					break
				}
			}
		}
	}

	fmt.Println(numSafe)
	fmt.Println((numSafe + numSafeDampened))
}
