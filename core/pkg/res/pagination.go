package res

import (
	"effect-mobile/pkg/postgres"
	"effect-mobile/pkg/utils"
	"fmt"
	"net/url"
)

type PaginationParams struct {
	Offset      int    `json:"page"`
	Limit       int    `json:"limit"`
	SortField   string `json:"sortField"`
	SortDir     string `json:"sortDir"`
	FilterField string `json:"filterField"`
	FilterValue string `json:"filterValue"`
}

type PaginationResponse[T any] struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Data   []T `json:"data"`
}

func NewPagResp[T any](pagParams *PaginationParams, data []T, total int) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		Total:  total,
		Offset: pagParams.Offset,
		Limit:  pagParams.Limit,
		Data:   data,
	}
}

func MapQueryToPagParams(query url.Values) *PaginationParams {

	limit := utils.ParseStringToIntOrDefault(query.Get("limit"), 10)
	offset := utils.ParseStringToIntOrDefault(query.Get("offset"), 0)
	sortField := query.Get("sortField")
	sortDir := query.Get("sortDir")
	filterField := query.Get("filterField")
	filterValue := query.Get("filterValue")

	if sortField == "" {
		sortField = "created_at"
		sortDir = "DESC"
	}

	return &PaginationParams{
		Limit:       limit,
		Offset:      offset,
		SortField:   sortField,
		SortDir:     sortDir,
		FilterField: filterField,
		FilterValue: filterValue,
	}

}

func (pagParams *PaginationParams) BuildSqlFromParams(table string, includeDeleted bool) (listingSql, totalSql string, args []any) {

	where := ""

	if !includeDeleted {
		where = "deleted_at IS NULL "

	}

	if pagParams.FilterField != "" && pagParams.FilterValue != "" {
		where += fmt.Sprintf("AND \"%s\" = $1 ", pagParams.FilterField)
		args = append(args, pagParams.FilterValue)
	}

	listingSql = postgres.
		NewSqlBuilder(table, table).
		SetLimit(pagParams.Limit, pagParams.Offset).
		SetOrderBy(fmt.Sprintf("%s %s", pagParams.SortField, pagParams.SortDir)).
		SetWhere(where).
		Build()

	totalSql = postgres.
		NewSqlBuilder(table, table).
		SetWhere(where).
		SetAggFn("COUNT", "total").
		Build()

	return
}
