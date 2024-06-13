package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tijanadmi/ddn_rdc/models"
)


type DatabaseRepo interface {
	Connection() *sql.DB
	Authenticate(ctx context.Context,username, testPassword string) error
	GetUserByUsername(ctx context.Context,username string) (*models.User, error)
	GetUserByID(ctx context.Context,id int) (*models.User, error)
	CreateSession( arg models.CreateSessionParams) (models.Session, error)
	GetSession( id uuid.UUID) (models.Session, error)
}
