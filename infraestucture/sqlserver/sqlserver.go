package sqlserver

import (
	"database/sql"
	"fmt"
	"strings"

	"api-holo/model"
)

// BuildSQLInsert builds a query INSERT of postgres
func BuildSQLInsert(table string, fields []string, fieldID string) string {
	var args, vals string

	for k, v := range fields {
		args += fmt.Sprintf("%s,", v)
		vals += fmt.Sprintf("@p%d,", k+1)
	}

	if len(fields) > 0 {
		args = args[:len(args)-1]
		vals = vals[:len(vals)-1]
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s, d_InsertDate", table, args, vals, fieldID)
}

// BuildSQLUpdateByID builds a query UPDATE of postgres
func BuildSQLUpdateByID(table string, fields []string, fieldID string) string {
	if len(fields) == 0 {
		return ""
	}

	var args string
	for k, v := range fields {
		args += fmt.Sprintf("%s = @p%d, ", v, k+1)
	}

	return fmt.Sprintf("UPDATE %s SET %sd_UpdateDate = GETDATE() WHERE %s = @p%d", table, args, fieldID, len(fields)+1)
}

// BuildSQLSelect builds a query SELECT of postgres
func BuildSQLSelect(table string, fields []string, fieldID string) string {
	if len(fields) == 0 {
		return ""
	}

	var args string
	for _, v := range fields {
		args += fmt.Sprintf("%s, ", v)
	}

	return fmt.Sprintf("SELECT %s, %sd_InsertDate, d_UpdateDate FROM %s", fieldID, args, table)
}

// BuildSQLSelectFields builds a query SELECT of postgres
func BuildSQLSelectFields(table string, fields []string) string {
	if len(fields) == 0 {
		return ""
	}

	var args string
	for _, v := range fields {
		args += fmt.Sprintf("%s, ", v)
	}

	return fmt.Sprintf("SELECT %s FROM %s", args[:len(args)-2], table)
}

// BuildSQLUpdateBy builds a query UPDATE by of postgres
func BuildSQLUpdateBy(table string, fields []string, by string) string {
	if len(fields) == 0 {
		return ""
	}

	var args string
	for k, v := range fields {
		args += fmt.Sprintf("%s = @p%d, ", v, k+1)
	}

	return fmt.Sprintf("UPDATE %s SET %sd_UpdateDate = now() WHERE %s = @p%d", table, args, by, len(fields)+1)
}

// BuildSQLWhere builds and returns a query WHERE of postgres and its arguments
func BuildSQLWhere(fields model.Fields) (string, []interface{}) {
	if fields.IsEmpty() {
		return "", []interface{}{}
	}

	query, length := "WHERE", len(fields)
	lastFieldIndex := length - 1
	nGroups := 0
	args := []interface{}{}

	paramSequence := 1
	for key, field := range fields {
		setChainingField(&field)
		setOperatorField(&field)
		setAliases(&field)
		setGroupOpen(&field)

		if field.GroupOpen {
			nGroups++
		}

		switch field.Operator {
		case model.In:
			query = fmt.Sprintf("%s %s", query, BuildIN(field))
		case model.IsNull, model.IsNotNull:
			query = fmt.Sprintf("%s %s %s",
				query,
				strings.ToLower(field.Name),
				field.Operator,
			)
		default:
			query = fmt.Sprintf("%s %s %s @p%d",
				query,
				strings.ToLower(field.Name),
				field.Operator,
				paramSequence,
			)
		}

		// Close the group
		if (nGroups > 0) && field.GroupClose {
			nGroups--
			query += ")"
		}

		// if exists still groups open, close them in the last field
		if (nGroups > 0) && (key == lastFieldIndex) {
			query += strings.Repeat(")", nGroups)
		}

		// Add chainingKey (OR, AND) except in the last field
		if key != lastFieldIndex {
			query = fmt.Sprintf("%s %s", query, field.ChainingKey)
		}

		// Add arguments of the parameters except "IN" operator
		if field.Operator != model.In && field.Operator != model.IsNull &&
			field.Operator != model.IsNotNull {
			args = append(args, sql.Named(fmt.Sprintf("p%d", paramSequence), field.Value))
			paramSequence++
		}
	}

	return query, args
}

// BuildSQLOrderBy builds and returns a query ORDER BY of postgres and its arguments
func BuildSQLOrderBy(sorts model.SortFields) string {
	if sorts.IsEmpty() {
		return ""
	}

	query, length := "ORDER BY", len(sorts)

	for key, sort := range sorts {
		setSortFieldOrder(&sort)
		query = fmt.Sprintf("%s %s %s",
			query,
			strings.ToLower(sort.Name),
			sort.Order,
		)
		if key != (length - 1) {
			query = fmt.Sprintf("%s,", query)
		}
	}

	return query
}

// BuildSQLPagination builds and returns a query OFFSET LIMIT of postgres for pagination
func BuildSQLPagination(pag model.Pagination) string {
	if pag.Limit == 0 && pag.Page == 0 {
		return ""
	}
	if pag.MaxLimit == 0 {
		pag.MaxLimit = 20
	}

	if pag.Limit == 0 || pag.Limit > pag.MaxLimit {
		pag.Limit = pag.MaxLimit
	}
	if pag.Page == 0 {
		pag.Page = 1
	}

	offset := pag.Page*pag.Limit - pag.Limit

	pagination := fmt.Sprintf("LIMIT %d OFFSET %d", pag.Limit, offset)

	return pagination
}

// ColumnsAliased return the column names with aliased of the table
func ColumnsAliased(fields []string, aliased string) string {
	if len(fields) == 0 {
		return ""
	}
	columns := ""
	for _, v := range fields {
		columns += fmt.Sprintf("%s.%s, ", aliased, v)
	}

	return fmt.Sprintf("%s.id, %s%s.d_InsertDate, %s.d_UpdateDate",
		aliased, columns, aliased, aliased)
}

// ColumnsAliasedWithDefault return the column names with aliased of the table
func ColumnsAliasedWithDefault(fields []string, aliased string) string {
	if len(fields) == 0 {
		return ""
	}
	columns := ""
	for _, v := range fields {
		columns += fmt.Sprintf("%s.%s, ", aliased, v)
	}

	return fmt.Sprintf("%s.id, %s%s.d_InsertDate, %s.d_UpdateDate",
		aliased, columns, aliased, aliased)
}

func CheckExpected(err error, expected error) error {
	if err != nil {
		if strings.Contains(err.Error(), "psql: expected") {
			return expected
		}
	}

	return err
}

func BuildIN(field model.Field) string {
	nameField := strings.ToLower(field.Name)
	// if the IN failed, return mistakeIN for not select nothing in the field
	mistakeIN := fmt.Sprintf("%s = 0", nameField)

	var args string
	switch items := field.Value.(type) {
	case []uint:
		if len(items) == 0 {
			return mistakeIN
		}

		for _, item := range items {
			args += fmt.Sprintf("%d,", item)
		}

		return fmt.Sprintf("%s IN (%s)", nameField, strings.TrimSuffix(args, ","))
	case []int:
		if len(items) == 0 {
			return mistakeIN
		}

		for _, item := range items {
			args += fmt.Sprintf("%d,", item)
		}

		return fmt.Sprintf("%s IN (%s)", nameField, strings.TrimSuffix(args, ","))
	case []string:
		if len(items) == 0 {
			return mistakeIN
		}

		for _, item := range items {
			args += fmt.Sprintf("'%s',", item)
		}

		return fmt.Sprintf("%s IN (%s)", nameField, strings.TrimSuffix(args, ","))
	default:
		return mistakeIN
	}
}

func setChainingField(field *model.Field) {
	if field.ChainingKey == "" {
		field.ChainingKey = model.And
	}
}

func setOperatorField(field *model.Field) {
	if field.Operator == "" {
		field.Operator = model.Equals
	}
}

func setAliases(field *model.Field) {
	if field.Source != "" {
		field.Name = fmt.Sprintf("%s.%s", field.Source, field.Name)
	}
}

func setGroupOpen(field *model.Field) {
	if field.GroupOpen {
		field.Name = fmt.Sprintf("(%s", field.Name)
	}
}

func setSortFieldOrder(sortField *model.SortField) {
	if sortField.Order == "" {
		sortField.Order = model.Asc
	}
}
