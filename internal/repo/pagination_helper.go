package repo

func (r *MillionaireRepo) GetTotalCount(whereClause string, args ...interface{}) (int, error) {
	query := countQuery + whereClause
	var total int
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}
