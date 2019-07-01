package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"jd.com/cc/jstack-cc-server/model"
	db2 "jd.com/jstack-common/db"
)

const (
	VniPoolBulkSize = 1000
	VniBulkSize     = 100
	VniUnused       = 0
	VniUsed         = 1
	VniRelease      = 2
)

var ErrIpUnavailable = errors.New("Ip address has been used")
var ErrVniUnavailable = errors.New("Vni has been used")

type DbOperator interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DbScanner interface {
	Scan(dest ...interface{}) error
}

func transferSQL(sql string, params []interface{}) (placeholder string) {
	if len(params) > 0 {
		placeholder = "?" + strings.Repeat(",?", len(params)-1)
		return fmt.Sprintf(sql, placeholder)
	}
	return
}

func clause(fields map[string]string, clause map[string]interface{}, sep string) (*string, *[]interface{}, error) {
	buffer := make([]string, 0, 10)
	values := make([]interface{}, 0, 10)
	for k, v := range clause {
		value, ok := fields[k]
		if !ok {
			return nil, nil, errors.New("fields is illegal")
		}
		buffer = append(buffer, fmt.Sprintf(" %v = ? ", value))
		values = append(values, v)
	}

	sql := strings.Join(buffer, sep)
	return &sql, &values, nil
}

func transferCountSql(table string, filter map[string]interface{}) (string, []interface{}) {
	filterStr, filterFv := tranferSqlFilters(filter)
	strSql := "select count(id) from " + table + filterStr
	return strSql, filterFv
}

func transferListSql(table string, filter map[string]interface{}, field []string, limit int, offset int, order string, od int) (string, []interface{}) {
	orders, ods := []string{}, []int{}
	if order != "" {
		orders, ods = append(orders, order), append(ods, od)
	}
	return transferListSqlWithOrders(table, filter, field, limit, offset, orders, ods)
}

func camelToUnix(s string) string {
	var tmp string
	for i, c := range s {
		if c >= 65 && c <= 90 {
			if i != 0 {
				tmp += "_" + string(c+32)
			} else {
				tmp += string(c + 32)
			}

		} else {
			tmp += string(c)
		}
	}
	return tmp
}

func UpdateDb(do DbOperator, tn string, id string, kv map[string]interface{}) error {
	if len(kv) == 0 {
		return nil
	}
	strSql := "update `" + tn + "`"
	var fk []string
	var fv []interface{}
	for k, v := range kv {
		tmpK := camelToUnix(k)
		switch tmpK {
		case "updated_at":
			fk = append(fk, "`updated_at`=now()")
		case "version":
			fk = append(fk, "`version`=`version`+1")
		default:
			switch v.(type) {
			case float64, int, string:
				fk = append(fk, "`"+tmpK+"`=?")
				fv = append(fv, v)
			default:
				return errors.New("Unsupport value type.")
			}
		}

	}
	strSql += " set " + strings.Join(fk, ",")
	strSql += " where `id`=?"
	fv = append(fv, id)

	_, err := db2.Cast2DbOperator(do).Exec(strSql, fv...)
	return err
}

func listBasicDb(do DbOperator, tn string, filter map[string]interface{}, limit int, offset int, order string, od int) ([]*model.ResourceBasic, error) {
	var ret []*model.ResourceBasic
	strSql, values := transferListSql(tn, filter, []string{"id", "version"}, limit, offset, order, od)
	result, err := db2.Cast2DbOperator(do).Query(strSql, values...)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		r := new(model.ResourceBasic)
		err = result.Scan(&(r.Id), &(r.Version))
		if err != nil {
			return ret, err
		}
		ret = append(ret, r)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

//gen sql with order by list
func transferListSqlWithOrders(table string, filter map[string]interface{}, field []string, limit int, offset int, orders []string, ods []int) (string, []interface{}) {
	fields := []string{}
	for _, f := range field {
		fields = append(fields, "`"+strings.Trim(f, "`")+"`")
	}
	filterStr, filterFv := tranferSqlFilters(filter)
	orderStr := transferSqlOrders(orders, ods)
	limitStr, limitFv := transferSqlLimit(limit, offset)
	strSql := "select " + strings.Join(fields, ",") + " from `" + table + "`" + filterStr + orderStr + limitStr
	return strSql, append(filterFv, limitFv...)
}

func tranferSqlFilters(filter map[string]interface{}) (string, []interface{}) {
	filterStr := ""
	var fk []string
	var fv []interface{}
	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string, bool:
			fk = append(fk, "`"+tmpK+"` = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		filterStr += " where " + strings.Join(fk, " and ")
	}
	return filterStr, fv
}

func transferSqlOrders(orders []string, ods []int) string {
	ordStr := ""
	if len(orders) > 0 && len(orders) == len(ods) {
		ordStr += " order by "
		getRangeStr := func(od int) string {
			if od == 0 {
				return "desc"
			} else if od == 1 {
				return "asc"
			} else {
				return ""
			}
		}
		for index, order := range orders {
			ordStr += fmt.Sprintf("%s %s,", order, getRangeStr(ods[index]))
		}
		ordStr = ordStr[:len(ordStr)-1]
	}
	return ordStr
}

func transferSqlLimit(limit int, offset int) (string, []interface{}) {
	var strSql string
	var fv []interface{}
	if limit >= 0 {
		if offset >= 0 {
			strSql += " LIMIT ?, ?"
			fv = append(fv, offset)
			fv = append(fv, limit)
		} else {
			strSql += " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return strSql, fv
}
