package awtree

import (
	"fmt"
	"sort"
)

// detectTables finds table patterns: rows with │ or | at consistent column positions.
func detectTables(g *Grid) []Element {
	// For each column, find rows where there's a vertical separator.
	colRows := make(map[int][]int)
	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			ch := g.At(row, col).Char
			if ch == '│' || ch == '|' {
				colRows[col] = append(colRows[col], row)
			}
		}
	}

	// Identify horizontal separator rows (rows filled mostly with ─, ━, -, =, ┼, ╋, +).
	isSepRow := make([]bool, g.Rows)
	for row := 0; row < g.Rows; row++ {
		sepCount := 0
		nonSpace := 0
		for col := 0; col < g.Cols; col++ {
			ch := g.At(row, col).Char
			if ch != ' ' && ch != 0 {
				nonSpace++
			}
			if isHorizSepChar(ch) {
				sepCount++
			}
		}
		// Row is a separator if most non-space chars are separator chars.
		if nonSpace > 0 && sepCount*100/nonSpace >= 80 {
			isSepRow[row] = true
		}
	}

	// Find separator columns: need 3+ data rows (excluding sep rows) with │/|.
	// A valid table separator must have non-space content on both sides in at
	// least one row — this filters out panel borders which only have content
	// on one side.
	var sepCols []int
	for col, rows := range colRows {
		if len(rows) < 2 {
			continue
		}
		// Check that at least one row has content on both sides of this column.
		hasContentBothSides := false
		for _, row := range rows {
			hasLeft := false
			for c := col - 1; c >= 0; c-- {
				ch := g.At(row, c).Char
				if ch != ' ' && ch != 0 && ch != '│' && ch != '|' {
					hasLeft = true
					break
				}
			}
			hasRight := false
			for c := col + 1; c < g.Cols; c++ {
				ch := g.At(row, c).Char
				if ch != ' ' && ch != 0 && ch != '│' && ch != '|' {
					hasRight = true
					break
				}
			}
			if hasLeft && hasRight {
				hasContentBothSides = true
				break
			}
		}
		if !hasContentBothSides {
			continue
		}
		if len(rows) >= 2 {
			sepCols = append(sepCols, col)
		}
	}

	if len(sepCols) == 0 {
		return nil
	}

	// Build row -> sep col set map.
	rowSepCols := make(map[int]map[int]bool)
	for col, rows := range colRows {
		for _, row := range rows {
			if rowSepCols[row] == nil {
				rowSepCols[row] = make(map[int]bool)
			}
			rowSepCols[row][col] = true
		}
	}

	// For each separator column, find runs of rows allowing separator row gaps.
	type tableCandidate struct {
		startRow int
		endRow   int   // exclusive
		dataRows int   // rows that are not separator lines
		sepCols  []int // columns with vertical separators
	}

	var best tableCandidate

	for _, pivotCol := range sepCols {
		rows := colRows[pivotCol]
		if len(rows) < 2 {
			continue
		}

		// Find runs allowing separator row gaps.
		runs := findRunsWithSepGaps(rows, isSepRow, g.Rows)
		for _, run := range runs {
			// Count data rows (non-separator).
			dataCount := 0
			for _, r := range run {
				if !isSepRow[r] {
					dataCount++
				}
			}
			if dataCount < 2 {
				continue
			}

			// Check which other sep cols are present in all data rows of this run.
			commonCols := []int{pivotCol}
			for _, otherCol := range sepCols {
				if otherCol == pivotCol {
					continue
				}
				allPresent := true
				for _, r := range run {
					if isSepRow[r] {
						continue // skip separator rows
					}
					if !rowSepCols[r][otherCol] {
						allPresent = false
						break
					}
				}
				if allPresent {
					commonCols = append(commonCols, otherCol)
				}
			}

			if dataCount > best.dataRows {
				best = tableCandidate{
					startRow: run[0],
					endRow:   run[len(run)-1] + 1,
					dataRows: dataCount,
					sepCols:  commonCols,
				}
			}
		}
	}

	if best.dataRows < 2 {
		return nil
	}

	numDataRows := best.dataRows
	numCols := len(best.sepCols) + 1
	totalHeight := best.endRow - best.startRow

	// Compute bounding box.
	sortedSepCols := make([]int, len(best.sepCols))
	copy(sortedSepCols, best.sepCols)
	sort.Ints(sortedSepCols)

	rightSep := sortedSepCols[len(sortedSepCols)-1]

	// Left edge: find leftmost non-space content across all table rows.
	minCol := rightSep
	for row := best.startRow; row < best.endRow; row++ {
		for col := 0; col < rightSep; col++ {
			ch := g.At(row, col).Char
			if ch != ' ' && ch != 0 {
				if col < minCol {
					minCol = col
				}
				break
			}
		}
	}

	// Right edge: find rightmost non-space content across all table rows.
	maxCol := rightSep
	for row := best.startRow; row < best.endRow; row++ {
		for col := g.Cols - 1; col > rightSep; col-- {
			ch := g.At(row, col).Char
			if ch != ' ' && ch != 0 {
				if col > maxCol {
					maxCol = col
				}
				break
			}
		}
	}

	return []Element{
		{
			Type:  ElementTable,
			Label: fmt.Sprintf("%dx%d", numDataRows, numCols),
			Bounds: Rect{
				Row:    best.startRow,
				Col:    minCol,
				Width:  maxCol - minCol + 1,
				Height: totalHeight,
			},
		},
	}
}

// isHorizSepChar returns true if the character is a horizontal separator.
func isHorizSepChar(ch rune) bool {
	switch ch {
	case '─', '━', '-', '=', '┼', '╋', '+', '┬', '┴', '╪', '╤', '╧':
		return true
	}
	return false
}

// findRunsWithSepGaps finds groups of rows from `dataRows` that are consecutive
// when allowing separator rows (isSepRow) as gaps between them.
func findRunsWithSepGaps(dataRows []int, isSepRow []bool, totalRows int) [][]int {
	if len(dataRows) == 0 {
		return nil
	}

	var runs [][]int
	current := []int{dataRows[0]}

	for i := 1; i < len(dataRows); i++ {
		prev := dataRows[i-1]
		curr := dataRows[i]

		// Check if all rows between prev and curr are separator rows.
		allSep := true
		for r := prev + 1; r < curr; r++ {
			if !isSepRow[r] {
				allSep = false
				break
			}
		}

		if curr == prev+1 || allSep {
			// Include any separator rows in between.
			for r := prev + 1; r < curr; r++ {
				current = append(current, r)
			}
			current = append(current, curr)
		} else {
			runs = append(runs, current)
			current = []int{curr}
		}
	}
	runs = append(runs, current)
	return runs
}
