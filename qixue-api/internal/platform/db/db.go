package db

import (
	"fmt"
	"net"
	"strconv"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"qixue-api/internal/platform/config"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	mc := mysqldriver.NewConfig()
	mc.User = cfg.DBUser
	mc.Passwd = cfg.DBPassword
	mc.Net = "tcp"
	mc.Addr = net.JoinHostPort(cfg.DBHost, strconv.Itoa(cfg.DBPort))
	mc.DBName = cfg.DBName
	mc.ParseTime = true
	mc.Loc = loc
	mc.Collation = "utf8mb4_unicode_ci"
	mc.Params = map[string]string{"charset": "utf8mb4"}

	conn, err := gorm.Open(mysql.Open(mc.FormatDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("连接 MySQL: %w", err)
	}
	return conn, nil
}
