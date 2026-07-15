package container

import (
	"context"
	"log/slog"
	"os"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

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
	"github.com/human9001/teams/pkg/closer"
)

type Container struct {
	Api    *apiv1.API
	Router *chi.Mux
}

func NewDIContainer() (*Container, error) {
	cfg := config.AppConfig()
	db, err := mysqlOpen()
	if err != nil {
		return nil, err
	}
	rdb := cache.NewRedisTaskClient(cfg.Cache.Addr, cfg.Cache.Password, cfg.Cache.Db)

	txManager := txManager(db)
	jwt := auth.NewJWTService(cfg.App.JWTSecret, cfg.App.JWTTTL)
	teamServservice := teamService.NewTeamService(mysql.NewTeamRepository(db, txManager))
	authServservice := authSrv.NewAuthService(mysql.NewUserRepository(db, txManager), jwt)
	taskCache := cache.NewTaskListCache(rdb, cfg.Cache.TTL)
	taskService := task.NewTaskService(mysql.NewTaskRepository(db, txManager), mysql.NewMembershipRepository(db), taskCache)
	taskHistoryService := taskhistory.NewTaskHistoryService(mysql.NewTaskHistoryRepository(db, txManager))
	commentService := comment.NewCommentService(mysql.NewCommentRepository(db, txManager))
	api := apiv1.NewAPI(teamServservice, authServservice, taskService, taskHistoryService, commentService)
	r := api.NewRouter(cfg.App.JWTSecret)

	return &Container{Api: api, Router: r}, nil
}

func mysqlOpen() (*sqlx.DB, error) {
	cfg := config.AppConfig()
	db, err := sqlx.Open("mysql", cfg.MySQL.DSN())
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return db, nil
}

func txManager(db *sqlx.DB) mysql.TxManager {
	txManager, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		slog.Error("ошибка создания txManager", "error", err)
		os.Exit(1)
	}

	return txManager
}
