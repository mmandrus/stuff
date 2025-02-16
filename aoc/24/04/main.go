package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	input := string(b)
	part1(input)
	part2(input)
}

func part1(input string) {
	rows := strings.Split(input, "\n")
	rows = rows[:len(rows)-1]
	rowLength := len(rows[0])
	numResults := 0

	var searchR func(i, j, di, dj int, s string)

	searchR = func(i, j, di, dj int, s string) {
		if len(s) == 0 {
			numResults++
			return
		}
		i, j = i+di, j+dj
		if i < len(rows) && j > -1 && j < rowLength && rows[i][j] == s[0] {
			searchR(i, j, di, dj, s[1:])
		}
	}

	search := func(i, j int, s string) {
		// look right
		searchR(i, j, 0, 1, s)
		// look down
		searchR(i, j, 1, 0, s)
		// look down right
		searchR(i, j, 1, 1, s)
		// look down left
		searchR(i, j, 1, -1, s)
	}

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for j := 0; j < len(row); j++ {
			char := row[j]
			if char == byte('X') {
				search(i, j, "MAS")
			} else if char == byte('S') {
				search(i, j, "AMX")
			}
		}
	}

	fmt.Println(numResults)
}

func part2(input string) {
	rows := strings.Split(input, "\n")
	rows = rows[:len(rows)-1]
	numResults := 0

	check := func(i, j int) {
		t := fmt.Sprintf("%c%c", rows[i-1][j-1], rows[i+1][j+1])
		if t != "MS" && t != "SM" {
			return
		}
		t = fmt.Sprintf("%c%c", rows[i-1][j+1], rows[i+1][j-1])
		if t != "MS" && t != "SM" {
			return
		}
		numResults++
	}

	for i := 1; i < len(rows)-1; i++ {
		row := rows[i]
		for j := 1; j < len(row)-1; j++ {
			char := row[j]
			if char == byte('A') {
				check(i, j)
			}
		}
	}

	fmt.Println(numResults)
}
