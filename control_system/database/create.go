package database

import (
	"control_system/models"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MyDb struct {
	Db *gorm.DB
}

func initDataBase(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error

	// 重试连接PostgreSQL，最多等待60秒
	maxRetries := 12
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// 测试连接
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					fmt.Printf("Database connection success\n")
					return db
				}
			}
		}

		fmt.Printf("Database connection attempt %d/%d failed: %v, retrying in %v\n", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	fmt.Printf("Failed to connect to database after %d attempts\n", maxRetries)
	panic(err)
}

func NewMyDb(dsn string) *MyDb {
	db := &MyDb{}
	db.Db = initDataBase(dsn)

	// 自动迁移数据库表
	err := db.autoMigrate()
	if err != nil {
		fmt.Printf("Auto migrate failed: %v\n", err)
		panic(err)
	}

	// 插入初始数据
	err = db.insertInitialData()
	if err != nil {
		fmt.Printf("Insert initial data failed: %v\n", err)
	}

	return db
}

// autoMigrate 自动创建数据库表
func (db *MyDb) autoMigrate() error {
	fmt.Println("Starting database auto migration...")

	// 自动迁移所有模型
	err := db.Db.AutoMigrate(
		&models.User{},
		&models.ConsortiumNode{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %v", err)
	}

	fmt.Println("Database auto migration completed successfully")
	return nil
}

// insertInitialData 插入初始数据
func (db *MyDb) insertInitialData() error {
	fmt.Println("Checking and inserting initial data...")

	// 检查是否已存在admin用户
	var count int64
	err := db.Db.Model(&models.User{}).Where("user_name = ?", "admin").Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check admin user: %v", err)
	}

	// 如果不存在admin用户，则创建
	if count == 0 {
		now := time.Now()
		adminUser := models.User{
			UserName:  "admin",
			Password:  "agri_chain",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = db.Db.Create(&adminUser).Error
		if err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}

		fmt.Println("Admin user created successfully")
	} else {
		fmt.Println("Admin user already exists, skipping creation")
	}

	return nil
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
