package main

import (
	"auth/iternal/server"
	"auth/pkg/config"
	"auth/pkg/database"
	"auth/pkg/database/models"
	"auth/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func preload(cfg *config.Config) error {
	//Миграции
	{
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable", cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password, cfg.Database.Port)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gLog.Default.LogMode(gLog.Silent),
		})
		if err != nil {
			return err
		}

		dbQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name)
		if err := db.Exec(dbQuery).Error; err != nil && !strings.Contains(err.Error(), "already exists") {
			return errors.New(err.Error())
		}

		db, err = database.Init(cfg)
		if err != nil {
			return err
		}

		err = db.AutoMigrate(&models.RefreshToken{})
		if err != nil {
			return errors.New("Ошибка при миграции: " + err.Error())
		}
	}

	return nil
}

func main() {
	//Инициализация конфигурации
	cfg := config.MustLoad()

	if err := logger.New(cfg); err != nil {
		panic(any(fmt.Errorf("Ошибка при инициализации логера: %v\n", err)))
	}

	if err := preload(cfg); err != nil {
		panic(any(err))
	}

	//Инициализация приложения
	srv := server.New(cfg)

	{
		wait := time.Second * 15

		// Запуск сервера в отдельном потоке
		go func() {
			logger.Info("Сервер запущен на адресе: %d", cfg.ServerConfig.Port)
			if err := srv.Listen(fmt.Sprintf(":%d", cfg.ServerConfig.Port)); err != nil {
				logger.Error("Ошибка при прослушивании сервера: %v", err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		_ = srv.ShutdownWithContext(ctx)
		database.CloseConnection()
		os.Exit(0)
	}
}
