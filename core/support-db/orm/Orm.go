package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type EntityInfo struct {
	TableInfo
	TableName string
	Entity interface{}
}

type TableInfo struct {
	Name string
	Fields []FieldInfo
	// table fields as key,entity fields as value
	MapperXML map[string] string
}
type FieldInfo struct {
	Name string
	IsPrimaryKey bool
	IsAutoGenerate bool
	Val reflect.Value
}

//表名
type TableName string
//表名类型
var typeTableName TableName
var tableNameType reflect.Type = reflect.TypeOf(typeTableName)
var ModelMapping = make(map[string] EntityInfo)

func Register(entity interface{}) {
	tbInfo, _ := getTableInfo(entity)
	ModelMapping[tbInfo.Name] = EntityInfo{TableName:tbInfo.Name, Entity:entity}
}

func getTableInfo(entity interface{}) (tableInfo *TableInfo, err error) {
	defer func() {
		if e := recover(); err != nil {
			tableInfo = nil
			err = e.(error)
		}
	}()

	if entity == nil {
		return nil, errors.New("entity is nil")
	}

	err = nil
	tableInfo = &TableInfo{
		MapperXML: make(map[string] string),
	}
	rt := reflect.TypeOf(entity)
	rv := reflect.ValueOf(entity)

	tableInfo.Name = rt.Name()
	fmt.Println(rt.Name())
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	for i, j := 0, rt.NumField(); i < j; i++ {
		rtf := rt.Field(i)
		rvf := rv.Field(i)

		if rtf.Type == tableNameType {
			tableInfo.Name = string(rtf.Tag.Get("name"))
			continue
		}
		if rtf.Tag == "-"{
			continue
		}
		//解析字段的tag
		var f FieldInfo
		//没有tag,表字段名与实体字段ing一致
		if rtf.Tag == ""{
			f = FieldInfo{Name:rtf.Name, IsAutoGenerate:false, IsPrimaryKey:false, Val:rvf}
			tableInfo.MapperXML[rtf.Name] = rtf.Name
		} else {
			strTag := string(rtf.Tag)
			if strings.Index(strTag, ":") == -1{
				//tag中没有":"时，表字段名与实体字段ing一致
				f = FieldInfo{Name:rtf.Name, IsAutoGenerate:false, IsPrimaryKey:false, Val:rvf}
				tableInfo.MapperXML[rtf.Name] = rtf.Name
			}else{
				//解析tag中的name值为表字段名
				strName := rtf.Tag.Get("name")
				if strName == ""{
					strName = rtf.Name
				}
				//解析tag中的PK
				isPk := false
				strIspk := rtf.Tag.Get("PK")
				if strIspk == "true"{
					isPk = true
				}
				//解析tag中的auto
				isAuto := false
				strIsauto := rtf.Tag.Get("auto")
				if strIsauto == "true"{
					isAuto = true
				}
				f = FieldInfo{Name:strName, IsPrimaryKey:isPk, IsAutoGenerate:isAuto, Val:rvf}
				tableInfo.MapperXML[strName] = rtf.Name
			}
		}
		tableInfo.Fields = append(tableInfo.Fields, f)
	}
	return
}


/*
  根据实体生成插入语句
*/
func GenerateInsertSql(model interface{})(string, []interface{}, *TableInfo, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, nil, err
	}
	if len(tbInfo.Fields) == 0 {
		return "", nil, nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}

	//根据字段信息拼Sql语句，以及参数值
	strSql := "insert into `" + tbInfo.Name + "`"
	strFileds := ""
	strValues := ""
	var params []interface{}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate {
			continue
		}
		strFileds += "`" + v.Name + "`,"
		strValues += "?,"
		params = append(params, v.Val.Interface())
	}
	if strFileds == ""{
		return "", nil, nil, errors.New(tbInfo.Name + "结构体中没有字段，或只有自增字段")
	}
	strFileds = strings.TrimRight(strFileds, ",")
	strValues = strings.TrimRight(strValues, ",")
	strSql += " (" + strFileds + ") values(" + strValues + ")"
	fmt.Println("sql: ",strSql)
	fmt.Println("params: ",params)
	return strSql, params, tbInfo, nil
}

