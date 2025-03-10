package repo

import "log/slog"

func (r *millionaireRepo) GetTotalCount(whereClause string, args ...interface{}) (int, error) {
	r.log.Debug("Executing count query", slog.String("query", countQuery+whereClause), slog.Any("args", args))

	query := countQuery + whereClause
	var total int
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		r.log.Error("Failed to get total count", slog.String("query", query), slog.Any("args", args), slog.Any("error", err))
		return 0, err
	}

	r.log.Info("Total count retrieved", slog.Int("total", total))
	return total, nil
}
