package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var stateStr string
	flag.StringVar(&stateStr, "state", "", "the initial sudoku state")
	flag.Parse()

	// build initial state
	initialState, err := getInitialState(stateStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Initial State is valid")

	printState(initialState)

	iterations := 0

	stack := [][]uint8{initialState}

	for len(stack) > 0 {
		iterations++

		state := stack[0]
		stack = stack[1:]

		if isComplete(state) {
			fmt.Println("Found a solution")
			printState(state)
			break
		}

		expandedStates := expand(state)

		stack = append(expandedStates, stack...)
	}

	fmt.Println("Done")
}

// getInitialState parses a string into a sudoku state
func getInitialState(state string) ([]uint8, error) {
	cells := strings.Split(state, ",")
	if len(cells) != 81 {
		return nil, errors.New("could not parse state")
	}

	var initialState []uint8
	for i := 0; i < 81; i++ {
		val, err := strconv.ParseInt(cells[i], 10, 8)
		if err != nil {
			return nil, errors.New("could not parse cell")
		}

		initialState = append(initialState, uint8(val))
	}

	return initialState, nil
}

func printState(state []uint8) error {
	if len(state) != 81 {
		return errors.New("state is not the right size")
	}

	// top edge
	printMajorHorizontalBoundary()
	for i := 0; i < 9; i++ {
		// each row
		printRow(state, i)
	}

	return nil
}

func printHorizontalBoundary() {
	fmt.Println("-------------------")
}

func printMajorHorizontalBoundary() {
	fmt.Println("===================")
}

func printRow(state []uint8, i int) {
	start := i * 9
	end := i*9 + 9
	r := state[start:end]
	row := getFormattedCells(r)
	fmt.Printf("\u2016%s|%s|%s\u2016%s|%s|%s\u2016%s|%s|%s\u2016\n", row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8])

	rem := i % 3
	if rem == 2 {
		printMajorHorizontalBoundary()
	} else {
		printHorizontalBoundary()
	}
}

func getFormattedCells(row []uint8) []string {
	var formatted []string

	for i := 0; i < 9; i++ {
		if cell := row[i]; cell == 0 {
			formatted = append(formatted, " ")
		} else {
			formatted = append(formatted, fmt.Sprintf("%d", cell))
		}
	}

	return formatted
}

func expand(state []uint8) [][]uint8 {
	expandIdx := getExpandIndex(state)
	return expandOn(state, expandIdx)
}

// getExpandIndex finds the next index upon which to expand
// it looks for the index with the lowest branching factor
func getExpandIndex(state []uint8) int {
	lowestBranchingFactor := 10 // default value larger than any possible branching factor
	bestIndex := 0

	for idx, val := range state {
		if val == 0 {
			if branches := getPossibleValues(state, idx); len(branches) < lowestBranchingFactor {
				lowestBranchingFactor = len(branches)
				bestIndex = idx
			}
		}

		// we can immediately fill in this value!
		if lowestBranchingFactor == 1 {
			break
		}
	}

	return bestIndex
}

func expandOn(state []uint8, idx int) [][]uint8 {
	var expandedStates [][]uint8

	expandValues := getPossibleValues(state, idx)

	for _, v := range expandValues {
		expandedState := make([]uint8, 81)
		copy(expandedState, state)
		expandedState[idx] = v
		expandedStates = append(expandedStates, expandedState)
	}
	return expandedStates
}

// simple completion check assuming the state is valid
func isComplete(state []uint8) bool {
	for _, val := range state {
		if val == 0 {
			return false
		}
	}

	return true
}

func getPossibleValues(state []uint8, idx int) []uint8 {
	// all possible values
	values := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// get square neighbours
	squareNeighbours := getSquareNeighbours(state, idx)
	// get row neighbours
	rowNeighbours := getRowNeighbours(state, idx)
	// get column neighbours
	columnNeighbours := getColumnNeighbours(state, idx)

	values = removeFromSlice(values, squareNeighbours)
	values = removeFromSlice(values, rowNeighbours)
	values = removeFromSlice(values, columnNeighbours)

	return values
}

func removeFromSlice(s []uint8, rem []uint8) []uint8 {
	var t []uint8

	for _, n := range s {
		found := false
		for _, r := range rem {
			if r == n {
				found = true
				break
			}
		}

		if !found {
			t = append(t, n)
		}
	}

	return t
}

func getSquareNeighbours(state []uint8, idx int) []uint8 {
	var neighbours []uint8

	// determine cell column
	cellCol := idx % 3

	// determine cell row
	cellRow := (idx / 9) % 3

	// determine top left corner of cell
	cellStart := idx - cellRow*9 - cellCol

	// iterate over rows
	for i := 0; i < 3; i++ {
		// iterate over cols
		for j := 0; j < 3; j++ {
			curr := cellStart + 9*i + j
			if curr != idx {
				if val := state[curr]; val != 0 {
					neighbours = append(neighbours, val)
				}
			}
		}
	}

	return neighbours
}

func getRowNeighbours(state []uint8, idx int) []uint8 {
	var neighbours []uint8

	// determine row
	row := idx / 9

	for i := 0; i < 9; i++ {
		curr := 9*row + i
		if curr != idx {
			if val := state[curr]; val != 0 {
				neighbours = append(neighbours, val)
			}
		}
	}

	return neighbours
}

func getColumnNeighbours(state []uint8, idx int) []uint8 {
	var neighbours []uint8

	// determine column
	col := idx % 9

	for i := 0; i < 9; i++ {
		curr := i*9 + col
		if curr != idx {
			if val := state[curr]; val != 0 {
				neighbours = append(neighbours, val)
			}
		}
	}

	return neighbours
}
