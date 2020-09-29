package strata

import "fmt"

// JoinType for JoinTable
type JoinType int

const (
	// Left join
	Left JoinType = iota
	// Right join
	Right JoinType = iota
	// Inner join
	Inner JoinType = iota
)

// SQL returns the SQL representation of the join type
func (jt *JoinType) SQL() string {
	if jt == nil {
		return ""
	}
	switch *jt {
	case Right:
		return "RIGHT"
	case Inner:
		return "INNER"
	case Left:
		return "LEFT"
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

	sql := ""
	for i, table := range *jt {
		if i > 0 {
			sql += ", "
		}
		sql += table.JoinType.SQL() + " JOIN "
		sql += table.LHSField.SQL() + " "
		sql += table.ComparisonType.SQL()
		sql += table.RHSField.SQL()
	}
	return sql, nil
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
