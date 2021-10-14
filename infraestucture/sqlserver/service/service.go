package service

import (
	"database/sql"

	"api-holo/infraestucture/sqlserver"
	sqlutil "api-holo/kit/sqlserver"
	"api-holo/model"
)

const Table = "service"

var Fields = []string{
	"v_ProtocolId",
	"v_PersonId",
	"i_ServiceStatusId",
	"d_ServiceDate",
	"i_AptitudeStatusId",
	"v_OrganizationId",
}

const fieldID = "v_ServiceId"

var (
	psqlInsert = sqlserver.BuildSQLInsert(Table, Fields, fieldID)
	psqlUpdate = sqlserver.BuildSQLUpdateByID(Table, Fields, fieldID)
	psqlGetAll = sqlserver.BuildSQLSelect(Table, Fields, fieldID)
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{db}
}

func (c Service) Create(m *model.Service) error {
	stmt, err := c.db.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.ProtocolID,
		m.PersonID,
		m.StatusID,
		m.ServiceDate,
		m.OrganizationID,
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) Update(m *model.Service) error {
	stmt, err := c.db.Prepare(psqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = sqlutil.ExecAffectingOneRow(
		stmt,
		m.ProtocolID,
		m.PersonID,
		m.StatusID,
		m.ServiceDate,
		m.OrganizationID,
		m.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) GetWhere(filter model.Fields, sort model.SortFields) (model.Service, error) {
	conditions, args := sqlserver.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	sorts := sqlserver.BuildSQLOrderBy(sort)
	query += " " + sorts

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return model.Service{}, err
	}
	defer stmt.Close()

	return c.scanRow(stmt.QueryRow(args...))
}

func (c Service) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Services, error) {
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

	ms := make(model.Services, 0)
	for rows.Next() {
		m, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (c Service) scanRow(s sqlutil.RowScanner) (model.Service, error) {
	m := model.Service{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&m.ID,
		&m.ProtocolID,
		&m.PersonID,
		&m.StatusID,
		&m.ServiceDate,
		&m.OrganizationID,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}
