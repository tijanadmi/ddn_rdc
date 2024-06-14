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

	DeleteDDNInterruptionOfDelivery(ctx context.Context,Id string) error
	GetDDNInterruptionOfDeliveryById(ctx context.Context,id int) (*models.DDNInterruptionOfDelivery, error)
	GetDDNInterruptionOfDelivery(ctx context.Context, arg models.ListInterruptionParams) ([]*models.DDNInterruptionOfDelivery, error)
	InsertDDNInterruptionOfDeliveryP(ctx context.Context, ddnintd models.DDNInterruptionOfDelivery) error
	UpdateDDNInterruptionOfDeliveryP(ctx context.Context, ddnintd models.DDNInterruptionOfDelivery) error

	GetMrcById(ctx context.Context,id int) (*models.SMrc, error)
	GetSMrc(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SMrc, error)
	GetSTipPrekById(ctx context.Context,id int) (*models.STipPrek, error)
	GetSTipPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.STipPrek, error)
	GetSVrPrekById(ctx context.Context,id int) (*models.SVrPrek, error) 
	GetSVrPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SVrPrek, error)
	GetSUzrokPrekById(ctx context.Context,id int) (*models.SUzrokPrek, error)
	GetSUzrokPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SUzrokPrek, error)
	GetSPoduzrokPrekById(ctx context.Context,id int) (*models.SPoduzrokPrek, error)
	GetSPoduzrokPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SPoduzrokPrek, error)
	GetSMernaMestaById(ctx context.Context,id int) (*models.SMernaMesta, error) 
	GetSMernaMesta(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SMernaMesta, error)

}
