package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// @project life
// @author arthur
// @created 13.05.2022

func main() {
	rand.Seed(time.Now().UnixNano())

	row := flag.Int("row", 15, "rows count")
	col := flag.Int("col", 15, "columns count")
	flag.Parse()

	gameField := CreateGameField(*row, *col)
	Shuffle(gameField)

	stdout := bufio.NewWriter(os.Stdout)

	start := time.Now()
	generation := 1
	ticker := time.NewTicker(100 * time.Millisecond)
	for t := range ticker.C {
		if err := stdout.Flush(); err != nil {
			log.Println(err)
		}

		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

		rows, cols, count := RowsColsCount(gameField)
		fmt.Fprintln(stdout, "-----", t.Format("15:04:05"),
			"Gen:", generation,
			"Rows:", rows,
			"Cols:", cols,
			"Count:", count,
			"Second:", time.Now().Sub(start).Seconds(),
			"-----",
		)

		drawed := DrawPlain(stdout, gameField)
		fmt.Fprintf(stdout, "\nDrawed: %d\n", drawed)
		generation++
		NewTick(gameField)
	}
}

func NewTick(field []*Cell) {
	// calc and save next state
	for _, row := range field {
		cell := row
		for cell != nil {
			switch cell.Check() {
			case 1:
				cell.inNextState = true
			case -1:
				cell.inNextState = false
			case 0:
				cell.inNextState = cell.alive
			}
			cell = cell.right
		}
	}
	// append next state to current
	for _, row := range field {
		cell := row
		for cell != nil {
			cell.alive = cell.inNextState
			cell = cell.right
		}
	}
}

type Cell struct {
	start              *Cell
	alive, inNextState bool

	row, column int

	left, right, top, down      *Cell
	diagLeftTop, diagRightTop   *Cell
	diagLeftDown, diagRightDown *Cell
}

// Check return 1,0,-1
// Reborn, living, dying
func (c *Cell) Check() int {
	liveNeighbors := 0
	chk := func(cells ...*Cell) {
		for _, cell := range cells {
			if cell != nil && cell.alive {
				liveNeighbors++
			}
		}
	}
	chk(c.left, c.diagLeftTop, c.top, c.diagRightTop, c.right, c.diagRightDown, c.down, c.diagLeftDown)
	if liveNeighbors < 2 || liveNeighbors > 3 {
		return -1 // die
	}
	if liveNeighbors == 3 {
		return 1
	}
	return 0
}

const Blank = "_"

func (c *Cell) String() string {
	if c == nil {
		return Blank
	}
	if c.alive {
		return "██"
	}
	return "░░"
}

func (c *Cell) StringNeighbor() string {
	return fmt.Sprintf("%s\t%s\t%s\n%s\t%s\t%s\n%s\t%s\t%s\n",
		c.diagLeftTop, c.top, c.diagRightTop,
		c.left, c, c.right,
		c.diagLeftDown, c.down, c.diagRightDown)
}

func CreateGameField(rowNum, columnNum int) []*Cell {
	field := InitRows(rowNum)
	PrepareColumns(field, columnNum)
	DiagonalLinking(field)
	return field
}

func Shuffle(field []*Cell) {
	for _, row := range field {
		cell := row
		for cell != nil {
			if rand.Intn(100) > 50 {
				cell.alive = true
			}
			cell = cell.right
		}
	}

}

func DrawPlain(w io.Writer, field []*Cell) int {
	drawed := 0
	for _, cell := range field {
		for cell != nil {
			drawed++
			if _, err := fmt.Fprintf(w, "%s", cell); err != nil {
				log.Println(err)
			}
			cell = cell.right
		}
		fmt.Fprint(w, "\n")
		if _, err := w.Write([]byte("\n")); err != nil {
			log.Println(err)
		}
	}
	return drawed
}

// InitRows return GameField with *Cell
// make top/down linking between Cells
func InitRows(rowNum int) []*Cell {
	var field = make([]*Cell, 0, rowNum)
	var rowCell *Cell

	for row := 0; row < rowNum; row++ {
		cell := &Cell{row: row, column: 0}
		if rowCell != nil {
			rowCell.down = cell
			cell.top = rowCell
			rowCell = rowCell.down
		} else {
			rowCell = cell
		}
		field = append(field, rowCell)
	}
	return field
}

// PrepareColumns fill Columns
// make left/right linking
func PrepareColumns(field []*Cell, colNum int) {
	var prevRow *Cell
	for row, rowCell := range field {
		cell := rowCell
		for column := 1; column < colNum; column++ {
			cell.right = &Cell{row: row, column: column, left: cell, start: rowCell}
			if prevRow != nil {
				cell.top = prevRow
				prevRow.down = cell
				// prevRow next
				prevRow = prevRow.right
			}
			cell = cell.right
		}
		if prevRow != nil { // prepare last element on column
			cell.top = prevRow
			prevRow.down = cell
		}
		prevRow = rowCell
	}
}

func DiagonalLinking(field []*Cell) {
	// diagonal neighborhoods
	for _, row := range field {
		cell := row
		for cell != nil {
			if cell.top != nil {
				// left top
				if cell.left != nil {
					// neighbors neighbors
					if cell.left.top == cell.top.left { // left top diagonal
						cell.diagLeftTop = cell.left.top
					}
				}
				if cell.right != nil {
					// neighbors neighbors
					rightTop, topRight := cell.right.top, cell.top.right
					if rightTop == topRight {
						cell.diagRightTop = rightTop
					}
				}
			}

			if cell.down != nil {
				if cell.left != nil {
					leftDown, downLeft := cell.left.down, cell.down.left
					if leftDown == downLeft {
						cell.diagLeftDown = leftDown
					}
				}
				if cell.right != nil {
					rightDown, downRight := cell.right.down, cell.down.right
					if rightDown == downRight {
						cell.diagRightDown = rightDown
					}
				}
			}

			cell = cell.right
		}
	}
}

func Lookup(field []*Cell, row, column int) *Cell {
	if len(field)-1 < row {
		return nil
	}
	rowCell := field[row]
	for rowCell != nil {
		if rowCell.column == column {
			return rowCell
		}
		rowCell = rowCell.right
	}
	return nil
}

func RowsColsCount(field []*Cell) (int, int, int) {
	rows, cols, count := 0, 0, 0
	for _, cell := range field {
		rows++
		for cell != nil {
			count++
			cell = cell.right
		}
	}

	for cell := field[0]; cell != nil; cell = cell.right {
		cols++
	}
	return rows, cols, count
}
