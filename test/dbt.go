package main

import (
	db2 "GoLibApp/test/dao"
	"flag"
	"fmt"
	"jd.com/cc/jstack-cc-server/model"
	"jd.com/jstack-common/util"
)

func main() {
	MysqlConfFile := flag.String("dbconf", "/root/db.conf", "Db config file name")
	dbConf, err := MysqlConfFac{Path: MysqlConfFile}.ParseConfig()
	if err != nil {
		fmt.Println("[Init] Parse mysql config error: ", err.Error())
		return
	}
	db, err := MysqlInstance{Conf: dbConf}.NewMysqlInstance()
	if err != nil {
		fmt.Println("[Init] Create Db connection error: ", err.Error())
		return
	}

	portBindingDao := &db2.PortBindingDao{}

	tx, err := db.Begin()
	for i := 0; i < 1; i++ {
		pb := &model.PortBinding{
			Id:     "pb-" + util.Uuid(),
			PortId: "port-" + util.Uuid(),
			HostId: "host-" + util.Uuid(),
		}
		portBindingDao.Create(tx, pb)
	}
	if err := portBindingDao.Delete(tx, "11111"); err != nil {
		fmt.Println(err.Error())
	}
}
