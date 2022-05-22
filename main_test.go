package main

import (
	"fmt"
	"testing"
)

// @project life
// @author arthur
// @created 13.05.2022

func TestInitRows_RowsLength(t *testing.T) {
	for i := 0; i < 100; i++ {
		field := InitRows(i)
		if len(field) != i {
			t.Errorf("Expected %d len, got %d ", i, len(field))
		}
	}
}

func TestInitRows_RowColumn(t *testing.T) {
	for n := 0; n < 100; n++ {
		field := InitRows(n)
		for i, row := range field {
			if row == nil {
				t.Fatal("Row is nil", "i=", i)
			}
			if row.row != i {
				t.Errorf("Incorrect row value, expected %d but given %d ", i, row.row)
			}
			if row.column != 0 {
				t.Errorf("Incorrect row value, expected %d but given %d ", 0, row.column)
			}
		}
	}

}

func TestInitRows_Linking(t *testing.T) {
	n := 3
	field := InitRows(n)
	var prev *Cell
	for _, row := range field {
		if prev != nil {
			if row.top == nil {
				t.Error("Expected top is not nil")
				continue
			}
			if row.top != prev {
				t.Error("Incorrect top neighbor")
			}
			if prev.down != row {
				t.Error("Top neighbor: incorrect down neighbor")
			}
		}
		prev = row
	}
}

func TestPrepareColumns_Length(t *testing.T) {
	field := InitRows(3)
	n := 10
	PrepareColumns(field, n)
	for _, row := range field {
		i := 0
		col := row
		for col != nil {
			i++
			col = col.right
		}
		if i != n {
			t.Errorf("Expected %d but given %d columns ", n, i)
		}
	}
}

func TestPrepareColumns_Linking(t *testing.T) {
	field := InitRows(10)
	n := 10
	PrepareColumns(field, n)
	for _, row := range field {
		var col = row
		var prev *Cell
		for col != nil {
			if col.left != prev {
				t.Errorf("Expected left is %s but given %s", prev, col.left)
			}
			if prev != nil {
				if prev.right != col {
					t.Errorf("Expected prev right %s should equal %s col ", prev.right, col)
				}
			}
			prev = col
			col = col.right
		}
	}

}

var table = []struct {
	input int
}{
	{input: 10},
	{input: 100},
	{input: 1000},
}

func BenchmarkCreateGameField(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("init_size %d ", v.input*v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CreateGameField(v.input, v.input)
			}
		})
	}
}

func BenchmarkInitRows(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("init_size %d ", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				InitRows(v.input)
			}
		})
	}
}

func BenchmarkPrepareColumns(b *testing.B) {
	for _, v := range table {
		field := InitRows(v.input)
		b.Run(fmt.Sprintf("init_size %d ", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				PrepareColumns(field, v.input)
			}
		})
	}
}

func BenchmarkDiagonalLinking(b *testing.B) {
	for _, v := range table {
		field := InitRows(v.input)
		PrepareColumns(field, v.input)
		b.Run(fmt.Sprintf("init_size %d ", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DiagonalLinking(field)
			}
		})
	}
}
