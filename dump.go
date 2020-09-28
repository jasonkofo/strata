package database

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/IMQS/devcon/config"
	"github.com/jackc/pgx"
)

type suggestionType int

type fieldMetadata struct {
	descriptions []pgx.FieldDescription
	rowValues    []interface{}
}

type nameIDPair struct {
	fieldDescIdx int
	name         string
}

const (
	searchTowns suggestionType = iota
	searchErven
)

func (t *suggestionType) Name() string {
	if t == nil {
		return ""
	}
	switch *t {
	case searchTowns:
		return "TownSuggestions"
	case searchErven:
		return "Erven"
	default:
		return "Unknown"
	}
}

// GetERFSuggestions inspects the main database for suggestions for ERFs
func (s *DBServer) GetERFSuggestions(searchTerm, limitingConditions string) ([]map[string]interface{}, error) {
	return s.getSuggestions(searchTerm, limitingConditions, searchErven)
}

// GetTownSuggestions inspects the main database for suggestions for Towns, based on configuration
func (s *DBServer) GetTownSuggestions(searchTerm, limitingConditions string) ([]map[string]interface{}, error) {
	return s.getSuggestions(searchTerm, limitingConditions, searchTowns)
}

func (s *DBServer) getSuggestions(searchTerm, limitingConditions string, searchType suggestionType) ([]map[string]interface{}, error) {
	if !s.allowExternalMainDB || s.mainDBSchema == nil {
		return nil, nil
	}

	sql, err := s.constructCadastralSQLFromConf(searchTerm, searchType, limitingConditions)
	if err != nil {
		return nil, fmt.Errorf("Invalid MainDB Schema: %v", err)
	}

	rows, err := s.mainPool.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("Could not get %v: %v", searchType.Name(), err)
	}

	fieldDesc := rows.FieldDescriptions()
	desc := fieldMetadata{
		descriptions: fieldDesc,
	}

	fields := []map[string]interface{}{}

	for rows.Next() {
		out := make(map[string]interface{})
		rowValues, err := rows.Values()
		desc.rowValues = rowValues

		if err != nil || rowValues == nil {
			s.log.Errorf("Could not read row Value:\nSQL: %v\nVALUES: %v", sql, rowValues)
			continue
		}

		if valueMatches(&desc, searchTerm) {
			for i, val := range rowValues {
				out[fieldDesc[i].Name] = FixValue(val)
			}
			fields = append(fields, out)
		}

	}

	return fields, nil
}

func FixString(value string) string {
	return strings.TrimLeft(value, "0")
}

func FixValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return FixString(v)
	default:
		return value
	}
}

func (s *DBServer) constructCadastralSQLFromConf(searchTerm string, sugType suggestionType, limitingConditions string) (string, error) {
	var (
		builder         strings.Builder
		scannableTables *[]config.MainSchemaColumns
	)

	builder.Reset()
	builder.Grow(256)
	defer builder.Reset()

	switch sugType {
	case searchTowns:
		scannableTables = &s.mainDBSchema.TownSuggestionTables
		break
	case searchErven:
		scannableTables = &s.mainDBSchema.ErfTables
		break
	default:
		return "", errors.New("Unrecognized suggestion type")
	}

	builder.WriteString("SELECT * FROM (")

	for i, dbTables := range *scannableTables {
		lc := resolveLimitingConditions(limitingConditions, &dbTables.Columns)
		if i > 0 {
			builder.WriteString(" UNION ALL ")
		}
		builder.WriteString("SELECT ")
		if dbTables.SelectDistinct {
			builder.WriteString("DISTINCT ")
		}
		builder.Write(cadastralSelectClause(dbTables.Columns))
		builder.WriteString(" FROM ")
		builder.WriteString(postgresQuoteSelector(dbTables.TableName))
		builder.WriteString(" WHERE ")
		builder.Write(cadastralWhereConditions(searchTerm, dbTables, lc))
		builder.WriteString(" LIMIT 25) " + string(i+65) + " ")
		if hasOrderClause, clause := cadastralOrder(dbTables.Columns); hasOrderClause {
			builder.WriteString(" ORDER BY ")
			builder.WriteString(clause)
		}
	}

	return builder.String(), nil
}

