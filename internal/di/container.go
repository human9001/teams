package container

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"

	"github.com/human9001/teams/config"
	"github.com/human9001/teams/internal/application/service/comment"
	taskhistory "github.com/human9001/teams/internal/application/service/history"
	"github.com/human9001/teams/internal/application/service/task"
	teamService "github.com/human9001/teams/internal/application/service/team"
	authSrv "github.com/human9001/teams/internal/application/service/user"
	"github.com/human9001/teams/internal/infrastructure/auth"
	"github.com/human9001/teams/internal/infrastructure/cache"
	"github.com/human9001/teams/internal/infrastructure/database/mysql"
	apiv1 "github.com/human9001/teams/internal/interfaces/http/api/v1"
	"github.com/human9001/teams/internal/logger"
	"github.com/human9001/teams/pkg/closer"
)

type Container struct {
	Api            *apiv1.API
	Router         *chi.Mux
	tx             *manager.Manager
	jwt            *auth.JWTService
	teamSrv        *teamService.TeamService
	authSrv        *authSrv.AuthService
	taskCache      *cache.TaskListCache
	taskSrv        *task.TaskService
	taskHistorySrv *taskhistory.TaskHistoryService
	commentSrv     *comment.CommentService
	db             *sqlx.DB
	rDb            *redis.Client
	middleware     []func(next http.Handler) http.Handler
}

func (d *Container) Build() {
	d.initLogger()
	d.initTracer()
	d.initMetrics()
	d.api()
	d.router()
}

// func NewDIContainer() (*Container, error) {
// 	cfg := config.AppConfig()
// 	db, err := mysqlOpen()
// 	if err != nil {
// 		return nil, err
// 	}

// 	txManager := txManager(db)
// 	jwt := auth.NewJWTService(cfg.App.JWTSecret, cfg.App.JWTTTL)
// 	teamServservice := teamService.NewTeamService(mysql.NewTeamRepository(db, txManager))
// 	authServservice := authSrv.NewAuthService(mysql.NewUserRepository(db, txManager), jwt)
// 	taskCache := cache.NewTaskListCache(rdb, cfg.Cache.TTL)
// 	taskService := task.NewTaskService(mysql.NewTaskRepository(db, txManager), mysql.NewMembershipRepository(db), taskCache)
// 	taskHistoryService := taskhistory.NewTaskHistoryService(mysql.NewTaskHistoryRepository(db, txManager))
// 	commentService := comment.NewCommentService(mysql.NewCommentRepository(db, txManager))
// 	api := apiv1.NewAPI(teamServservice, authServservice, taskService, taskHistoryService, commentService)

// 	r := api.NewRouter(cfg.App.JWTSecret, initMetrics())

//		return &Container{Api: api, Router: r}, nil
//	}
func (d *Container) router() *chi.Mux {
	if d.Router == nil {
		cfg := config.AppConfig()
		d.Router = d.api().NewRouter(cfg.App.JWTSecret, d.middleware...)
	}
	return d.Router
}

func (d *Container) api() *apiv1.API {
	if d.Api == nil {
		d.Api = apiv1.NewAPI(d.teamService(), d.authService(), d.taskService(), d.taskHistoryService(), d.commentService())
	}
	return d.Api
}

func (d *Container) jwtService() *auth.JWTService {
	if d.jwt == nil {
		cfg := config.AppConfig()
		d.jwt = auth.NewJWTService(cfg.App.JWTSecret, cfg.App.JWTTTL)

	}
	return d.jwt
}

