package strata

import (
	"fmt"
	"reflect"
)

// Table is an abstraction of the table type
type Table struct {
	Name            string
	Schema          string
	Alias           *string
	LHS             string
	Fields          TableFields
	WhereConditions Wheres
}

// Tables is a collection of table
type Tables []Table

// IsJoinTable returns whether or not the derived object is a join table using reflection
func (t *Table) IsJoinTable() bool {
	jt := JoinTable{}
	if t == nil {
		return false
	}
	return reflect.TypeOf(t).Elem() == reflect.TypeOf(jt)
}

// AddFields adds all the desired string fields
func (t *Table) AddFields(fields ...TableField) {
	for i := 0; i < len(fields); i++ {
		fields[i].Alias = t.Alias
	}
	t.Fields.append(fields...)
}

// AddWhereCondition adds a where limitation to the table
func (t *Table) AddWhereCondition(lhs *TableField, rhs interface{}, comparisonType ComparisonType) error {

	return nil
}

// AddSimpleStringFields adds all the desired string fields
func (t *Table) AddSimpleStringFields(fields ...string) {
	t.Fields.addSimpleStringFields(t.Alias, fields...)
}

// AddSimpleNumberFields adds all the desired string fields
func (t *Table) AddSimpleNumberFields(fields ...string) {
	t.Fields.addSimpleNumberFields(t.Alias, fields...)
}

// AddSimpleDateFields adds all the desired string fields
func (t *Table) AddSimpleDateFields(fields ...string) {
	t.Fields.addSimpleDateFields(t.Alias, fields...)
}

// AddGeometryField adds a single complex field type
func (t *Table) AddGeometryField(name, friendlyName, formattedName string) {
	t.Fields.addGeometryField(t.Alias, name, friendlyName, formattedName)
}

// AddStringField adds a single complex field type
func (t *Table) AddStringField(name, friendlyName, formattedName string) {
	t.Fields.addStringField(t.Alias, name, friendlyName, formattedName)
}

// AddNumberField adds a single complex field type
func (t *Table) AddNumberField(name, friendlyName, formattedName string) {
	t.Fields.addNumberField(t.Alias, name, friendlyName, formattedName)
}

// AddDateField adds a single complex field type
func (t *Table) AddDateField(name, friendlyName, formattedName string) {
	t.Fields.addDateField(t.Alias, name, friendlyName, formattedName)
}

// AddFieldByProperties adds a single field to the dataset
func (t *Table) AddFieldByProperties(name, friendlyName, formattedName, _type string) {
	t.Fields.addFieldByProperties(t.Alias, name, friendlyName, formattedName, _type)
}

// FieldByName returns a field by the name
func (t *Table) FieldByName(name string) *TableField {
	return t.Fields.fieldByName(name)
}

// SetWhereConditions creates a where condition on the table object on an existing object in the object.
// It first inspects if the field with the given name exists in the object, if not, returning an error
// Then it attempts to create a where condition given the predicate and add it to the table object
func (t *Table) SetWhereConditions(fieldName string, comparisonType ComparisonType, rhs interface{}) error {
	where := t.FieldByName(fieldName).Where(comparisonType, rhs)
	if where == nil {
		return fmt.Errorf("Could not find field %v in the Table object", fieldName)
	}

	t.WhereConditions = *where
	return nil
}

// SQL returns the name of the table object represented as an SQL selector
func (t *Table) SQL() string {
	sql := ""
	if t.Schema != "" {
		sql += chainSelector(t.Schema, t.Name)
	} else {
		sql += insertDoubleQuotes(t.Name)
	}

	if t.Alias != nil && *t.Alias != "" {
		sql += " " + insertDoubleQuotes(*t.Alias)
	}
	return sql
}

func (t *Table) fixFields() {
	if len(t.Fields) == 0 {
		t.Fields = nil
	}

	if len(t.WhereConditions.ToNativeSlice()) == 0 {
		t.WhereConditions = Wheres{}
	}
}
