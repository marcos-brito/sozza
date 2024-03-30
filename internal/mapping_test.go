package internal

import (
	"reflect"
	"testing"
)

func TestParsingMappingFile(t *testing.T) {
	tests := []struct {
		input    string
		expected *Mapping
	}{
		{
			`table1:
                - field1: csv1
                  field2: csv2
                  field3: csv3
                - field1: csv1
                  field2: csv2
                  field3: csv3
            `,
			&Mapping{
				"table1": []map[string]string{
					{
						"field1": "csv1",
						"field2": "csv2",
						"field3": "csv3",
					},
					{
						"field1": "csv1",
						"field2": "csv2",
						"field3": "csv3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		got, err := ReadMapping([]byte(tt.input))

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Expected %v, but got %v", tt.expected, got)
		}
	}

}