func (d *Container) teamService() *teamService.TeamService {
	if d.teamSrv == nil {
		d.teamSrv = teamService.NewTeamService(mysql.NewTeamRepository(d.mysqlOpen(), d.txM()))

	}
	return d.teamSrv
}
func (d *Container) authService() *authSrv.AuthService {
	if d.authSrv == nil {
		d.authSrv = authSrv.NewAuthService(mysql.NewUserRepository(d.mysqlOpen(), d.txM()), d.jwtService())
	}
	return d.authSrv
}
func (d *Container) taskService() *task.TaskService {
	if d.taskSrv == nil {
		d.taskSrv = task.NewTaskService(mysql.NewTaskRepository(d.mysqlOpen(), d.txM()), mysql.NewMembershipRepository(d.mysqlOpen()), d.taskCacheService())
	}
	return d.taskSrv
}
func (d *Container) taskHistoryService() *taskhistory.TaskHistoryService {
	if d.taskHistorySrv == nil {
		d.taskHistorySrv = taskhistory.NewTaskHistoryService(mysql.NewTaskHistoryRepository(d.mysqlOpen(), d.txM()))

	}
	return d.taskHistorySrv
}
func (d *Container) commentService() *comment.CommentService {
	if d.commentSrv == nil {
		d.commentSrv = comment.NewCommentService(mysql.NewCommentRepository(d.mysqlOpen(), d.txM()))
	}
	return d.commentSrv
}
func (d *Container) taskCacheService() *cache.TaskListCache {
	if d.taskCache == nil {
		cfg := config.AppConfig()
		c := cache.NewTaskListCache(d.redisDb(), cfg.Cache.TTL)
		d.taskCache = c
	}
	return d.taskCache
}

func (d *Container) redisDb() *redis.Client {
	if d.rDb == nil {
		cfg := config.AppConfig()
		r := cache.NewRedisTaskClient(cfg.Cache.Addr, cfg.Cache.Password, cfg.Cache.Db)
		d.rDb = r
	}
	return d.rDb
}

func (d *Container) mysqlOpen() *sqlx.DB {
	if d.db == nil {
		cfg := config.AppConfig()
		db, err := sqlx.Open("mysql", cfg.MySQL.DSN())
		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)
		db.SetConnMaxIdleTime(cfg.MySQL.ConnMaxIdleTime)

		closer.Add("Mysql", func(_ context.Context) error {
			err := db.Close()
			if err != nil {
				slog.Error("DB close", "error", err)
			}
			return nil
		})

		if err := db.Ping(); err != nil {
			panic(err)
		}

		d.db = db
	}
	return d.db

}

func (d *Container) txM() mysql.TxManager {
	if d.tx == nil {
		tx, err := manager.New(trmsqlx.NewDefaultFactory(d.mysqlOpen()))
		if err != nil {
			slog.Error("ошибка создания txManager", "error", err)
			os.Exit(1)
		}
		d.tx = tx
	}

	return d.tx
}

func (d *Container) initLogger() {

	ctx := context.Background()

	lp, err := logger.InitLogger(ctx)
	if err != nil {
		panic(err)
	}

	closer.Add("Otel Logger", func(_ context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return lp.Shutdown(shutdownCtx)
	})

	otelLogger := slog.New(otelslog.NewHandler("teams-service", otelslog.WithLoggerProvider(lp)))
	slog.SetDefault(otelLogger)

}

func (d *Container) initTracer() {
	ctx := context.Background()
	tp, err := logger.InitTracer(ctx)
	if err != nil {
		panic(err)
	}

	closer.Add("Otel Tracer", func(_ context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return tp.Shutdown(shutdownCtx)
	})
}

func (d *Container) initMetrics() {
	ctx := context.Background()

	mp, err := logger.InitMetrics(ctx)
	if err != nil {
		panic(err)
	}

	closer.Add("Otel Metrics", func(_ context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mp.Shutdown(shutdownCtx)
	})

	meter := otel.Meter("teams-service")
	requestsTotal, _ := meter.Int64Counter("http_requests_total")
	requestDuration, _ := meter.Float64Histogram("http_request_duration_ms")

	d.addMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			fmt.Println("measuring time ...")
			next.ServeHTTP(w, r)
			requestsTotal.Add(r.Context(), 1)
			requestDuration.Record(r.Context(), float64(time.Since(start).Milliseconds()))
		})
	})
}

func (d *Container) addMiddleware(f func(next http.Handler) http.Handler) {
	d.middleware = append(d.middleware, f)
}
