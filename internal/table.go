package internal

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
