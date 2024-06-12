package oraclerepo

import "database/sql"

type OracleDBRepo struct {
	DB *sql.DB
}

func (m *OracleDBRepo) Connection() *sql.DB {
	return m.DB
}