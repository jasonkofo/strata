package strata

import "strings"

// TableField is an abstraction of the field type
type TableField struct {
	Name          string    `json:"name"`          // quoted name
	Alias         *string   `json:"alias"`         // alias for reference within the quer
	FormattedName string    `json:"formattedName"` // unquoted provision for custom names (perhaps using formulas) - i.e. SUBSTRING(\"fieldName\" FROM '[A-Za-z]+_([A-Za-z]+[A-Z.])').
	FriendlyName  string    `json:"friendlyName"`
	Type          FieldType `json:"type"`
}

// TableFields is an array of table fields
type TableFields []TableField

func (tf *TableField) pickFriendlyName() string {
	if tf.FriendlyName != "" {
		return " as " + tf.FriendlyName
	}
	return ""
}

func (tf *TableField) transformsField() bool {
	return containsQuoted(tf.FormattedName, tf.Name)
}

// fixFormattedName tries to detect if the selector name is in the Forma ttedName
// field.
func (tf *TableField) fixFormattedName() string {
	n := tf.Name
	if tf.Name == "" || !strings.Contains(tf.FormattedName, tf.Name) || !tf.transformsField() {
		return tf.FormattedName
	}

	createdField := delimitDot(*tf.Alias, tf.FormattedName)

	return strings.Replace(tf.FormattedName, n, createdField, 1)
}

func (tf *TableField) pickSelectorName() string {
	if tf.FormattedName != "" {
		return tf.fixFormattedName()
	}

	if tf.Alias != nil && *tf.Alias != "" && tf.Name != "" {
		return chainSelector(*tf.Alias, tf.Name)
	}
	return insertDoubleQuotes(tf.Name)
}

func (tf *TableFields) append(field ...TableField) {
	*tf = append(*tf, field...)
}

func (tf *TableFields) fieldByName(name string) *TableField {
	if tf == nil || len(*tf) == 0 {
		return nil
	}

	if idx := tf.getIndex(name); idx != -1 {
		return tf.getField(idx)
	}

	return nil
}

// getField returns the
func (tf *TableFields) getField(i int) *TableField {
	if tf == nil || i < 0 || i > len(*tf)-1 {
		return nil
	}
	return &(*tf)[i]
}

// getIndex returns the index of the first element having
// the given name
func (tf *TableFields) getIndex(name string) int {
	for i, field := range *tf {
		if field.Name == name {
			return i
		}
	}
	return -1
}

// SQL returns the SQL representation of the field
func (tf *TableField) SQL() string {
	if tf == nil {
		return ""
	}
	sql := tf.pickSelectorName()
	if suffix := tf.pickFriendlyName(); suffix != "" {
		sql += " " + suffix
	}
	return sql
}

// ToSlice is a convenient method for returning a slice from
// a single element
func (tf *TableField) ToSlice() TableFields {
	return TableFields{*tf}
}

// IsString returns whether the field type is a string
func (tf *TableField) IsString() bool {
	return tf.Type == String
}

// IsNumber returns whether the field type is a number
func (tf *TableField) IsNumber() bool {
	return tf.Type == String
}

// IsDate returns whether the field type is a date
func (tf *TableField) IsDate() bool {
	return tf.Type == Date
}

// IsGeometry returns whether the field type is a geometry field
func (tf *TableField) IsGeometry() bool {
	return tf.Type == Geometry
}

func (tf *TableFields) addField(alias *string, name, friendlyName, formattedName string, _type FieldType) {
	tf.append(makeField(alias, name, friendlyName, formattedName, _type))
}

func makeField(alias *string, name, friendlyName, formattedName string, _type FieldType) TableField {
	return TableField{
		Name:          name,
		Alias:         alias,
		FriendlyName:  friendlyName,
		FormattedName: formattedName,
		Type:          _type,
	}
}

func field(name string, _type FieldType) TableField {
	return TableField{
		Name: name,
		Type: _type,
	}
}

// Field returns a TableField of parsed type
func Field(name string, _type string) TableField {
	return field(name, ParseFieldType(_type))
}

// StringField returns a TableField of string type
func StringField(name string) TableField {
	return field(name, String)
}

// NumberField returns a TableField of number type
func NumberField(name string) TableField {
	return field(name, Number)
}

// DateField returns a TableField of date type
func DateField(name string) TableField {
	return field(name, Date)
}

// GeometryField returns a TableField of Geometry type
func GeometryField(name string) TableField {
	return field(name, Geometry)
}

func (tf *TableFields) addStringField(alias *string, name, friendlyName, formattedName string) {
	tf.addField(alias, name, friendlyName, formattedName, String)
}

func (tf *TableFields) addNumberField(alias *string, name, friendlyName, formattedName string) {
	tf.addField(alias, name, friendlyName, formattedName, Number)
}

func (tf *TableFields) addDateField(alias *string, name, friendlyName, formattedName string) {
	tf.addField(alias, name, friendlyName, formattedName, Date)
}

func (tf *TableFields) addGeometryField(alias *string, name, friendlyName, formattedName string) {
	tf.addField(alias, name, friendlyName, formattedName, Geometry)
}

// addFieldByProperties adds a field from primitive types
func (tf *TableFields) addFieldByProperties(alias *string, name, formattedName, friendlyName, _type string) {
	tf.addField(alias, name, formattedName, friendlyName, ParseFieldType(_type))
}

func (tf *TableFields) addSimpleStringFields(alias *string, fields ...string) {
	for _, field := range fields {
		tf.addField(alias, field, "", "", String)
	}
}

func (tf *TableFields) addSimpleNumberFields(alias *string, fields ...string) {
	for _, field := range fields {
		tf.addField(alias, field, "", "", Number)
	}
}

func (tf *TableFields) addSimpleDateFields(alias *string, fields ...string) {
	for _, field := range fields {
		tf.addField(alias, field, "", "", Date)
	}
}

func (tf *TableFields) addGeometryFields(alias *string, fields ...string) {
	for _, field := range fields {
		tf.addField(alias, field, "", "", Geometry)
	}
}

func (tf *TableFields) setAlias(alias *string) {
	for i := 0; tf != nil && i < len(*tf); i++ {
		(*tf)[i].Alias = alias
	}
}

// SQL returns the SQL representation of the table fields
func (tf *TableFields) SQL() string {
	sql := ""
	for i, field := range *tf {
		if i > 0 && i < len(*tf) {
			sql += ", "
		}
		sql += field.SQL()
	}

	return sql
}
