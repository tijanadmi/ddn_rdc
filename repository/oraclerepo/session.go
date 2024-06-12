package oraclerepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tijanadmi/ddn_rdc/models"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  :1, :2, :3, :4, :5, :6, :7
) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
`


func (m *OracleDBRepo) CreateSession(arg models.CreateSessionParams) (models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions
WHERE id = :1 LIMIT 1
`

func (m *OracleDBRepo) GetSession( id uuid.UUID) (models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, getSession, id)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}