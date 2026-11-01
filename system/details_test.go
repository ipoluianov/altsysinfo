package system

import (
	"reflect"
	"testing"
)

func TestDropEmptyColumns(t *testing.T) {
	table := newDetailTable("Test", "A", "B", "C")
	table.Section("Section")
	table.Add("a1", "", "c1")
	table.Add("a2", "", "")
	table.DropEmptyColumns()

	if want := []string{"A", "C"}; !reflect.DeepEqual(table.Columns, want) {
		t.Errorf("columns = %q, want %q", table.Columns, want)
	}
	want := [][]string{{"Section"}, {"a1", "c1"}, {"a2", ""}}
	if !reflect.DeepEqual(table.Rows, want) {
		t.Errorf("rows = %q, want %q", table.Rows, want)
	}
}
