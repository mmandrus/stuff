package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	input := string(b)

	// expression to find all `mul(dd,dd)` instances and their indices
	expr := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	matches := expr.FindAllString(input, -1)
	matchIndices := expr.FindAllStringIndex(input, -1)

	// expression to find all `d()` and `don't()` indices
	expr = regexp.MustCompile(`don't\(\)|do\(\)`)
	conditionalIndices := expr.FindAllStringIndex(input, -1)

	allResults, conditionalResults := 0, 0
	for i, j := 0, -1; i < len(matches); i++ {
		m := matches[i]
		mIdx := matchIndices[i]
		// advance the conditional cursor until just before this match
		for {
			if j+1 != len(conditionalIndices) && conditionalIndices[j+1][0] < mIdx[0] {
				j++
				continue
			}
			break
		}
		isDo := j == -1 || conditionalIndices[j][1]-conditionalIndices[j][0] == 4 // length of "do()"
		digits := strings.Split(m[len("mul("):len(m)-1], ",")
		d1, _ := strconv.Atoi(digits[0])
		d2, _ := strconv.Atoi(digits[1])
		allResults += d1 * d2
		if isDo {
			conditionalResults += d1 * d2
		}
	}
	fmt.Println(allResults)
	fmt.Println(conditionalResults)
}
