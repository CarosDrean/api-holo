package cie10

import (
	"database/sql"

	"api-holo/infraestucture/sqlserver"
	sqlutil "api-holo/kit/sqlserver"
	"api-holo/model"
)

const Table = "cie10"

var Fields = []string{
	"v_CIE10Description1",
	"v_CIE10Description2",
}

const fieldID = "v_CIE10Id"

var (
	psqlInsert = sqlserver.BuildSQLInsert(Table, Fields, fieldID)
	psqlUpdate = sqlserver.BuildSQLUpdateByID(Table, Fields, fieldID)
	psqlGetAll = sqlserver.BuildSQLSelect(Table, Fields, fieldID)
)

type Cie10 struct {
	db *sql.DB
}

func New(db *sql.DB) Cie10 {
	return Cie10{db}
}

func (c Cie10) Create(m *model.Cie10) error {
	stmt, err := c.db.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.Description,
		sqlutil.StringToNull(m.DescriptionTwo),
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c Cie10) Update(m *model.Cie10) error {
	stmt, err := c.db.Prepare(psqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = sqlutil.ExecAffectingOneRow(
		stmt,
		m.Description,
		sqlutil.StringToNull(m.DescriptionTwo),
		m.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c Cie10) GetWhere(filter model.Fields, sort model.SortFields) (model.Cie10, error) {
	conditions, args := sqlserver.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	sorts := sqlserver.BuildSQLOrderBy(sort)
	query += " " + sorts

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return model.Cie10{}, err
	}
	defer stmt.Close()

	return c.scanRow(stmt.QueryRow(args...))
}

func (c Cie10) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Cie10s, error) {
	conditions, args := sqlserver.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	sorts := sqlserver.BuildSQLOrderBy(sort)
	query += " " + sorts

	pagination := sqlserver.BuildSQLPagination(pag)
	query += " " + pagination

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(model.Cie10s, 0)
	for rows.Next() {
		m, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (c Cie10) scanRow(s sqlutil.RowScanner) (model.Cie10, error) {
	m := model.Cie10{}
	descriptionTwoNull := sql.NullString{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&m.ID,
		&m.Description,
		&descriptionTwoNull,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.DescriptionTwo = descriptionTwoNull.String
	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}
