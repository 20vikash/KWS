package store

import (
	"context"
	"kws/kws/models"
	"kws/kws/models/web"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Auth interface {
		CreateUser(ctx context.Context, first_name, last_name, email, password, user_name string) error
		GenerateToken(ctx context.Context, email string) string
		VerifyUser(ctx context.Context, email string) error
		LoginUser(ctx context.Context, userName, password string) (*models.User, error)
	}

	InMemory interface {
		SetEmailToken(ctx context.Context, email string, token string) error
		GetEmailFromToken(ctx context.Context, token string) string
		DeleteEmailToken(ctx context.Context, token string) error
		PushFreeIp(ctx context.Context, ip int, key string) error
		PopFreeIp(ctx context.Context, key string) (int, error)
		PutDeployResult(ctx context.Context, userName, jobID, password, ip string, success bool, containerIP string) error
		GetDeployResult(ctx context.Context, jobID string) (bool, *web.Instance, error)
		PutStopResult(ctx context.Context, result bool, jobID string) error
		GetStopResult(ctx context.Context, jobID string) (bool, bool, error)
		PutKillResult(ctx context.Context, result bool, jobID string) error
		GetKillResult(ctx context.Context, jobID string) (bool, bool, error)
	}

	Instance interface {
		CreateInstance(ctx context.Context, uid int, userName, insUser, insPassword string) error
		RemoveInstance(ctx context.Context, uid int) error
		StopInstance(ctx context.Context, uid int) error
		StartInstance(ctx context.Context, uid int) error
		Exists(ctx context.Context, uid int) (bool, error)
		GetData(ctx context.Context, uid int) (*web.InsData, error)
	}

	MessageQueue interface {
		PushMessageInstance(ctx context.Context, message *QueueMessage) error
	}

	Wireguard interface {
		AddPeer(ctx context.Context, uid int, wgType *models.WireguardType) error
		RemovePeer(ctx context.Context, pubKey string, uid int) (int, error)
		GetDevices(ctx context.Context, uid int) ([]models.WireguardType, error)
		AllocateNextFreeIP(ctx context.Context, maxHostNumber int, uid int, wgType *models.WireguardType) (int, error)
	}

	PgService interface {
		GetPassword(ctx context.Context, pid int) (string, error)
		GetDatabases(ctx context.Context, pid, uid int) (int, []web.Database, error)
		GetUsers(ctx context.Context, uid int) ([]web.User, error)
		AddUser(ctx context.Context, pgUser *models.PGServiceUser) (int, error)
		AddDatabase(ctx context.Context, pgUser *models.PGServiceUser, pgDatabase *models.PGServiceDatabase) error
		RemoveUser(ctx context.Context, pgUser *models.PGServiceUser) error
		RemoveDatabase(ctx context.Context, pgUser *models.PGServiceUser, pgDatabase *models.PGServiceDatabase) error
	}

	Domains interface {
		AddDomain(ctx context.Context, domain *models.Domain) error
		RemoveDomain(ctx context.Context, domain *models.Domain) error
		GetUserDomains(ctx context.Context, domain *models.Domain) (*[]web.Domain, error)
	}
}

func NewStore(pg *pgxpool.Pool, redis *redis.Client, mq *MQ) *Storage {
	return &Storage{
		Auth: &AuthStore{
			db: pg,
		},
		InMemory: &RedisStore{
			Ds: redis,
		},
		Instance: &InstanceStore{
			Db: pg,
		},
		MessageQueue: &MQ{
			Ch:    mq.Ch,
			Queue: mq.Queue,
		},
		Wireguard: &WireguardStore{
			Con: pg,
		},
		PgService: &PgServiceStore{
			Con: pg,
		},
		Domains: &Domain{
			Con: pg,
		},
	}
}
