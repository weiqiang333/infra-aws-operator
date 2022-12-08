package databases

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBConns struct {
	DBCRUD *gorm.DB
}

func NewDBConns() DBConns {
	return DBConns{
		DBCRUD: nil,
	}
}

func (dbs *DBConns) ConnsMysql(user, password, address, dbName string, connmaxlifetime time.Duration, maxopenconns int) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		password,
		address,
		dbName,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed DBConns Open err: %s", err.Error())
	}
	DBs, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed DBConns Conns err: %s", err.Error())
	}
	DBs.SetConnMaxLifetime(connmaxlifetime)
	DBs.SetMaxOpenConns(maxopenconns)
	dbs.DBCRUD = db
	return nil
}

//func (dbs *DBConns) ConnsSqlite(dbFile string, connmaxlifetime time.Duration, maxopenconns int) error {
//	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//	})
//	if err != nil {
//		return fmt.Errorf("Failed DBConns Open err: %s", err.Error())
//	}
//	DBs, err := db.DB()
//	if err != nil {
//		return fmt.Errorf("Failed DBConns Conns err: %s", err.Error())
//	}
//	DBs.SetConnMaxLifetime(connmaxlifetime)
//	DBs.SetMaxOpenConns(maxopenconns)
//	dbs.DBCRUD = db
//	return nil
//}

func (dbs *DBConns) Close() error {
	sqlDB, err := dbs.DBCRUD.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
