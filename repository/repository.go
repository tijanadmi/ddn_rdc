package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tijanadmi/ddn_rdc/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	Authenticate(ctx context.Context, username, testPassword string) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateSession(arg models.CreateSessionParams) (models.Session, error)
	GetSession(id uuid.UUID) (models.Session, error)

	DeleteDDNInterruptionOfDelivery(ctx context.Context, Id string) error
	GetDDNInterruptionOfDeliveryById(ctx context.Context, id int) (*models.DDNInterruptionOfDelivery, error)
	
	GetDDNInterruptionOfDeliveryByPage(ctx context.Context, arg models.ListInterruptionWithPaginationParams) ([]*models.DDNInterruptionOfDelivery, int, error)
	GetAllDDNInterruptionOfDelivery(ctx context.Context, arg models.ListInterruptionParams) ([]*models.DDNInterruptionOfDelivery, int, error)
	InsertDDNInterruptionOfDeliveryP(ctx context.Context, ddnintd models.CreateDDNInterruptionOfDeliveryPParams) (int, error)
	UpdateDDNInterruptionOfDeliveryP(ctx context.Context, id int, version int, ddnintd models.CreateDDNInterruptionOfDeliveryPParams) error

	GetMrcById(ctx context.Context, id int) (*models.SMrc, error)
	GetSMrc(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SMrc, error)
	GetSTipPrekById(ctx context.Context, id int) (*models.STipPrek, error)
	GetSTipPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.STipPrek, error)
	GetSVrPrekById(ctx context.Context, id int) (*models.SVrPrek, error)
	GetSVrPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SVrPrek, error)
	GetSUzrokPrekById(ctx context.Context, id int) (*models.SUzrokPrek, error)
	GetSUzrokPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SUzrokPrek, error)
	GetSPoduzrokPrekById(ctx context.Context, id int) (*models.SPoduzrokPrek, error)
	GetSPoduzrokPrek(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SPoduzrokPrek, error)
	GetSMernaMestaById(ctx context.Context, id int) (*models.SMernaMesta, error)
	GetSMernaMesta(ctx context.Context, arg models.ListLimitOffsetParams) ([]*models.SMernaMesta, error)
	GetObjById(ctx context.Context, id int) (*models.ObjLOV, error)
	GetObjTSRP(ctx context.Context, arg models.ListObjectLimitOffsetParams) ([]*models.ObjLOV, error)
	GetObjHETEVE(ctx context.Context, arg models.ListObjectLimitOffsetParams) ([]*models.ObjLOV, error)
	GetPoljaGE(ctx context.Context, arg models.ListPoljaLimitOffsetParams) ([]*models.PoljaLOV, error)
	GetPoljeGEById(ctx context.Context, id int) (*models.PoljaLOV, error)

	GetPiMMByParams(ctx context.Context, arg models.ListPiMMParams) ([]*models.PiMM, int, error)
	GetPiMMByParamsByPage(ctx context.Context, arg models.ListPiMMParamsByPage) ([]*models.PiMM, int, error)

	GetPiMMT4ByParams(ctx context.Context, arg models.ListPiMMT4Params) ([]*models.PiMMT4, int, error)
	GetPiMMT4ByParamsByPage(ctx context.Context, arg models.ListPiMMT4ParamsByPage) ([]*models.PiMMT4, int, error)

	GetPiDDByParams(ctx context.Context, arg models.ListPiDDParams) ([]*models.PiDD, int, error)
	GetPiDDByParamsByPage(ctx context.Context, arg models.ListPiDDParamsByPage) ([]*models.PiDD, int, error)
	GetPiDDT4ByParams(ctx context.Context, arg models.ListPiDDT4Params) ([]*models.PiMMT4, int, error)
}
