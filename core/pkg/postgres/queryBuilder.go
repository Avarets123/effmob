package postgres

import (
	"fmt"
	"iter"
	"reflect"
	"strings"
)

type SqlQueryBuilder struct {
	table, alias, joins, where, selects, orderBy, aggFn, aggAlias string
	limit, offset                                                 int
	args                                                          any
}

func NewSqlBuilder(table, alias string) *SqlQueryBuilder {
	return &SqlQueryBuilder{
		table: table,
		alias: alias,
	}
}

func (builder *SqlQueryBuilder) SetWhere(where string) *SqlQueryBuilder {
	builder.where = where
	return builder
}

func (builder *SqlQueryBuilder) SetJoins(joins string) *SqlQueryBuilder {
	builder.joins = joins
	return builder
}

func (builder *SqlQueryBuilder) SetSelects(selects string) *SqlQueryBuilder {
	builder.selects = selects
	return builder
}

func (builder *SqlQueryBuilder) SetOrderBy(orderBy string) *SqlQueryBuilder {
	builder.orderBy = orderBy
	return builder
}

func (builder *SqlQueryBuilder) SetLimit(limit, offset int) *SqlQueryBuilder {
	builder.limit = limit
	builder.offset = offset
	return builder
}

func (builder *SqlQueryBuilder) SetAggFn(aggFn, aggAlias string) *SqlQueryBuilder {
	builder.aggFn = aggFn
	builder.aggAlias = aggAlias
	return builder
}

func (builder *SqlQueryBuilder) SetArgs(args ...any) *SqlQueryBuilder {
	builder.args = args
	return builder
}

func (builder *SqlQueryBuilder) Build() string {

	query := "SELECT "

	if builder.selects == "" {
		builder.selects = "* "
	}

	if builder.aggFn == "" {
		query += builder.selects
	} else {
		query = fmt.Sprintf("%s %s(*) as %s ", query, builder.aggFn, builder.aggAlias)
	}

	query += " FROM " + builder.table + " as " + builder.alias

	if builder.joins != "" {
		query += " " + builder.joins
	}

	if builder.where != "" {
		query += " WHERE " + builder.where

	}

	if builder.orderBy != "" {
		query += " ORDER BY " + builder.orderBy
	}

	if builder.limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, builder.limit, builder.offset)

	}

	return query

}

func GetInsertSqlFromModel[T any](table string, model T) (string, []any) {

	args := []any{}

	sql := "INSERT INTO " + table

	fields := ""
	values := ""

	t := reflect.TypeOf(model)
	v := reflect.ValueOf(&model).Elem()

	for i := 0; i < t.NumField(); i++ {

		if v.Field(i).IsZero() {
			continue
		}

		field := t.Field(i)
		tagValue := field.Tag.Get("db")

		if tagValue == "" {
			continue
		}

		fields += fmt.Sprintf("\"%s\", ", tagValue)

		data := v.Field(i).Interface()
		values += fmt.Sprintf("$%d, ", len(args)+1)

		args = append(args, data)

	}

	if len(args) == 0 {
		return "", nil
	}

	fields = trimSuffix(fields, ", ")
	values = trimSuffix(values, ", ")

	sql += fmt.Sprintf(" (%s) VALUES (%s)", fields, values)

	return sql, args
}

func GetUpdateSqlFromModel[T any](table, whereField, whereValue string, model T) (string, []any) {

	args := []any{}
	sql := "UPDATE " + table + " SET "
	sets := ""

	t := reflect.TypeOf(model)
	v := reflect.ValueOf(&model).Elem()

	for i := 0; i < t.NumField(); i++ {

		if v.Field(i).IsZero() {
			continue
		}

		field := t.Field(i)
		tagValue := field.Tag.Get("db")

		if tagValue == "" {
			continue
		}

		sets += fmt.Sprintf("\"%s\" = $%d, ", tagValue, len(args)+1)

		data := v.Field(i).Interface()
		args = append(args, data)

	}

	if len(args) == 0 {
		return "", nil
	}

	sets = trimSuffix(sets, ", ")

	sql += fmt.Sprintf("%s WHERE \"%s\" = $%d ", sets, whereField, len(args)+1)

	args = append(args, whereValue)

	return sql, args

}

func GetInsertBatchSqlFromModels[T any](table string, models []T) (string, []any) {

	if len(models) == 0 {
		return "", nil
	}

	args := []any{}
	model := models[0]

	sql := "INSERT INTO " + table

	fields := ""
	values := ""

	t := reflect.TypeOf(model)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("db")
		if tagValue == "" || tagValue == "id" {
			continue
		}
		fields += fmt.Sprintf("\"%s\", ", tagValue)
	}

	for val, params := range getBatchValueFromModel(models) {
		values += fmt.Sprintf("(%s), ", val)
		args = params
	}

	fields = trimSuffix(fields, ", ")
	values = trimSuffix(values, ", ")

	sql += fmt.Sprintf(" (%s) VALUES %s", fields, values)

	fmt.Println(len(args))
	fmt.Println(args)

	return sql, args

}

func getBatchValueFromModel[T any](models []T) iter.Seq2[string, []any] {
	args := []any{}

	return func(yield func(string, []any) bool) {
		for i := 0; i < len(models); i++ {
			model := models[i]
			values := ""
			t := reflect.TypeOf(model)
			v := reflect.ValueOf(&model).Elem()

			for f := 0; f < t.NumField(); f++ {
				field := t.Field(f)

				tagValue := field.Tag.Get("db")
				if tagValue == "" || tagValue == "id" {
					continue
				}

				if v.Field(f).IsZero() {
					args = append(args, nil)
					continue
				}
				data := v.Field(f).Interface()
				args = append(args, data)
				values += fmt.Sprintf("$%d, ", len(args))
			}

			yield(trimSuffix(values, ", "), args)

		}
	}
}

func trimSuffix(str, suff string) string {
	return strings.TrimSuffix(str, suff)
}
