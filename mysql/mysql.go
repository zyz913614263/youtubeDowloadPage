package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"zyz.com/m/config"
)

var DefaultDB *sql.DB

func InitMysql() {
	if config.DefaultConfig.NouseMysql {
		return
	}
	// 设置数据库连接信息
	username := config.DefaultConfig.MysqlUsername
	password := config.DefaultConfig.MysqlPassword
	host := config.DefaultConfig.MysqlHost
	port := config.DefaultConfig.MysqlPort
	dbname := config.DefaultConfig.MysqlDBName

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	//defer db.Close()

	// 检查数据库是否存在
	err = checkDatabaseExists(db, dbname)
	if err != nil {
		log.Fatalf("Error checking database existence: %v", err)
	}

	log.Println("Database check completed")
	// 更新连接字符串，指定要使用的数据库
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname)
	DefaultDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}
	//defer db.Close()

	// 检查是否存在 user 表
	if err = checkUserTableExists(DefaultDB); err != nil {
		panic(err.Error())
	}
}

// checkUserTableExists 检查是否存在 user 表，如果不存在则创建
func checkUserTableExists(db *sql.DB) error {
	// 查询数据库中是否存在 user 表
	query := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'user'`
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}

	// 如果 user 表不存在，则创建表
	if count == 0 {
		createTableQuery := `
			CREATE TABLE user (
				id INT AUTO_INCREMENT PRIMARY KEY,
				username VARCHAR(255) NOT NULL UNIQUE,
				password VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL UNIQUE,
				is_member TINYINT(1) DEFAULT 0,           -- 是否是会员，0表示非会员，1表示会员
				membership_expiry_date DATE DEFAULT NULL, -- 会员到期时间
				free_trial_remaining INT DEFAULT 0         -- 免费次数
			)
		`
		if _, err := db.Exec(createTableQuery); err != nil {
			return err
		}
		log.Println("User table created successfully")
	} else {
		log.Println("User table already exists")
	}
	return nil
}

// checkDatabaseExists 检查数据库是否存在，如果不存在则创建数据库
func checkDatabaseExists(db *sql.DB, dbname string) error {
	// 查询数据库是否存在
	var dbExist bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)", dbname).Scan(&dbExist)
	if err != nil {
		return err
	}

	// 如果数据库不存在，则创建数据库
	if !dbExist {
		_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
		if err != nil {
			return err
		}
		log.Println("Database created successfully")
	} else {
		log.Println("Database already exists")
	}
	return nil
}
