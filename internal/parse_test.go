package internal

import (
	"math/rand"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestParseFieldValue(t *testing.T) {
	binarys := []string{
		getRandomBinaryPath(),
		getRandomBinaryPath(),
	}

	tests := []struct {
		value      string
		expected   Insertable
		shouldFail bool
	}{
		{
			binarys[0] + ", param1, param2",
			&FormatedInput{
				scriptPath: binarys[0],
				params:     []string{"param1", "param2"},
			},
			false,
		},
		{
			binarys[1] + ", param 1, param2, param 3",
			&FormatedInput{
				scriptPath: binarys[1],
				params:     []string{"param 1", "param2", "param 3"},
			},
			false,
		},
		{
			"__table__",
			&TableReference{
				referenceTable: "table",
				insertion:      0,
			},
			false,
		},
		{
			"__table__ ,4",
			&TableReference{
				referenceTable: "table",
				insertion:      4,
			},
			false,
		},
		{
			"__table",
			&RegularInsertion{value: "__table"},
			false,
		},
		{
			"_table__",
			&RegularInsertion{value: "_table__"},
			false,
		},
	}

	for _, tt := range tests {
		parser := NewParser(nil)

		got, err := parser.parseFieldValue(tt.value)

		if tt.shouldFail && err == nil {
			t.Errorf("Expected \"%s\", to fail but got %v", tt.value, got)
			return
		}

		if !reflect.DeepEqual(tt.expected, got) && !tt.shouldFail {
			t.Errorf("Expected %v, but got %v: %s", tt.expected, got, tt.value)
		}
	}
}

func TestParseRegularInsertion(t *testing.T) {
	tests := []struct {
		value    string
		expected *RegularInsertion
	}{
		{
			"csvData",
			&RegularInsertion{
				value: "csvData",
			},
		},
		{
			"1234",
			&RegularInsertion{
				value: "1234",
			},
		},
		{
			"whatever",
			&RegularInsertion{
				value: "whatever",
			},
		},
	}

	for _, tt := range tests {
		parser := NewParser(nil)

		got, _ := parser.parseRegularInsertion(tt.value)

		if !reflect.DeepEqual(tt.expected, got) {
			t.Errorf("Expected %v, but got %v", tt.expected, got)
		}
	}

}
func TestParseFormattedInput(t *testing.T) {
	binarys := []string{
		getRandomBinaryPath(),
		getRandomBinaryPath(),
		getRandomBinaryPath(),
	}

	tests := []struct {
		value      string
		expected   *FormatedInput
		shouldFail bool
	}{
		{
			binarys[0] + ", param1, param2",
			&FormatedInput{
				scriptPath: binarys[0],
				params:     []string{"param1", "param2"},
			},
			false,
		},
		{
			binarys[1] + ", param 1, param2, param 3",
			&FormatedInput{
				scriptPath: binarys[1],
				params:     []string{"param 1", "param2", "param 3"},
			},
			false,
		},
		{
			binarys[1] + ", param1",
			&FormatedInput{
				scriptPath: binarys[1],
				params:     []string{"param1"},
			},
			false,
		},
		{
			binarys[2],
			&FormatedInput{
				scriptPath: binarys[2],
				params:     []string{},
			},
			false,
		},
		{
			"fakebin" + ", param1",
			nil,
			true,
		},
		{
			"fake/bin" + ", param1, param 2",
			nil,
			true,
		},
		{
			"evenfakier",
			nil,
			true,
		},
		{
			"evenfakier" + " param 1 345",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		parser := NewParser(nil)

		got, err := parser.parseFormatedInput(tt.value)

		if tt.shouldFail && err == nil {
			t.Errorf("Expected %v, to fail but got %v", tt.expected, got)
			return
		}

		if !reflect.DeepEqual(tt.expected, got) && !tt.shouldFail {
			t.Errorf("Expected %v, but got %v", tt.expected, got)
		}
	}

}

func TestParseTableReference(t *testing.T) {
	tests := []struct {
		value string
		// It may return a RegularInsertion
		expected   Insertable
		shouldFail bool
	}{
		{
			"__table",
			&RegularInsertion{value: "__table"},
			false,
		},
		{
			"_table__",
			&RegularInsertion{value: "_table__"},
			false,
		},
		{
			"__table__",
			&TableReference{
				referenceTable: "table",
				insertion:      0,
			},
			false,
		},
		{
			"__table__ ,4",
			&TableReference{
				referenceTable: "table",
				insertion:      4,
			},
			false,
		},
		{
			"__table__, 789",
			&TableReference{
				referenceTable: "table",
				insertion:      789,
			},
			false,
		},
		{
			"__table__, random",
			nil,
			true,
		},
		{
			"__table__, 45.13",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		parser := NewParser(nil)

		got, err := parser.parseTableReference(tt.value)

		if tt.shouldFail && err == nil {
			t.Errorf("Expected \"%s\", to fail but got %v", tt.value, got)
			return
		}

		if !reflect.DeepEqual(tt.expected, got) && !tt.shouldFail {
			t.Errorf("Expected %v, but got %v: %s", tt.expected, got, tt.value)
		}
	}
}

// TODO: handle the errors and make it work with other platforms
// this gotta be the ugliest funtion a ever wrote btw
func getRandomBinaryPath() string {
	pathValue, _ := os.LookupEnv("PATH")
	possiblePaths := strings.Split(pathValue, ":")
	dirPath := possiblePaths[rand.Intn(len(possiblePaths))]
	entries, _ := os.ReadDir(dirPath)
	bin := entries[rand.Intn(len(entries))]

	return path.Join(dirPath, bin.Name())
}

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		candidate string
		expected  bool
	}{
		{
			path.Join(os.TempDir()),
			true,
		},
		{
			path.Join(os.TempDir(), "/dir"),
			true,
		},
		{
			path.Join(os.TempDir(), "/dir/file2"),
			true,
		},
		{
			path.Join(os.TempDir(), "path/to/something"),
			false,
		},
		{
			path.Join(os.TempDir(), "path"),
			false,
		},
	}

	for _, tt := range tests {
		got := isValidPath(tt.candidate)

		if got != tt.expected {
			t.Errorf("Expected %t, but got %t: %s", tt.expected, got, tt.candidate)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		candidate string
		expected  bool
	}{
		{
			"123",
			true,
		},
		{
			"1",
			true,
		},
		{
			"phrase",
			false,
		},
		{
			"1.4",
			false,
		},
		{
			".45",
			false,
		},
		{
			"12.44$ab",
			false,
		},
	}

	for _, tt := range tests {
		got := isNumeric(tt.candidate)

		if got != tt.expected {
			t.Errorf("Expected %t, but got %t", tt.expected, got)
		}
	}
}
