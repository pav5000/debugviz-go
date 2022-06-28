package table

import "fmt"

type Block struct {
	caption string
	rows    []Row
}

type RowClass string

const (
	RowClassTitle RowClass = "title"
)

type Row struct {
	Class RowClass `json:",omitempty"`
	Cells []Cell
}

type Cell struct {
	Colspan int `json:",omitempty"`
	Value   string
}

func New(caption string) *Block {
	return &Block{
		caption: caption,
	}
}

func (b *Block) Type() string {
	return "table"
}

func (b *Block) Data() interface{} {
	if b == nil {
		return nil
	}
	return struct {
		Caption string
		Rows    []Row
	}{
		Caption: b.caption,
		Rows:    b.rows,
	}
}

func (b *Block) addRow(row Row) {
	if b == nil {
		return
	}
	b.rows = append(b.rows, row)
}

func newRow(values ...string) Row {
	row := Row{
		Cells: make([]Cell, 0, len(values)),
	}
	for _, value := range values {
		row.Cells = append(row.Cells, Cell{
			Value: value,
		})
	}
	return row
}

func (b *Block) AddRow(values ...interface{}) *Block {
	if b == nil {
		return nil
	}
	strs := make([]string, 0, len(values))
	for _, value := range values {
		strs = append(strs, fmt.Sprint(value))
	}
	b.addRow(newRow(strs...))
	return b
}

func (b *Block) AddTitleRow(values ...string) *Block {
	if b == nil {
		return nil
	}
	row := newRow(values...)
	row.Class = RowClassTitle
	b.addRow(row)
	return b
}
