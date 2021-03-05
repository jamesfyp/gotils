package mysqldb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func initDB(dsn string) (db *gorm.DB) {

	if dsn == "" {
		panic("dsn must not empty")

	}

	mysqldb := mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	})

	db, err := gorm.Open(mysqldb, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Warn, // Log level
			},
		),
	})

	rawDb, _ := db.DB()

	rawDb.SetMaxIdleConns(10)
	rawDb.SetMaxOpenConns(200)
	rawDb.SetConnMaxLifetime(3 * 60 * time.Second)

	if err != nil {
		panic("mysql connect error " + err.Error())

	}

	if db.Error != nil {
		panic("database error " + db.Error.Error())

	}

	return db
}
