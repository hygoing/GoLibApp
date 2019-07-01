package db

import (
	"fmt"
	"github.com/google/uuid"
	"jd.com/cc/jstack-cc-server/model"
	db2 "jd.com/jstack-common/db"
)

type PortBindingDao struct {
}

var PortBindingFields = []string{"id", "host_id", "port_id"}

func (dao *PortBindingDao) Create(do db2.DbOperator, portBinding *model.PortBinding) error {
	strSql := "insert into port_binding(id, port_id, host_id) values(?, ?, ?)"
	values := []interface{}{portBinding.Id, portBinding.PortId, portBinding.HostId}

	_, err := db2.Cast2DbOperator(do).Exec(strSql, values...)
	if err != nil {
		return err
	}
	return nil
}

func (dao *PortBindingDao) Delete(do db2.DbOperator, id string) error {
	strSql := "delete from port_binding where id=?"
	_, err := db2.Cast2DbOperator(do).Exec(strSql, id)
	switch {
	case err != nil:
		return err
	}
	return nil
}


func (dao *PortBindingDao) List(do db2.DbOperator, params map[string]interface{}) ([]*model.PortBinding, error) {
	var portBindings []*model.PortBinding
	var keyMap = map[string]string{"PortId": "port_id", "HostId": "host_id"}
	strSql := "select id, port_id, host_id from port_binding"
	var values []interface{}
	var first bool = true
	for key, value := range params {
		if dbKey, ok := keyMap[key]; ok {
			if first {
				strSql += " where " + dbKey + "=?"
				first = false
			} else {
				strSql += " and " + dbKey + "=?"
			}
			values = append(values, value)
		}
	}

	rows, err := db2.Cast2DbOperator(do).Query(strSql, values...)
	if err != nil {
		return portBindings, err
	}

	defer rows.Close()
	err = nil
	for rows.Next() {
		var portBinding = &model.PortBinding{}
		err := rows.Scan(&(portBinding.Id),
			&(portBinding.PortId),
			&(portBinding.HostId))
		if err != nil {
		} else {
			portBindings = append(portBindings, portBinding)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return portBindings, nil
}

func (dao *PortBindingDao) UpdatePortBinding(do db2.DbOperator, p string, hostId string) error {
	strSql := "update port_binding set host_id=? where port_id=?"
	values := []interface{}{hostId, p}

	_, err := db2.Cast2DbOperator(do).Exec(strSql, values...)
	if err != nil {
		return err
	}
	return nil
}

func (dao *PortBindingDao) UpdateHostIdById(do db2.DbOperator, ids []string, hostId string) error {
	strSql := "update port_binding set host_id=? where id in (%s)"
	var idStr string
	values := []interface{}{hostId}
	for idx, v := range ids {
		values = append(values, v)
		if idx == 0 {
			idStr += "?"
		} else {
			idStr += ",?"
		}
	}
	strSql = fmt.Sprintf(strSql, idStr)
	_, err := db2.Cast2DbOperator(do).Exec(strSql, values...)
	if err != nil {
		return err
	}
	return nil
}

func (dao *PortBindingDao) BindingPort(do db2.DbOperator, id, host string) error {
	var portBinding = &model.PortBinding{PortId: id, HostId: host}
	portBinding.Id = uuid.New().String()
	return dao.Create(do, portBinding)
}

func (dao *PortBindingDao) UnbindPortWithAllHost(do db2.DbOperator, p string) error {
	strSql := "delete from `port_binding` where port_id=?"
	values := []interface{}{p}

	_, err := db2.Cast2DbOperator(do).Exec(strSql, values...)
	if err != nil {
		return err
	}
	return nil
}
