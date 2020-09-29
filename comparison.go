package strata

// ComparisonType is a type of comparison
type ComparisonType int

const (
	// Exact comparisons will insert equality comparisons
	Exact ComparisonType = iota
	// NotEqual comparisons will insert not equal comparisons
	NotEqual
	// Like comparison will insert LIKE comparator phrase
	Like
	// ILike comparison will insert ILIKE comparator phrase
	ILike
	// IsNotNull is a comparison that is not null
	IsNotNull
	// IsNull is a comparison that is null
	IsNull
	// LTreeSubsists prints out the SQL for the operator
	LTreeSubsists
	// Remember to add changes to function GetComparisonOperator()
)

// SQL operator returns the string representation of the equality type
func (t *ComparisonType) SQL() string {
	if t == nil {
		return "IS NOT NULL"
	}
	switch *t {
	case NotEqual:
		return "<>"
	case ILike:
		return "ILIKE"
	case Like:
		return "LIKE"
	case IsNull:
		return "IS NULL"
	case LTreeSubsists:
		return "@>"
	case Exact:
		return ""
	case IsNotNull:
		fallthrough
	default:
		return "IS NOT NULL"
	}
}

// NeedsRHS indicates whether a RHS field is needed to make a valid comparison
func (t *ComparisonType) NeedsRHS() bool {
	if t == nil {
		return false
	}

	if t == nil {
		return true
	}

	return (*t) != IsNotNull || *t != IsNull
}
