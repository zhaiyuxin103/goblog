package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"goblog/pkg/logger"
	"time"
)

// DB 数据库对象
var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	initDB()
	createTables()
}

func initDB() {

	var err error

	// 设置数据库连接信息
	config := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池
	/**
	 * DSN：全称 `Data Source Name`，表示数据库连接源，用于定义数据库的连接信息
	 * 不同数据库的 DSN 格式不同，MySQL 的 DSN 格式如下：
	 * [用户名[:密码]@][协议(数据库服务器地址)]]/数据库名称?参数列表
	 * [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	 */
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大链接数
	DB.SetMaxOpenConns(100)
	// 设置最大空闲链接数
	DB.SetMaxIdleConns(25)
	// 设置每个链接当过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
    );`

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
