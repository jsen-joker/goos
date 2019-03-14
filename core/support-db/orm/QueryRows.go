package orm

import (
	"database/sql"
	"reflect"
	"strconv"
)

type QueryRows struct{
	* sql.Rows
	Values map[string]interface{} //表字段和值的映射
	ColumnNames []string //表字段名集合
}

/*
  获取数据
*/
func (this *QueryRows)Next()bool{
	bResult := this.Rows.Next()
	if bResult{
		//获取表字段名称集合
		if this.ColumnNames == nil || len(this.ColumnNames) == 0{
			this.ColumnNames, _ = this.Rows.Columns()
		}
		//初始化表字段和值的映射
		if this.Values == nil{
			this.Values = make(map[string]interface{})
		}
		//调用scan函数的参数
		scanArgs := make([]interface{}, len(this.ColumnNames))
		//scan函数的值
		values := make([][]byte, len(this.ColumnNames))
		for i := range values{
			scanArgs[i] = &values[i]
		}
		this.Rows.Scan(scanArgs...)
		//将结果存放到Values中
		for i := 0; i < len(this.ColumnNames); i++{
			this.Values[this.ColumnNames[i]] = values[i]
		}
	}
	return bResult
}

/*
  将数据映射到实体切片
  tbname：U对应的数据表名
*/
func (this *QueryRows)To(tbname string) ([]interface{},error) {
	mi := ModelMapping[tbname]
	ti, err := getTableInfo(mi.Entity)
	if err != nil {
		return nil, err
	}
	var models []interface{}
	for this.Next(){
		v := reflect.New(reflect.TypeOf(mi.Entity).Elem()).Elem()
		for k, val := range this.Values{
			f := v.FieldByName(ti.MapperXML[k])
			var strVal string
			if bt, ok := val.([]byte); ok{
				strVal = string(bt)
				switch f.Type().Name(){
				case "int":
					i, _ := strconv.ParseInt(strVal, 10, 64)
					f.SetInt(i)
					break
				case "int64":
					i, _ := strconv.ParseInt(strVal, 10, 64)
					f.SetInt(i)
					break
				case "string":
					f.SetString(strVal)
					break
				}
			}
		}
		models = append(models, v.Interface())
	}
	return models, nil
}