package main

import (
	"fmt"
	"os"
	"strings"
)

type move struct {
	di int
	dj int
}
type guardCoord struct {
	i int
	j int
	m move
}

func (c guardCoord) string() string {
	return fmt.Sprintf("%d|%d|%d|%d", c.i, c.j, c.m.di, c.m.dj)
}

var (
	right = move{0, 1}
	down  = move{1, 0}
	left  = move{0, -1}
	up    = move{-1, 0}
)
var moveLookup = map[rune]move{
	'>': right,
	'v': down,
	'<': left,
	'^': up,
}
var turnLookup = map[move]move{
	right: down,
	down:  left,
	left:  up,
	up:    right,
}

func copy(s map[string]bool) map[string]bool {
	c := make(map[string]bool, len(s))
	for k, v := range s {
		c[k] = v
	}
	return c
}

func main() {
	b, _ := os.ReadFile("input.txt")
	rows := strings.Split(string(b), "\n")
	rows = rows[:len(rows)-1] // trim blank line
	numCols := len(rows[0])

	// note where obstacles are with x+y coordinate keys
	obstacles := make(map[string]bool)
	var guardPos guardCoord
	for i := 0; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j++ {
			c := rune(rows[i][j])
			if c == '#' {
				obstacles[fmt.Sprintf("%d|%d", i, j)] = true
			} else if c != '.' {
				guardPos = guardCoord{
					i: i,
					j: j,
					m: moveLookup[c],
				}
			}
		}
	}
	distinctPositions := make(map[string]bool, 0)    // part 1
	distinctNewObstacles := make(map[string]bool, 0) // part 2

	var walkPath func(distinctStates map[string]bool, guardPos guardCoord, searchForObstance bool)

	walkPath = func(distinctStates map[string]bool, guardPos guardCoord, searchForObstacle bool) {
		if ok := distinctStates[guardPos.string()]; ok {
			// we've been at this exact state and we are obstacle searching, and the obstacle wouldn't be out of bounds, increment
			op := []int{guardPos.i + guardPos.m.di, guardPos.j + guardPos.m.dj}
			if searchForObstacle {
				distinctNewObstacles[fmt.Sprintf("%d|%d", op[0], op[1])] = true
				return
			}
		} else {
			distinctStates[guardPos.string()] = true
		}
		if !searchForObstacle {
			distinctPositions[fmt.Sprintf("%d|%d", guardPos.i, guardPos.j)] = true
		}

		nextCoord := guardCoord{
			i: guardPos.i + guardPos.m.di,
			j: guardPos.j + guardPos.m.dj,
			m: guardPos.m,
		}
		if nextCoord.i < 0 || nextCoord.j < 0 || nextCoord.i == len(rows) || nextCoord.j == numCols {
			// out of bounds
			return
		}

		if obstacles[fmt.Sprintf("%d|%d", nextCoord.i, nextCoord.j)] {
			guardPos.m = turnLookup[guardPos.m]
		} else {
			guardPos = nextCoord
			// for each unique state, turn right and see if following that path would get us to a loop
			if !searchForObstacle && !obstacles[fmt.Sprintf("%d|%d", guardPos.i+guardPos.m.di, guardPos.j+guardPos.m.dj)] {
				walkPath(copy(distinctStates), guardCoord{guardPos.i, guardPos.j, turnLookup[guardPos.m]}, true)
			}
		}

		walkPath(distinctStates, guardPos, searchForObstacle)
	}

	distinctStates := make(map[string]bool)
	walkPath(distinctStates, guardPos, false)

	fmt.Println(len(distinctPositions))
	fmt.Println(len(distinctNewObstacles))
}
