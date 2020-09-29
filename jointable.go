package strata

import (
	"bytes"
	"fmt"
)

// JoinType for JoinTable
type JoinType int

const (
	// LeftJoin join
	LeftJoin JoinType = iota
	// RightJoin type
	RightJoin
	// InnerJoin type
	InnerJoin
	// OuterJoin type
	OuterJoin
)

// SQL returns the SQL representation of the join type
func (jt *JoinType) SQL() string {
	if jt == nil {
		return ""
	}
	switch *jt {
	case RightJoin:
		return "RIGHT"
	case InnerJoin:
		return "INNER"
	case LeftJoin:
		return "LEFT"
	case OuterJoin:
		return "OUTER"
	default:
		return ""
	}
}

// JoinTable is a table with extra properties, will be appended to
// the from clause of this statement as a Join
type JoinTable struct {
	Table
	// When adding a new property,
	JoinType       JoinType
	ComparisonType ComparisonType
	LHSField       *TableField
	RHSField       *TableField
}

func (jt *JoinTable) assert() error {
	if jt.LHSField == nil {
		return fmt.Errorf("LHSField of join table %v undefined", jt.SQL())
	}
	if jt.RHSField == nil {
		return fmt.Errorf("RHSField of join table %v undefined", jt.SQL())
	}
	return nil
}

// JoinTables is a collection of join tables
type JoinTables []JoinTable

func (jt *JoinTables) assert() error {
	if jt == nil {
		return nil
	}
	for _, _jt := range *jt {
		if err := _jt.assert(); err != nil {
			return err
		}
	}

	return nil
}

func (jt *JoinTables) fields() TableFields {
	fields := TableFields{}
	for _, table := range *jt {
		fields.append(table.Fields...)
	}
	return fields
}

func (jt *JoinTables) wheres() []Wheres {
	wheres := []Wheres{}
	for _, table := range *jt {
		wheres = append(wheres, table.WhereConditions)
	}
	return wheres
}

// SQL returns the SQL representation of the join tables object
func (jt *JoinTables) SQL() (string, error) {
	if err := jt.assert(); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.Grow(150)
	for i, table := range *jt {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(" " + table.JoinType.SQL() + " JOIN ")
		buf.WriteString(table.SQL())
		buf.WriteString(" ON ")
		buf.WriteString(table.LHSField.SQL() + " ")
		buf.WriteString(table.ComparisonType.SQL() + " ")
		if !table.ComparisonType.IsExact() {
			buf.WriteString("'%' + ")
		}
		buf.WriteString(table.RHSField.SQL() + " ")
		if !table.ComparisonType.IsExact() {
			buf.WriteString(" + '%' ")
		}
	}
	return buf.String(), nil
}

func (jt *JoinTables) append(tables ...JoinTable) {
	*jt = append(*jt, tables...)
}

func (jt *JoinTable) fixFields() {
	if len(jt.Fields) == 0 {
		jt.Fields = nil
	}

	if len(jt.WhereConditions.ToNativeSlice()) == 0 {
		jt.WhereConditions = Wheres{}
	}
}

func (jt *JoinTables) fixFields() {
	for _, tables := range *jt {
		tables.fixFields()
	}
}

func makeJoinTable(name, schema string, _type JoinType) *JoinTable {
	return &JoinTable{
		Table: Table{
			Name: name, Schema: schema,
		},
		JoinType: _type,
	}
}

// MakeLeftJoinTable returns a join table with
func MakeLeftJoinTable(name, schema string) *JoinTable {
	return makeJoinTable(name, schema, LeftJoin)
}

// MakeRightJoinTable returns a join table with
func MakeRightJoinTable(name, schema string) *JoinTable {
	return makeJoinTable(name, schema, RightJoin)
}

// MakeInnerJoinTable returns a join table with
func MakeInnerJoinTable(name, schema string) *JoinTable {
	return makeJoinTable(name, schema, InnerJoin)
}

// MakeOuterJoinTable returns a join table with
func MakeOuterJoinTable(name, schema string) *JoinTable {
	return makeJoinTable(name, schema, OuterJoin)
}

func (jt *JoinTable) setRHSField(field *TableField) {
	if field == nil {
		return
	}
	jt.RHSField = field
}

// SetEqualTo is some syntactic sugar for setting the comparison operator to equality
func (jt *JoinTable) SetEqualTo(field *TableField) *JoinTable {
	jt.ComparisonType = Equal
	jt.setRHSField(field)
	return jt
}

// SetNotEqualTo is some syntactic sugar for setting the comparison operator to not equality
func (jt *JoinTable) SetNotEqualTo(field *TableField) *JoinTable {
	jt.ComparisonType = NotEqual
	jt.setRHSField(field)
	return jt
}

// SetLHSField is some syntactic sugar for setting the lhs field
func (jt *JoinTable) SetLHSField(name string) *JoinTable {
	f := jt.FieldByName(name)
	if f == nil {
		return jt
	}
	jt.LHSField = f
	return jt
}

// SetILike is some syntactic sugar for setting the comparison operator to ILike
func (jt *JoinTable) SetILike(field *TableField) *JoinTable {
	jt.ComparisonType = ILike
	jt.setRHSField(field)
	return jt
}

// SetLike is some syntactic sugar for setting the comparison operator to Like
func (jt *JoinTable) SetLike(field *TableField) *JoinTable {
	jt.ComparisonType = Like
	jt.setRHSField(field)
	return jt
}

// SetNotILike is some syntactic sugar for setting the comparison operator to NotILike
func (jt *JoinTable) SetNotILike(field *TableField) *JoinTable {
	jt.ComparisonType = NotILike
	jt.setRHSField(field)
	return jt
}

// SetNotLike is some syntactic sugar for setting the comparison operator to NotLike
func (jt *JoinTable) SetNotLike(field *TableField) *JoinTable {
	jt.ComparisonType = NotLike
	jt.setRHSField(field)
	return jt
}

// WithFields is a join table wrapper that adds the given fields to the
func (jt *JoinTable) WithFields(fields ...TableField) *JoinTable {
	jt.AddFields(fields...)
	return jt
}