func valueMatches(meta *fieldMetadata, searchTerm string) bool {
	// value is the property, according to configuration, that is going to be
	// returned to the front end and rendered by the autocomplete component
	// This means that we treat any field named "value" as

	if meta == nil || (*meta).descriptions == nil || (*meta).rowValues == nil {
		return false
	}

	idx := getFieldDescIdxByName(&meta.descriptions, "value")
	if idx == -1 {
		return false
	}

	if (*meta).descriptions[idx].Name != "value" {
		return false
	}

	rString, ok := (*meta).rowValues[idx].(string)
	if !ok {
		return false
	}

	return isContained(rString, searchTerm)
}

func resolveLimitingConditions(limitingConditions string, columns *[]config.MainSchemaColumn) *map[string]nameIDPair {
	if limitingConditions == "" || columns == nil {
		return nil
	}

	var (
		out  = make(map[string]nameIDPair)
		cArr = strings.Split(limitingConditions, ";")
	)

	for _, c := range cArr {
		cond := strings.Split(c, "=")
		if len(cond) != 2 {
			continue
		}
		out[cond[0]] = nameIDPair{
			name:         cond[1],
			fieldDescIdx: getColumnIndexByName(columns, cond[0]),
		}
	}
	return &out
}

func cadastralSelectClause(columns []config.MainSchemaColumn) []byte {
	var buffer bytes.Buffer
	buffer.Reset()
	buffer.Grow(100)
	defer buffer.Reset()

	for i, column := range columns {
		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(createSelectClause(column))
	}
	return buffer.Bytes()
}

func cadastralWhereConditions(searchTerm string, dbTable config.MainSchemaColumns, limitingConditions *map[string]nameIDPair) []byte {
	var buffer bytes.Buffer
	buffer.Reset()
	buffer.Grow(256)
	defer buffer.Reset()

	writeSearchCondition := func(name, searchTerm string) {
		buffer.WriteString(postgresQuoteSelector(name))
		buffer.WriteString(` ILIKE '%` + searchTerm + `%'`)

	}

	for i, column := range dbTable.Columns {
		if !column.IncludeInConditionals {
			continue
		}

		if i > 0 && limitingConditions != nil {
			buffer.WriteString(surroundWithSpaces(dbTable.GetLogicalOperator()))
		} else if i > 0 {
			buffer.WriteString(surroundWithSpaces("OR"))
		}

		if !column.AllowExternalConditions || limitingConditions == nil {
			writeSearchCondition(column.Name, searchTerm)
			continue
		}

		nameKey := pickFriendlyNameOrSelector(column)
		lc, ok := (*limitingConditions)[nameKey]
		if !ok {
			writeSearchCondition(column.Name, searchTerm)
			continue
		}

		writeSearchCondition(column.Name, lc.name)
	}

	return buffer.Bytes()
}

// cadastralOrder only caters for a single column to be ordered, for performance reasons
func cadastralOrder(columns []config.MainSchemaColumn) (bool, string) {
	var buffer bytes.Buffer
	buffer.Reset()
	buffer.Grow(256)
	defer buffer.Reset()

	for _, column := range columns {
		if column.Order {
			return true, postgresQuoteSelector(*column.FriendlyName) + " ASC"
		}
	}

	return false, ""
}

func getColumnIndexByName(columns *[]config.MainSchemaColumn, name string) int {
	if columns == nil {
		return -1
	}

	for i, column := range *columns {
		if column.FriendlyName != nil && *column.FriendlyName == name {
			return i
		}
	}

	return -1
}

func getFieldDescIdxByName(fieldDesc *[]pgx.FieldDescription, name string) int {
	if fieldDesc == nil {
		return -1
	}

	for i, desc := range *fieldDesc {
		if strings.ToLower(desc.Name) == name {
			return i
		}
	}

	return -1
}