/*
  根据实体生成修改的sql语句
*/
func GenerateUpdateSql(model interface{})(string, []interface{}, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, err
	}
	if len(tbInfo.Fields) == 0 {
		return "", nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}
	//根据字段信息拼Sql语句，以及参数值
	strSql := "update `" + tbInfo.Name + "` set "
	strFileds := ""
	strWhere := ""
	var p interface{}
	var params []interface{}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate && !v.IsPrimaryKey{
			continue
		}
		if v.IsPrimaryKey{
			strWhere += "`" + v.Name + "`=?"
			p = v.Val.Interface()
			continue
		}
		vv := v.Val.Interface()
		if vv == nil || vv == "" {
			continue
		}
		tp := reflect.TypeOf(vv).String()
		if tp == "string" && vv == "" {
			continue
		} else if tp == "int" && vv.(int) == 0 {
			continue
		} else if tp == "int64" && vv.(int64) == 0 {
			continue
		}
		strFileds += "`" + v.Name + "`=?,"
		params = append(params, vv)
	}
	params = append(params, p)
	strFileds = strings.TrimRight(strFileds, ",")
	strSql += strFileds + " where " + strWhere
	fmt.Println("update sql: ", strSql)
	fmt.Println("update params: ", params)
	return strSql, params, nil
}

func GenerateUpdateBySql(model interface{}, where interface{})(string, []interface{}, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, err
	}
	whereTbInfo, err := getTableInfo(where)
	if err != nil{
		return "", nil, err
	}
	if len(tbInfo.Fields) == 0 {
		return "", nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}
	if len(whereTbInfo.Fields) == 0 {
		return "", nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}
	//根据字段信息拼Sql语句，以及参数值
	strSql := "update `" + tbInfo.Name + "` set "
	strFileds := ""
	strWhere := ""
	var params []interface{}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate && !v.IsPrimaryKey{
			continue
		}
		if v.IsPrimaryKey{
			continue
		}
		vv := v.Val.Interface()
		if vv == nil || vv == "" {
			continue
		}
		tp := reflect.TypeOf(vv).String()
		if tp == "string" && vv == "" {
			continue
		} else if tp == "int" && vv.(int) == 0 {
			continue
		} else if tp == "int64" && vv.(int64) == 0 {
			continue
		}
		strFileds += "`" + v.Name + "`=?,"
		params = append(params, vv)
	}

	var p interface{}

	for _, v := range whereTbInfo.Fields {
		if v.IsAutoGenerate && !v.IsPrimaryKey{
			continue
		}
		vv := v.Val.Interface()
		if vv == nil || vv == "" {
			continue
		}
		tp := reflect.TypeOf(vv).String()
		if tp == "string" && vv == "" {
			continue
		} else if tp == "int" && vv.(int) == 0 {
			continue
		} else if tp == "int64" && vv.(int64) == 0 {
			continue
		}

		strWhere += "`" + v.Name + "`=?"
		p = vv
	}


	params = append(params, p)
	strFileds = strings.TrimRight(strFileds, ",")
	strSql += strFileds + " where " + strWhere
	fmt.Println("update sql: ", strSql)
	fmt.Println("update params: ", params)
	return strSql, params, nil
}

/*
  自动生成删除的sql语句，以主键为删除条件
*/
func GenerateDeleteSql(model interface{})(string, []interface{}, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, err
	}
	//根据字段信息拼Sql语句，以及参数值
	strSql := "delete from " + tbInfo.Name + " where "
	var idVal interface{}
	for _, v := range tbInfo.Fields{
		if v.IsPrimaryKey{
			strSql += v.Name + "=?"
			idVal = v.Val.Interface()
		}
	}
	params := []interface{}{idVal}
	fmt.Println("update sql: ", strSql)
	fmt.Println("update params: ", params)
	return strSql, params, nil
}

/*
  设置自增长字段的值
*/
func SetAuto(result sql.Result, tbInfo *TableInfo)(err error){
	defer func(){
		if e := recover(); e != nil{
			err = e.(error)
		}
	}()
	id, err := result.LastInsertId()
	if id == 0{
		return
	}
	if err != nil{
		return
	}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate && v.Val.CanSet(){
			v.Val.SetInt(id)
			break
		}
	}
	return
}

