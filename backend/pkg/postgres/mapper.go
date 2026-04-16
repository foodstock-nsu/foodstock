package postgres

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNumericIsInvalid       = errors.New("numeric value is invalid")
	ErrNumericValueIsTooLarge = errors.New("numeric value is too large to store in this type")
)

func Float64ToNumeric(f float64, precision int) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	s := strconv.FormatFloat(f, 'f', precision, 64)

	err := n.Scan(s)
	if err != nil {
		return pgtype.Numeric{}, err
	}

	return n, nil
}

func NumericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0.0, ErrNumericIsInvalid
	}
	f, _ := n.Float64Value()
	return f.Float64, nil
}

func Int64ToNumeric(val int64, exp int32) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(val),
		Exp:   exp,
		Valid: true,
	}
}

func NumericToInt64(n pgtype.Numeric, exp int32) (int64, error) {
	if !n.Valid {
		return 0, ErrNumericIsInvalid
	}

	acc := new(big.Int).Set(n.Int)

	diff := n.Exp - exp
	if diff > 0 {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil)
		acc.Mul(acc, mul)
	} else if diff < 0 {
		div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-diff)), nil)
		acc.Div(acc, div)
	}

	if !acc.IsInt64() {
		return 0, ErrNumericValueIsTooLarge
	}

	return acc.Int64(), nil
}
