package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contents := string(b)
	lines := strings.Split(contents, "\n")

	l1, l2 := make([]int, len(lines)), make([]int, len(lines))
	for i, l := range lines {
		if len(l) == 0 {
			continue
		}
		vals := strings.Split(l, "   ")
		v1, err := strconv.Atoi(vals[0])
		if err != nil {
			panic(err)
		}
		v2, err := strconv.Atoi(vals[1])
		if err != nil {
			panic(err)
		}
		l1[i], l2[i] = v1, v2
	}

	// sorting facilitates both parts
	sortTime := time.Now()
	slices.Sort(l1)
	slices.Sort(l2)
	sortDuration := time.Since(sortTime)

	// part 1
	p1Time := time.Now()
	total := 0.0
	for i := 0; i < len(l1); i++ {
		total += math.Abs(float64(l1[i] - l2[i]))
	}
	p1Duration := time.Since(p1Time)
	fmt.Println(int(total))

	// part 2
	p2Time := time.Now()
	total2 := 0
	rIndex := 0
outer:
	for _, v1 := range l1 {
		// find number of occurrences of v1 in l2 and advance cursor to following value
		numOccurrences := 0
	search:
		for {
			v2 := l2[rIndex]
			if v2 > v1 {
				// no occurrences
				break
			}
			if v2 == v1 {
				// continue until we are at a different number
				for {
					numOccurrences++
					rIndex++
					if l2[rIndex] != v1 {
						break search
					}
				}
			}
			if v2 < v1 {
				rIndex++
			}
			if rIndex == len(l1) {
				break outer
			}
		}
		total2 += numOccurrences * v1
	}
	p2Duration := time.Since(p2Time)
	fmt.Println(total2)
	fmt.Printf("sort time: %dμs\npart 1 time: %dμs\npart 2 time: %dμs\n", sortDuration.Microseconds(), p1Duration.Microseconds(), p2Duration.Microseconds())
}
