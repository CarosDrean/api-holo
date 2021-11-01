package user

import (
	"database/sql"

	"api-holo/infraestucture/sqlserver"
	sqlutil "api-holo/kit/sqlserver"
	"api-holo/model"
)

const table = "systemuser"

var fields = []string{
	"v_UserName",
	"v_Pasword",
	"i_SystemUserTypeId",
}

const fieldID = "i_SystemUserId"

var (
	psqlInsert = sqlserver.BuildSQLInsert(table, fields, fieldID)
	psqlUpdate = sqlserver.BuildSQLUpdateByID(table, fields, fieldID)
	psqlGetAll = sqlserver.BuildSQLSelect(table, fields, fieldID)
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{db}
}

func (c Service) Create(m *model.User) error {
	stmt, err := c.db.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.UserName,
		m.Password,
		m.Type,
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) Update(m *model.User) error {
	stmt, err := c.db.Prepare(psqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = sqlutil.ExecAffectingOneRow(
		stmt,
		m.UserName,
		m.Password,
		m.Type,
		m.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) GetWhere(filter model.Fields, sort model.SortFields) (model.User, error) {
	conditions, args := sqlserver.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	sorts := sqlserver.BuildSQLOrderBy(sort)
	query += " " + sorts

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return model.User{}, err
	}
	defer stmt.Close()

	return c.scanRow(stmt.QueryRow(args...))
}

func (c Service) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error) {
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

	ms := make(model.Users, 0)
	for rows.Next() {
		m, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (c Service) scanRow(s sqlutil.RowScanner) (model.User, error) {
	m := model.User{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&m.ID,
		&m.UserName,
		&m.Password,
		&m.Type,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}
