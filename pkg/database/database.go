package database

//https://gorm.io/ru_RU/docs/

import (
	"auth/pkg/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var dbase *gorm.DB

// Init - Инициализация базы данных
func Init(cfg *config.Config) (*gorm.DB, error) {
	dbCfg := cfg.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbCfg.Host,
		dbCfg.User,
		dbCfg.Password, dbCfg.Name, dbCfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}

// GetDB - Получение ссылки на экземпляр базы данных
func GetDB() *gorm.DB {
	if dbase == nil {
		cfg := config.GlobalConfig
		dbase, _ = Init(cfg)
		sleep := time.Duration(1)
		for dbase == nil {
			sleep *= 2
			fmt.Printf("Не удалось подключиться к базе данных, повторное подключение через %d секунд", sleep)
			time.Sleep(sleep * time.Second)
			dbase, _ = Init(cfg)
		}
	}
	return dbase
}

func CloseConnection() {
	db := GetDB()
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	fmt.Println("Закрытие соединения с базой данных.")
}
