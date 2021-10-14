package sqlserver

import (
	"database/sql"
	"errors"
	"time"
)

// ErrExpectedOneRow is used for indicate that expected 1 row
var ErrExpectedOneRow = errors.New("expected 1 row")

// RowScanner utilidad para leer los registros de un Query
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// ExecAffectingOneRow ejecuta una sentencia (statement),
// esperando una sola fila afectada.
func ExecAffectingOneRow(stmt *sql.Stmt, args ...interface{}) error {
	r, err := stmt.Exec(args...)
	if err != nil {
		return err
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return ErrExpectedOneRow
	}

	return nil
}

// TimeToNull devuelve una estructura nil si la fecha está en valor (zero)
func TimeToNull(t time.Time) sql.NullTime {
	r := sql.NullTime{}
	r.Time = t

	if !t.IsZero() {
		r.Valid = true
	}

	return r
}

// ParseDateToTime devuelve una estructura nil si la hora está en valor (zero)
func ParseDateToTime(s string) sql.NullTime {
	format := "15:04:05"
	t, _ := time.Parse(format, s)

	return TimeToNull(t)
}

// Int64ToNull devuelve una estructura nil si el entero es (zero)
func Int64ToNull(i int64) sql.NullInt64 {
	r := sql.NullInt64{}
	r.Int64 = i

	if i > 0 {
		r.Valid = true
	}

	return r
}

// StringToNull devuelve una estructura nil si la cadena de texto está vacia
func StringToNull(s string) sql.NullString {
	r := sql.NullString{}
	r.String = s

	if s != "" {
		r.Valid = true
	}

	return r
}

// FloatToNull devuelve una estructura nil si el float es (zero)
func FloatToNull(i float64) sql.NullFloat64 {
	r := sql.NullFloat64{}
	r.Float64 = i

	if i > 0 {
		r.Valid = true
	}

	return r
}
