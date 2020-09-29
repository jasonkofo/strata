package strata

import (
	"bytes"
	"fmt"
	"strconv"
)

// SQLElement definition
type SQLElement interface {
	SQL() string
}

// Query is the second highest level of abstraction for the SQL result set -
// they are joined together using unions
type Query struct {
	baseTable  Table
	joinTables JoinTables
	Limit      int
}

// NestedFields returns all the that are in the query object (i.e.
// in the base table and the join tables) as a single set of
// TableFields. This is used to create the select statement
func (q *Query) NestedFields() string {
	fields := TableFields{}
	fields.append(q.baseTable.Fields...)
	fields.append(q.joinTables.fields()...)
	return fields.SQL()
}

// NestedWheres returns the nested where information
func (q *Query) NestedWheres() (string, error) {
	wheres := WhereSet{}
	q.baseTable.fixFields()
	q.joinTables.fixFields()
	wheres.append(q.baseTable.WhereConditions)
	wheres.append(q.joinTables.wheres()...)
	return wheres.SQL()
}

// NestedTables definition
func (q *Query) NestedTables() (string, error) {
	var (
		tables = ""
		jt     = ""
		err    error
	)
	q.baseTable.fixFields()
	tables += q.baseTable.SQL()

	if jt, err = q.joinTables.SQL(); err != nil {
		return "", err
	}

	tables += jt
	return tables, nil
}

// SQL returns the sql representation of the Query hierarchy of
// objects
func (q *Query) SQL() (string, error) {
	var buf bytes.Buffer
	buf.Grow(300)
	if q == nil {
		return "", fmt.Errorf("Query object is undefined - cannot create a union")
	}

	var (
		nf         = q.NestedFields()
		sql        = delimitSpace("SELECT", nf)
		tables, e1 = q.NestedTables()
		where, e2  = q.NestedWheres()
	)
	if e2 != nil || e1 != nil {
		return "", fmt.Errorf("The following errors were encountered:\n(1)\tTABLES:\t%v\n(2)\tWHERE:\t%v", e1, e2)
	}

	sql = delimitSpace(sql, "FROM", tables)

	if where != "" {
		fmt.Println(where)
		sql = delimitSpace(sql, "WHERE", where)
	}

	if q.Limit != 0 {
		sql = delimitSpace(sql, "LIMIT", strconv.Itoa(q.Limit))
	}

	return sql, nil
}

// SetBaseTable definition
func (q *Query) SetBaseTable(bt Table) {
	one := randomString(4)
	bt.Alias = &one
	bt.Fields.setAlias(&one)
	q.baseTable = bt
}

// SetBaseTableFromProperties sets the base table
func (q *Query) SetBaseTableFromProperties(name, schema string) {
	q.baseTable = Table{
		Name:   name,
		Schema: schema,
	}
}

// AddJoinTables appends JoinTables into the Query object
func (q *Query) AddJoinTables(tables ...JoinTable) {
	for _, table := range tables {
		one := randomString(4)
		table.Alias = &one
		table.Fields.setAlias(&one)
		q.joinTables.append(table)
	}
}

// Union is an abstraction of multiple queries
type Union []Query

// SQL returns the SQL for the given ResultSet, which
// is an object that is composed of tables, fields and other types
func (u *Union) SQL() (string, error) {
	if u == nil {
		return "", fmt.Errorf("Union object is undefined - cannot create a union")
	}

	sql := ""
	for i, query := range *u {
		if i > 0 {
			sql += " UNION ALL "
		}
		q, err := query.SQL()
		if err != nil {
			return "", err
		}
		sql += "(\n" + q + "\n)"
	}
	return sql, nil
}
