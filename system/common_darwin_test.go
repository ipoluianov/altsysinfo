package system

import (
	"reflect"
	"testing"
)

func TestProfilerValue(t *testing.T) {
	tests := []struct {
		key  string
		v    any
		want string
	}{
		{"spdisplays_pixelresolution", "spdisplays_2560x1664Retina", "2560 x 1664 Retina"},
		{"USBKeyLocationID", "0x01000000", "0x01000000"},
		{"coreaudio_device_transport", "coreaudio_device_type_builtin", "Builtin"},
		{"spdisplays_main", "spdisplays_yes", "Yes"},
	}
	for _, test := range tests {
		if got := profilerValue(test.key, test.v); got != test.want {
			t.Errorf("profilerValue(%q, %q) = %q, want %q", test.key, test.v, got, test.want)
		}
	}
}

func TestAddProfilerItemGroups(t *testing.T) {
	table := newKeyValueTable("Power")
	addProfilerSections(table, []map[string]any{{
		"_name":         "sppower_information",
		"AC Power":      map[string]any{"Wake On LAN": "Yes"},
		"Battery Power": map[string]any{"Wake On LAN": "No"},
	}})
	want := [][]string{
		{"Power settings"},
		{"AC Power"}, {"Wake On LAN", "Yes"},
		{"Battery Power"}, {"Wake On LAN", "No"},
	}
	if !reflect.DeepEqual(table.Rows, want) {
		t.Errorf("rows = %q, want %q", table.Rows, want)
	}
}
