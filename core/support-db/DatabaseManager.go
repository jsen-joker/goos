package support_db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsen-joker/goos/core/env"
	"github.com/jsen-joker/goos/core/support-db/orm"
)

var defaultDB *sql.DB

func Init()  {
	CreateDefaultDatabase(env.GoosDatabase)
}
func Close()  {
	defer defaultDB.Close()
}

func CreateDefaultDatabase(url string)  {
	db, err := sql.Open("mysql", url)

	if err != nil {
		panic(err)
		return
	}
	defaultDB = db

	//Query("select username from account where id=?", "1")
}


func SqlInsert(sql string, args ...interface{}) (eff int64, err error)  {
	result, err := defaultDB.Exec(sql, args)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Insert(entity interface{}) (eff int64, err error)  {
	sqlStr, args, tbInfo, err := orm.GenerateInsertSql(entity)
	if err != nil {
		return 0, err
	}
	result, err := defaultDB.Exec(sqlStr, args...)
	if err != nil {
		return 0, err
	}
	_ = orm.SetAuto(result, tbInfo)

	return result.RowsAffected()
}

func SqlQuery(sql string, args ...interface{}) (result *orm.QueryRows, err error) {

	rows, err := defaultDB.Query(sql, args...)

	if err != nil {
		return nil, err
	}
	return &orm.QueryRows{Rows: rows, Values: make(map[string]interface{}),}, nil

	//var models []interface{}
	//modelName := reflect.TypeOf(mType)
	//fmt.Println(modelName)
	//
	//for rows.Next()  {
	//	var id int64
	//	var username string
	//	var password string
	//	if err := rows.Scan(&id, &username, &password); err != nil {
	//		log.Fatal(err)
	//	}
	//	models = append(models, Test{
	//		Id: id,
	//		Username: username,
	//		Password: password,
	//	})
	//}
	//
	//return models

}

func Query(q *orm.QueryWrapper) (result []interface{}, err error) {
	sqlStr :=  q.GetSelectSql()
	fmt.Println("sql: ", sqlStr)
	vals := q.GetValues()
	fmt.Println("params: ", vals)
	myrows, err := SqlQuery(sqlStr, vals...)
	if err != nil{
		return nil, err
	}
	return myrows.To(q.GetFrom())

}


func QueryOne(q *orm.QueryWrapper) (result interface{}, err error) {
	sqlStr :=  q.GetSelectSql()
	fmt.Println("sql: ", sqlStr)
	vals := q.GetValues()
	fmt.Println("params: ", vals)
	myrows, err := SqlQuery(sqlStr, vals...)
	if err != nil{
		return nil, err
	}
	r, e := myrows.To(q.GetFrom())
	if e != nil {
		return nil, e
	}
	if len(r) == 1 {
		return r[0], nil
	} else if len(r) == 0 {
		return nil, nil
	} else {
		return nil, errors.New("more than one result be queried")
	}
}

func Count(q *orm.QueryWrapper) (eff int64, err error) {
	sqlStr := q.GetCountSql()
	fmt.Println("sql: ", sqlStr)
	vals := q.GetValues()
	fmt.Println("params: ", vals)
	if rows, eff := defaultDB.Query(sqlStr, vals...); eff == nil {
		var total int64 = 0
		for rows.Next()  {
			if err := rows.Scan(&total); err != nil {
				return 0, err
			}
		}
		return total, nil
	} else {
		return 0, err
	}
}

func SqlUpdate(sql string, args ...interface{}) (eff int64, err error)  {
	result, err := defaultDB.Exec(sql, args)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Update(entity interface{}) (eff int64, err error)  {
	sqlStr, args, err := orm.GenerateUpdateSql(entity)
	if err != nil {
		return 0, err
	}
	result, err := defaultDB.Exec(sqlStr, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateBy(entity interface{}, where interface{}) (eff int64, err error)  {
	sqlStr, args, err := orm.GenerateUpdateBySql(entity, where)
	if err != nil {
		return 0, err
	}
	result, err := defaultDB.Exec(sqlStr, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func SqlDelete(sql string, args ...interface{}) (eff int64, err error) {
	result, err := defaultDB.Exec(sql, args)

	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Delete(entity interface{}) (eff int64, err error)  {
	sqlStr, args, err := orm.GenerateDeleteSql(entity)
	if err != nil {
		return 0, err
	}
	result, err := defaultDB.Exec(sqlStr, args)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//func Query1(sql string, args ...interface{}) {
//
//	rows, err := defaultDB.Query(sql, args...)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for rows.Next()  {
//		var username string
//		if err := rows.Scan(&username); err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("%s is\n", username)
//	}
//
//	if err := rows.Err(); err != nil {
//		log.Fatal(err)
//	}
//
//}

