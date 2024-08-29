package database

import (
	"fmt"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MyDb struct {
	Db *gorm.DB
}

func initDataBase(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func NewMyDb(dsn string) *MyDb {
	db := &MyDb{}
	db.Db = initDataBase(dsn)
	return db
}

type ReDb struct {
	Db *redis.Client
}

func NewReDb(ip, port string) *ReDb {

	db := ReDb{}
	db.Db = redis.NewClient(&redis.Options{
		Addr:     ip + ":" + port,
		Password: "",
		DB:       0,
	})
	if db.Db == nil {
		panic(fmt.Errorf("open redis failed"))
	}
	return &db
}
