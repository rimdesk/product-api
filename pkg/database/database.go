package database

import (
	"fmt"
	"github.com/rimdesk/product-api/pkg/data/entities"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"time"
)

type Connection interface {
	GetEngine() interface{}
	ConnectDB()
	HandleMigration()
	SetConfig(*gorm.Config)
}

type gormDB struct {
	config *gorm.Config
	engine *gorm.DB
}

func (db *gormDB) SetConfig(cfg *gorm.Config) {
	db.config = cfg
}

func (db *gormDB) GetEngine() interface{} {
	return db.engine
}

func NewGormDatabase() Connection {
	return &gormDB{}
}

func (db *gormDB) ConnectDB() {
	dbType := os.Getenv("DB_TYPE")
	log.Println("Loading db config for: ", dbType)
	var err error
	switch dbType {
	case "mysql":
		db.engine, err = db.loadMysqlDB()
		if err != nil {
			log.Fatalln("unable to load sqlite db", err.Error())
		}
	case "sqlite":
		db.engine, err = db.LoadSqliteDB()
		if err != nil {
			log.Fatalln("unable to load sqlite db", err.Error())
		}
	case "postgres":
		log.Println("loading config for postgres")
		db.engine, err = db.loadPostgresDB()
		if err != nil {
			log.Fatalln("unable to load postgres db", err.Error())
		}

		log.Println("postgres connected successfully")
	}

	db.HandleMigration()
}

func (db *gormDB) LoadSqliteDB() (*gorm.DB, error) {
	dbName := fmt.Sprintf("%s.db", os.Getenv("DB_NAME"))
	conn, err := gorm.Open(sqlite.Open(dbName), db.config)

	if err != nil {
		panic("failed to connect database")
		return nil, err
	}
	log.Println("connected to sqlite db")

	return conn, nil
}

func (db *gormDB) loadPostgresDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB.PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	var counts int64
	var backOff = 1 * time.Second
	var connection *gorm.DB
	for {
		c, err := gorm.Open(postgres.Open(dsn), db.config)
		if err != nil {
			fmt.Println("Postgres DB not yet ready to connect!")
			counts++
		} else {
			fmt.Println("Connected to postgres db!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
	}

	return connection, nil
}

func (db *gormDB) loadMysqlDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return nil, err
	}

	return conn, nil
}

func (db *gormDB) HandleMigration() {
	err := db.engine.AutoMigrate(&entities.Product{})
	if err != nil {
		log.Println("migration failed:", err.Error())
	}
	log.Println("tables migrated successfully!")
}
