package internal

import (
	"fmt"
	"strings"
)

type Table struct {
	name   string
	fields map[string]Insertable
	order  []string
}

func newTable(name string, fields map[string]Insertable) *Table {
	order := []string{}

	for field := range fields {
		order = append(order, field)
	}

	return &Table{name: name, fields: fields, order: order}
}

func (t *Table) createStatment() string {
	placeholders := strings.Repeat("?, ", len(t.fields))
	placeholders, _ = strings.CutSuffix(placeholders, ", ")

	statment := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		t.name,
		strings.Join(t.order, ", "),
		placeholders,
	)

	return statment
}

func (t *Table) buildValues(context InsertContext) ([]any, error) {
	values := []any{}

	for _, field := range t.order {
		insertable := t.fields[field]
		value, err := insertable.generateValue(context)

		if err != nil {
			return nil, fmt.Errorf("Error generating value for %s:%s: %s", t.name, field, err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (t *Table) findReferences() []TableReference {
	references := []TableReference{}
	for _, insertable := range t.fields {
		switch value := insertable.(type) {
		case *TableReference:
			references = append(references, *value)
		}
	}

	return references
}
