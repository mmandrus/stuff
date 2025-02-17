package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	b, _ := os.ReadFile("rules.txt")
	rules := strings.Split(string(b), "\n")
	b, _ = os.ReadFile("updates.txt")
	updates := strings.Split(string(b), "\n")

	// part 1
	// map of int to the things that come after it
	pagesProcessed := make(map[string][]string)
	// some dp to make things more performance
	validTransitions := make(map[string]bool)

	// part 2
	// knowing what comes *after* what is more helpful as it turns out
	afterMap := make(map[string]map[string]bool)

	// figure out what comes before what
	for i := 0; i < len(rules); i++ {
		r := rules[i]
		if len(r) == 0 {
			continue
		}
		p1, p2 := r[:2], r[3:]

		// part 1
		if l, ok := pagesProcessed[p1]; ok {
			pagesProcessed[p1] = append(l, p2)
		} else {
			pagesProcessed[p1] = []string{p2}
		}

		// part 2
		after, ok := afterMap[p2]
		if !ok {
			after = make(map[string]bool, 0)
		}
		after[p1] = true
		afterMap[p2] = after
	}

	var check func([]string) bool

	// recursively move through each update and check status along the way
	check = func(s []string) bool {
		if len(s) == 1 {
			// everything up until now has been valid, so it doesn't matter what this number is
			return true
		}
		p1, p2 := s[0], s[1]
		if validTransitions[p1+p2] {
			// we've already seen that this is a valid transition so continue on
			return check(s[1:])
		}

		// figure out if this is a valid transition and save the result
		isValid := true
		if entry, ok := pagesProcessed[p2]; ok && slices.Contains(entry, p1) {
			isValid = false
		}
		validTransitions[p1+p2] = isValid
		if !isValid {
			return false
		}

		return check(s[1:])
	}

	sum := 0
	incorrect := make([][]string, 0)
	for i := 0; i < len(updates); i++ {
		u := strings.TrimSpace(updates[i])
		if len(u) == 0 {
			continue
		}
		pages := strings.Split(u, ",")
		if ok := check(pages); ok {
			v, _ := strconv.Atoi(pages[len(pages)/2])
			sum += v
		} else {
			incorrect = append(incorrect, pages)
		}
	}
	fmt.Println(sum)

	p2Sum := 0
	for _, page := range incorrect {
		correctOrder := make([]string, len(page))
		for i, p := range page {
			after := afterMap[p]
			idx := 0
			// calculate where this should go in the list by checking how many items in the list should be before it
			for j := 0; j < len(page); j++ {
				if i != j && after[page[j]] {
					idx++
				}
			}
			correctOrder[idx] = p
		}
		v, _ := strconv.Atoi(correctOrder[len(correctOrder)/2])
		p2Sum += v
	}
	fmt.Println(p2Sum)
}
