package repo

import "fmt"

func BuildWhereClause(filter MillionaireFilter) (string, []interface{}) {
	var where string
	var args []interface{}
	var conditions []string

	if filter.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("last_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.LastName+"%")
	}
	if filter.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("first_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.FirstName+"%")
	}
	if filter.MiddleName != "" {
		conditions = append(conditions, fmt.Sprintf("middle_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.MiddleName+"%")
	}
	if filter.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.Country+"%")
	}

	if len(conditions) > 0 {
		where = " WHERE " + JoinConditions(conditions, " AND ")
	}
	return where, args
}

func JoinConditions(conditions []string, sep string) string {
	var result string
	for i, c := range conditions {
		if i > 0 {
			result += sep
		}
		result += c
	}
	return result
}
