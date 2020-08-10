package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

var dns string

// Init 初始化 MySQL 链接
func Init() (err error) {
	storageConfig := viper.GetStringMapString("mysql")
	dns = genDNS(storageConfig)

	// 启动时就开打数据库链接
	if err := initEngine(); err != nil {
		return err
	}

	// 测试数据库链接是否正常
	if err := db.Ping(); err != nil {
		return err
	}
	return
}

func genDNS(storageConfig map[string]string) string {
	// "user:password@tcp(host:port)/dbname"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		storageConfig["user"],
		storageConfig["password"],
		storageConfig["host"],
		storageConfig["port"],
		storageConfig["dbname"],
		storageConfig["charset"])
}

func initEngine() (err error) {

	driver := viper.GetString("mysql.driver")
	db, err = sqlx.Connect(driver, dns)
	if err != nil {
		return
	}

	// 设置最大连接数
	viper.SetDefault("mysql.max_conn", 200)
	maxConn := viper.GetInt("mysql.max_conn")
	db.SetMaxOpenConns(maxConn)

	// 设置最大空闲连接数
	viper.SetDefault("mysql.max_idle", 20)
	maxIdle := viper.GetInt("mysql.max_idle")
	db.SetMaxIdleConns(maxIdle)

	return
}

// Close 程序退出时释放 MySQL 链接
// 不直接对外暴露 db 变量, 而是对外暴露一个 Close 方法
func Close() {
	err := db.Close()
	fmt.Println(err)
}

func ClearTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		zap.L().Error("mysql transaction failed, 事务回滚rollback也失败了", zap.Error(err))
	}
}
