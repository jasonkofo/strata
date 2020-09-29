package strata

import (
	"fmt"
)

// Where is the abstraction of a where condition
type Where struct {
	LHSField       *TableField
	RHSField       interface{}
	ComparisonType ComparisonType
}

// Wheres is a slice of where conditions
type Wheres struct {
	Wheres      []Where
	IsInclusive bool
}

func (w *Where) rightFieldSQL() string {
	if w.RHSField == nil {
		return ""
	}
	switch rhs := w.RHSField.(type) {
	case *TableField:
		return rhs.SQL()
	case string:
		islike := !w.ComparisonType.IsExact()
		return insertStringLiterals(rhs, islike, islike)
	case int:
		return fmt.Sprintf("%v", rhs)
	case int64:
		return fmt.Sprintf("%v", rhs)
	default:
		return ""
	}
}

// WhereSet is a set of Where objects - they will be marshalled into a
// set of larger where conditions using the OR keywords
type WhereSet []Wheres

// SQL returns the SQL for the given where condition
func (w *Where) SQL() (string, error) {
	// Do absolutely nothing if there is no left hand side of the comparison
	if w.LHSField == nil {
		return "", fmt.Errorf("No left hand field object provided")
	}

	rhs := w.rightFieldSQL()
	if w.ComparisonType.NeedsRHS() && rhs == "" {
		return "", fmt.Errorf("No right hand field object provided - is necessary for comparison type")
	}

	return delimitSpace(
		w.LHSField.SQL(),
		w.ComparisonType.SQL(),
		rhs,
	), nil
}

// Append appends a where condition to the where object
func (ws *Wheres) Append(w ...Where) {
	if w != nil {
		ws.Wheres = append(ws.Wheres, w...)
	}
}

// ToNativeSlice returns the underlying golang kind
func (ws *Wheres) ToNativeSlice() []Where {
	if len(ws.Wheres) == 0 {
		return nil
	}

	return []Where(ws.Wheres)
}

// SQL returns the SQL representation of the where conditions
func (ws *Wheres) SQL() (string, error) {
	sql := ""
	for i, w := range ws.Wheres {
		if i > 0 {
			sql += ws.inclusiveSQL()
		}
		s, err := w.SQL()
		if err != nil {
			return "", err
		}
		sql += s
	}
	return sql, nil
}

func (ws *Wheres) inclusiveSQL() string {
	if ws.IsInclusive {
		return "AND"
	}
	return "OR"
}

// add adds where conditions to the Where
func (ws *Wheres) add(wheres ...Where) {
	if ws == nil {
		ws = &Wheres{}
	}
	ws.Append(wheres...)
}

// Append appends a where object to the whereset
func (ws *WhereSet) append(whereSet ...Wheres) {
	*ws = append(*ws, whereSet...)
}

// SQL return an SQL representation of the WhereSet. It assumes that
// the Wheresets will be included with "OR" conditions
func (ws *WhereSet) SQL(whereSet ...Wheres) (string, error) {
	sql := ""
	for i, where := range *ws {
		if where.Wheres == nil || len(where.Wheres) == 0 {
			break
		}
		if i > 0 {
			sql += " OR "
		}
		_w, err := where.SQL()
		if err != nil {
			return "", err
		}
		sql += _w
	}
	return sql, nil
}
