package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
	mapping Mapping
}

type FormatedInput struct {
	scriptPath string
	params     []string
}

type TableReference struct {
	referenceTable string
	insertion      int
}

type RegularInsertion struct {
	value string
}

func NewParser(mapping Mapping) *Parser {
	return &Parser{
		mapping: mapping,
	}
}

func (p *Parser) parse() ([]Table, error) {
	tables := []Table{}

	for tableName, insertions := range p.mapping {
		for idx, insertion := range insertions {
			fields := map[string]Insertable{}

			for field, value := range insertion {
				log.Debugf("Parsing %s:%d:%s", tableName, idx, field)
				insertable, err := p.parseFieldValue(value)

				if err != nil {
					return nil, fmt.Errorf("Error parsing %s:%d:%s: %s", tableName, idx, field, err)
				}

				fields[field] = insertable
			}

			tables = append(tables, *newTable(tableName, fields))
		}
	}

	return tables, nil
}

func (p *Parser) parseFieldValue(value string) (Insertable, error) {
	if strings.HasPrefix(value, "__") {
		return p.parseTableReference(value)
	}

	if candidate := strings.Split(value, ",")[0]; isValidPath(candidate) {
		return p.parseFormatedInput(value)
	}

	return p.parseRegularInsertion(value)
}

func (p *Parser) parseTableReference(value string) (Insertable, error) {
	values := []string{}

	for _, v := range strings.Split(value, ",") {
		values = append(values, strings.TrimSpace(v))
	}

	if !(strings.HasPrefix(values[0], "__") && strings.HasSuffix(values[0], "__")) &&
		len(values) == 1 {
		return p.parseRegularInsertion(value)
	}

	if len(values) == 1 {
		if !strings.HasSuffix(values[0], "__") {
			return nil, fmt.Errorf(
				"Unexpected text after %s. It should be a single number",
				values[0],
			)
		}

		return &TableReference{
			referenceTable: values[0][2 : len(values[0])-2],
			insertion:      0,
		}, nil
	}

	if isNumeric(values[1]) && len(values) <= 2 {
		// We already know it is a valid number. No need to handle the error
		insertion, _ := strconv.Atoi(values[1])

		return &TableReference{
			referenceTable: values[0][2 : len(values[0])-2],
			insertion:      insertion,
		}, nil
	}

	return nil, fmt.Errorf("Unexpected text after %s. It should be a single number", values[0])
}

func (p *Parser) parseFormatedInput(value string) (Insertable, error) {
	values := []string{}

	for _, v := range strings.Split(value, ",") {
		values = append(values, strings.TrimSpace(v))
	}

	fileInfo, err := os.Stat(values[0])

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid file", values[0])
	}

	if fileInfo.Mode().IsRegular() {
		return &FormatedInput{
			scriptPath: values[0],
			params:     values[1:],
		}, nil
	}

	return nil, fmt.Errorf("%s is not a valid file", values[0])
}

func (p *Parser) parseRegularInsertion(value string) (Insertable, error) {
	return &RegularInsertion{
		value: value,
	}, nil
}

func isValidPath(candidate string) bool {
	_, err := os.Stat(filepath.Clean(candidate))

	return err == nil
}

func isNumeric(candidate string) bool {
	for _, c := range candidate {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
