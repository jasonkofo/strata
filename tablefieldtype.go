package strata

// FieldType is the enumerated fieldtype
type FieldType int

const (
	// Nil is a field type that is unknown
	Nil FieldType = iota
	// String is part of the enum for field types
	String
	// Number is part of the enum for field types
	Number
	// Date is part of the enum for field types
	Date
	// Geometry is part of the enum for field types
	Geometry
)

func (ft *FieldType) String() string {
	if ft == nil {
		return ""
	}
	switch *ft {
	case String:
		return "String"
	case Number:
		return "Number"
	case Date:
		return "Date"
	case Geometry:
		return "Geometry"
	case Nil:
		fallthrough
	default:
		return ""
	}
}

func isDate(name string) bool {
	name = cleanString(name)
	return name == "time" || name == "timestamp" || name == "date"
}

func isNumber(name string) bool {
	name = cleanString(name)
	return name == "int" || name == "biginteger" || name == "integer" || name == "number" || name == "num"
}

func isString(name string) bool {
	name = cleanString(name)
	return name == "string" || name == "varchar" || name == "char"
}

func isGeometry(name string) bool {
	name = cleanString(name)
	return name == "geo" || name == "geom" || name == "geometry"
}

// ParseFieldType returns a FieldType
func ParseFieldType(_type string) FieldType {
	if isString(_type) {
		return String
	}

	if isNumber(_type) {
		return Number
	}

	if isDate(_type) {
		return Date
	}

	if isGeometry(_type) {
		return Geometry
	}

	return Nil
}
