package repository

import (
	"fmt"
	"strings"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
)

// queryBuilder помогает строить SQL-запросы с параметрами.
type queryBuilder struct {
	conditions []string
	args       []any
	argNum     int
}

func newQueryBuilder() *queryBuilder {
	return &queryBuilder{argNum: 1}
}

// addInCondition добавляет условие IN с множеством значений.
func (qb *queryBuilder) addInCondition(column string, values []any) {
	if len(values) == 0 {
		return
	}

	placeholders := make([]string, len(values))

	for i, v := range values {
		placeholders[i] = fmt.Sprintf("$%d", qb.argNum)
		qb.args = append(qb.args, v)
		qb.argNum++
	}

	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
}

// whereClause возвращает WHERE часть запроса.
func (qb *queryBuilder) whereClause() string {
	if len(qb.conditions) == 0 {
		return ""
	}

	return " WHERE " + strings.Join(qb.conditions, " AND ")
}

// addPagination добавляет LIMIT и OFFSET.
func (qb *queryBuilder) addPagination(pagination *models.Pagination) string {
	if pagination == nil {
		return ""
	}

	clause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", qb.argNum, qb.argNum+1)

	qb.args = append(qb.args, pagination.Limit, pagination.Offset)
	qb.argNum += 2

	return clause
}

func toAnySlice[T any](s []T) []any {
	result := make([]any, len(s))

	for i, v := range s {
		result[i] = v
	}

	return result
}
