package internal

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"slices"
)

type Inserter struct {
	database            *sql.DB
	mapping             *Mapping
	csvPath             string
	insertionReferences map[string][]int64
	header              map[string]int
}

// Data to be passed to a Insertable
type InsertContext struct {
	header              map[string]int
	insertionReferences map[string][]int64
	csvContent          []string
}

type Insertable interface {
	generateValue(context InsertContext) (string, error)
}

func newInserter(database *sql.DB, mapping *Mapping, csvPath string) (*Inserter, error) {
	inserter := &Inserter{
		database:            database,
		mapping:             mapping,
		csvPath:             csvPath,
		insertionReferences: map[string][]int64{},
	}

	header, err := inserter.CreateHeader()

	if err != nil {
		return nil, err
	}

	inserter.header = header

	return inserter, nil
}

func (i *Inserter) CreateHeader() (map[string]int, error) {
	file, err := os.Open(i.csvPath)

	if err != nil {
		return nil, fmt.Errorf("Could not open the csv file: %s", err)
	}

	reader := csv.NewReader(file)
	line, err := reader.Read()

	if err != nil {
		return nil, fmt.Errorf("Could not read the csv file: %s", err)
	}

	header := map[string]int{}
	for idx, field := range line {
		header[field] = idx
	}

	return header, nil
}

func createStatements(tables []Table, transaction *sql.Tx) ([]*sql.Stmt, error) {
	statements := []*sql.Stmt{}

	for _, table := range tables {
		log.Debugf("Creating prepared statement for %s", table.name)
		statement, err := transaction.Prepare(table.createStatment())

		if err != nil {
			return nil, fmt.Errorf("Could not create prepared statment for %s: %s", table.name, err)
		}

		log.Tracef("Created %s", table.createStatment())
		statements = append(statements, statement)
	}
	return statements, nil
}

func (f *FormatedInput) generateValue(context InsertContext) (string, error) {
	params := []string{}

	for _, param := range f.params {
		index, ok := context.header[param]

		if !ok {
			return "", fmt.Errorf("Csv file does not have a %s field", param)
		}

		params = append(params, context.csvContent[index])
	}

	out, err := exec.Command(f.scriptPath, params...).Output()
	if err != nil {
		return "", fmt.Errorf("Error executing formatting script: %s", err)
	}

	return string(out), nil
}

func (t *TableReference) generateValue(context InsertContext) (string, error) {
	references, ok := context.insertionReferences[t.referenceTable]

	if !ok {
		return "", fmt.Errorf("Tried to reference %s, but it was never inserted", t.referenceTable)
	}

	if t.insertion > len(references)-1 {
		return "", fmt.Errorf(
			"%s had %d insertions, but tried to get the %d nth",
			t.referenceTable,
			len(references),
			t.insertion,
		)
	}

	return fmt.Sprint(references[t.insertion]), nil
}

func (t *RegularInsertion) generateValue(context InsertContext) (string, error) {
	index, ok := context.header[t.value]

	if !ok {
		return "", fmt.Errorf("Csv file does not have a %s field", t.value)
	}

	value := context.csvContent[index]

	return value, nil
}
