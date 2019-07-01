package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfFac struct {
	Path *string
}

type MysqlConf struct {
	Ip            string
	Port          int
	User          string
	Passwd        string
	DB            string
	Timeout       int
	MaxConnection int
	MaxLifetime   int
}

func (fac MysqlConfFac) ParseConfig() (*MysqlConf, error) {
	data, err := ioutil.ReadFile(*fac.Path)
	if err != nil {
		return nil, err
	}
	config := &MysqlConf{
		Ip:            "127.0.0.1",
		Port:          3306,
		DB:            "cc",
		Timeout:       3000,
		MaxConnection: 300,
		MaxLifetime:   1000,
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type MysqlInstance struct {
	Conf   *MysqlConf
}

func (ins MysqlInstance) NewMysqlInstance() (*sql.DB, error) {
	strConn := "%s:%s@tcp(%s:%d)/%s?autocommit=true&parseTime=true&timeout=%dms&loc=Asia%%2FShanghai&tx_isolation='READ-COMMITTED'"
	url := fmt.Sprintf(strConn, ins.Conf.User, ins.Conf.Passwd,
		ins.Conf.Ip, ins.Conf.Port, ins.Conf.DB, ins.Conf.Timeout)
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(ins.Conf.MaxConnection)
	db.SetMaxIdleConns(ins.Conf.MaxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(ins.Conf.MaxLifetime))

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
