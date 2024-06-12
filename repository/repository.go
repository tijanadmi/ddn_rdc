package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/tijanadmi/ddn_rdc/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	Authenticate(username, testPassword string) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateSession( arg models.CreateSessionParams) (models.Session, error)
	GetSession( id uuid.UUID) (models.Session, error)
}