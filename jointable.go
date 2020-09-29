package sqlgen

// JoinTable is a table with extra properties, will be appended to
// the from clause of this statement as a Join
type JoinTable struct {
	Table
	// When adding a new property,
	JoinType       JoinType
	ComparisonType ComparisonType
	LHSJoinField   *TableField
	RHSJoinField   *TableField
}

// JoinTables is a collection of join tables
type JoinTables []JoinTable

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
func (jt *JoinTables) SQL() string {
	sql := ""
	for i, table := range *jt {
		if i > 0 {
			sql += ", "
		}
		sql += table.SQL()
	}
	return sql
}

func (jt *JoinTables) append(tables ...JoinTable) {
	*jt = append(*jt, tables...)
}

func (t *JoinTable) fixFields() {
	if len(t.Fields) == 0 {
		t.Fields = nil
	}

	if len(t.WhereConditions.ToNativeSlice()) == 0 {
		t.WhereConditions = Wheres{}
	}
}

func (jt *JoinTables) fixFields() {
	for _, tables := range *jt {
		tables.fixFields()
	}
}
