package config

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/rimdesk/product-api/pkg/types"
	"golang.org/x/net/http2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)




type appConfig struct {
}

func (cfg *appConfig) Http2() *http2.Server {
	return &http2.Server{}
}

func (cfg *appConfig) GetServerAddress() string {
	return fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
}

func New() types.GlobalConfig {
	return &appConfig{}
}

func (*appConfig) LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load .env")
	}
	log.Println(".env loaded successfully!")
}


func (*appConfig) DatabaseConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}
}
