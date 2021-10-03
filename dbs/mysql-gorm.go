package dbs

import (
	"time"

	"github.com/jinzhu/gorm"
)

//MySQLGORM mysql gorm
type MySQLGORM struct {
	_dba *gorm.DB
}

//Initial initial mysql gorm
func (slf *MySQLGORM) Initial(dsn string, maxConn, maxIdleConn, lifeSec int) error {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	slf._dba = db

	slf._dba.DB().SetMaxOpenConns(maxConn)
	slf._dba.DB().SetMaxIdleConns(maxIdleConn)
	slf._dba.DB().SetConnMaxLifetime(time.Duration(lifeSec) * time.Millisecond)
	if err := slf._dba.DB().Ping(); err != nil {
		return err
	}
	return nil
}

//DB Returns Gorm db
func (slf *MySQLGORM) DB() *gorm.DB {
	/*if err := slf._dba.DB().Ping(); err != nil {
		//TODO: 提供一个重连方法
	}*/

	return slf._dba
}

//Close close gorm db
func (slf *MySQLGORM) Close() {
	slf._dba.Close()
}
