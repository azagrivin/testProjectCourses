package main

import (
	"context"
	"fmt"
	syslog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azagrivin/testProjectCourses/config"
	"github.com/azagrivin/testProjectCourses/internal/consumer"
	router "github.com/azagrivin/testProjectCourses/internal/delivery/http"
	"github.com/azagrivin/testProjectCourses/internal/logger"
	"github.com/azagrivin/testProjectCourses/internal/repository"
	"github.com/azagrivin/testProjectCourses/internal/services"
	"github.com/azagrivin/testProjectCourses/internal/services/btcusdt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate make --directory=../../ easyjson
func main() {

	syslog.Print("Conf")
	var cfg *config.Config
	{
		cfg = config.Get()
	}

	syslog.Print("Logger")
	var log logger.HttpLogger
	{
		var zapLog *zap.SugaredLogger
		{
			if cfg.App.Debug {
				zapLog = logger.NewDevelopment().Sugar()
			} else {
				zapLog = logger.NewProduction().Sugar()
			}
		}

		log = logger.New(zapLog)
	}

	var err error

	log.Info("DB")
	var db *sqlx.DB
	{
		dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DB.Driver,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
		)

		db, err = sqlx.Open(cfg.Driver, dsn)
		if err != nil {
			log.Errorf("connect to database error, %v", err)
		}
	}

	var (
		repo = repository.NewRepository(db)

		services = &services.Services{
			BtcUsdt: btcusdt.NewService(repo),
		}
	)

	currenciesProcessor := consumer.NewConsumer(time.Second*10, log)
	btcusdtProcessor := consumer.NewConsumer(time.Hour*24, log)

	go func() {
		consumFn := func(ctx context.Context) error {
			return services.BtcUsdt.GetCurrencies(ctx, time.Now())
		}

		if err = currenciesProcessor.CatchAndServe(consumFn); err != nil {
			log.Errorf("consumer error, %v", err)
		}
	}()
	go func() {
		consumFn := func(ctx context.Context) error {
			return services.BtcUsdt.GetBtcusdt(ctx)
		}

		if err = btcusdtProcessor.CatchAndServe(consumFn); err != nil {
			log.Errorf("consumer error, %v", err)
		}
	}()

	go router.NewRouter(cfg, services, log).Run()

	terminateCh := make(chan os.Signal)
	defer close(terminateCh)

	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)
	<-terminateCh

	fmt.Println()
	log.Info("Shutdown")

	ctx, cancel := context.WithCancel(context.Background())
	go currenciesProcessor.Shutdown(ctx)
	go btcusdtProcessor.Shutdown(ctx)

	go func() {
		<-terminateCh
		log.Info("Terminating")
		signal.Stop(terminateCh)
	}()

	cancel()
	time.Sleep(2 * time.Second)
	os.Exit(1)
}
