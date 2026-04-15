package mapper

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func toNumeric(f float64, precision int) pgtype.Numeric {
	var n pgtype.Numeric
	s := strconv.FormatFloat(f, 'f', precision, 64)

	err := n.Scan(s)
	if err != nil {
		return pgtype.Numeric{}
	}

	return n
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
